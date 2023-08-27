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
		var key = "key1";
		var i =123;
		var j = |2-5- (10%90)|;
		println(j);

		var a[] = [1,key,i,88];
		println(a[-1]);
		var ja[] = [1,20,3,4,5,6,7,8,9,10];
		var x{} = {key:i, "jojo":1, "jaja": ja};
		var s[] = [0..100, x,90];
		println(">>>>",s[4..9, 15, 8..13, 101]);
		println("s 101 :",s[101]["jaja"][1]);

		var b{} = {key:i, "jojo":[1,4,7], "jaja": 3};
		println("b>",b);

		b["jojo"] = 2;
		println("b>",b);

		

		var frase1[] = ["hola1", "mundo1"];
		var frase2[] = ["hola2", "mundo2"];
		var frases[] = [frase1,frase2, [8,9], {"pp" : [1,0,0,{"ko":{"p":0}}]}];

		println("aa>:",frases);
		frases[3]["pp"]= 10;		
		println("aa>:",frases);

		println("aa:",frases[2][1]);
		println("aa:",frases[3]["pp"][3]["ko"]);
		println("aa:",frases[0][1]);

		println(b["key1","jaja"]);
		

		`
		run(source)
	}

}
