package coati2lang

import (
	"fmt"
)

func init() {
	GlobalFx["println"] = Println{}
	GlobalFx["print"] = Print{}
	GlobalFx["fprint"] = Fprint{}
	GlobalFx["sprint"] = Sprint{}
	GlobalFx["len"] = Len{}
}

type Println struct {
}

func (c Println) Call(interpreter *Interpreter, arguments []interface{}) interface{} {
	fmt.Println(arguments...)
	return nil
}

func (c Println) Arity() int {
	return -1
}

type Print struct {
}

func (c Print) Call(interpreter *Interpreter, arguments []interface{}) interface{} {
	n, _ := fmt.Print(arguments...)
	return n
}

func (c Print) Arity() int {
	return -1
}

type Fprint struct {
}

func (c Fprint) Call(interpreter *Interpreter, arguments []interface{}) interface{} {
	n, _ := fmt.Println(arguments...)
	return n
}

func (c Fprint) Arity() int {
	return -1
}

type Sprint struct {
}

func (c Sprint) Call(interpreter *Interpreter, arguments []interface{}) interface{} {
	return fmt.Sprint(arguments...)
}

func (c Sprint) Arity() int {
	return -1
}

type Len struct {
}

func (c Len) Call(interpreter *Interpreter, arguments []interface{}) interface{} {
	return len(arguments)
}

func (c Len) Arity() int {
	return 1
}
