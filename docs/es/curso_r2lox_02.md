# Curso go-r2lox: Módulo 2 - Análisis Léxico y Tokens

## Objetivos del Módulo
- Entender qué es el análisis léxico y por qué es importante
- Explorar cómo go-r2lox convierte código fuente en tokens
- Analizar la implementación del scanner
- Experimentar con diferentes tipos de tokens

---

## ¿Qué es el Análisis Léxico?

### Introducción
El análisis léxico es la primera fase en el procesamiento de cualquier lenguaje de programación. Su función es convertir una secuencia de caracteres (el código fuente) en una secuencia de tokens significativos.

```
Código fuente: "var x = 42;"
       ↓
   Análisis Léxico
       ↓
Tokens: [VAR] [IDENTIFIER:"x"] [EQUAL] [NUMBER:42] [SEMICOLON]
```

### ¿Por qué es Necesario?
Sin análisis léxico, el parser tendría que trabajar directamente con caracteres individuales, lo que sería extremadamente complejo. Los tokens proporcionan una abstracción que simplifica las fases posteriores.

---

## Anatomía de un Token

### Estructura de Token en go-r2lox
```go
type Token struct {
    Type    TokenType     // Tipo de token (VAR, NUMBER, etc.)
    Lexeme  string        // Texto original del código fuente
    Literal interface{}   // Valor procesado (para números, strings)
    Line    int          // Número de línea (para error reporting)
}
```

### Ejemplo Práctico
```lox
var precio = 19.99;
```

Se convierte en tokens:
```
Token{Type: VAR,        Lexeme: "var",   Literal: nil,    Line: 1}
Token{Type: IDENTIFIER, Lexeme: "precio", Literal: nil,    Line: 1}
Token{Type: EQUAL,      Lexeme: "=",     Literal: nil,    Line: 1}
Token{Type: NUMBER,     Lexeme: "19.99", Literal: 19.99,  Line: 1}
Token{Type: SEMICOLON,  Lexeme: ";",     Literal: nil,    Line: 1}
Token{Type: EOF,        Lexeme: "",      Literal: nil,    Line: 1}
```

---

## Tipos de Tokens en Lox

### 1. Tokens de Puntuación
```lox
( ) { } [ ] , . ; : ?
```

Ejemplos en contexto:
```lox
fun saludar(nombre) {        // ( ) {
    var mensaje = "Hola";    // ;
    var array = [1, 2, 3];   // [ , ]
    return mensaje;          // ;
}                            // }
```

### 2. Operadores
```lox
// Aritméticos
+ - * / ** %

// Comparación  
> >= < <= == !=

// Lógicos
and or not

// Asignación
=

// Incremento/Decremento
++ --
```

### 3. Literales
```lox
// Números
42
3.14159
-17
1.23e4

// Cadenas
"Hola mundo"
"Línea 1\nLínea 2"
"""
Cadena
multilínea
"""

// Booleanos
true
false

// Nil
nil
```

### 4. Identificadores y Palabras Clave
```lox
// Palabras clave reservadas
var let const fun class if else while for return

// Identificadores válidos
variable
_privado
camelCase
snake_case
Variable123
```

---

## Implementación del Scanner

### Estructura Principal
```go
type Scanner struct {
    Source  string    // Código fuente completo
    Tokens  []Token   // Lista de tokens generados
    Start   int       // Inicio del token actual
    Current int       // Posición actual en el código
    Line    int       // Línea actual
}
```

### Proceso de Scanning

#### 1. Bucle Principal
```go
func (s *Scanner) ScanTokens() []Token {
    for !s.isAtEnd() {
        s.Start = s.Current  // Marcar inicio del próximo token
        s.scanToken()        // Procesar el token
    }
    
    s.Tokens = append(s.Tokens, NewToken(EOF, "", nil, s.Line))
    return s.Tokens
}
```

