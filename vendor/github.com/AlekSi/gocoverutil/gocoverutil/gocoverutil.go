// Package gocoverutil implements merging of go cover profiles and running go test -cover with correct flags.
package gocoverutil

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/tools/cover"
)

type byLines []cover.ProfileBlock

func (bl byLines) Len() int          { return len(bl) }
func (bl byLines) Swap(i int, j int) { bl[i], bl[j] = bl[j], bl[i] }
func (bl byLines) Less(i int, j int) bool {
	return bl[i].StartLine < bl[j].StartLine || bl[i].StartLine == bl[j].StartLine && bl[i].StartCol < bl[j].StartCol
}

// Merge merges several coverage files into single file.
// All input files are fully read before output file is written,
// so it may be one of the input files.
func Merge(inputFiles []string, outputFile string) error {
	blocks := make(map[string][]cover.ProfileBlock)
	var mode string
	for _, f := range inputFiles {
		profiles, err := cover.ParseProfiles(f)
		if err != nil {
			return err
		}
		for _, p := range profiles {
			if mode == "" {
				mode = p.Mode
			}
			if mode != p.Mode {
				return fmt.Errorf("different modes: %s and %s", mode, p.Mode)
			}

			blocks[p.FileName] = append(blocks[p.FileName], p.Blocks...)
		}
	}

	// sort files
	inputFiles = make([]string, 0, len(blocks))
	for file := range blocks {
		inputFiles = append(inputFiles, file)
	}
	sort.Strings(inputFiles)

	for _, file := range inputFiles {
		sort.Sort(byLines(blocks[file]))

		// merge blocks
		var newBlocks []cover.ProfileBlock
		var prev cover.ProfileBlock
		for _, b := range blocks[file] {
			// skip full duplicate
			if prev == b {
				continue
			}

			// change count inside previous block if only count changed
			prev.Count = b.Count
			if prev == b {
				if mode == "set" {
					newBlocks[len(newBlocks)-1].Count = 1
				} else {
					newBlocks[len(newBlocks)-1].Count += b.Count
				}
				prev = b
				continue
			}

			newBlocks = append(newBlocks, b)
			prev = b
		}
		blocks[file] = newBlocks
	}

	f, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err = f.WriteString(fmt.Sprintf("mode: %s\n", mode)); err != nil {
		return err
	}
	for _, file := range inputFiles {
		for _, b := range blocks[file] {
			// encoding/base64/base64.go:34.44,37.40 3 1
			// where the fields are: name.go:line.column,line.column numberOfStatements count
			l := fmt.Sprintf("%s:%d.%d,%d.%d %d %d\n", file, b.StartLine, b.StartCol, b.EndLine, b.EndCol, b.NumStmt, b.Count)
			if _, err = f.WriteString(l); err != nil {
				return err
			}
		}
	}
	return nil
}

// list uses `go list` command to expand packages list to a sorted list without duplicates.
func list(packages []string) ([]string, error) {
	args := append([]string{"list"}, packages...)
	cmd := exec.Command("go", args...)
	cmd.Stderr = os.Stderr
	b, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var res []string
	lines := strings.Split(string(b), "\n")
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l != "" {
			res = append(res, l)
		}
	}
	return res, nil
}

// Test runs `go test -cover` with correct flags for all packages in flagSet, and merges coverage files.
// Returned error may be *exec.ExitError if tests failed.
func Test(flagSet *flag.FlagSet, ignore []string, outputFile string, logger *log.Logger) error {
	args, err := list(flagSet.Args())
	if err != nil {
		return err
	}
	packages := args

	// handle ignore slice only if it is really given to avoid expanding it with `go list` to the same package
	if len(ignore) > 0 {
		ignore, err = list(ignore)
		if err != nil {
			return err
		}
		packages = nil
		for _, a := range args {
			var skip bool
			for _, i := range ignore {
				if a == i {
					skip = true
					break
				}
			}
			if !skip {
				packages = append(packages, a)
			}
		}
	}

	if len(packages) == 0 {
		return fmt.Errorf("nothing to test, all packages are ignored")
	}

	// copy flags from flagSet, add -coverpkg with all packages
	var flags []string
	flagSet.Visit(func(f *flag.Flag) {
		flags = append(flags, fmt.Sprintf("-%s=%s", f.Name, f.Value.String()))
	})
	flags = append(flags, fmt.Sprintf("-coverpkg=%s", strings.Join(packages, ",")))

	// create temporary directory
	f, err := ioutil.TempFile("", "gocoverutil-")
	if err != nil {
		return err
	}
	if err = f.Close(); err != nil {
		return err
	}
	dir := f.Name()
	if err = os.Remove(dir); err != nil {
		return err
	}
	if err = os.Mkdir(dir, 0777); err != nil {
		return err
	}

	files := make([]string, 0, len(packages))
	for _, p := range packages {
		// get temporary file name
		if f, err = ioutil.TempFile(dir, filepath.Base(p)+"-"); err != nil {
			return err
		}
		files = append(files, f.Name())
		if err = f.Close(); err != nil {
			return err
		}

		// run go test with added -coverprofile
		args := append([]string{"test"}, flags...)
		args = append(args, fmt.Sprintf("-coverprofile=%s", f.Name()))
		args = append(args, p)
		cmd := exec.Command("go", args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		logger.Printf(strings.Join(cmd.Args, " "))
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	// merge and remove files
	if err = Merge(files, outputFile); err != nil {
		return err
	}
	return os.RemoveAll(dir)
}
