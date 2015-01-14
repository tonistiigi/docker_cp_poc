package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path"

	"github.com/tonistiigi/docker_cp_poc/archive"
)

func main() {
	flag.Parse()

	argv := flag.Args()

	if len(argv) < 1 {
		log.Fatalf("Usage: %s path/to/unpack < content.tar", os.Args[0])
	}

	if err := unpack(os.Stdin, argv[0]); err != nil {
		if err != io.EOF {
			log.Fatal(err)
		}
	}

}

func unpack(content io.ReadCloser, dest string) error {
	dest = path.Clean(dest) // there is no sfx magic here
	if fi, err := os.Stat(dest); err != nil {
		return err
	} else if !fi.IsDir() {
		return fmt.Errorf("%s is not a directory", dest)
	}

	// this does not currently handle file/dir conflicts properly
	return archive.Untar(content, dest, nil)
}