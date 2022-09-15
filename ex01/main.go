package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	flag.Parse()

	if flag.Arg(0) == "" {
		fmt.Fprintf(os.Stderr, "error: empty argument\n")
		os.Exit(1)
	} else if flag.Arg(1) != "" {
		fmt.Fprintf(os.Stderr, "error: multiple arguments\n")
		os.Exit(1)
	}
	fmt.Println(flag.Arg(0))
}
