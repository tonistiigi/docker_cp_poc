package main

import (
	"flag"
	"io"
	"log"
	"os"
	"os/exec"
)

func main() {
	flag.Parse()

	argv := flag.Args()

	if len(argv) < 2 {
		log.Fatalf("Usage: %s from to", os.Args[0])
	}

	from := argv[0]
	to := argv[1]

	pack := exec.Command("docker-pack", from)
	log.Println("docker-pack", from)
	packout, err := pack.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	packStderr, err := pack.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}
	go io.Copy(os.Stderr, packStderr)

	// todo: Exceptions for:
	// /foo -> /bar
	// /foo/ -> /bar

	unpack := exec.Command("docker-unpack", to)
	log.Println("docker-unpack", to)
	unpackin, err := unpack.StdinPipe()
	unpackStderr, err := unpack.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}
	go io.Copy(os.Stderr, unpackStderr)
	if err != nil {
		log.Fatal(err)
	}
	if err := unpack.Start(); err != nil {
		log.Fatal(err)
	}

	_, err = io.Copy(unpackin, packout)
	if err != nil {
		log.Fatal(err) // this doesn't check dest errors properly. Broken pipe.
	}

}
