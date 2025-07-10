# Roadmap de Desarrollo: go-r2lox

## Resumen Ejecutivo

Este roadmap define la hoja de ruta para la evolución del intérprete go-r2lox, priorizando estabilidad, funcionalidad y rendimiento. Cada fase incluye estimaciones de tiempo y dependencias claras.

## Metodología de Priorización

### Criterios de Evaluación
- **Estabilidad**: Impacto en la confiabilidad del sistema
- **Usabilidad**: Mejora en la experiencia del desarrollador
- **Funcionalidad**: Nuevas capacidades del lenguaje
- **Rendimiento**: Optimizaciones de velocidad y memoria
- **Mantenibilidad**: Facilidad de desarrollo futuro

### Niveles de Prioridad
- 🔥 **Crítico**: Bloquea el desarrollo o uso básico
- ⚡ **Alto**: Mejora significativa en core functionality
- 📈 **Medio**: Características importantes pero no críticas
- 💡 **Bajo**: Mejoras nice-to-have

---

## Fase 1: Estabilización Core (Mes 1)
**Objetivo**: Crear una base sólida y confiable

### 🔥 Prioridad Crítica

#### 1.1 Sistema de Manejo de Errores Robusto
**Estimación**: 7 días  
**Complejidad**: Alta  
**Ubicación**: `coati2lang/r2loxerrors.go`, `coati2lang/interpreter.go`

**Tareas específicas:**
- [ ] Crear tipos de error estructurados (2 días)
  ```go
  type LoxError struct {
      Type     ErrorType
      Message  string
      Line     int
      Column   int
      Context  string
      Stack    []CallFrame
  }
  ```
- [ ] Implementar manejo graceful de errores (2 días)
- [ ] Agregar stack traces detallados (2 días)
- [ ] Crear sistema de recovery para REPL (1 día)

**Dependencias**: Ninguna  
**Bloqueadores**: Uso actual de panic/log.Fatalln

#### 1.2 Suite de Pruebas Básica
**Estimación**: 8 días  
**Complejidad**: Media  
**Ubicación**: Nuevos archivos `*_test.go`

**Tareas específicas:**
- [ ] Tests unitarios para Scanner (2 días)
  - Tokenización correcta
  - Manejo de edge cases
  - Números, strings, identificadores
- [ ] Tests unitarios para Parser (3 días)
  - Expresiones básicas
  - Statements
  - Manejo de errores
- [ ] Tests unitarios para Interpreter (3 días)
  - Evaluación de expresiones
  - Variables y scope
  - Funciones básicas

**Dependencias**: Sistema de errores mejorado  
**Cobertura objetivo**: >70%

#### 1.3 Refactoring de Arquitectura Base
**Estimación**: 5 días  
**Complejidad**: Media  
**Ubicación**: `main.go`, `coati2lang/interpreter.go`

**Tareas específicas:**
- [ ] Separar responsabilidades en main.go (1 día)
- [ ] Crear interfaces claras entre componentes (2 días)
- [ ] Implementar patrón Result para manejo de errores (2 días)

**Dependencias**: Sistema de errores

### ⚡ Prioridad Alta

#### 1.4 Optimizaciones Críticas de Rendimiento
**Estimación**: 6 días  
**Complejidad**: Alta  
**Ubicación**: `coati2lang/enviroment.go`, `coati2lang/interpreter.go`

**Tareas específicas:**
- [ ] Optimizar lookups de variables con maps (2 días)
- [ ] Reducir allocaciones innecesarias (2 días)
- [ ] Implementar pooling de objetos frecuentes (2 días)

**Dependencias**: Tests unitarios  
**Objetivo**: 2x mejora en benchmarks básicos

---

## Fase 2: Funcionalidad Extendida (Mes 2)
**Objetivo**: Expandir capacidades del lenguaje

### ⚡ Prioridad Alta

#### 2.1 Clases y Orientación a Objetos
**Estimación**: 15 días  
**Complejidad**: Muy Alta  
**Ubicación**: `coati2lang/parser.go`, `coati2lang/interpreter.go`

**Tareas específicas:**
- [ ] Diseño de sintaxis de clases (2 días)
  ```lox
  class Person {
      init(name, age) {
          this.name = name;
          this.age = age;
      }
      
      greet() {
          return "Hello, I'm " + this.name;
      }
  }
  ```
- [ ] Implementar parsing de clases (3 días)
- [ ] Crear sistema de instancias (4 días)
- [ ] Implementar métodos y propiedades (3 días)
- [ ] Agregar herencia simple (3 días)

**Dependencias**: Suite de pruebas, sistema de errores  
**Bloqueadores**: Complejidad del sistema de scope

#### 2.2 Sistema de Módulos Básico
**Estimación**: 12 días  
**Complejidad**: Alta  
**Ubicación**: Nuevos archivos `coati2lang/modules.go`

