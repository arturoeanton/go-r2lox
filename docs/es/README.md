# go-r2lox: Documentación Completa

## Índice de Documentación

### 📚 Documentación Principal
- **[README.md](README.md)** - Esta página (introducción y navegación)
- **[implementation_book.md](implementation_book.md)** - Guía completa de implementación
- **[mejoras.md](mejoras.md)** - Análisis de mejoras futuras

### 🏗️ Arquitectura y Análisis
- **[arquitectura.md](arquitectura.md)** - Arquitectura profunda del sistema
- **[analisis_general.md](analisis_general.md)** - Análisis general del proyecto
- **[analisis_funcional.md](analisis_funcional.md)** - Análisis funcional detallado
- **[analisis_tecnico.md](analisis_tecnico.md)** - Análisis técnico profundo
- **[analisis_deep.md](analisis_deep.md)** - Análisis profundo de implementación

### 🗺️ Planificación y Desarrollo
- **[roadmap.md](roadmap.md)** - Hoja de ruta del proyecto
- **[issues.md](issues.md)** - Issues categorizados por prioridad
- **[brainstorm.md](brainstorm.md)** - Ideas y conceptos para el futuro

### 🎓 Curso y Tutoriales
- **[curso_r2lox_01.md](curso_r2lox_01.md)** - Introducción a Lox y go-r2lox
- **[curso_r2lox_02.md](curso_r2lox_02.md)** - Análisis léxico y tokens
- **[curso_r2lox_03.md](curso_r2lox_03.md)** - Análisis sintáctico y AST
- **[curso_r2lox_04.md](curso_r2lox_04.md)** - Intérprete y evaluación
- **[curso_r2lox_05.md](curso_r2lox_05.md)** - Funciones y closures
- **[curso_r2lox_06.md](curso_r2lox_06.md)** - Optimización y futuro

---

## ¿Qué es go-r2lox?

go-r2lox es una implementación en Go del lenguaje de programación Lox, inspirada en la segunda implementación del libro "Crafting Interpreters" de Robert Nystrom. Este proyecto demuestra cómo construir un intérprete completo usando Go, siguiendo las mejores prácticas de diseño de lenguajes de programación.

## Características Principales

### 🚀 Intérprete Completo
- **Análisis Léxico**: Tokenización completa del código fuente
- **Análisis Sintáctico**: Parser recursivo descendente
- **Evaluación**: Intérprete basado en árbol de sintaxis abstracta (AST)
- **Gestión de Memoria**: Manejo automático de variables y scope

### 💻 Funcionalidades del Lenguaje
- **Tipos de Datos**: Números, cadenas, booleanos, nil
- **Variables**: Declaración, asignación y scoping
- **Funciones**: Declaración, llamadas y closures
- **Estructuras de Control**: if/else, while, for
- **Operadores**: Aritméticos, lógicos, de comparación
- **Colecciones**: Arrays y mapas con sintaxis intuitiva

### 🛠️ Herramientas de Desarrollo
- **REPL**: Modo interactivo para pruebas rápidas
- **Ejecución de Archivos**: Intérprete de scripts completos
- **Manejo de Errores**: Reportes detallados con números de línea
- **Extensibilidad**: Arquitectura modular para agregar características

## Arquitectura del Sistema

### Flujo de Ejecución
1. **Entrada** → Código fuente Lox
2. **Scanner** → Tokens
3. **Parser** → AST (Árbol de Sintaxis Abstracta)
4. **Interpreter** → Evaluación y ejecución
5. **Salida** → Resultados

### Componentes Principales
- **`main.go`**: Punto de entrada y coordinación
- **`scanner.go`**: Análisis léxico
- **`parser.go`**: Análisis sintáctico
- **`interpreter.go`**: Evaluación y ejecución
- **`enviroment.go`**: Gestión de variables y scope

## Comenzando

### Instalación
```bash
git clone https://github.com/arturoeanton/go-r2lox.git
cd go-r2lox
go build
```