#### 2. Reconocimiento de Tokens
```go
func (s *Scanner) scanToken() {
    c := s.advance()  // Obtener siguiente carácter
    
    switch c {
    case '(':
        s.addToken(LEFT_PAREN, "(")
    case ')':
        s.addToken(RIGHT_PAREN, ")")
    case '+':
        if s.match('+') {  // Verificar si es ++
            s.addToken(PLUS_PLUS, "++")
        } else {
            s.addToken(PLUS, "+")
        }
    case '"':
        s.string()  // Procesar cadena
    default:
        if s.isDigit(c) {
            s.number()      // Procesar número
        } else if s.isAlpha(c) {
            s.identifier()  // Procesar identificador
        } else {
            Errors(s.Line, "Unexpected character.")
        }
    }
}
```

### Casos Especiales

#### Números
```go
func (s *Scanner) number() interface{} {
    // Consumir dígitos enteros
    for s.isDigit(s.peek()) {
        s.advance()
    }
    
    // Buscar parte decimal
    if s.peek() == '.' && s.isDigit(s.peekNext()) {
        s.advance()  // Consumir el '.'
        
        for s.isDigit(s.peek()) {
            s.advance()
        }
    }
    
    // Convertir a float64
    value, err := strconv.ParseFloat(s.Source[s.Start:s.Current], 64)
    if err != nil {
        Errors(s.Line, "Error parsing number.")
    }
    
    s.addToken(NUMBER, value)
    return value
}
```

#### Cadenas
```go
func (s *Scanner) string() {
    for s.peek() != '"' && !s.isAtEnd() {
        if s.peek() == '\n' {
            s.Line++  // Contar líneas en strings multilínea
        }
        s.advance()
    }
    
    if s.isAtEnd() {
        Errors(s.Line, "Unterminated string.")
        return
    }
    
    s.advance()  // Consumir " final
    
    // Extraer contenido (sin comillas)
    value := s.Source[s.Start+1 : s.Current-1]
    
    // Procesar secuencias de escape
    value, err := strconv.Unquote("\"" + value + "\"")
    if err != nil {
        Errors(s.Line, "Error parsing string.")
    }
    
    s.addToken(STRING, value)
}
```

#### Identificadores y Palabras Clave
```go
func (s *Scanner) identifier() {
    for s.isAlphaNumeric(s.peek()) {
        s.advance()
    }
    
    text := s.Source[s.Start:s.Current]
    
    // Verificar si es palabra clave
    tokenType, ok := keywords[text]
    if !ok {
        tokenType = IDENTIFIER
    }
    
    s.addToken(tokenType, text)
}

// Mapa de palabras clave
var keywords = map[string]TokenType{
    "and":    AND,
    "class":  CLASS,
    "else":   ELSE,
    "false":  FALSE,
    "for":    FOR,
    "fun":    FUN,
    "if":     IF,
    "nil":    NIL,
    "or":     OR,
    "return": RETURN,
    "true":   TRUE,
    "var":    VAR,
    "while":  WHILE,
}
```

---

## Experimentos Prácticos

### Herramienta de Análisis de Tokens

Puedes crear un programa simple para ver cómo se tokeniza el código:

```go
// token_analyzer.go
package main

import (
    "fmt"
    "os"
    "./coati2lang"
)

func main() {
    if len(os.Args) != 2 {
        fmt.Println("Usage: go run token_analyzer.go 'código lox'")
        return
    }
    
    source := os.Args[1]
    scanner := coati2lang.NewScanner(source)
    tokens := scanner.ScanTokens()
    
    fmt.Printf("Código: %s\n\n", source)
    fmt.Println("Tokens:")
    for i, token := range tokens {
        fmt.Printf("%d: %s\n", i+1, tokenString(token))
    }
}

func tokenString(token coati2lang.Token) string {
    if token.Literal != nil {
        return fmt.Sprintf("%-15s %-10s %v", 
            token.Type, token.Lexeme, token.Literal)
    }
    return fmt.Sprintf("%-15s %s", token.Type, token.Lexeme)
}
```

### Ejemplos de Tokenización

#### Expresión Aritmética
```lox
var resultado = (5 + 3) * 2;
```

