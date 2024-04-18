// brc compresses or decompresses using brotli format.
package main

import (
	"fmt"
	"io"
	"os"

	"github.com/andybalholm/brotli"
)

type config struct {
	compress      bool
	compressLevel int

	help func()
}

func run(conf config) error {
	switch {
	case conf.compress:
		w := brotli.NewWriterLevel(os.Stdout, conf.compressLevel)
		_, err := io.Copy(w, os.Stdin)
		if err != nil {
			return fmt.Errorf("compress: %w", err)
		}
		err = w.Close()
		if err != nil {
			return fmt.Errorf("compress closing: %w", err)
		}

	case !conf.compress:
		r := brotli.NewReader(os.Stdin)
		_, err := io.Copy(os.Stdout, r)
		if err != nil {
			return fmt.Errorf("decompress: %w", err)
		}
	}

	return nil
}

func main() {
	conf, err := parseArgs(os.Args[1:])
	if err != nil {
		die(2, err)
	}
	if conf.help != nil {
		conf.help()
		os.Exit(0)
	}

	err = run(conf)
	if err != nil {
		die(1, err)
	}
}

func die(code int, err error) {
	fmt.Fprintln(os.Stderr, "brc:", err)
	os.Exit(code)
}
