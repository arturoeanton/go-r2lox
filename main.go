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
		
		fun sum (a, b) {
			print  a+b;
			return a+b;
		}
		print sum(2,3);
		
		print clock();
		
		var j = 0;
		while ( j < 10 ) {
			print j;
			j = j + 1;
		}
		for (var i = 0; i < 10; i = i + 1){
			print i;
		} 

		if (true or false) {
			print "if true";
		}else {
			print "no if true";
		}


		print "hi" or 2; // hi
		print nil or "yes"; // yes
		
		print "hi" and 2;  // 2
		print nil and "yes"; // nil

		print "-----------------";
		print "Table OR";
		print true or true; // true
		print true or false; // true
		print  false or true; // true
		print false or false; // false


		print "Table AND"; 
		print true and true; // true
		print true and false; // false
		print  false and true; // false
		print false and false; // false

		print "-----------------";


		

		var a = "global a"; 
var b = "global b"; 
var c = "global c"; 
{
    var a = "outer a"; 
    var b = "outer b";
    {
        var a = "inner a"; 
        print a;
        print b;
        print c;
    }
    print a; 
    print b; 
    print c;
}
print a; 
print b; 
print c;


			`
		run(source)
	}

}
