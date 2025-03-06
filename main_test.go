package main

import (
	"testing"

	"github.com/rogpeppe/go-internal/testscript"
)

func TestMain(m *testing.M) {
	testscript.Main(m, map[string]func(){
		"tablizer": main,
	})
}

func TestTablizer(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir: "t",
	})
}
