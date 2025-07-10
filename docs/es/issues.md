# Issues Categorizados: go-r2lox

## Resumen Ejecutivo

Este documento categoriza todos los issues identificados en el proyecto go-r2lox, organizados por tema, prioridad, complejidad y ubicación específica en el código. Cada issue incluye estimación de tiempo y análisis de impacto.

---

## 🔥 Issues Críticos (Prioridad 1)

### Error Handling (Manejo de Errores)

#### ERR-001: Sistema de Errores Frágil
**Prioridad**: 🔥 Crítica  
**Complejidad**: Alta  
**Estimación**: 7 días  
**Ubicación**: `coati2lang/r2loxerrors.go:5-15`, `coati2lang/interpreter.go:319-321`

**Problema**:
```go
// Línea 319 en interpreter.go
log.Fatalln("Undefined variable '" + expr.Name.Lexeme + "'.")
```

**Impacto**:
- El programa termina abruptamente en cualquier error
- REPL se cierra completamente por errores menores
- No hay información de contexto para debugging
- Imposible recuperarse de errores

**Solución Propuesta**:
```go
type LoxError struct {
    Type     ErrorType
    Message  string
    Line     int
    Column   int
    Token    Token
    Context  string
    Stack    []CallFrame
}

func (i *Interpreter) reportError(err LoxError) {
    // Manejo graceful sin terminar programa
}
```

**Archivos Afectados**:
- `coati2lang/r2loxerrors.go` (refactoring completo)
- `coati2lang/interpreter.go` (múltiples líneas)
- `coati2lang/parser.go` (líneas 228, 244, 301)
- `main.go` (manejo de errores global)

#### ERR-002: Panics No Controlados en Parser
**Prioridad**: 🔥 Crítica  
**Complejidad**: Media  
**Estimación**: 3 días  
**Ubicación**: `coati2lang/parser.go:228-229`, `coati2lang/parser.go:298-303`

**Problema**:
```go
// Línea 228 en parser.go
panic(message1)

// Líneas 298-303
defer func() {
    if r := recover(); r != nil {
        fmt.Fprintf(os.Stderr, "%s", r)
        os.Exit(2)
    }
}()
```

**Impacto**:
- Parser crash termina todo el programa
- No hay sincronización de errores
- Múltiples errores no se reportan

**Solución**:
- Implementar error recovery en parser
- Agregar sincronización de tokens
- Colectar múltiples errores antes de abortar

### Memory Management (Gestión de Memoria)

#### MEM-001: Leaks de Memoria en Entornos
**Prioridad**: 🔥 Crítica  
**Complejidad**: Media  
**Estimación**: 4 días  
**Ubicación**: `coati2lang/enviroment.go:15-25`, `coati2lang/interpreter.go:149-162`

**Problema**:
```go
// enviroment.go - sin limpieza de referencias
type Enviroment struct {
    values    map[string]interface{}
    enclosing *Enviroment  // Referencia circular potencial
}
```

**Impacto**:
- Memoria creciente en ejecuciones largas
- Referencias circulares en closures profundos
- No hay garbage collection explícito

**Solución**:
- Implementar weak references donde sea apropiado
- Agregar cleanup explícito de entornos
- Optimizar reutilización de entornos

#### MEM-002: Allocaciones Excesivas en Evaluación
**Prioridad**: 🔥 Crítica  
**Complejidad**: Alta  
**Estimación**: 5 días  
**Ubicación**: `coati2lang/interpreter.go:447-473`

**Problema**:
```go
// Línea 447-473 - full_evaluate crea nuevos slices/maps frecuentemente
var values []interface{} = make([]interface{}, len(expr_a))
var values map[interface{}]interface{} = make(map[interface{}]interface{})
```

**Impacto**:
- Alto uso de memoria en operaciones simples
- Pressure en garbage collector
- Rendimiento degradado

---

