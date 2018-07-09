package main

import (
	"bufio"
	"io"
	"os"

	"github.com/dlsniper/test2teamcity/stdlib"
)

func toTeamCity(r io.Reader, w io.Writer, process func(string, io.Writer)) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) != 0 {
			process(line, w)
		} else {
			break
		}
	}
}

func main() {
	toTeamCity(os.Stdin, os.Stdout, stdlib.ProcessStdLib)
}