**Tareas específicas:**
- [ ] Diseñar sintaxis de import/export (2 días)
  ```lox
  // math.lox
  export fun sqrt(x) { return x ** 0.5; }
  
  // main.lox
  import "math" as math;
  print math.sqrt(16);
  ```
- [ ] Implementar loader de módulos (4 días)
- [ ] Crear sistema de namespaces (3 días)
- [ ] Manejar dependencias circulares (3 días)

**Dependencias**: Sistema de errores robusto

### 📈 Prioridad Media

#### 2.3 Mejoras al REPL
**Estimación**: 6 días  
**Complejidad**: Media  
**Ubicación**: `main.go`

**Tareas específicas:**
- [ ] Implementar historial de comandos (2 días)
- [ ] Agregar auto-completado básico (2 días)
- [ ] Crear comandos especiales (:help, :exit, :vars) (2 días)

**Dependencias**: Sistema de errores

#### 2.4 Biblioteca Estándar Básica
**Estimación**: 8 días  
**Complejidad**: Media  
**Ubicación**: Nuevos archivos `coati2lang/stdlib/`

**Tareas específicas:**
- [ ] Módulo math (2 días)
  - sin, cos, tan, log, exp
- [ ] Módulo string (2 días)
  - split, join, trim, replace
- [ ] Módulo array (2 días)
  - push, pop, slice, map, filter
- [ ] Módulo io básico (2 días)
  - readFile, writeFile

**Dependencias**: Sistema de módulos

---

## Fase 3: Características Avanzadas (Mes 3)
**Objetivo**: Agregar funcionalidades modernas al lenguaje

### 📈 Prioridad Media

#### 3.1 Sintaxis Extendida
**Estimación**: 10 días  
**Complejidad**: Media  
**Ubicación**: `coati2lang/scanner.go`, `coati2lang/parser.go`

**Tareas específicas:**
- [ ] Operador ternario (2 días)
  ```lox
  var result = condition ? "yes" : "no";
  ```
- [ ] Operadores compuestos (2 días)
  ```lox
  x += 5;
  str *= 3;
  ```
- [ ] Destructuring assignment (3 días)
  ```lox
  var [a, b] = [1, 2];
  var {name, age} = person;
  ```
- [ ] String interpolation (3 días)
  ```lox
  var message = `Hello ${name}, you are ${age} years old`;
  ```

**Dependencias**: Parser robusto

#### 3.2 Async/Await Básico
**Estimación**: 14 días  
**Complejidad**: Muy Alta  
**Ubicación**: Nuevos archivos `coati2lang/async.go`

**Tareas específicas:**
- [ ] Diseñar modelo de concurrencia (3 días)
- [ ] Implementar Promises/Futures (4 días)
- [ ] Agregar sintaxis async/await (4 días)
- [ ] Crear scheduler básico (3 días)

**Dependencias**: Clases, sistema de errores avanzado

### 💡 Prioridad Baja

#### 3.3 Macros y Metaprogramación
**Estimación**: 12 días  
**Complejidad**: Muy Alta  
**Ubicación**: Nuevos archivos `coati2lang/macros.go`

**Tareas específicas:**
- [ ] Diseñar sistema de macros (4 días)
- [ ] Implementar expansion en parse time (4 días)
- [ ] Crear macros básicas predefinidas (4 días)

**Dependencias**: Parser muy estable

---

## Fase 4: Optimización y Herramientas (Mes 4)
**Objetivo**: Optimizar rendimiento y agregar herramientas de desarrollo

### ⚡ Prioridad Alta

#### 4.1 Variable Resolution y Optimizaciones
**Estimación**: 10 días  
**Complejidad**: Alta  
**Ubicación**: Nuevos archivos `coati2lang/resolver.go`

**Tareas específicas:**
- [ ] Implementar resolver pass (4 días)
- [ ] Optimizar variable lookups (3 días)
- [ ] Agregar constant folding (3 días)

**Dependencias**: Suite de pruebas comprehensiva

#### 4.2 Debugging Tools
**Estimación**: 12 días  
**Complejidad**: Alta  
**Ubicación**: Nuevos archivos `coati2lang/debugger.go`

**Tareas específicas:**
- [ ] Implementar breakpoints (4 días)
- [ ] Crear step-through debugging (4 días)
- [ ] Agregar variable inspection (2 días)
- [ ] Crear debugger REPL (2 días)

**Dependencias**: Sistema de errores avanzado

### 📈 Prioridad Media

#### 4.3 Profiling y Análisis
**Estimación**: 8 días  
**Complejidad**: Media  
**Ubicación**: Nuevos archivos `coati2lang/profiler.go`

**Tareas específicas:**
- [ ] Implementar function call profiling (3 días)
- [ ] Crear memory profiling (3 días)
- [ ] Agregar reporting tools (2 días)

