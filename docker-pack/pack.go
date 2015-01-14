package main

import (
	"flag"
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
		log.Fatalf("Usage: %s path/to/archive", os.Args[0])
	}

	if tar, err := pack(argv[0]); err != nil {
		log.Fatal(err)
	} else {
		io.Copy(os.Stdout, tar)
	}
}

func pack(src string) (content io.ReadCloser, err error) {
	if _, err = os.Stat(src); err != nil {
		return
	}

	// todo: think when to clean path

	var filter []string
	dir, base := path.Split(src)

	if len(base) == 0 {
		dir, base = path.Split(src[:len(src)-1])
	}

	if base != "." {
		filter = append(filter, base)
	}

	return archive.TarWithOptions(dir, &archive.TarOptions{
		Compression:  archive.Uncompressed,
		IncludeFiles: filter,
	})

}
