package main

type ProgramArgs struct {
	FileName string `arg:"positional, required"`
	NoPath   bool   `arg:"-n, --no-path" default:"false" help:"Include directories in user's $PATH."`
	Ugly     bool   `arg:"-u, --ugly" default:"false" help:"Remove colorized output. Yes it's ugly."`
}
