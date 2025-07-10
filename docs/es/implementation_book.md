# Libro de Implementación: go-r2lox

## Índice

1. [Introducción](#introducción)
2. [Arquitectura General](#arquitectura-general)
3. [Análisis Léxico (Scanner)](#análisis-léxico-scanner)
4. [Análisis Sintáctico (Parser)](#análisis-sintáctico-parser)
5. [Árbol de Sintaxis Abstracta (AST)](#árbol-de-sintaxis-abstracta-ast)
6. [Intérprete y Evaluación](#intérprete-y-evaluación)
7. [Gestión de Variables y Scope](#gestión-de-variables-y-scope)
8. [Funciones y Closures](#funciones-y-closures)
9. [Estructuras de Control](#estructuras-de-control)
10. [Tipos de Datos y Operadores](#tipos-de-datos-y-operadores)
11. [Manejo de Errores](#manejo-de-errores)
12. [Optimizaciones](#optimizaciones)
13. [Casos de Uso Avanzados](#casos-de-uso-avanzados)
14. [Debugging y Profiling](#debugging-y-profiling)

---

## Introducción

Este documento detalla la implementación completa del intérprete go-r2lox, basado en la segunda implementación del libro "Crafting Interpreters". La implementación sigue un enfoque tree-walking interpreter, construyendo un AST y evaluándolo directamente.

### Filosofía de Diseño

1. **Simplicidad**: Código claro y mantenible
2. **Correctitud**: Implementación fiel al estándar Lox
3. **Extensibilidad**: Arquitectura modular para futuras mejoras
4. **Rendimiento**: Optimizaciones donde no comprometan la claridad

### Estructura del Proyecto

```
go-r2lox/
├── main.go                 # Punto de entrada
├── coati2lang/            # Paquete principal del intérprete
│   ├── scanner.go         # Análisis léxico
│   ├── parser.go          # Análisis sintáctico
│   ├── interpreter.go     # Evaluación
│   ├── enviroment.go      # Gestión de variables
│   ├── expr.go            # Definiciones AST
│   ├── tokens.go          # Definiciones de tokens
│   ├── globals.go         # Funciones globales
│   └── r2loxerrors.go     # Manejo de errores
├── docs/                  # Documentación
└── script.lox            # Script de ejemplo
```

---

## Arquitectura General

### Flujo de Ejecución

```
Código Fuente → Scanner → Tokens → Parser → AST → Interpreter → Resultado
```

### Componentes Principales

#### 1. **Scanner (Análisis Léxico)**
- **Responsabilidad**: Convierte código fuente en tokens
- **Entrada**: String con código fuente
- **Salida**: Slice de tokens
- **Archivo**: `coati2lang/scanner.go`

#### 2. **Parser (Análisis Sintáctico)**
- **Responsabilidad**: Convierte tokens en AST
- **Entrada**: Slice de tokens
- **Salida**: AST (statements)
- **Archivo**: `coati2lang/parser.go`

#### 3. **Interpreter (Evaluación)**
- **Responsabilidad**: Ejecuta el AST
- **Entrada**: AST
- **Salida**: Valores resultantes
- **Archivo**: `coati2lang/interpreter.go`

---

## Análisis Léxico (Scanner)

### Propósito
El scanner convierte el código fuente en una secuencia de tokens que el parser puede entender.

### Implementación Detallada

#### Estructura del Scanner

```go
type Scanner struct {
    Source  string    // Código fuente completo
    Tokens  []Token   // Tokens generados
    Start   int       // Inicio del token actual
    Current int       // Posición actual en el código
    Line    int       // Línea actual (para errores)
}
```

#### Proceso de Tokenización

1. **Inicialización**
   ```go
   func NewScanner(source string) *Scanner {
       return &Scanner{
           Source:  source,
           Tokens:  []Token{},
           Start:   0,
           Current: 0,
           Line:    1,
       }
   }
   ```

2. **Bucle Principal**
   ```go
   func (s *Scanner) ScanTokens() []Token {
       for !s.isAtEnd() {
           s.Start = s.Current
           s.scanToken()
       }
       s.Tokens = append(s.Tokens, NewToken(EOF, "", nil, s.Line))
       return s.Tokens
   }
   ```

#### Manejo de Tokens Especiales

##### Números
```go
func (s *Scanner) number() interface{} {
    // Dígitos enteros
    for s.isDigit(s.peek()) {
        s.advance()
    }
    
    // Parte fraccionaria
    if s.peek() == '.' && s.isDigit(s.peekNext()) {
        s.advance() // Consume el '.'
        for s.isDigit(s.peek()) {
            s.advance()
        }
    }
    
    value, err := strconv.ParseFloat(s.Source[s.Start:s.Current], 64)
    if err != nil {
        Errors(s.Line, "Error parsing number.")
    }
    
    s.addToken(NUMBER, value)
    return value
}
```

##### Cadenas
```go
func (s *Scanner) string() {
    for s.peek() != '"' && !s.isAtEnd() {
        if s.peek() == '\n' {
            s.Line++
        }
        s.advance()
    }
    s.advance()
    
    if s.isAtEnd() {
        Errors(s.Line, "Unterminated string.")
        return
    }
    
    value := s.Source[s.Start+1 : s.Current-1]
    value, err := strconv.Unquote("\"" + value + "\"")
    if err != nil {
        Errors(s.Line, "Error parsing string.")
    }
    s.addToken(STRING, value)
}
```

##### Cadenas Multilínea
```go
func (s *Scanner) multiline_string() {
    s.advance() // Consume segundo "
    s.advance() // Consume tercer "
    
    for s.peek() != '"' && !s.isAtEnd() {
        if s.peek() == '\n' {
            s.Line++
        }
        s.advance()
    }
    
    // Consume las tres comillas finales
    s.advance()
    s.advance()
    s.advance()
    
    value := s.Source[s.Start+3 : s.Current-3]
    s.addToken(STRING, value)
}
```

##### Identificadores y Palabras Clave
```go
func (s *Scanner) identifier() {
    for s.isAlphaNumeric(s.peek()) {
        s.advance()
    }
    
    text := s.Source[s.Start:s.Current]
    
    tokenType, ok := keywords[text]
    if !ok {
        tokenType = IDENTIFIER
    }
    
    s.addToken(tokenType, text)
}
```

#### Mapa de Palabras Clave

```go
var keywords = map[string]TokenType{
    "and":        AND,
    "class":      CLASS,
    "else":       ELSE,
    "false":      FALSE,
    "for":        FOR,
    "fun":        FUN,
    "if":         IF,
    "nil":        NIL,
    "or":         OR,
    "return":     RETURN,
    "true":       TRUE,
    "var":        VAR,
    "while":      WHILE,
    "let":        LET,
    // ... más palabras clave
}
```

### Casos Especiales

#### Operadores Compuestos
```go
case '*':
    if s.match('*') {
        s.addToken(STAR_STAR, "**")  // Exponenciación
    } else {
        s.addToken(STAR, "*")
    }
```

#### Comentarios
```go
case '/':
    if s.match('/') {
        // Comentario hasta el final de línea
        for s.peek() != '\n' && !s.isAtEnd() {
            s.advance()
        }
    } else {
        s.addToken(SLASH, "/")
    }
```

---

## Análisis Sintáctico (Parser)

### Propósito
El parser convierte la secuencia de tokens en un Árbol de Sintaxis Abstracta (AST) que representa la estructura del programa.

### Gramática de Lox

```
program        → declaration* EOF ;
declaration    → funDecl | varDecl | statement ;
funDecl        → "fun" function ;
function       → IDENTIFIER "(" parameters? ")" block ;
varDecl        → "var" IDENTIFIER ( "=" expression )? ";" ;
statement      → exprStmt | ifStmt | returnStmt | whileStmt | forStmt | block ;
block          → "{" declaration* "}" ;
exprStmt       → expression ";" ;
ifStmt         → "if" "(" expression ")" statement ("else" statement)? ;
returnStmt     → "return" expression? ";" ;
whileStmt      → "while" "(" expression ")" statement ;
forStmt        → "for" "(" (varDecl | exprStmt | ";") expression? ";" expression? ")" statement ;
expression     → assignment ;
assignment     → IDENTIFIER "=" assignment | or ;
or             → and ("or" and)* ;
and            → equality ("and" equality)* ;
equality       → comparison (("!=" | "==") comparison)* ;
comparison     → term ((">" | ">=" | "<" | "<=") term)* ;
term           → factor (("-" | "+") factor)* ;
factor         → unary (("/" | "*" | "**" | "%") unary)* ;
unary          → ("!" | "-") unary | call ;
call           → primary ("(" arguments? ")")* ;
primary        → NUMBER | STRING | "true" | "false" | "nil" | "(" expression ")" | IDENTIFIER ;
```

### Implementación del Parser

#### Estructura Principal

```go
type Parser struct {
    Tokens  []Token
    Current int
    Start   int
}
```

#### Métodos de Parsing

##### Expresiones Primarias
```go
func (p *Parser) primary() Expr {
    if p.match(FALSE) {
        return Literal{Value: false}
    }
    if p.match(TRUE) {
        return Literal{Value: true}
    }
    if p.match(NIL) {
        return Literal{Value: nil}
    }
    
    if p.match(NUMBER, STRING, MULTILINE_STRING, TEMPLATE_STRING) {
        return Literal{Value: p.previous().Literal}
    }
    
    if p.match(IDENTIFIER) {
        name := p.previous()
        
        // Manejo de selectores (arrays/maps)
        selectors := [][]Expr{}
        for p.check(LEFT_BRACKET) || p.check(DOT) {
            if p.match(LEFT_BRACKET) {
                array := p.Array()
                selectors = append(selectors, array)
            } else if p.match(DOT) {
                name := p.consume(IDENTIFIER, "Expect property name after '.'.")
                selectors = append(selectors, []Expr{Literal{Value: name.Lexeme}})
            }
        }
        
        return Var{Name: name, Selectors: selectors}
    }
    
    if p.match(LEFT_PAREN) {
        expr := p.Equality()
        p.consume(RIGHT_PAREN, "Expect ')' after expression.")
        return Grouping{Expression: expr}
    }
    
    return nil
}
```

##### Operadores Binarios
```go
func (p *Parser) Equality() Expr {
    expr := p.Comparison()
    
    for p.match(BANG_EQUAL, EQUAL_EQUAL) {
        operator := p.previous()
        right := p.Comparison()
        expr = Binary{Left: expr, Operator: operator, Right: right}
    }
    
    return expr
}

func (p *Parser) Comparison() Expr {
    expr := p.Term()
    
    for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
        operator := p.previous()
        right := p.Term()
        expr = Binary{Left: expr, Operator: operator, Right: right}
    }
    
    return expr
}
```

##### Llamadas a Funciones
```go
func (p *Parser) finishCall(callee Expr) Expr {
    arguments := []Expr{}
    
    var this Var
    if self, ok := callee.(Var); ok {
        if len(self.Selectors) > 0 {
            this = self
            this.Selectors = this.Selectors[:len(this.Selectors)-1]
        }
    }
    
    if !p.check(RIGHT_PAREN) {
        for {
            if len(arguments) >= 255 {
                Errors(p.peek().Line, "Can't have more than 255 arguments.")
            }
            arguments = append(arguments, p.Expression())
            if !p.match(COMMA) {
                break
            }
        }
    }
    
    paren := p.consume(RIGHT_PAREN, "Expect ')' after arguments.")
    
    return Call{Callee: callee, Paren: paren, Arguments: arguments, This: this}
}
```

#### Manejo de Arrays y Mapas

##### Arrays
```go
func (p *Parser) Array() []Expr {
    initializer := []Expr{}
    
    if !p.check(RIGHT_BRACKET) {
        for {
            if len(initializer) >= 255 {
                Errors(p.peek().Line, "Can't have more than 255 arguments.")
            }
            expr := p.Expression()
            
            // Rangos (1..10)
            if p.match(DOT) {
                p.consume(DOT, "Expect '.' after arguments.")
                literal := p.Expression()
                start := int(expr.(Literal).Value.(float64))
                end := int(literal.(Literal).Value.(float64))
                for i := start; i <= end; i++ {
                    initializer = append(initializer, Literal{Value: float64(i)})
                }
            } else {
                initializer = append(initializer, expr)
            }
            
            if !p.match(COMMA) {
                break
            }
        }
    }
    
    p.consume(RIGHT_BRACKET, "Expect ']' after arguments.")
    return initializer
}
```

##### Mapas
```go
func (p *Parser) Map() []ItemVar {
    initializer := []ItemVar{}
    
    if !p.check(RIGHT_BRACE) {
        for {
            if len(initializer) >= 255 {
                Errors(p.peek().Line, "Can't have more than 255 arguments.")
            }
            key := p.Expression()
            
            if key_identifier, ok := key.(Var); ok {
                key = Literal{Value: key_identifier.Name.Lexeme}
            }
            
            if p.check(ARROW) {
                // Funciones lambda con =>
                p.consume(ARROW, "Expect '=>' after key.")
                name := Token{Type: IDENTIFIER, Lexeme: "subfx-" + uuid.NewString()}
                subfx := p.Function("method", name)
                subVar := Var{Name: name, Sub: true, InitializerFx: subfx.(Function)}
                initializer = append(initializer, ItemVar{Key: key, Value: subVar})
            } else {
                p.consume(COLON, "Expect ':' after key.")
                value := p.Expression()
                initializer = append(initializer, ItemVar{Key: key, Value: value})
            }
            
            if !p.match(COMMA) {
                break
            }
        }
    }
    
    p.consume(RIGHT_BRACE, "Expect '}' after arguments.")
    return initializer
}
```

#### Manejo de Errores y Sincronización

```go
func (p *Parser) synchronize() {
    p.advance()
    
    for !p.isAtEnd() {
        if p.previous().Type == SEMICOLON {
            return
        }
        
        switch p.peek().Type {
        case CLASS, FUN, VAR, FOR, IF, WHILE, RETURN:
            return
        }
        
        p.advance()
    }
}
```

---

## Árbol de Sintaxis Abstracta (AST)

### Propósito
El AST representa la estructura sintáctica del programa en un formato que el intérprete puede evaluar eficientemente.

### Definiciones de Nodos

#### Expresiones (Expr)

```go
type Expr interface {
    AcceptExpr(visitor ExprVisitor) interface{}
}

type ExprVisitor interface {
    VisitBinaryExpr(expr Binary) interface{}
    VisitUnaryExpr(expr Unary) interface{}
    VisitLiteralExpr(expr Literal) interface{}
    VisitGroupingExpr(expr Grouping) interface{}
    VisitVariableExpr(expr Var) interface{}
    VisitAssignExpr(expr Assign) interface{}
    VisitCallExpr(expr Call) interface{}
    VisitLogicalExpr(expr Logical) interface{}
    VisitGroupingABSExpr(expr GroupingABS) interface{}
}
```

#### Tipos de Expresiones

##### Literales
```go
type Literal struct {
    Value interface{}
}

func (l Literal) AcceptExpr(visitor ExprVisitor) interface{} {
    return visitor.VisitLiteralExpr(l)
}
```

##### Operadores Binarios
```go
type Binary struct {
    Left     Expr
    Operator Token
    Right    Expr
}

func (b Binary) AcceptExpr(visitor ExprVisitor) interface{} {
    return visitor.VisitBinaryExpr(b)
}
```

##### Variables
```go
type Var struct {
    Name             Token
    Selectors        [][]Expr      // Para arrays/maps
    InitializerVal   Expr          // Para declaraciones
    InitializerArray []Expr        // Para arrays
    InitializerMap   []ItemVar     // Para mapas
    InitializerFx    interface{}   // Para funciones
    SizeArrayInit    int          // Tamaño de array
    Sub              bool         // Es subvariable
}

func (v Var) AcceptExpr(visitor ExprVisitor) interface{} {
    return visitor.VisitVariableExpr(v)
}
```

##### Llamadas a Funciones
```go
type Call struct {
    Callee    Expr
    Paren     Token
    Arguments []Expr
    This      Var        // Para métodos
    SubCall   *Call      // Para llamadas encadenadas
}

func (c Call) AcceptExpr(visitor ExprVisitor) interface{} {
    return visitor.VisitCallExpr(c)
}
```

#### Sentencias (Stmt)

```go
type Stmt interface {
    AcceptStmt(visitor StmtVisitor) interface{}
}

type StmtVisitor interface {
    VisitExpressionStmt(stmt Expression) interface{}
    VisitVar(stmt Var) interface{}
    VisitBlockStmt(stmt Block) interface{}
    VisitIfStmt(stmt If) interface{}
    VisitWhileStmt(stmt While) interface{}
    VisitFunctionStmt(stmt Function) interface{}
    VisitReturnStmt(stmt Return) interface{}
}
```

##### Tipos de Sentencias

###### Expresiones
```go
type Expression struct {
    Expression Expr
}

func (e Expression) AcceptStmt(visitor StmtVisitor) interface{} {
    return visitor.VisitExpressionStmt(e)
}
```

###### Bloques
```go
type Block struct {
    Statements []Stmt
}

func (b Block) AcceptStmt(visitor StmtVisitor) interface{} {
    return visitor.VisitBlockStmt(b)
}
```

###### Funciones
```go
type Function struct {
    Name       Token
    Parameters []Token
    Body       []Stmt
    Closure    *Enviroment
}

func (f Function) AcceptStmt(visitor StmtVisitor) interface{} {
    return visitor.VisitFunctionStmt(f)
}
```

### Patrón Visitor

El patrón Visitor permite agregar operaciones a los nodos del AST sin modificar sus definiciones.

#### Implementación en el Intérprete

```go
type Interpreter struct {
    Stmts      []Stmt
    enviroment *Enviroment
}

func (i *Interpreter) VisitBinaryExpr(expr Binary) interface{} {
    left := i.evaluate(expr.Left)
    right := i.evaluate(expr.Right)
    
    switch expr.Operator.Type {
    case MINUS:
        return left.(float64) - right.(float64)
    case PLUS:
        if _, ok := left.(float64); ok {
            return left.(float64) + right.(float64)
        }
        if _, ok := left.(string); ok {
            return left.(string) + right.(string)
        }
        return nil
    case STAR:
        return left.(float64) * right.(float64)
    case SLASH:
        return left.(float64) / right.(float64)
    case STAR_STAR:
        return math.Pow(left.(float64), right.(float64))
    case PERCENT:
        return (left.(float64) * right.(float64) / 100.0)
    // ... más operadores
    }
    
    return nil
}
```

---

## Intérprete y Evaluación

### Propósito
El intérprete recorre el AST y ejecuta las operaciones representadas por cada nodo.

### Estructura del Intérprete

```go
type Interpreter struct {
    Stmts      []Stmt
    enviroment *Enviroment
}

func NewInterpreter(stmts []Stmt) *Interpreter {
    global := NewEnviroment(nil)
    
    // Registrar funciones globales
    for k, v := range GlobalFx {
        global.Define(k, v)
    }
    
    return &Interpreter{
        Stmts:      stmts,
        enviroment: global,
    }
}
```

### Evaluación de Expresiones

#### Literales
```go
func (i *Interpreter) VisitLiteralExpr(expr Literal) interface{} {
    return expr.Value
}
```

#### Variables
```go
func (i *Interpreter) VisitVariableExpr(expr Var) interface{} {
    value, ok := i.enviroment.Get(expr.Name.Lexeme)
    
    if !ok {
        if expr.Sub {
            expr.Sub = false
            i.VisitVar(expr)
            value, ok = i.enviroment.Get(expr.Name.Lexeme)
        }
        if !ok {
            log.Fatalln("Undefined variable '" + expr.Name.Lexeme + "'.")
        }
    }
    
    // Manejo de selectores (arrays/maps)
    if len(expr.Selectors) > 0 {
        for _, arraySelector := range expr.Selectors {
            if array, ok := value.([]interface{}); ok {
                values := make([]interface{}, len(arraySelector))
                for index, selExpr := range arraySelector {
                    selector := i.evaluate(selExpr)
                    pos := int(selector.(float64))
                    if pos < 0 {
                        pos = len(array) + pos
                    }
                    values[index] = array[pos]
                }
                if len(values) == 1 {
                    value = values[0]
                    continue
                }
                value = values
            }
            
            if m, ok := value.(map[interface{}]interface{}); ok {
                values := make(map[interface{}]interface{})
                var selector interface{}
                for _, selExpr := range arraySelector {
                    selector = i.evaluate(selExpr)
                    values[selector] = m[selector]
                }
                if len(values) == 1 {
                    value = values[selector]
                    continue
                }
                value = values
            }
        }
    }
    
    return value
}
```

#### Asignaciones
```go
func (i *Interpreter) VisitAssignExpr(expr Assign) interface{} {
    value := i.full_evaluate(expr.Value)
    old, ok := i.enviroment.Get(expr.Name.Lexeme)
    if !ok {
        log.Fatalln("Undefined variable '" + expr.Name.Lexeme + "'.")
    }
    
    // Asignación con selectores
    path_var := make([]interface{}, len(expr.Selectors))
    if len(expr.Selectors) > 0 {
        for index, arraySelector := range expr.Selectors {
            for _, selExpr := range arraySelector {
                path_var[index] = i.full_evaluate(selExpr)
            }
        }
        new, err := i.setByPath(old, path_var, value)
        if err != nil {
            log.Fatalln(err)
        }
        if expr.Name.Lexeme != "this" {
            i.enviroment.Assign(expr.Name.Lexeme, new)
        }
        return value
    } else {
        i.enviroment.Assign(expr.Name.Lexeme, value)
    }
    return value
}
```

#### Llamadas a Funciones
```go
func (i *Interpreter) VisitCallExpr(expr Call) interface{} {
    callee := i.evaluate(expr.Callee)
    
    // Manejo de métodos de string
    if variable, ok := expr.Callee.(Var); ok {
        value := i.evaluate(variable)
        if str, ok := value.(string); ok {
            fx := variable.Selectors[len(variable.Selectors)-1][0].(Literal).Value.(string)
            return STRING_FX_MAP[fx](str, i, expr.Arguments)
        }
    }
    
    var arguments []interface{}
    for _, argument := range expr.Arguments {
        arguments = append(arguments, i.evaluate(argument))
    }
    
    callable, ok := callee.(LoxCallable)
    if !ok {
        fmt.Println("Can only call functions and classes.")
        return nil
    }
    
    if callable.Arity() != -1 && callable.Arity() != len(arguments) {
        fmt.Println("Expected", callable.Arity(), "arguments but got", len(arguments), ".")
        return nil
    }
    
    // Llamada con contexto this
    if expr.This.Name.Lexeme != "" {
        this := i.full_evaluate(expr.This)
        return callable.Call(i, arguments, this)
    }
    
    return callable.Call(i, arguments, nil)
}
```

### Evaluación de Sentencias

#### Expresiones como Sentencias
```go
func (i *Interpreter) VisitExpressionStmt(stmt Expression) interface{} {
    return i.evaluate(stmt.Expression)
}
```

#### Declaraciones de Variables
```go
func (i *Interpreter) VisitVar(stmt Var) interface{} {
    var value interface{}
    
    if stmt.InitializerVal != nil {
        value = i.full_evaluate(stmt.InitializerVal)
    }
    
    if stmt.InitializerArray != nil {
        var values []interface{} = make([]interface{}, stmt.SizeArrayInit)
        for index, value := range stmt.InitializerArray {
            values[index] = i.full_evaluate(value)
        }
        value = values
    }
    
    if stmt.InitializerMap != nil {
        var values map[interface{}]interface{} = make(map[interface{}]interface{})
        for _, item := range stmt.InitializerMap {
            key := i.full_evaluate(item.Key)
            value := i.full_evaluate(item.Value)
            values[key] = value
        }
        value = values
    }
    
    if stmt.InitializerFx != nil {
        i.VisitFunctionStmt(stmt.InitializerFx.(Function))
        value = i.full_evaluate(Var{Name: stmt.Name, Sub: false})
    }
    
    i.enviroment.Define(stmt.Name.Lexeme, value)
    return nil
}
```

#### Bloques
```go
func (i *Interpreter) VisitBlockStmt(stmt Block) interface{} {
    return i.executeBlock(stmt.Statements, *NewEnviroment(i.enviroment))
}

func (i *Interpreter) executeBlock(stmts []Stmt, enviroment Enviroment) interface{} {
    previous := i.enviroment
    defer func() {
        i.enviroment = previous
    }()
    
    i.enviroment = &enviroment
    for _, stmt := range stmts {
        if stmt == nil {
            continue
        }
        i.execute(stmt)
    }
    return nil
}
```

### Funciones Especiales

#### Evaluación Completa
```go
func (i *Interpreter) full_evaluate(expr Expr) interface{} {
    value := i.evaluate(expr)
    
    // Manejo de expresiones anidadas
    expr_v, ok := value.(Expr)
    if ok {
        value := i.full_evaluate(expr_v)
        return value
    }
    
    // Manejo de arrays de expresiones
    expr_a, ok := value.([]Expr)
    if ok {
        var values []interface{} = make([]interface{}, len(expr_a))
        for index, value := range expr_a {
            values[index] = i.full_evaluate(value)
        }
        return values
    }
    
    // Manejo de mapas de expresiones
    expr_m, ok := value.([]ItemVar)
    if ok {
        var values map[interface{}]interface{} = make(map[interface{}]interface{})
        for _, value := range expr_m {
            values[i.full_evaluate(value.Key)] = i.full_evaluate(value.Value)
        }
        return values
    }
    
    return value
}
```

#### Asignación por Ruta
```go
func (i *Interpreter) setByPath(target interface{}, path []interface{}, value interface{}) (interface{}, error) {
    if len(path) == 0 {
        return nil, errors.New("path is too short")
    }
    
    switch t := target.(type) {
    case []interface{}:
        index := int(path[0].(float64))
        
        // Extensión automática del slice
        for len(t) <= index {
            t = append(t, nil)
        }
        
        if len(path) == 1 {
            t[index] = value
            return t, nil
        }
        return i.setByPath(t[index], path[1:], value)
        
    case map[interface{}]interface{}:
        key := path[0]
        if len(path) == 1 {
            t[key] = value
            return t, nil
        }
        
        if _, exists := t[key]; !exists {
            return nil, errors.New("key not found")
        }
        
        if _, ok := t[key].([]interface{}); ok {
            new, err := i.setByPath(t[key], path[1:], value)
            if err != nil {
                return nil, err
            }
            t[key] = new
            return t, nil
        }
        
        return i.setByPath(t[key], path[1:], value)
        
    default:
        return nil, errors.New("unsupported type")
    }
}
```

---

## Gestión de Variables y Scope

### Propósito
El sistema de entornos (environments) gestiona las variables y sus ámbitos de visibilidad.

### Estructura del Entorno

```go
type Enviroment struct {
    values    map[string]interface{}
    enclosing *Enviroment
}

func NewEnviroment(enclosing *Enviroment) *Enviroment {
    return &Enviroment{
        values:    make(map[string]interface{}),
        enclosing: enclosing,
    }
}
```

### Operaciones Principales

#### Definición de Variables
```go
func (e *Enviroment) Define(name string, value interface{}) {
    e.values[name] = value
}
```

#### Obtención de Variables
```go
func (e *Enviroment) Get(name string) (interface{}, bool) {
    if value, ok := e.values[name]; ok {
        return value, true
    }
    
    if e.enclosing != nil {
        return e.enclosing.Get(name)
    }
    
    return nil, false
}
```

#### Asignación de Variables
```go
func (e *Enviroment) Assign(name string, value interface{}) error {
    if _, ok := e.values[name]; ok {
        e.values[name] = value
        return nil
    }
    
    if e.enclosing != nil {
        return e.enclosing.Assign(name, value)
    }
    
    return errors.New("undefined variable '" + name + "'")
}
```

### Cadena de Entornos

Los entornos forman una cadena para implementar el scoping léxico:

```
Global Environment
    │
    ├── Function Environment
    │   │
    │   ├── Block Environment
    │   │   │
    │   │   └── Nested Block Environment
    │   │
    │   └── Another Block Environment
    │
    └── Another Function Environment
```

### Implementación de Scoping

#### Alcance Global
```go
func NewInterpreter(stmts []Stmt) *Interpreter {
    global := NewEnviroment(nil)
    
    // Funciones globales
    for k, v := range GlobalFx {
        global.Define(k, v)
    }
    
    return &Interpreter{
        Stmts:      stmts,
        enviroment: global,
    }
}
```

#### Alcance de Función
```go
func (f Function) Call(interpreter *Interpreter, arguments []interface{}, this interface{}) interface{} {
    // Crear nuevo entorno para la función
    enviroment := NewEnviroment(f.Closure)
    
    // Vincular parámetros
    for i, param := range f.Parameters {
        enviroment.Define(param.Lexeme, arguments[i])
    }
    
    // Ejecutar cuerpo de la función
    defer func() {
        if r := recover(); r != nil {
            if ret, ok := r.(Return); ok {
                return ret.Result
            }
            panic(r)
        }
    }()
    
    interpreter.executeBlock(f.Body, *enviroment)
    return nil
}
```

#### Alcance de Bloque
```go
func (i *Interpreter) executeBlock(stmts []Stmt, enviroment Enviroment) interface{} {
    previous := i.enviroment
    defer func() {
        i.enviroment = previous
    }()
    
    i.enviroment = &enviroment
    for _, stmt := range stmts {
        if stmt == nil {
            continue
        }
        i.execute(stmt)
    }
    return nil
}
```

---

## Funciones y Closures

### Propósito
Implementar funciones como valores de primera clase con soporte para closures.

### Interface de Funciones Llamables

```go
type LoxCallable interface {
    Call(interpreter *Interpreter, arguments []interface{}, this interface{}) interface{}
    Arity() int
}
```

### Implementación de Funciones

#### Estructura de Función
```go
type Function struct {
    Name       Token
    Parameters []Token
    Body       []Stmt
    Closure    *Enviroment
}
```

#### Implementación de LoxCallable
```go
func (f Function) Call(interpreter *Interpreter, arguments []interface{}, this interface{}) interface{} {
    // Crear entorno para la función
    enviroment := NewEnviroment(f.Closure)
    
    // Vincular parámetros
    for i, param := range f.Parameters {
        if i < len(arguments) {
            enviroment.Define(param.Lexeme, arguments[i])
        }
    }
    
    // Vincular 'this' si existe
    if this != nil {
        enviroment.Define("this", this)
    }
    
    // Ejecutar cuerpo con manejo de return
    defer func() {
        if r := recover(); r != nil {
            if ret, ok := r.(Return); ok {
                // Return statement encontrado
                return
            }
            panic(r) // Re-panic si no es un return
        }
    }()
    
    interpreter.executeBlock(f.Body, *enviroment)
    return nil
}

func (f Function) Arity() int {
    return len(f.Parameters)
}
```

### Declaración de Funciones

```go
func (i *Interpreter) VisitFunctionStmt(stmt Function) interface{} {
    function := Function{
        Name:       stmt.Name,
        Parameters: stmt.Parameters,
        Body:       stmt.Body,
        Closure:    i.enviroment, // Captura el entorno actual
    }
    i.enviroment.Define(stmt.Name.Lexeme, function)
    return nil
}
```

### Closures

Los closures capturan el entorno léxico donde se definen:

```go
// Ejemplo de closure
fun makeCounter() {
    var count = 0;
    
    fun counter() {
        count = count + 1;
        return count;
    }
    
    return counter;
}

var counter = makeCounter();
print counter(); // 1
print counter(); // 2
```

La implementación mantiene una referencia al entorno de captura:

```go
function := Function{
    Name:       stmt.Name,
    Parameters: stmt.Parameters,
    Body:       stmt.Body,
    Closure:    i.enviroment, // Entorno de captura
}
```

### Funciones Globales

#### Función Clock
```go
type Clock struct {}

func (c Clock) Call(interpreter *Interpreter, arguments []interface{}, this interface{}) interface{} {
    return time.Now().Unix()
}

func (c Clock) Arity() int {
    return 0
}
```

#### Registro de Funciones Globales
```go
var GlobalFx = make(map[string]LoxCallable)

func init() {
    GlobalFx["clock"] = Clock{}
}
```

### Return Statements

```go
type Return struct {
    Keyword Token
    Value   Expr
    Result  interface{}
}

func (i *Interpreter) VisitReturnStmt(stmt Return) interface{} {
    var result interface{}
    if stmt.Value != nil {
        result = i.evaluate(stmt.Value)
    }
    panic(Return{Keyword: stmt.Keyword, Value: stmt.Value, Result: result})
}
```

---

## Estructuras de Control

### Propósito
Implementar las estructuras de control del lenguaje: if/else, while, for.

### If/Else Statements

#### Estructura AST
```go
type If struct {
    Condition  Expr
    ThenBranch Stmt
    ElseBranch Stmt
}
```

#### Implementación
```go
func (i *Interpreter) VisitIfStmt(stmt If) interface{} {
    if i.isTruthy(i.evaluate(stmt.Condition)) {
        return i.execute(stmt.ThenBranch)
    } else if stmt.ElseBranch != nil {
        return i.execute(stmt.ElseBranch)
    }
    return nil
}
```

#### Parsing
```go
func (p *Parser) IfStatement() Stmt {
    p.consume(LEFT_PAREN, "Expect '(' after 'if'.")
    condition := p.Expression()
    p.consume(RIGHT_PAREN, "Expect ')' after if condition.")
    
    thenBranch := p.Statement()
    var elseBranch Stmt
    if p.match(ELSE) {
        elseBranch = p.Statement()
    }
    
    return If{Condition: condition, ThenBranch: thenBranch, ElseBranch: elseBranch}
}
```

### While Loops

#### Estructura AST
```go
type While struct {
    Condition Expr
    Body      Stmt
}
```

#### Implementación
```go
func (i *Interpreter) VisitWhileStmt(stmt While) interface{} {
    for i.isTruthy(i.evaluate(stmt.Condition)) {
        i.execute(stmt.Body)
    }
    return nil
}
```

#### Parsing
```go
func (p *Parser) WhileStatement() Stmt {
    p.consume(LEFT_PAREN, "Expect '(' after 'while'.")
    condition := p.Expression()
    p.consume(RIGHT_PAREN, "Expect ')' after while condition.")
    body := p.Statement()
    return While{Condition: condition, Body: body}
}
```

### For Loops

Los bucles for se implementan como azúcar sintáctico para while:

```go
func (p *Parser) ForStatement() Stmt {
    p.consume(LEFT_PAREN, "Expect '(' after 'for'.")
    
    // Inicializador
    var initializer Stmt
    if p.match(SEMICOLON) {
        initializer = nil
    } else if p.match(VAR) {
        initializer = p.VarDeclaration()
    } else {
        initializer = p.ExpressionStatement()
    }
    
    // Condición
    var condition Expr
    if !p.check(SEMICOLON) {
        condition = p.Expression()
    }
    p.consume(SEMICOLON, "Expect ';' after loop condition.")
    
    // Incremento
    var increment Expr
    if !p.check(RIGHT_PAREN) {
        increment = p.Expression()
    }
    p.consume(RIGHT_PAREN, "Expect ')' after for clauses.")
    
    body := p.Statement()
    
    // Desugaring a while loop
    if increment != nil {
        body = Block{Statements: []Stmt{body, Expression{Expression: increment}}}
    }
    
    if condition == nil {
        condition = Literal{Value: true}
    }
    body = While{Condition: condition, Body: body}
    
    if initializer != nil {
        body = Block{Statements: []Stmt{initializer, body}}
    }
    
    return body
}
```

### Operadores Lógicos

#### Estructura AST
```go
type Logical struct {
    Left     Expr
    Operator Token
    Right    Expr
}
```

#### Implementación con Short-Circuit
```go
func (i *Interpreter) VisitLogicalExpr(expr Logical) interface{} {
    left := i.evaluate(expr.Left)
    
    if expr.Operator.Type == OR {
        if i.isTruthy(left) {
            return left
        }
    } else {
        if !i.isTruthy(left) {
            return left
        }
    }
    
    return i.evaluate(expr.Right)
}
```

### Función de Truthiness

```go
func (i *Interpreter) isTruthy(object interface{}) bool {
    if object == nil {
        return false
    }
    if object == false {
        return false
    }
    return true
}
```

---

## Tipos de Datos y Operadores

### Tipos de Datos Soportados

#### Números (float64)
```go
// Ejemplos
42
3.14159
-17.5
```

#### Cadenas (string)
```go
// Ejemplos
"Hello, World!"
"Line 1\nLine 2"
"""
Multiline
string
"""
```

#### Booleanos (bool)
```go
// Ejemplos
true
false
```

#### Nil (nil)
```go
// Ejemplo
nil
```

#### Arrays ([]interface{})
```go
// Ejemplos
[1, 2, 3]
["a", "b", "c"]
[1, "hello", true]
[1..10]  // Rango
```

#### Mapas (map[interface{}]interface{})
```go
// Ejemplos
{"name": "John", "age": 30}
{1: "one", 2: "two"}
```

### Operadores Aritméticos

#### Implementación
```go
func (i *Interpreter) VisitBinaryExpr(expr Binary) interface{} {
    left := i.evaluate(expr.Left)
    right := i.evaluate(expr.Right)
    
    switch expr.Operator.Type {
    case MINUS:
        return left.(float64) - right.(float64)
    case PLUS:
        if _, ok := left.(float64); ok {
            return left.(float64) + right.(float64)
        }
        if _, ok := left.(string); ok {
            return left.(string) + right.(string)
        }
        return nil
    case SLASH:
        return left.(float64) / right.(float64)
    case STAR:
        return i.handleMultiplication(left, right)
    case STAR_STAR:
        return math.Pow(left.(float64), right.(float64))
    case PERCENT:
        return (left.(float64) * right.(float64) / 100.0)
    }
    
    return nil
}
```

#### Multiplicación Especial
```go
func (i *Interpreter) handleMultiplication(left, right interface{}) interface{} {
    _, rightIsString := right.(string)
    _, leftIsString := left.(string)
    _, rightIsNumber := right.(float64)
    _, leftIsNumber := left.(float64)
    
    if rightIsString && leftIsNumber {
        var result string
        for i := 0; i < int(left.(float64)); i++ {
            result += right.(string)
        }
        return result
    }
    
    if rightIsNumber && leftIsString {
        var result string
        for i := 0; i < int(right.(float64)); i++ {
            result += left.(string)
        }
        return result
    }
    
    if rightIsNumber && leftIsNumber {
        return left.(float64) * right.(float64)
    }
    
    return nil
}
```

### Operadores de Comparación

```go
case GREATER:
    return left.(float64) > right.(float64)
case GREATER_EQUAL:
    return left.(float64) >= right.(float64)
case LESS:
    return left.(float64) < right.(float64)
case LESS_EQUAL:
    return left.(float64) <= right.(float64)
case BANG_EQUAL:
    return !i.isEqual(left, right)
case EQUAL_EQUAL:
    return i.isEqual(left, right)
```

### Función de Igualdad

```go
func (i *Interpreter) isEqual(a interface{}, b interface{}) bool {
    if a == nil && b == nil {
        return true
    }
    if a == nil {
        return false
    }
    return a == b
}
```

### Operadores Unarios

```go
func (i *Interpreter) VisitUnaryExpr(expr Unary) interface{} {
    value := i.evaluate(expr.Value)
    
    switch expr.Operator.Type {
    case MINUS:
        return -(value.(float64))
    case PLUS_PLUS:
        return value.(float64) + 1
    case MINUS_MINUS:
        return value.(float64) - 1
    case BANG:
        return !i.isTruthy(value)
    default:
        return nil
    }
}
```

### Operadores de Incremento/Decremento

Implementados como azúcar sintáctico:

```go
// i++ se convierte en i = i + 1
for p.match(PLUS_PLUS) {
    op := p.previous()
    operator := Token{Type: PLUS, Lexeme: "++", Literal: nil, Line: op.Line}
    right := Literal{Value: 1.0}
    
    var name string
    if expr, ok := expr.(Var); ok {
        name = expr.Name.Lexeme
    }
    
    expr = Assign{
        Name: Token{Type: IDENTIFIER, Lexeme: name, Literal: nil, Line: op.Line},
        Value: Binary{Left: expr, Operator: operator, Right: right}
    }
}
```

### Métodos de String

```go
var STRING_FX_MAP = map[string]func(string, *Interpreter, []Expr) interface{}{
    "length": func(s string, i *Interpreter, args []Expr) interface{} {
        return float64(len(s))
    },
    "substring": func(s string, i *Interpreter, args []Expr) interface{} {
        start := int(i.evaluate(args[0]).(float64))
        end := int(i.evaluate(args[1]).(float64))
        return s[start:end]
    },
    "charAt": func(s string, i *Interpreter, args []Expr) interface{} {
        index := int(i.evaluate(args[0]).(float64))
        return string(s[index])
    },
    "indexOf": func(s string, i *Interpreter, args []Expr) interface{} {
        substr := i.evaluate(args[0]).(string)
        return float64(strings.Index(s, substr))
    },
    "toUpperCase": func(s string, i *Interpreter, args []Expr) interface{} {
        return strings.ToUpper(s)
    },
    "toLowerCase": func(s string, i *Interpreter, args []Expr) interface{} {
        return strings.ToLower(s)
    },
}
```

---

## Manejo de Errores

### Sistema Actual

#### Errores de Parsing
```go
func (p *Parser) consume(t TokenType, message string) Token {
    line := p.peek().Line
    if p.check(t) {
        return p.advance()
    }
    
    message1 := fmt.Sprint(message, "\n\tLine -> ", fmt.Sprint(line), "\n\tNear ->"+p.peek().Lexeme)
    panic(message1)
}
```

#### Errores de Runtime
```go
func (i *Interpreter) VisitVariableExpr(expr Var) interface{} {
    value, ok := i.enviroment.Get(expr.Name.Lexeme)
    
    if !ok {
        log.Fatalln("Undefined variable '" + expr.Name.Lexeme + "'.")
    }
    
    return value
}
```

### Problemas Identificados

1. **Uso excesivo de panic/log.Fatalln**: Termina abruptamente el programa
2. **Falta de stack traces**: Dificulta el debugging
3. **Errores no recuperables**: En REPL, un error cierra todo
4. **Información limitada**: Los errores no proporcionan contexto suficiente

### Mejoras Propuestas

#### Sistema de Errores Estructurado
```go
type LoxError struct {
    Type    ErrorType
    Message string
    Line    int
    Column  int
    Token   Token
    Context string
}

type ErrorType int

const (
    SYNTAX_ERROR ErrorType = iota
    RUNTIME_ERROR
    TYPE_ERROR
    REFERENCE_ERROR
)
```

#### Manejo Graceful
```go
func (i *Interpreter) safeEvaluate(expr Expr) (interface{}, error) {
    defer func() {
        if r := recover(); r != nil {
            // Convertir panic a error
        }
    }()
    
    return i.evaluate(expr), nil
}
```

#### Stack Traces
```go
type CallStack struct {
    frames []CallFrame
}

type CallFrame struct {
    Function string
    Line     int
    Column   int
}

func (cs *CallStack) Push(function string, line, column int) {
    cs.frames = append(cs.frames, CallFrame{function, line, column})
}

func (cs *CallStack) Pop() {
    if len(cs.frames) > 0 {
        cs.frames = cs.frames[:len(cs.frames)-1]
    }
}
```

---

## Optimizaciones

### Optimizaciones Actuales

#### Lazy Evaluation en Operadores Lógicos
```go
func (i *Interpreter) VisitLogicalExpr(expr Logical) interface{} {
    left := i.evaluate(expr.Left)
    
    if expr.Operator.Type == OR {
        if i.isTruthy(left) {
            return left // Short-circuit
        }
    } else {
        if !i.isTruthy(left) {
            return left // Short-circuit
        }
    }
    
    return i.evaluate(expr.Right)
}
```

#### Reutilización de Entornos
```go
func (i *Interpreter) executeBlock(stmts []Stmt, enviroment Enviroment) interface{} {
    previous := i.enviroment
    defer func() {
        i.enviroment = previous // Restaurar entorno
    }()
    
    i.enviroment = &enviroment
    // ... ejecutar statements
}
```

### Optimizaciones Propuestas

#### Variable Resolution
```go
type Resolver struct {
    interpreter *Interpreter
    scopes      []map[string]bool
}

func (r *Resolver) resolve(stmt Stmt) {
    // Resolver variables durante parse time
    // Evitar lookups dinámicos durante runtime
}
```

#### Constant Folding
```go
func (p *Parser) optimizeBinary(left, right Expr, operator Token) Expr {
    // Si ambos operandos son literales, evaluar en parse time
    if leftLit, ok := left.(Literal); ok {
        if rightLit, ok := right.(Literal); ok {
            return p.foldConstants(leftLit, rightLit, operator)
        }
    }
    return Binary{Left: left, Operator: operator, Right: right}
}
```

#### Interning de Strings
```go
var stringPool = make(map[string]string)

func internString(s string) string {
    if interned, exists := stringPool[s]; exists {
        return interned
    }
    stringPool[s] = s
    return s
}
```

---

## Casos de Uso Avanzados

### Programación Funcional

#### Funciones como Valores
```lox
fun makeAdder(x) {
    fun adder(y) {
        return x + y;
    }
    return adder;
}

var add5 = makeAdder(5);
print add5(3); // 8
```

#### Higher-Order Functions
```lox
fun forEach(array, callback) {
    for (var i = 0; i < array.length; i++) {
        callback(array[i]);
    }
}

forEach([1, 2, 3], fun(x) { print x; });
```

### Estructuras de Datos Complejas

#### Arrays Multidimensionales
```lox
var matrix = [[1, 2], [3, 4]];
print matrix[0][1]; // 2
```

#### Mapas Anidados
```lox
var config = {
    "database": {
        "host": "localhost",
        "port": 5432
    },
    "cache": {
        "enabled": true,
        "ttl": 3600
    }
};

print config["database"]["host"]; // localhost
```

### Metaprogramación

#### Funciones Dinámicas
```lox
fun createGetter(property) {
    return fun(obj) {
        return obj[property];
    };
}

var getName = createGetter("name");
var person = {"name": "John", "age": 30};
print getName(person); // John
```

#### Reflection Básica
```lox
fun typeof(value) {
    if (value == nil) return "nil";
    // Implementar detección de tipos
}
```

---

## Debugging y Profiling

### Herramientas de Debugging Actuales

#### Print Debugging
```lox
fun debug(message, value) {
    print "[DEBUG] " + message + ": " + value;
}

debug("Variable x", x);
```

#### Información de Errores
```go
func Errors(line int, message string) {
    fmt.Fprintf(os.Stderr, "[line %d] Error: %s\n", line, message)
    HasError = true
}
```

### Herramientas Propuestas

#### Breakpoints
```go
type Debugger struct {
    breakpoints map[int]bool
    stepMode    bool
    callStack   []CallFrame
}

func (d *Debugger) checkBreakpoint(line int) {
    if d.breakpoints[line] || d.stepMode {
        d.interactive(line)
    }
}
```

#### Variable Inspection
```go
func (d *Debugger) inspect(enviroment *Enviroment) {
    fmt.Println("Variables in scope:")
    for name, value := range enviroment.values {
        fmt.Printf("  %s = %v\n", name, value)
    }
}
```

#### Performance Profiling
```go
type Profiler struct {
    functionCalls map[string]int
    executionTime map[string]time.Duration
}

func (p *Profiler) startFunction(name string) {
    // Iniciar medición
}

func (p *Profiler) endFunction(name string) {
    // Finalizar medición
}
```

---

## Conclusión

Este documento detalla la implementación completa del intérprete go-r2lox, desde el análisis léxico hasta la evaluación final. La arquitectura modular permite futuras extensiones y optimizaciones manteniendo la claridad del código.

### Puntos Clave

1. **Arquitectura Clara**: Separación de responsabilidades entre scanner, parser e interpreter
2. **Patrón Visitor**: Permite extensibilidad sin modificar definiciones de AST
3. **Gestión de Scope**: Sistema de entornos para variables y closures
4. **Tipos Dinámicos**: Soporte para múltiples tipos de datos con coerción automática
5. **Estructuras de Control**: Implementación completa de if/else, while, for
6. **Funciones Avanzadas**: Soporte para closures y funciones como valores

### Próximos Pasos

1. Implementar mejoras en manejo de errores
2. Agregar suite de pruebas comprehensiva
3. Optimizar rendimiento del intérprete
4. Expandir biblioteca estándar
5. Agregar soporte para clases y herencia

Este intérprete sirve como base sólida para futuras extensiones del lenguaje Lox y como ejemplo educativo de construcción de intérpretes en Go.