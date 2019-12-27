package go_vars_helper

import (
	"go/build"
	"os"
)

var GOROOT string
var GOPATH string

// Package init
func init() {
	GOROOT = os.Getenv("GOROOT")
	if GOROOT == "" {
		GOROOT = build.Default.GOROOT
	}

	GOPATH = os.Getenv("GOPATH")
	if GOPATH == "" {
		GOPATH = build.Default.GOPATH
	}
}
