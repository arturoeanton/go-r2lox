# Curso go-r2lox: Módulo 1 - Introducción a Lox y go-r2lox

## Bienvenidos al Curso de go-r2lox

### Objetivos del Módulo
- Entender qué es Lox y por qué estudiarlo
- Familiarizarse con go-r2lox y su instalación
- Escribir y ejecutar los primeros programas en Lox
- Comprender la filosofía detrás del diseño del lenguaje

---

## ¿Qué es Lox?

### Introducción al Lenguaje
Lox es un lenguaje de programación dinámico, interpretado y orientado a objetos, diseñado específicamente para el libro "Crafting Interpreters" de Robert Nystrom. Aunque es un lenguaje educativo, Lox incorpora características de lenguajes modernos como:

- **Tipado dinámico**: Los tipos se determinan en tiempo de ejecución
- **Gestión automática de memoria**: No necesitas manejar memoria manualmente
- **Funciones de primera clase**: Las funciones pueden ser valores
- **Closures**: Funciones que capturan su entorno léxico
- **Clases y herencia**: Programación orientada a objetos (en desarrollo)

### Filosofía de Diseño
```lox
// Lox prioriza la simplicidad y claridad
print "¡Hola, mundo!";

// Variables simples
var nombre = "Juan";
var edad = 25;

// Funciones intuitivas
fun saludar(persona) {
    return "¡Hola, " + persona + "!";
}

print saludar(nombre);
```

---

## ¿Qué es go-r2lox?

### Una Implementación en Go
go-r2lox es una implementación completa del intérprete Lox escrita en Go. Sigue la segunda implementación del libro "Crafting Interpreters", usando un intérprete basado en árbol (tree-walking interpreter).

### Ventajas de go-r2lox
- **Rendimiento**: Go ofrece mejor rendimiento que Java o Python
- **Simplicidad**: Código claro y fácil de entender
- **Portabilidad**: Ejecuta en cualquier plataforma donde corra Go
- **Extensibilidad**: Fácil de modificar y extender

---

## Instalación y Configuración

### Prerequisitos
```bash
# Verificar instalación de Go (versión 1.18+)
go version

# Si no tienes Go instalado, descárgalo de:
# https://golang.org/download/
```

### Instalación de go-r2lox
```bash
# Clonar el repositorio
git clone https://github.com/arturoeanton/go-r2lox.git
cd go-r2lox

# Construir el intérprete
go build

# Verificar la instalación
./go-r2lox --help
```

### Primer Programa
```bash
# Crear un archivo de prueba
echo 'print "¡Hola desde Lox!";' > hola.lox

# Ejecutar el programa
./go-r2lox hola.lox
```

---

## Conceptos Fundamentales de Lox

### 1. Variables y Tipos

#### Números
```lox
var entero = 42;
var decimal = 3.14159;
var negativo = -17;
var cientifico = 1.23e4;  // 12300

print entero + decimal;   // 45.14159
```

#### Cadenas
```lox
var saludo = "Hola";
var nombre = "María";
var mensaje = saludo + ", " + nombre + "!";

print mensaje;  // "Hola, María!"

// Cadenas multilínea
var poema = """
Roses are red,
Violets are blue,
Lox is awesome,
And so are you!
""";
```

#### Booleanos y Nil
```lox
var verdadero = true;
var falso = false;
var vacio = nil;

print verdadero and falso;  // false
print verdadero or falso;   // true
print not verdadero;        // false
```

### 2. Operadores

#### Aritméticos
```lox
print 5 + 3;     // 8
print 10 - 4;    // 6
print 6 * 7;     // 42
print 15 / 3;    // 5
print 2 ** 3;    // 8 (potencia)
print 50 % 20;   // 10 (porcentaje: 50 * 20 / 100)
```

#### Comparación
```lox
print 5 > 3;     // true
print 2 <= 8;    // true
print 4 == 4;    // true
print 5 != 3;    // true
```

#### Lógicos
```lox
print true and false;   // false
print true or false;    // true
print not true;         // false

// Short-circuit evaluation
print nil or "valor";   // "valor"
print "primero" or "segundo";  // "primero"
```

### 3. Estructuras de Control

#### Condicionales
```lox
var edad = 18;

if (edad >= 18) {
    print "Eres mayor de edad";
} else {
    print "Eres menor de edad";
}

// Múltiples condiciones
var calificacion = 85;

if (calificacion >= 90) {
    print "Excelente";
} else if (calificacion >= 80) {
    print "Muy bien";
} else if (calificacion >= 70) {
    print "Bien";
} else {
    print "Necesitas mejorar";
}
```

#### Bucles
```lox
// While loop
var i = 1;
while (i <= 5) {
    print "Número: " + i;
    i = i + 1;
}

// For loop
for (var j = 0; j < 3; j = j + 1) {
    print "Iteración: " + j;
}
```

---

## Primer Proyecto: Calculadora Simple

Vamos a crear una calculadora básica que demuestre los conceptos aprendidos:

```lox
// calculadora.lox
print "=== Calculadora Simple ===";

// Función para sumar
fun sumar(a, b) {
    return a + b;
}

// Función para restar
fun restar(a, b) {
    return a - b;
}

// Función para multiplicar
fun multiplicar(a, b) {
    return a * b;
}

// Función para dividir
fun dividir(a, b) {
    if (b == 0) {
        print "Error: División por cero";
        return nil;
    }
    return a / b;
}

// Función principal de la calculadora
fun calculadora(operacion, a, b) {
    if (operacion == "suma") {
        return sumar(a, b);
    } else if (operacion == "resta") {
        return restar(a, b);
    } else if (operacion == "multiplicacion") {
        return multiplicar(a, b);
    } else if (operacion == "division") {
        return dividir(a, b);
    } else {
        print "Operación no válida";
        return nil;
    }
}

// Ejemplos de uso
var resultado1 = calculadora("suma", 10, 5);
print "10 + 5 = " + resultado1;

var resultado2 = calculadora("multiplicacion", 7, 8);
print "7 * 8 = " + resultado2;

var resultado3 = calculadora("division", 20, 4);
print "20 / 4 = " + resultado3;

// Prueba de división por cero
var resultado4 = calculadora("division", 10, 0);
```

### Ejecutar la Calculadora
```bash
# Guardar el código en calculadora.lox
./go-r2lox calculadora.lox
```

---

## Modo REPL (Read-Eval-Print Loop)

El REPL permite experimentar con Lox de forma interactiva:

```bash
# Iniciar REPL
./go-r2lox

# Ahora puedes escribir código Lox línea por línea:
> var x = 42;
> var y = 8;
> print x + y;
50
> fun doble(n) { return n * 2; }
> print doble(21);
42
> // Salir con Ctrl+C
```

### Ejemplos de Sesión REPL
```lox
// Experimento con variables
> var nombre = "Ana";
> var apellido = "García";
> var completo = nombre + " " + apellido;
> print completo;
Ana García

// Experimento con funciones
> fun factorial(n) {
    if (n <= 1) return 1;
    return n * factorial(n - 1);
}
> print factorial(5);
120

// Experimento con closures
> fun contador() {
    var cuenta = 0;
    return fun() {
        cuenta = cuenta + 1;
        return cuenta;
    };
}
> var miContador = contador();
> print miContador();
1
> print miContador();
2
```

---

## Ejercicios Prácticos

### Ejercicio 1: Variables y Operadores
Crea un programa que:
1. Declare variables para tu nombre, edad y ciudad
2. Calcule tu año de nacimiento (año actual - edad)
3. Imprima una presentación personalizada

```lox
// Tu solución aquí
var nombre = "Tu nombre";
var edad = 25;
var ciudad = "Tu ciudad";
var anoActual = 2024;

// Completar...
```

### Ejercicio 2: Condicionales
Escribe un programa que determine si un número es:
- Positivo, negativo o cero
- Par o impar
- Mayor, menor o igual a 100

### Ejercicio 3: Bucles
Crea un programa que:
1. Imprima los números del 1 al 10
2. Calcule la suma de los números del 1 al N
3. Imprima la tabla de multiplicar de un número

### Ejercicio 4: Funciones
Implementa las siguientes funciones:
```lox
// Función que determine si un número es primo
fun esPrimo(numero) {
    // Tu implementación
}

// Función que calcule el máximo de tres números
fun maximo(a, b, c) {
    // Tu implementación
}

// Función que convierta temperatura Celsius a Fahrenheit
fun celsiusAFahrenheit(celsius) {
    // Tu implementación
}
```

---

## Recursos Adicionales

### Documentación
- [README principal](README.md) - Visión general del proyecto
- [Guía de implementación](implementation_book.md) - Detalles técnicos
- [Referencia de arquitectura](arquitectura.md) - Diseño del sistema

### Enlaces Útiles
- [Crafting Interpreters](http://craftinginterpreters.com/) - El libro original
- [Repositorio go-r2lox](https://github.com/arturoeanton/go-r2lox)
- [Documentación de Go](https://golang.org/doc/)

### Próximo Módulo
En el [Módulo 2](curso_r2lox_02.md) aprenderemos sobre:
- Análisis léxico en profundidad
- Tokens y palabras clave
- Manejo de errores lexicales
- Implementación del scanner

---

## Resumen del Módulo 1

### Lo que Aprendimos
✅ Qué es Lox y su filosofía de diseño  
✅ Instalación y configuración de go-r2lox  
✅ Tipos de datos básicos (números, cadenas, booleanos, nil)  
✅ Operadores aritméticos, de comparación y lógicos  
✅ Estructuras de control (if/else, while, for)  
✅ Funciones básicas  
✅ Uso del REPL para experimentación  

### Conceptos Clave
- **Tipado dinámico**: Los tipos se determinan en tiempo de ejecución
- **Sintaxis familiar**: Similar a C/Java/JavaScript
- **Funciones como valores**: Base para conceptos avanzados
- **REPL**: Herramienta para experimentación rápida

### Preparación para el Módulo 2
Asegúrate de:
- Tener go-r2lox funcionando correctamente
- Haber completado al menos 2 de los ejercicios prácticos
- Entender la diferencia entre variables, expresiones y sentencias
- Estar cómodo usando el REPL

¡Felicidades por completar el Módulo 1! En el próximo módulo profundizaremos en cómo go-r2lox procesa el código que escribes.