### Uso Básico
```bash
# Ejecutar un archivo
./go-r2lox script.lox

# Modo REPL
./go-r2lox

# Con parámetros
./go-r2lox -script=mi_archivo.lox
```

### Primer Programa
```lox
// Hola mundo
print "¡Hola, mundo desde Lox!";

// Variables
var nombre = "Arturo";
var edad = 30;
print "Mi nombre es " + nombre + " y tengo " + edad + " años.";

// Funciones
fun saludar(nombre) {
    return "¡Hola, " + nombre + "!";
}

print saludar("Go-R2Lox");
```

## Documentación Detallada

### Para Desarrolladores
- **[implementation_book.md](implementation_book.md)**: Guía completa de implementación
- **[arquitectura.md](arquitectura.md)**: Arquitectura profunda del sistema
- **[analisis_tecnico.md](analisis_tecnico.md)**: Análisis técnico detallado

### Para Usuarios
- **[curso_r2lox_01.md](curso_r2lox_01.md)**: Introducción al lenguaje Lox
- **[curso_r2lox_02.md](curso_r2lox_02.md)**: Sintaxis y características básicas
- **[curso_r2lox_03.md](curso_r2lox_03.md)**: Funciones y programación avanzada

### Para Contribuidores
- **[roadmap.md](roadmap.md)**: Hoja de ruta del proyecto
- **[issues.md](issues.md)**: Issues priorizados y categorizados
- **[mejoras.md](mejoras.md)**: Análisis de mejoras futuras

## Estado del Proyecto

### ✅ Implementado
- Análisis léxico completo
- Parser recursivo descendente
- Intérprete funcional
- Variables y funciones
- Estructuras de control
- Arrays y mapas básicos
- REPL interactivo

### 🚧 En Desarrollo
- Mejoras en manejo de errores
- Optimizaciones de rendimiento
- Biblioteca estándar expandida
- Sistema de módulos

### 📋 Planeado
- Clases y orientación a objetos
- Debugging tools
- Compilación a bytecode
- Optimizaciones avanzadas

## Contribuir

### Cómo Contribuir
1. Lee la documentación completa
2. Revisa [issues.md](issues.md) para tareas disponibles
3. Consulta [roadmap.md](roadmap.md) para el plan futuro
4. Sigue las guías en [analisis_tecnico.md](analisis_tecnico.md)

### Áreas que Necesitan Ayuda
- **Testing**: Expansión de la suite de pruebas
- **Optimización**: Mejoras de rendimiento
- **Documentación**: Ejemplos y tutoriales
- **Características**: Nuevas funcionalidades del lenguaje

## Recursos Adicionales

### Links Útiles
- [Crafting Interpreters](http://craftinginterpreters.com/) - Libro original
- [Lox Language Reference](http://craftinginterpreters.com/appendix-i.html) - Especificación del lenguaje
- [Go Documentation](https://golang.org/doc/) - Documentación de Go

### Comunidad
- **Issues**: Reporta bugs y solicita características
- **Discussions**: Discusiones sobre implementación
- **Wiki**: Documentación colaborativa

---

## Navegación Rápida

### 🎯 Nuevo en el Proyecto
1. [curso_r2lox_01.md](curso_r2lox_01.md) - Empieza aquí
2. [implementation_book.md](implementation_book.md) - Implementación detallada
3. [arquitectura.md](arquitectura.md) - Arquitectura del sistema

### 👨‍💻 Desarrollador Experimentado
1. [analisis_tecnico.md](analisis_tecnico.md) - Análisis técnico
2. [analisis_deep.md](analisis_deep.md) - Análisis profundo
3. [mejoras.md](mejoras.md) - Mejoras futuras

### 🚀 Contribuidor
1. [roadmap.md](roadmap.md) - Plan de desarrollo
2. [issues.md](issues.md) - Tareas disponibles
3. [brainstorm.md](brainstorm.md) - Ideas y conceptos

---

*Esta documentación está en constante evolución. Si encuentras errores o tienes sugerencias, por favor crea un issue o contribuye directamente.*