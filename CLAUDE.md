# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is `go-r2lox`, a Go implementation of the Lox programming language interpreter based on the second implementation from "Crafting Interpreters" by Robert Nystrom. The project implements a tree-walking interpreter that builds an Abstract Syntax Tree (AST) from source code.

## Build and Run Commands

### Build the interpreter:
```bash
go build
```

### Run a Lox script:
```bash
./go-r2lox <script.lox>
```

### Run in REPL mode:
```bash
./go-r2lox
```

### Run with custom script flag:
```bash
./go-r2lox -script=path/to/script.lox
```

### Standard Go commands:
```bash
go run main.go <script.lox>
go test ./...
go fmt ./...
go vet ./...
```

## Architecture

The interpreter follows a classic compiler architecture with these main phases:

1. **Lexical Analysis** (`coati2lang/scanner.go`): Tokenizes source code into tokens
2. **Parsing** (`coati2lang/parser.go`): Builds an AST from tokens using recursive descent parsing
3. **Interpretation** (`coati2lang/interpreter.go`): Executes the AST using the visitor pattern

### Key Components

- **`main.go`**: Entry point that handles file reading and orchestrates the interpretation pipeline
- **`coati2lang/` package**: Contains the core interpreter implementation
  - `scanner.go`: Lexical analyzer that converts source text to tokens
  - `parser.go`: Recursive descent parser that builds AST from tokens
  - `interpreter.go`: Tree-walking interpreter that executes the AST
  - `tokens.go`: Token definitions and types
  - `expr.go`: Expression AST node definitions
  - `enviroment.go`: Variable environment/scope management
  - `globals.go`: Global function definitions
  - `r2loxerrors.go`: Error handling utilities

### Language Features Implemented

- Variables and assignment (including arrays and maps)
- Functions with closures
- Control flow (if/else, while, for loops)
- Expressions (arithmetic, logical, comparison)
- Built-in functions (clock)
- String methods and operations
- Array/map indexing and manipulation

### Error Handling

The interpreter uses Go's panic/recover mechanism for error handling during parsing and runtime errors. Parse errors are caught and reported with line numbers.

## Development Notes

- The project uses Go modules with `go.mod` requiring Go 1.20+
- External dependency: `github.com/google/uuid` for generating unique identifiers
- The interpreter supports both single-file execution and REPL mode
- No formal test suite is present in the current codebase
- The `script.lox` file serves as a default test/example script

## Code Style

- Follow standard Go conventions
- Package `coati2lang` contains the interpreter core
- Use visitor pattern for AST traversal
- Error reporting includes line numbers for debugging