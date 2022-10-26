package main

import "github.com/alexflint/go-arg"

func main() {
	var args ProgramArgs
	arg.MustParse(&args)
}
