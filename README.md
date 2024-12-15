# go-r2lox: Un intérprete de Lox en Go

Este proyecto es una implementación en Go del lenguaje de programación Lox, tal como se presenta en el libro "Crafting Interpreters" de Robert Nystrom. Se utiliza la segunda implementación del libro, basada en un árbol de sintaxis abstracta (AST).

## ¿Qué es Lox?

Lox es un lenguaje de scripting dinámico, imperativo y orientado a objetos. Está diseñado para ser incrustado en aplicaciones y para aprender sobre la construcción de lenguajes.

## ¿Qué es "Crafting Interpreters"?

"Crafting Interpreters" (http://craftinginterpreters.com/) es un libro que guía al lector a través de la creación de un intérprete para un lenguaje llamado Lox, presentando dos implementaciones: una en Java (jlox) y otra en C (clox). Este proyecto se basa en los conceptos de la segunda implementación (clox), adaptándolos a Go.

## Características

*   **Intérprete basado en AST:** Se construye un árbol de sintaxis abstracta a partir del código fuente, lo que permite un análisis más profundo y facilita futuras optimizaciones.
*   **Implementado en Go:** Aprovecha las características de Go, como su concurrencia y rendimiento.
*   **Siguiendo "Crafting Interpreters":** Se adhiere a la estructura y los conceptos presentados en el libro.

## Cómo usar go-r2lox

### Prerrequisitos

*   Go instalado (versión 1.18 o superior probablemente funcione)

### Instalación

Para clonar el repositorio y construir el intérprete, ejecuta los siguientes comandos:
´´´
    git clone https://github.com/arturoeanton/go-r2lox.git
    cd go-r2lox
    go build
´´´
### Ejecución

Para ejecutar un archivo Lox, usa el siguiente comando:
´´´
    ./go-r2lox <ruta_al_archivo.lox>
´´´
Por ejemplo, para ejecutar un archivo llamado `test.lox`, ejecuta:
´´´
    ./go-r2lox test.lox
´´´
También puedes ejecutar el intérprete en modo REPL (Read-Eval-Print Loop) ejecutando simplemente:
´´´
    ./go-r2lox
´´´
Esto te permitirá ingresar expresiones Lox directamente en la terminal.

### Ejemplos de código Lox

Aquí hay un ejemplo simple de código Lox:
´´´
    print "Hola, mundo!";

    var a = 10;
    var b = 20;
    print a + b;

    fun add(x, y) {
        return x + y;
    }

    print add(5, 3);
´´´
## Estado del proyecto

Este proyecto se encuentra en desarrollo y aún no implementa todas las características de Lox. Se planea seguir avanzando según lo propuesto en "Crafting Interpreters".

## Contribuciones

Las contribuciones son bienvenidas. Si encuentras errores o tienes sugerencias de mejora, por favor abre un issue o un pull request.