## ⚡ Issues de Alta Prioridad (Prioridad 2)

### Performance (Rendimiento)

#### PERF-001: Variable Lookups Ineficientes
**Prioridad**: ⚡ Alta  
**Complejidad**: Media  
**Estimación**: 4 días  
**Ubicación**: `coati2lang/enviroment.go:27-35`

**Problema**:
```go
// Búsqueda lineal en cadena de entornos
func (e *Enviroment) Get(name string) (interface{}, bool) {
    if value, ok := e.values[name]; ok {
        return value, true
    }
    if e.enclosing != nil {
        return e.enclosing.Get(name)  // Recursión costosa
    }
    return nil, false
}
```

**Impacto**:
- O(n) en profundidad de scope para cada variable
- Penalización en funciones anidadas
- Escalabilidad pobre

**Solución**:
- Implementar variable resolution en parse-time
- Cache de variable locations
- Hash maps optimizados para lookups frecuentes

#### PERF-002: AST Traversal Redundante
**Prioridad**: ⚡ Alta  
**Complejidad**: Alta  
**Estimación**: 6 días  
**Ubicación**: `coati2lang/interpreter.go:475-477`

**Problema**:
```go
// evaluate llamado múltiples veces para mismos nodos
func (i *Interpreter) evaluate(expr Expr) interface{} {
    return expr.AcceptExpr(i)  // Sin caching
}
```

**Impacto**:
- Re-evaluación de sub-expresiones constantes
- Trabajo duplicado en loops
- CPU usage innecesario

### Architecture (Arquitectura)

#### ARCH-001: Acoplamiento Fuerte entre Componentes
**Prioridad**: ⚡ Alta  
**Complejidad**: Alta  
**Estimación**: 8 días  
**Ubicación**: `main.go:20-47`, `coati2lang/interpreter.go:40-52`

**Problema**:
```go
// main.go - lógica mezclada
func run(source string) {
    tokens := coati2lang.ScanTokens(source)
    parse := coati2lang.NewParser(tokens)
    expr := parse.Parse()
    interp := coati2lang.NewInterpreter(expr)
    interp.Interpret()
}
```

**Impacto**:
- Difícil testing de componentes individuales
- No hay separación de responsabilidades
- Inflexibilidad para diferentes use cases

**Solución**:
- Crear interfaces claras entre componentes
- Implementar dependency injection
- Separar pipeline de ejecución

#### ARCH-002: Estado Global en Scanner
**Prioridad**: ⚡ Alta  
**Complejidad**: Media  
**Estimación**: 3 días  
**Ubicación**: `coati2lang/scanner.go:44-60`

**Problema**:
```go
// Estado mutable compartido
type Scanner struct {
    Source  string
    Tokens  []Token
    Start   int      // Estado mutable
    Current int      // Estado mutable
    Line    int      // Estado mutable
}
```

**Impacto**:
- Scanner no es thread-safe
- No reutilizable para múltiples fuentes
- Testing complicado

---

## 📈 Issues de Prioridad Media (Prioridad 3)

### Functionality (Funcionalidad)

#### FUNC-001: Arrays de Tamaño Fijo Ineficientes
**Prioridad**: 📈 Media  
**Complejidad**: Media  
**Estimación**: 3 días  
**Ubicación**: `coati2lang/parser.go:357-375`

**Problema**:
```go
// Declaración de arrays con tamaño fijo innecesario
var array[10] = [1, 2, 3];  // Desperdicia memoria
```

**Impacto**:
- Uso ineficiente de memoria
- Limitaciones artificiales
- Sintaxis confusa

#### FUNC-002: Métodos de String Limitados
**Prioridad**: 📈 Media  
**Complejidad**: Baja  
**Estimación**: 2 días  
**Ubicación**: `coati2lang/strings_methods.go`

**Problema**:
Solo métodos básicos implementados (length, substring, charAt)

**Impacto**:
- Funcionalidad limitada para string processing
- Necesidad de workarounds

