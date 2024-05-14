package zygo

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

var ShellCmd string = "/bin/bash"

func init() {
	SetShellCmd()
}

// set ShellCmd as used by SystemFunction
func SetShellCmd() {
	if runtime.GOOS == "windows" {
		ShellCmd = os.Getenv("COMSPEC")
		return
	}
	try := []string{"/usr/bin/bash"}
	if !FileExists(ShellCmd) {
		for i := range try {
			b := try[i]
			if FileExists(b) {
				ShellCmd = b
				return
			}
		}
	}
}

// sys is a builder. shell out, return the combined output.
func SystemBuilder(env *Zlisp, name string, args []Sexp) (Sexp, error) {
	//P("SystemBuilder called with args='%#v'", args)
	return SystemFunction(env, name, args)
}

func SystemFunction(env *Zlisp, name string, args []Sexp) (Sexp, error) {
	if len(args) == 0 {
		return SexpNull, WrongNargs
	}

	flat, err := flattenToWordsHelper(args)
	if err != nil {
		return SexpNull, fmt.Errorf("flatten on '%#v' failed with error '%s'", args, err)
	}
	if len(flat) == 0 {
		return SexpNull, WrongNargs
	}

	joined := strings.Join(flat, " ")
	cmd := ShellCmd

	var out []byte
	if runtime.GOOS == "windows" {
		out, err = exec.Command(cmd, "/c", joined).CombinedOutput()
	} else {
		out, err = exec.Command(cmd, "-c", joined).CombinedOutput()
	}
	if err != nil {
		return SexpNull, fmt.Errorf("error from command: '%s'. Output:'%s'", err, string(Chomp(out)))
	}
	return &SexpStr{S: string(Chomp(out))}, nil
}

// given strings/lists of strings with possible whitespace
// flatten out to a array of SexpStr with no internal whitespace,
// suitable for passing along to (system) / exec.Command()
func FlattenToWordsFunction(env *Zlisp, name string, args []Sexp) (Sexp, error) {
	if len(args) == 0 {
		return SexpNull, WrongNargs
	}
	stringArgs, err := flattenToWordsHelper(args)
	if err != nil {
		return SexpNull, err
	}

	// Now convert to []Sexp{SexpStr}
	res := make([]Sexp, len(stringArgs))
	for i := range stringArgs {
		res[i] = &SexpStr{S: stringArgs[i]}
	}
	return env.NewSexpArray(res), nil
}

func flattenToWordsHelper(args []Sexp) ([]string, error) {
	stringArgs := []string{}

	for i := range args {
		switch c := args[i].(type) {
		case *SexpStr:
			many := strings.Split(c.S, " ")
			stringArgs = append(stringArgs, many...)
		case *SexpSymbol:
			stringArgs = append(stringArgs, c.name)
		case *SexpPair:
			carry, err := ListToArray(c)
			if err != nil {
				return []string{}, fmt.Errorf("tried to convert list of strings to array but failed with error '%s'. Input was type %T / val = '%#v'", err, c, c)
			}
			moreWords, err := flattenToWordsHelper(carry)
			if err != nil {
				return []string{}, err
			}
			stringArgs = append(stringArgs, moreWords...)
		default:
			return []string{}, fmt.Errorf("arguments to system must be strings; instead we have %T / val = '%#v'", c, c)
		}
	} // end i over args
	// INVAR: stringArgs has our flattened list.
	return stringArgs, nil
}

func Chomp(by []byte) []byte {
	if len(by) > 0 {
		n := len(by)
		if by[n-1] == '\n' {
			return by[:n-1]
		}
	}
	return by
}
