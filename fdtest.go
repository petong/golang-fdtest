package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

func main() {

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	open := func(paths []string) []*os.File {
		rc := make([]*os.File, 0, len(paths))
		for _, path := range paths {
			fh, err := os.Create(path)
			if err != nil {
				log.Printf("error: %v", err)
				continue
			}
			rc = append(rc, fh)
		}
		return rc
	}

	closers := func(fh []*os.File) []io.Closer {
		rc := make([]io.Closer, 0, len(fh))
		for _, f := range fh {
			rc = append(rc, f)
		}
		return rc
	}

	writers := func(fh []*os.File) []io.Writer {
		rc := make([]io.Writer, 0, len(fh))
		for _, f := range fh {
			rc = append(rc, f)
		}
		return rc
	}

	closefiles := func(f ...io.Closer) {
		for _, fh := range f {
			fh.Close()
		}
	}

	numberoffiles := makeRange(0, 2048)
	var files []string

	for _, name := range numberoffiles[1:] {
		files = append(files, "tmp/"+strconv.Itoa(name))
	}

	if _, err := os.Stat("tmp"); os.IsNotExist(err) {
		os.Mkdir("tmp", 0755)
	}

	//fmt.Println(files)

	tmpfiles := open(files)

	mw := io.MultiWriter(writers(tmpfiles)...)

	r := strings.NewReader("foo\n")
	io.Copy(mw, r)

	// block until we get a signal to exit
	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()
	fmt.Println("awaiting signal")
	<-done
	fmt.Println("cleaning up")
	closefiles(closers(tmpfiles)...)
}

func makeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}
