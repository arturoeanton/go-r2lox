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
	switch first.(type) {
	case []interface{}:
		return cloneArray(first.([]interface{}))
	case map[interface{}]interface{}:
		return cloneMap(first.(map[interface{}]interface{}))
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
		switch v.(type) {
		case []interface{}:
			clone[i] = cloneArray(v.([]interface{}))
		case map[interface{}]interface{}:
			clone[i] = cloneMap(v.(map[interface{}]interface{}))
		default:
			clone[i] = v
		}
	}
	return clone
}

func cloneMap(m map[interface{}]interface{}) map[interface{}]interface{} {
	clone := make(map[interface{}]interface{})
	for k, v := range m {
		switch v.(type) {
		case []interface{}:
			clone[k] = cloneArray(v.([]interface{}))
		case map[interface{}]interface{}:
			clone[k] = cloneMap(v.(map[interface{}]interface{}))
		default:
			clone[k] = v
		}
	}
	return clone
}
