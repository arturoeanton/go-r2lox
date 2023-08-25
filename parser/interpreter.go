package parser

import (
	"fmt"
	"strings"
	"time"

	"github.com/arturoeanton/go-r2lox/lexer"
)

type LoxCallable interface {
	Call(interpreter *Interpreter, arguments []interface{}) interface{}
	Arity() int
}

type Interpreter struct {
	Stmts      []Stmt
	enviroment *Enviroment
}

type Clock struct {
}

func (c Clock) Call(interpreter *Interpreter, arguments []interface{}) interface{} {
	return time.Now().Unix()
}

func (c Clock) Arity() int {
	return 0
}

func NewInterpreter(stmts []Stmt) *Interpreter {
	global := NewEnviroment(nil)

	global.Define("clock", Clock{})

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
	var value interface{}
	if stmt.Value != nil {
		value = i.evaluate(stmt.Value)
		return value
	}
	panic("return")
}

func (i *Interpreter) VisitFunctionStmt(stmt Function) interface{} {
	function := Function{
		Name:       stmt.Name,
		Parameters: stmt.Parameters,
		Body:       stmt.Body,
	}
	i.enviroment.Define(stmt.Name.Lexeme, function)
	return nil
}

func (i *Interpreter) VisitCallExpr(expr Call) interface{} {
	callee := i.evaluate(expr.Callee)
	var arguments []interface{}
	for _, argument := range expr.Arguments {
		arguments = append(arguments, i.evaluate(argument))
	}
	function, ok := callee.(LoxCallable)
	if !ok {
		fmt.Println("Can only call functions and classes.")
		return nil
	}
	if len(arguments) != function.Arity() {
		fmt.Println("Expected", function.Arity(), "arguments but got", len(arguments), ".")
		return nil
	}

	return function.Call(i, arguments)

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
	if expr.Operator.Type == lexer.OR {
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
	i.enviroment = &enviroment
	var result interface{}
	for _, stmt := range stmts {
		if stmt == nil {
			continue
		}
		result = i.execute(stmt)
	}
	i.enviroment = previous
	return result
}

func (i *Interpreter) VisitBlockStmt(stmt Block) interface{} {
	return i.executeBlock(stmt.Statements, *NewEnviroment(i.enviroment))
}

func (i *Interpreter) VisitExpressionStmt(stmt Expression) interface{} {
	return i.evaluate(stmt.Expression)
}

func (i *Interpreter) VisitPrint(stmt Print) interface{} {
	value := i.evaluate(stmt.Expression)
	fmt.Println(stringify(value))
	return value
}

func stringify(object interface{}) string {
	if object == nil {
		return "nil"
	}
	switch value := object.(type) {
	case float64:
		text := fmt.Sprintf("%v", value)

		if ok := strings.HasSuffix(text, ".0"); ok {
			text = text[:len(text)-2]
		}
		return text
	default:
		return fmt.Sprintf("%v", object)
	}
}

func (i *Interpreter) VisitBinaryExpr(expr Binary) interface{} {

	left := i.evaluate(expr.Left)
	right := i.evaluate(expr.Right)

	switch expr.Operator.Type {
	case lexer.MINUS:
		return left.(float64) - right.(float64)
	case lexer.PLUS:
		{
			if _, ok := left.(float64); ok {
				return left.(float64) + right.(float64)
			}
			if _, ok := left.(string); ok {
				return left.(string) + right.(string)
			}
			return nil
		}
	case lexer.SLASH:
		return left.(float64) / right.(float64)
	case lexer.STAR:
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
	case lexer.GREATER:
		return left.(float64) > right.(float64)
	case lexer.GREATER_EQUAL:
		return left.(float64) >= right.(float64)
	case lexer.LESS:
		return left.(float64) < right.(float64)
	case lexer.LESS_EQUAL:
		return left.(float64) <= right.(float64)
	case lexer.BANG_EQUAL:
		return !i.isEqual(left, right)
	case lexer.EQUAL_EQUAL:
		return i.isEqual(left, right)
	default:
		return nil
	}
}

func (i *Interpreter) VisitGroupingExpr(expr Grouping) interface{} {
	return i.evaluate(expr.Expression)
}

func (i *Interpreter) VisitLiteralExpr(expr Literal) interface{} {
	return expr.Value
}

func (i *Interpreter) VisitUnaryExpr(expr Unary) interface{} {
	right := i.evaluate(expr.Right)

	switch expr.Operator.Type {
	case lexer.MINUS:
		return -(right.(float64))
	case lexer.BANG:
		return !(i.isTruthy(right))
	default:
		return nil
	}

}

func (i *Interpreter) VisitVar(stmt Var) interface{} {
	var value interface{}
	if stmt.Initializer != nil {
		value = i.evaluate(stmt.Initializer)
	}
	i.enviroment.Define(stmt.Name.Lexeme, value)
	return nil
}

func (i *Interpreter) VisitVariableExpr(expr Var) interface{} {
	return i.enviroment.Get(expr.Name.Lexeme)
}

func (i *Interpreter) VisitAssignExpr(expr Assign) interface{} {
	value := i.evaluate(expr.Value)
	i.enviroment.Assign(expr.Name.Lexeme, value)
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
