package main

import (
	"flag"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/docker/docker/vendor/src/code.google.com/p/go/src/pkg/archive/tar"
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

	// todo: paths should be cleaned with suffixes kept

	// <magic>

	_, fbase := path.Split(from)
	_, tbase := path.Split(to)

	// /foo -> /bar
	// /foo/ -> /bar
	if fbase != "." && tbase != "" && tbase != "." {
		log.Println("check-destination", to)
		if !destinationExists(to) {
			log.Println("destination not found")
			to = path.Dir(to)

			if fbase == "" {
				fbase = path.Base(from)
			}

			packout = rename(packout, fbase, tbase)
			log.Println("rename", fbase, tbase)
		}
	}

	// </magic>

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

	if err := pack.Start(); err != nil {
		log.Fatal(err)
	}
	if err := unpack.Start(); err != nil {
		log.Fatal(err)
	}

	_, err = io.Copy(unpackin, packout)
	if err != nil {
		log.Fatal(err) // this doesn't check dest errors properly. Broken pipe.
	}

	unpackin.Close()
	if err = unpack.Wait(); err != nil {
		log.Fatal(err)
	}

}

// normally this would check for 404
func destinationExists(dest string) bool {
	unpack := exec.Command("docker-unpack", dest)
	err := unpack.Run()
	return err == nil
}

func rename(in io.ReadCloser, from, to string) io.ReadCloser {
	r, w := io.Pipe()

	// should wait until someone actually reads
	go func() {
		reader := tar.NewReader(in)
		writer := tar.NewWriter(w)

		for {
			hdr, err := reader.Next()
			if err == io.EOF {
				w.Close()
				return
			}
			if err != nil {
				w.CloseWithError(err)
				return
			}

			hdr.Name = path.Clean(hdr.Name)
			if hdr.Name == from || strings.HasPrefix(hdr.Name, from+"/") {
				hdr.Name = strings.Replace(hdr.Name, from, to, 1)
			}

			err = writer.WriteHeader(hdr)
			if err != nil {
				w.CloseWithError(err)
				return
			}

			if hdr.Size != 0 {
				_, err = io.Copy(writer, reader)
				if err != nil {
					w.CloseWithError(err)
					return
				}
			}
		}
	}()

	return r
}
