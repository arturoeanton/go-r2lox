package coati2lang

import (
	"fmt"
)

func init() {
	GlobalFx["clone"] = Clone{}
}

type Clone struct {
}

func (c Clone) Call(interpreter *Interpreter, arguments []interface{}, this interface{}) interface{} {
	var first interface{} = arguments[0]
	switch cast_element := first.(type) {
	case []interface{}:
		return cloneArray(cast_element)
	case map[interface{}]interface{}:
		return cloneMap(cast_element)
	default:
		return fmt.Errorf("clone: type %T not supported", first)
	}
}

func (c Clone) Arity() int {
	return 1
}

func cloneArray(array []interface{}) []interface{} {
	clone := make([]interface{}, len(array))
	for i, v := range array {
		switch cast_element := v.(type) {
		case []interface{}:
			clone[i] = cloneArray(cast_element)
		case map[interface{}]interface{}:
			clone[i] = cloneMap(cast_element)
		default:
			clone[i] = v
		}
	}
	return clone
}

func cloneMap(m map[interface{}]interface{}) map[interface{}]interface{} {
	clone := make(map[interface{}]interface{})
	for k, v := range m {
		switch cast_element := v.(type) {
		case []interface{}:
			clone[k] = cloneArray(cast_element)
		case map[interface{}]interface{}:
			clone[k] = cloneMap(cast_element)
		default:
			clone[k] = v
		}
	}
	return clone
}
