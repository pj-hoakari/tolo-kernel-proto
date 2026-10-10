package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	kernelproto "github.com/pj-hoakari/tolo-kernel-proto"
)

func main() {
	dir := flag.String("write-proto", "", "include root to write the bundled .proto files into (required)")

	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	if *dir == "" || flag.NArg() != 0 {
		flag.Usage()
		os.Exit(2)
	}

	written, err := kernelproto.WriteProto(*dir)
	if err != nil {
		logger.Error("write proto failed", "error", err)
		os.Exit(1)
	}

	for _, filename := range written {
		fmt.Println(filename)
	}
}