Tokens resultantes:
```
VAR         var
IDENTIFIER  resultado
EQUAL       =
LEFT_PAREN  (
NUMBER      5           5
PLUS        +
NUMBER      3           3
RIGHT_PAREN )
STAR        *
NUMBER      2           2
SEMICOLON   ;
EOF
```

#### Función con Condicional
```lox
fun maximo(a, b) {
    if (a > b) {
        return a;
    } else {
        return b;
    }
}
```

Tokens resultantes:
```
FUN         fun
IDENTIFIER  maximo
LEFT_PAREN  (
IDENTIFIER  a
COMMA       ,
IDENTIFIER  b
RIGHT_PAREN )
LEFT_BRACE  {
IF          if
LEFT_PAREN  (
IDENTIFIER  a
GREATER     >
IDENTIFIER  b
RIGHT_PAREN )
LEFT_BRACE  {
RETURN      return
IDENTIFIER  a
SEMICOLON   ;
RIGHT_BRACE }
ELSE        else
LEFT_BRACE  {
RETURN      return
IDENTIFIER  b
SEMICOLON   ;
RIGHT_BRACE }
RIGHT_BRACE }
EOF
```

---

## Manejo de Errores Léxicos

### Tipos de Errores

#### 1. Caracteres Inesperados
```lox
var x = 42@;  // '@' no es válido en Lox
```
Error: "Unexpected character '@' at line 1"

#### 2. Cadenas No Terminadas
```lox
var mensaje = "Hola mundo;  // Falta comilla de cierre
```
Error: "Unterminated string at line 1"

#### 3. Números Malformados
```lox
var numero = 123.456.789;  // Múltiples puntos decimales
```
Error: "Invalid number format at line 1"

### Estrategias de Recuperación
```go
func (s *Scanner) scanToken() {
    c := s.advance()
    
    switch c {
    // ... casos válidos ...
    default:
        if s.isDigit(c) {
            s.number()
        } else if s.isAlpha(c) {
            s.identifier()
        } else {
            // Reportar error pero continuar scanning
            Errors(s.Line, "Unexpected character.")
            // No terminar el programa, seguir procesando
        }
    }
}
```

---

## Optimizaciones del Scanner

### 1. Reconocimiento de Operadores Compuestos
```go
func (s *Scanner) scanOperator(c rune) {
    switch c {
    case '+':
        if s.match('+') {
            s.addToken(PLUS_PLUS, "++")
        } else if s.match('=') {
            s.addToken(PLUS_EQUAL, "+=")  // Futuro
        } else {
            s.addToken(PLUS, "+")
        }
    case '=':
        if s.match('=') {
            s.addToken(EQUAL_EQUAL, "==")
        } else if s.match('>') {
            s.addToken(ARROW, "=>")
        } else {
            s.addToken(EQUAL, "=")
        }
    }
}
```

### 2. Manejo de Comentarios
```go
case '/':
    if s.match('/') {
        // Comentario de línea: consumir hasta \n
        for s.peek() != '\n' && !s.isAtEnd() {
            s.advance()
        }
    } else if s.match('*') {
        // Comentario de bloque (futuro)
        s.blockComment()
    } else {
        s.addToken(SLASH, "/")
    }
```

### 3. Optimización de Memoria
```go
// Reutilizar strings para tokens comunes
var commonTokens = map[string]string{
    "(": "(",
    ")": ")",
    "{": "{",
    "}": "}",
    // ... más tokens comunes
}

func (s *Scanner) addToken(tokenType TokenType, lexeme string) {
    if common, exists := commonTokens[lexeme]; exists {
        lexeme = common  // Reutilizar string común
    }
    
    token := Token{
        Type:    tokenType,
        Lexeme:  lexeme,
        Literal: nil,
        Line:    s.Line,
    }
    s.Tokens = append(s.Tokens, token)
}
```

---

## Ejercicios Prácticos

### Ejercicio 1: Analizar Tokens
Usa la herramienta de análisis para examinar estos fragmentos de código:

