// brc compresses or decompresses using brotli format.
package main

import (
	"fmt"
	"io"
	"os"

	"github.com/andybalholm/brotli"
)

const fileExt = ".br"
const discard = "∅"

type action struct {
	fileIn        string
	fileOut       string
	compress      bool
	compressLevel int
	force         bool

	help func()
}

var defaultAction = action{
	compress:      true,
	compressLevel: 6,
}

func openIn(path string) (*os.File, error) {
	if path == "-" {
		return os.Stdin, nil
	}
	return os.Open(path)
}

func openOut(path string, force bool) (io.Writer, error) {
	if path == discard {
		return io.Discard, nil
	}
	if path == "-" {
		return os.Stdout, nil
	}
	var wflag int
	if force {
		wflag = os.O_TRUNC
	} else {
		wflag = os.O_EXCL
	}
	return os.OpenFile(path, os.O_WRONLY|os.O_CREATE|wflag, 0644)
}

func run(a action) (rerr error) {
	in, err := openIn(a.fileIn)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := openOut(a.fileOut, a.force)
	if err != nil {
		return err
	}
	defer safeCloseWriter(out, &rerr)

	switch {
	case a.compress:
		w := brotli.NewWriterLevel(out, a.compressLevel)
		_, err = io.Copy(w, in)
		if err != nil {
			return fmt.Errorf("compress: %w", err)
		}
		err = w.Close()
		if err != nil {
			return fmt.Errorf("compress closing: %w", err)
		}

	case !a.compress:
		r := brotli.NewReader(in)
		_, err := io.Copy(out, r)
		if err != nil {
			return fmt.Errorf("decompress: %w", err)
		}
	}

	return nil
}

func main() {
	a, err := parseArgs(os.Args[1:])
	if err != nil {
		die(2, err)
	}
	if a.help != nil {
		a.help()
		os.Exit(0)
	}

	err = run(a)
	if err != nil {
		die(1, err)
	}
}

func safeClose(c io.Closer, errp *error) {
	cerr := c.Close()
	if cerr != nil && *errp == nil {
		*errp = cerr
	}
}

func safeCloseWriter(w io.Writer, errp *error) {
	if c, ok := w.(io.Closer); ok {
		safeClose(c, errp)
	}
}

func die(code int, err error) {
	fmt.Fprintln(os.Stderr, "brc:", err)
	os.Exit(code)
}
