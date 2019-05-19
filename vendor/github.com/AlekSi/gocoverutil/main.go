package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/AlekSi/gocoverutil/gocoverutil"
)

var (
	coverprofileF = flag.String("coverprofile", "cover.out", "Output file.")
	ignoreF       = flag.String("ignore", "", "comma-separated list of packages to ignore; may contain '...' patterns")

	mergeFlagSet = flag.NewFlagSet("merge", flag.ExitOnError)

	testFlagSet = flag.NewFlagSet("test", flag.ExitOnError)

	// go build flags, in order of "go build -h"
	aF    = testFlagSet.Bool("a", false, "force rebuilding of packages that are already up-to-date.")
	nF    = testFlagSet.Bool("n", false, "print the commands but do not run them.")
	pF    = testFlagSet.Int("p", 1, "ignored for compatibility with go build")
	raceF = testFlagSet.Bool("race", false, "enable data race detection.")
	msanF = testFlagSet.Bool("msan", false, "enable interoperation with memory sanitizer.")
	workF = testFlagSet.Bool("work", false, "print the name of the temporary work directory and do not delete it when exiting.")
	xF    = testFlagSet.Bool("x", false, "print the commands.")
	tagsF = testFlagSet.String("tags", "", "a list of build tags to consider satisfied during the build.")
	// -v is redefined below

	// test binary flags (without "test." prefix), in order of "./gocoverutil-fizzbuzz.test -h"
	shortF   = testFlagSet.Bool("short", false, "tell long-running tests to shorten their run time.")
	timeoutF = testFlagSet.Duration("timeout", 0, "if a test runs longer than t, panic.")
	vF       = testFlagSet.Bool("v", false, "verbose output: log all tests as they are run.")
	// -coverprofile is defined by main command

	// go test flags, excluding go build and test binary flags, in order of "go test -h"
	covermodeF = testFlagSet.String("covermode", "", "set the mode for coverage analysis for the package[s] being tested.")
	// -coverpkg is set by Test method

	// TODO add more flags
)

func main() {
	log.SetFlags(0)
	mergeFlagSet.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s merge command merges several go coverage profiles into a single file.\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Usage:\n\n")
		fmt.Fprintf(os.Stderr, "  %s [global flags] merge [files]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Example:\n\n")
		fmt.Fprintf(os.Stderr, "  %s -coverprofile=cover.out merge internal/test/package1/package1.out internal/test/package2/package2.out\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Global flags:\n")
		flag.PrintDefaults()
	}
	testFlagSet.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s test command runs go test -cover with correct flags and merges profiles.\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Packages list is passed as arguments; they may contain `...` patterns.\n")
		fmt.Fprintf(os.Stderr, "The list is expanded, sorted and duplicates and ignored packages are removed.\n")
		fmt.Fprintf(os.Stderr, "`go test -coverpkg` flag is set automatically to the same list.\n")
		fmt.Fprintf(os.Stderr, "Only a single package is passed at once to `go test`, so it always acts as if `-p 1` is passed.\n")
		fmt.Fprintf(os.Stderr, "If tests are failing, %s exits with a correct exit code.\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Usage:\n\n")
		fmt.Fprintf(os.Stderr, "  %s [global flags] test [test flags] [packages]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Example:\n\n")
		fmt.Fprintf(os.Stderr, "  %s -coverprofile=cover.out test -v -covermode=count github.com/AlekSi/gocoverutil/internal/test/...\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Global flags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nTest flags:\n")
		testFlagSet.PrintDefaults()
	}
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s contains two commands: merge and test.\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Merge command merges several go coverage profiles into a single file.\n")
		fmt.Fprintf(os.Stderr, "Run `%s merge -h` for usage information.\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Test command runs go test -cover with correct flags and merges profiles.\n")
		fmt.Fprintf(os.Stderr, "Run `%s test -h` for usage information.\n", os.Args[0])
	}
	flag.Parse()

	var err error
	switch flag.Arg(0) {
	case "merge":
		mergeFlagSet.Parse(flag.Args()[1:])
		err = gocoverutil.Merge(mergeFlagSet.Args(), *coverprofileF)

	case "test":
		testFlagSet.Parse(flag.Args()[1:])
		logger := log.New(ioutil.Discard, "", 0)
		if *nF || *xF || *vF {
			logger.SetOutput(os.Stderr)
		}
		var ignore []string
		if len(*ignoreF) > 0 {
			ignore = strings.Split(*ignoreF, ",")
		}
		err = gocoverutil.Test(testFlagSet, ignore, *coverprofileF, logger)

	default:
		flag.Usage()
		log.Fatalf("\nUnexpected command '%s'.", flag.Arg(0))
	}

	if err != nil {
		log.Print(err)
		if eErr, ok := err.(*exec.ExitError); ok {
			if ws, ok := eErr.Sys().(*syscall.WaitStatus); ok {
				os.Exit(ws.ExitStatus())
			}
		}
		os.Exit(1)
	}
}