#### FUNC-003: Sin Soporte para Comentarios Multilínea
**Prioridad**: 📈 Media  
**Complejidad**: Baja  
**Estimación**: 1 día  
**Ubicación**: `coati2lang/scanner.go:169-177`

**Problema**:
```go
// Solo soporta comentarios de línea
if s.match('/') {
    // A comment goes until the end of the line.
    for s.peek() != '\n' && !s.isAtEnd() {
        s.advance()
    }
}
```

**Impacto**:
- Limitaciones en documentación de código
- Inconveniente para debugging

### Testing (Pruebas)

#### TEST-001: Ausencia Total de Tests
**Prioridad**: 📈 Media  
**Complejidad**: Alta  
**Estimación**: 15 días  
**Ubicación**: Todo el proyecto

**Problema**:
No existen archivos `*_test.go` en el proyecto

**Impacto**:
- No hay garantías de correctitud
- Refactoring riesgoso
- Regresiones no detectadas

**Solución**:
- Suite completa de unit tests
- Integration tests
- Benchmark tests

#### TEST-002: Sin CI/CD Pipeline
**Prioridad**: 📈 Media  
**Complejidad**: Baja  
**Estimación**: 2 días  
**Ubicación**: Archivo `.github/workflows/`

**Problema**:
No hay automatización de tests/builds

**Impacto**:
- Testing manual propenso a errores
- No validación automática de PRs

---

## 💡 Issues de Baja Prioridad (Prioridad 4)

### UX/DX (Experiencia de Usuario/Desarrollador)

#### UX-001: REPL Básico Sin Características Modernas
**Prioridad**: 💡 Baja  
**Complejidad**: Media  
**Estimación**: 5 días  
**Ubicación**: `main.go:30-46`

**Problema**:
REPL muy básico sin historial, auto-completado, o sintaxis highlighting

**Impacto**:
- Experiencia de desarrollo pobre
- Productividad reducida
- No competitivo con otros REPLs

#### UX-002: Mensajes de Error Poco Informativos
**Prioridad**: 💡 Baja  
**Complejidad**: Baja  
**Estimación**: 3 días  
**Ubicación**: `coati2lang/r2loxerrors.go:5-15`

**Problema**:
```go
func Errors(line int, message string) {
    fmt.Fprintf(os.Stderr, "[line %d] Error: %s\n", line, message)
}
```

**Impacto**:
- Difficult debugging para usuarios
- No sugerencias de fix

### Documentation (Documentación)

#### DOC-001: Comentarios de Código Insuficientes
**Prioridad**: 💡 Baja  
**Complejidad**: Baja  
**Estimación**: 4 días  
**Ubicación**: Todo el proyecto

**Problema**:
Muchas funciones sin documentación GoDoc

**Impacto**:
- Dificulta onboarding de contributors
- Mantenimiento más difícil

#### DOC-002: Falta de Ejemplos de Código
**Prioridad**: 💡 Baja  
**Complejidad**: Baja  
**Estimación**: 2 días  
**Ubicación**: `README.md`, documentation

**Problema**:
Solo un ejemplo básico en `script.lox`

**Impacto**:
- Learning curve empinada
- Adopción lenta

---

## 🏗️ Issues de Refactoring (Mejoras Técnicas)

### Code Quality (Calidad de Código)

#### QUAL-001: Funciones Demasiado Largas
**Prioridad**: 🏗️ Refactoring  
**Complejidad**: Media  
**Estimación**: 6 días  
**Ubicación**: `coati2lang/parser.go:403-516`, `coati2lang/interpreter.go:172-241`

**Problema**:
```go
// parser.go Array() - 113 líneas
func (p *Parser) Array() []Expr {
    // Función muy larga y compleja
}
```

**Impacto**:
- Difícil entendimiento y mantenimiento
- Testing complicado
- Alta complejidad ciclomática

