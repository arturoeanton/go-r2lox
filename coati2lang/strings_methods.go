package coati2lang

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

var (
	STRING_FX_MAP = map[string]func(string, ...any) interface{}{
		"len":        len1,
		"template":   template1,
		"number":     number1,
		"lower":      lower1,
		"upper":      upper1,
		"trim":       trim1,
		"trimleft":   trimleft1,
		"trimright":  trimright1,
		"trimprefix": trimprefix1,
		"trimsuffix": trimsuffix1,
		"split":      split1,
	}
)

func split1(s string, args ...interface{}) interface{} {
	interpreter := args[0].(*Interpreter)

	arge := args[1].([]Expr)
	argv := make([]interface{}, len(arge))
	for i, arg := range arge {
		argv[i] = arg.AcceptExpr(interpreter)
	}
	p := argv[0].(string)
	return strings.Split(s, p)
}

func lower1(s string, args ...interface{}) interface{} {
	return strings.ToLower(s)
}

func upper1(s string, args ...interface{}) interface{} {
	return strings.ToUpper(s)
}

func trim1(s string, args ...interface{}) interface{} {
	return strings.TrimSpace(s)
}

func trimleft1(s string, args ...interface{}) interface{} {
	return strings.TrimLeft(s, args[0].(string))
}

func trimright1(s string, args ...interface{}) interface{} {
	return strings.TrimRight(s, args[0].(string))
}

func trimprefix1(s string, args ...interface{}) interface{} {
	return strings.TrimPrefix(s, args[0].(string))
}

func trimsuffix1(s string, args ...interface{}) interface{} {
	return strings.TrimSuffix(s, args[0].(string))
}

func len1(s string, args ...interface{}) interface{} {
	return len(s)
}

func number1(s string, args ...interface{}) interface{} {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}
	return f
}

func template1(s string, args ...interface{}) interface{} {
	if len(args) == 0 {
		return s
	}
	interpreter := args[0].(*Interpreter)

	result, err := replaceNestedKeys(s, interpreter)
	if err != nil {
		return err
	}
	return result
}

func replaceNestedKeys(str string, interpreter *Interpreter) (string, error) {
	var result strings.Builder
	var key strings.Builder
	isKey := false

	for i, c := range str {
		if c == '$' && i+1 < len(str) && rune(str[i+1]) == '{' {
			isKey = true
			i++
		} else if isKey && c == '}' {
			isKey = false
			parts := strings.Split(key.String()[1:], "(")
			var keys []string
			if len(parts) > 1 {
				keys = strings.Split(parts[0], ".")
				keys[len(keys)-1] = keys[len(keys)-1] + "(" + parts[1]
			} else {
				keys = strings.Split(parts[0], ".")
			}

			this, val, args, err := getNestedValue(keys, interpreter.enviroment)
			if err != nil {
				return "", err
			}
			if fx, ok := val.(LoxCallable); ok {
				for i, arg := range args {
					if arg == "this" {
						args[i] = this
					}
					if args[i] == "true" {
						args[i] = true
					}
					if args[i] == "false" {
						args[i] = false
					}

					if args[i] == "nil" {
						args[i] = nil
					}

					if str, ok := args[i].(string); ok && len(str) > 0 && str[0] == '"' && str[len(str)-1] == '"' {
						args[i] = str[1 : len(str)-1]
					}

					if f, err := strconv.ParseFloat(args[i].(string), 64); err == nil {
						args[i] = f
					}

					if str, ok := args[i].(string); ok {
						split_path := strings.Split(str, ".")
						if len(split_path) > 1 {
							_, val, _, _ := getNestedValue(split_path, interpreter.enviroment)
							args[i] = val
						} else {
							args[i], _ = interpreter.enviroment.Get(str)
						}
					}

				}

				val = fx.Call(interpreter, args, this)
			}
			result.WriteString(fmt.Sprintf("%v", val))
			key.Reset()
		} else if isKey {
			key.WriteRune(c)
		} else {
			result.WriteRune(c)
		}
	}

	return result.String(), nil
}

func getNestedValue(keys []string, env *Enviroment) (interface{}, interface{}, []interface{}, error) {
	val, _ := env.Get(keys[0])
	var this interface{}
	var argsret []interface{}
	for i, key := range keys {
		this = val
		if i == 0 {
			continue
		}
		// TODO: handler arguments
		parser_function := strings.Split(key, "(")
		if len(parser_function) > 1 {
			key = parser_function[0]
			args := parser_function[1][:len(parser_function[1])-1]
			agrs_split := strings.Split(args, ",")
			argsret = make([]interface{}, len(agrs_split))
			for i, arg := range agrs_split {
				argsret[i] = strings.TrimSpace(arg)
			}

		}

		v := reflect.ValueOf(val)
		switch v.Kind() {
		case reflect.Map:
			keyv := reflect.ValueOf(key)
			vi := v.MapIndex(keyv)
			val = vi.Interface()
		case reflect.Slice, reflect.Array:
			index, err := strconv.Atoi(key)
			if err != nil {
				return this, nil, nil, errors.New("invalid array index")
			}
			val = v.Index(index).Interface()
		default:
			return this, nil, nil, errors.New("type not supported for key: " + key)
		}
	}

	return this, val, argsret, nil
}
