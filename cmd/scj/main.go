package main

import (
	"fmt"
	"os"
	"os/exec"
	_ "embed"
)

//go:embed scj.sh
var scj string

func main() {
	out, err := exec.Command("sh", append([]string{"-c", scj}, os.Args...)...).CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", out)
		os.Exit(1)
	} else {
		fmt.Fprintf(os.Stdout, "%s", out)
	}
}