```lox
// Fragmento 1
for (var i = 0; i < 10; i++) {
    print i;
}

// Fragmento 2  
fun fibonacci(n) {
    if (n <= 1) return n;
    return fibonacci(n-1) + fibonacci(n-2);
}

// Fragmento 3
var array = [1, 2, 3];
var objeto = {"clave": "valor"};
```

### Ejercicio 2: Identificar Errores
¿Qué errores léxicos tienen estos fragmentos?

```lox
// Error 1
var nombre = "Juan;

// Error 2  
var numero = 123.45.67;

// Error 3
var especial = 42#;

// Error 4
fun test() {
    /* comentario no soportado */
}
```

### Ejercicio 3: Extensión del Scanner
Diseña cómo agregarías soporte para:
1. Comentarios de bloque `/* ... */`
2. Strings con comillas simples `'texto'`
3. Números binarios `0b1010`
4. Números hexadecimales `0xFF`

---

## Comparación con Otros Lenguajes

### Python
```python
# Python tokeniza esto:
def suma(a, b):
    return a + b

# En tokens similares a:
[NAME:'def', NAME:'suma', OP:'(', NAME:'a', OP:',', NAME:'b', OP:')', 
 OP:':', NEWLINE, INDENT, NAME:'return', NAME:'a', OP:'+', NAME:'b', 
 NEWLINE, DEDENT]
```

### JavaScript
```javascript
// JavaScript tokeniza esto:
function suma(a, b) {
    return a + b;
}

// Con herramientas como Babel parser
```

### Ventajas del Enfoque de Lox
- **Simplicidad**: Menos tipos de tokens que lenguajes complejos
- **Claridad**: Cada token tiene un propósito claro
- **Extensibilidad**: Fácil agregar nuevos tipos de tokens

---

## Herramientas de Debugging

### Visualización de Tokens
```go
func printTokens(tokens []Token) {
    fmt.Println("╔══════════════════════════════════════════╗")
    fmt.Println("║               TOKEN ANALYSIS             ║")
    fmt.Println("╠═══════════╦══════════════╦═══════════════╣")
    fmt.Println("║   TYPE    ║    LEXEME    ║    LITERAL    ║")
    fmt.Println("╠═══════════╬══════════════╬═══════════════╣")
    
    for _, token := range tokens {
        fmt.Printf("║ %-9s ║ %-12s ║ %-13v ║\n",
            token.Type, token.Lexeme, token.Literal)
    }
    
    fmt.Println("╚═══════════╩══════════════╩═══════════════╝")
}
```

### Análisis de Performance
```go
func benchmarkScanner(source string) {
    start := time.Now()
    scanner := NewScanner(source)
    tokens := scanner.ScanTokens()
    duration := time.Since(start)
    
    fmt.Printf("Scanned %d tokens in %v\n", len(tokens), duration)
    fmt.Printf("Rate: %.2f tokens/ms\n", 
        float64(len(tokens))/float64(duration.Nanoseconds())*1000000)
}
```

---

## Resumen del Módulo 2

### Lo que Aprendimos
✅ Qué es el análisis léxico y por qué es importante  
✅ Estructura y tipos de tokens en Lox  
✅ Implementación del scanner en go-r2lox  
✅ Manejo de casos especiales (números, cadenas, identificadores)  
✅ Estrategias de manejo de errores léxicos  
✅ Optimizaciones y consideraciones de performance  

### Conceptos Clave
- **Token**: Unidad básica de significado en el código
- **Lexema**: Texto original del código fuente
- **Scanner**: Componente que realiza análisis léxico
- **Palabras clave**: Tokens reservados con significado especial

### Preparación para el Módulo 3
En el [Módulo 3](curso_r2lox_03.md) aprenderemos sobre:
- Análisis sintáctico y gramáticas
- Construcción del AST (Árbol de Sintaxis Abstracta)
- Implementación del parser
- Manejo de precedencia de operadores

Asegúrate de entender bien cómo funciona el scanner antes de continuar, ya que el parser depende completamente de los tokens que produce el scanner.