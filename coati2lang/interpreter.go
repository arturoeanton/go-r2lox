package coati2lang

import (
	"errors"
	"fmt"
	"log"
	"math"
	"time"
)

var (
	GlobalFx = make(map[string]LoxCallable)
)

func init() {
	GlobalFx["clock"] = Clock{}
}

type LoxCallable interface {
	Call(interpreter *Interpreter, arguments []interface{}, this interface{}) interface{}
	Arity() int
}

type Interpreter struct {
	Stmts      []Stmt
	enviroment *Enviroment
}

type Clock struct {
}

func (c Clock) Call(interpreter *Interpreter, arguments []interface{}, this interface{}) interface{} {
	return time.Now().Unix()
}

func (c Clock) Arity() int {
	return 0
}

func NewInterpreter(stmts []Stmt) *Interpreter {
	global := NewEnviroment(nil)

	for k, v := range GlobalFx {
		global.Define(k, v)
	}

	return &Interpreter{
		Stmts:      stmts,
		enviroment: global,
	}

}

func (i *Interpreter) Interpret() interface{} {
	result := i.executeBlock(i.Stmts, *i.enviroment)
	return result
}

func (i *Interpreter) execute(stmt Stmt) interface{} {
	return stmt.AcceptStmt(i)
}

func (i *Interpreter) VisitReturnStmt(stmt Return) interface{} {
	var result interface{}
	if stmt.Value != nil {
		result = i.evaluate(stmt.Value)

	}
	panic(Return{Keyword: stmt.Keyword, Value: stmt.Value, Result: result})
}

func (i *Interpreter) VisitFunctionStmt(stmt Function) interface{} {
	function := Function{
		Name:       stmt.Name,
		Parameters: stmt.Parameters,
		Body:       stmt.Body,
		Closure:    i.enviroment,
	}
	i.enviroment.Define(stmt.Name.Lexeme, function)
	return nil
}

func (i *Interpreter) VisitCallExpr(expr Call) interface{} {
	callee := i.evaluate(expr.Callee)

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

	//if _, ok := expr.(Var); ok {
	if expr.This.Name.Lexeme != "" {
		this := i.full_evaluate(expr.This)
		return callable.Call(i, arguments, this)
	}
	//}

	return callable.Call(i, arguments, nil)

}
func (i *Interpreter) VisitWhileStmt(stmt While) interface{} {
	for i.isTruthy(i.evaluate(stmt.Condition)) {
		i.execute(stmt.Body)
	}
	return nil
}

func (i *Interpreter) VisitIfStmt(stmt If) interface{} {
	if i.isTruthy(i.evaluate(stmt.Condition)) {
		return i.execute(stmt.ThenBranch)
	} else if stmt.ElseBranch != nil {
		return i.execute(stmt.ElseBranch)
	}
	return nil
}

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

func (i *Interpreter) VisitBlockStmt(stmt Block) interface{} {
	return i.executeBlock(stmt.Statements, *NewEnviroment(i.enviroment))
}

func (i *Interpreter) VisitExpressionStmt(stmt Expression) interface{} {
	return i.evaluate(stmt.Expression)
}

func (i *Interpreter) VisitBinaryExpr(expr Binary) interface{} {

	left := i.evaluate(expr.Left)
	right := i.evaluate(expr.Right)

	switch expr.Operator.Type {
	case MINUS:
		return left.(float64) - right.(float64)
	case PERCENT:
		return (left.(float64) * right.(float64) / 100.0)
	case PLUS:
		{
			if _, ok := left.(float64); ok {
				return left.(float64) + right.(float64)
			}
			if _, ok := left.(string); ok {
				return left.(string) + right.(string)
			}
			return nil
		}
	case SLASH:
		return left.(float64) / right.(float64)
	case STAR_STAR:
		return math.Pow(left.(float64), right.(float64))
	case STAR:
		// validate rigth is string
		{
			_, right_is_string := right.(string)
			_, left_is_string := left.(string)
			_, right_is_number := right.(float64)
			_, left_is_number := left.(float64)

			if right_is_string && left_is_number {
				var result string
				for i := 0; i < int(left.(float64)); i++ {
					result += right.(string)
				}
				return result
			}

			if right_is_number && left_is_string {
				var result string
				for i := 0; i < int(right.(float64)); i++ {
					result += left.(string)
				}
				return result
			}

			if right_is_number && left_is_number {
				return left.(float64) * right.(float64)
			}

			return nil
		}
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
	default:
		return nil
	}
}

func (i *Interpreter) VisitGroupingABSExpr(expr GroupingABS) interface{} {
	value := i.evaluate(expr.Expression)
	if value.(float64) < 0 {
		return -value.(float64)
	}
	return value
}

func (i *Interpreter) VisitGroupingExpr(expr Grouping) interface{} {
	return i.evaluate(expr.Expression)
}

func (i *Interpreter) VisitLiteralExpr(expr Literal) interface{} {
	return expr.Value
}

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
		return !(i.isTruthy(value))
	default:
		return nil
	}

}

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

		if len(expr.Selectors) == 1 {
			if array, ok := value.([]interface{}); ok {
				if len(array) == 1 {
					value = array[0]
				}
			}
		}
	}

	return value
}

func (i *Interpreter) setByPath(target interface{}, path []interface{}, value interface{}) (interface{}, error) {
	// Si no hay más elementos en la path, simplemente asigna el valor
	if len(path) == 0 {
		return nil, errors.New("path is too short")
	}

	switch t := target.(type) {
	case []interface{}:
		// Trata target como un slice
		index := int(path[0].(float64))

		// Si el índice está fuera de rango, extiende el slice
		for len(t) <= index {
			t = append(t, nil)
		}

		// Si esta es la última parte de la path, asigna el valor
		if len(path) == 1 {
			t[index] = value
			return t, nil
		}
		return i.setByPath(t[index], path[1:], value)

	case map[interface{}]interface{}:
		// Trata target como un mapa
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

func (i *Interpreter) VisitAssignExpr(expr Assign) interface{} {
	value := i.full_evaluate(expr.Value)
	old, ok := i.enviroment.Get(expr.Name.Lexeme)
	if !ok {
		log.Fatalln("Undefined variable '" + expr.Name.Lexeme + "'.")
	}

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
			//TODO: verificar que funcione en todos los casos de uso
			i.enviroment.Assign(expr.Name.Lexeme, new)
		}
		return value
	} else {
		i.enviroment.Assign(expr.Name.Lexeme, value)
	}
	return value
}
func (i *Interpreter) full_evaluate(expr Expr) interface{} {
	value := i.evaluate(expr)
	expr_v, ok := value.(Expr)
	if ok {
		value := i.full_evaluate(expr_v)
		return value
	}
	expr_a, ok := value.([]Expr)
	if ok {
		var values []interface{} = make([]interface{}, len(expr_a))
		for index, value := range expr_a {
			values[index] = i.full_evaluate(value)
		}
		return values
	}

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

func (i *Interpreter) evaluate(expr Expr) interface{} {
	return expr.AcceptExpr(i)
}

func (i *Interpreter) isTruthy(object interface{}) bool {
	if object == nil {
		return false
	}
	if object == false {
		return false
	}
	return true
}

func (i *Interpreter) isEqual(a interface{}, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return false
	}
	return a == b
}
