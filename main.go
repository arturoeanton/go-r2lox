package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/arturoeanton/go-r2lox/coati2lang"
)

func runFile(path string) (string, error) {
	archivo, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	contenido := string(archivo)
	return contenido, nil
}

func run(source string) {
	tokens := coati2lang.ScanTokens(source)
	parse := coati2lang.NewParser(tokens)

	expr := parse.Parse()
	interp := coati2lang.NewInterpreter(expr)
	interp.Interpret()

}

func main() {
	var arg_script string
	flag.StringVar(&arg_script, "script", "script.lox", "script to run")
	flag.Parse()

	source, err := runFile(arg_script)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(coati2lang.ERROR_FILE_NOT_FOUND)
	}
	run(source)
	if coati2lang.HasError {
		os.Exit(coati2lang.ERROR_SYNTAX)
	}

}