**Dependencias**: Debugging tools

#### 4.4 Compilación a Bytecode (Experimental)
**Estimación**: 20 días  
**Complejidad**: Muy Alta  
**Ubicación**: Nuevos archivos `coati2lang/compiler.go`, `coati2lang/vm.go`

**Tareas específicas:**
- [ ] Diseñar set de instrucciones (5 días)
- [ ] Implementar compilador AST→bytecode (8 días)
- [ ] Crear VM para ejecutar bytecode (7 días)

**Dependencias**: Resolver, optimizaciones

---

## Fase 5: Ecosistema y Comunidad (Mes 5-6)
**Objetivo**: Crear un ecosistema robusto alrededor del lenguaje

### 📈 Prioridad Media

#### 5.1 Package Manager
**Estimación**: 15 días  
**Complejidad**: Alta  
**Ubicación**: Nuevos archivos `cmd/loxpkg/`

**Tareas específicas:**
- [ ] Diseñar formato de paquetes (3 días)
- [ ] Implementar registry básico (5 días)
- [ ] Crear CLI para gestión de paquetes (4 días)
- [ ] Agregar dependency resolution (3 días)

**Dependencias**: Sistema de módulos maduro

#### 5.2 Language Server Protocol
**Estimación**: 18 días  
**Complejidad**: Muy Alta  
**Ubicación**: Nuevos archivos `cmd/lox-lsp/`

**Tareas específicas:**
- [ ] Implementar LSP básico (8 días)
- [ ] Agregar syntax highlighting (3 días)
- [ ] Crear auto-completion avanzado (4 días)
- [ ] Implementar go-to-definition (3 días)

**Dependencias**: Parser muy estable, resolver

### 💡 Prioridad Baja

#### 5.3 Transpilación a JavaScript/WASM
**Estimación**: 25 días  
**Complejidad**: Extrema  
**Ubicación**: Nuevos archivos `coati2lang/transpiler.go`

**Tareas específicas:**
- [ ] Diseñar target JavaScript (8 días)
- [ ] Implementar transpilador (12 días)
- [ ] Crear runtime JS (5 días)

**Dependencias**: Compilador a bytecode

---

## Cronograma Detallado

### Mes 1: Fundación Sólida
```
Semana 1: Sistema de errores + inicio tests
Semana 2: Completar tests unitarios
Semana 3: Refactoring arquitectura
Semana 4: Optimizaciones críticas + buffer
```

### Mes 2: Expansión del Lenguaje
```
Semana 1-2: Implementar clases básicas
Semana 3: Herencia y métodos avanzados
Semana 4: Sistema de módulos básico
```

### Mes 3: Características Modernas
```
Semana 1: Sintaxis extendida
Semana 2: REPL mejorado + stdlib
Semana 3-4: Async/await experimental
```

### Mes 4: Optimización
```
Semana 1: Variable resolution
Semana 2: Debugging tools básico
Semana 3: Profiling
Semana 4: Inicio bytecode compiler
```

---

## Métricas de Éxito

### Fase 1
- [ ] 0 crashes en suite de pruebas
- [ ] >70% cobertura de código
- [ ] 2x mejora en benchmarks
- [ ] REPL no se cierra por errores

### Fase 2
- [ ] Clases funcionales con herencia
- [ ] Sistema de módulos operativo
- [ ] >50 funciones en stdlib
- [ ] REPL con auto-completado

### Fase 3
- [ ] Sintaxis moderna completa
- [ ] Async básico funcional
- [ ] Debugging básico operativo

### Fase 4
- [ ] 5x mejora en rendimiento
- [ ] Debugger completo
- [ ] Bytecode compiler funcional

---

## Gestión de Riesgos

### Riesgos Técnicos
- **Complejidad de OOP**: Mitigar con diseño incremental
- **Rendimiento de bytecode**: Benchmark continuo
- **Async complexity**: Implementación experimental primero

### Riesgos de Recursos
- **Tiempo de desarrollo**: Priorizar características core
- **Mantenimiento**: Documentar todas las decisiones
- **Testing**: Automatizar desde el inicio

### Planes de Contingencia
- **Retraso en OOP**: Posponer herencia
- **Problemas de rendimiento**: Implementar profiling temprano
- **Complejidad de async**: Mover a fase posterior

---

## Conclusiones

Este roadmap balancea ambición con pragmatismo, priorizando:

1. **Estabilidad primero**: Base sólida antes de características avanzadas
2. **Incrementalidad**: Cada fase construye sobre la anterior
3. **Flexibilidad**: Posibilidad de ajustar prioridades según feedback
4. **Calidad**: Testing y documentación integrados en cada fase

La implementación de este roadmap convertirá go-r2lox de un intérprete básico a un lenguaje moderno y robusto, manteniendo la simplicidad y elegancia del diseño original de Lox.