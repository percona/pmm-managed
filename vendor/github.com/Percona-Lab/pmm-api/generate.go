// +build ignore

package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var directories = []string{"agent", "inventory"}

func main() {
	log.SetFlags(0)
	flag.Parse()

	protoc, err := exec.LookPath("protoc")
	if err != nil {
		log.Fatal(err)
	}

	for _, dir := range directories {
		files, err := filepath.Glob(dir + "/*.pb.go")
		if err != nil {
			log.Fatal(err)
		}
		for _, f := range files {
			if err = os.Remove(f); err != nil {
				log.Fatal(err)
			}
		}

		files, err = filepath.Glob(dir + "/*.proto")
		if err != nil {
			log.Fatal(err)
		}

		args := []string{"--proto_path=vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis"}
		args = append(args, "--proto_path="+dir)
		args = append(args, files...)
		args = append(args, "--go_out=plugins=grpc:"+dir)

		cmd := exec.Command(protoc, args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		log.Print(strings.Join(cmd.Args, " "))
		if err = cmd.Run(); err != nil {
			log.Fatal(err)
		}
	}
}
