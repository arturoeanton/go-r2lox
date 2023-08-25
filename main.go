package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/arturoeanton/go-r2lox/globals"
	"github.com/arturoeanton/go-r2lox/lexer"
	"github.com/arturoeanton/go-r2lox/parser"
)

func runFile(path string) (string, error) {
	archivo, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	contenido := string(archivo)
	return contenido, nil
}

func runPrompt() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			// Podrías manejar el error aquí si lo deseas
			break
		}

		run(line)
		globals.HasError = false
	}
}

func run(source string) {
	tokens := lexer.ScanTokens(source)
	parse := parser.NewParser(tokens)

	expr := parse.Parse()
	interp := parser.NewInterpreter(expr)
	interp.Interpret()

}

func main() {
	var arg_script string
	flag.StringVar(&arg_script, "script", "", "script to run")
	flag.Parse()
	if arg_script != "" {
		source, err := runFile(arg_script)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			os.Exit(globals.ERROR_FILE_NOT_FOUND)
		}
		run(source)
		if globals.HasError {
			os.Exit(globals.ERROR_SYNTAX)
		}
	} else {
		//runPrompt()
		source := `
	

		fun fib(n) {
			print n;
			if (n <= 1) return n;
			var r = fib(n - 1) + fib(n - 2);
			return r;
		}
		var f = fib(6); //13
		print f;
`
		run(source)
	}

}