#### QUAL-002: Magic Numbers y Strings
**Prioridad**: 🏗️ Refactoring  
**Complejidad**: Baja  
**Estimación**: 2 días  
**Ubicación**: `coati2lang/parser.go:332`, `coati2lang/interpreter.go:104`

**Problema**:
```go
if len(arguments) >= 255 {  // Magic number
    Errors(p.peek().Line, "Can't have more than 255 arguments.")
}
```

**Impacto**:
- Código menos mantenible
- Límites hardcodeados

#### QUAL-003: Duplicación de Código
**Prioridad**: 🏗️ Refactoring  
**Complejidad**: Media  
**Estimación**: 4 días  
**Ubicación**: `coati2lang/parser.go:422-434`, `coati2lang/parser.go:478-490`

**Problema**:
Lógica similar repetida en múltiples lugares

**Impacto**:
- Mantenimiento duplicado
- Inconsistencias potenciales

---

## Matriz de Priorización

| Issue | Prioridad | Complejidad | Tiempo | Impacto | Riesgo |
|-------|-----------|-------------|--------|---------|--------|
| ERR-001 | 🔥 | Alta | 7d | Alto | Alto |
| ERR-002 | 🔥 | Media | 3d | Alto | Medio |
| MEM-001 | 🔥 | Media | 4d | Alto | Alto |
| MEM-002 | 🔥 | Alta | 5d | Alto | Medio |
| PERF-001 | ⚡ | Media | 4d | Medio | Bajo |
| PERF-002 | ⚡ | Alta | 6d | Medio | Medio |
| ARCH-001 | ⚡ | Alta | 8d | Alto | Medio |
| ARCH-002 | ⚡ | Media | 3d | Medio | Bajo |
| TEST-001 | 📈 | Alta | 15d | Alto | Bajo |
| FUNC-001 | 📈 | Media | 3d | Bajo | Bajo |

---

## Plan de Resolución por Fases

### Fase 1 (Semana 1-2): Issues Críticos
1. **ERR-002** (3 días) - Más simple, base para ERR-001
2. **ERR-001** (7 días) - Sistema de errores completo
3. **MEM-001** (4 días) - Fundamental para estabilidad

### Fase 2 (Semana 3-4): Performance Core
1. **MEM-002** (5 días) - Optimización de memoria
2. **PERF-001** (4 días) - Variable lookups
3. **ARCH-002** (3 días) - Estado del scanner

### Fase 3 (Semana 5-6): Arquitectura
1. **ARCH-001** (8 días) - Refactoring mayor
2. **PERF-002** (6 días) - AST optimizations

### Fase 4 (Semana 7-9): Testing y Funcionalidad
1. **TEST-001** (15 días) - Suite de pruebas
2. **FUNC-001, FUNC-002, FUNC-003** (6 días total)

---

## Métricas de Seguimiento

### Métricas de Calidad
- **Error Rate**: % de ejecuciones que terminan en error
- **Memory Usage**: Memoria máxima en benchmark estándar
- **Test Coverage**: % de código cubierto por tests
- **Cyclomatic Complexity**: Complejidad promedio por función

### Métricas de Performance
- **Variable Lookup Time**: Nanosegundos promedio
- **Parse Time**: Tiempo para parsear 1000 líneas
- **Execution Time**: Tiempo para ejecutar script estándar
- **Memory Allocations**: Número de allocations por operación

### Objetivos de Mejora
- ✅ Error rate < 1% en suite de pruebas
- ✅ Test coverage > 80%
- ✅ Variable lookup < 100ns promedio
- ✅ Memory usage estable en ejecuciones largas

---

## Conclusión

La resolución sistemática de estos issues transformará go-r2lox de un prototipo funcional a un intérprete robusto y eficiente. La priorización se basa en impacto directo en estabilidad, rendimiento y experiencia del usuario, asegurando que cada mejora construya sobre una base sólida.