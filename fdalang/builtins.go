package fdalang

import (
	"fmt"
	"math"
)

func AbsInt64(n int64) int64 {
	y := n >> 63
	return (n ^ y) - y
}

const (
	BuiltinPrint    = "print"
	BuiltinEmpty    = "empty"
	BuiltinLength   = "length"
	BuiltinAbsInt   = "absInt"
	BuiltinAbsFloat = "absFloat"
)

func (e *ExecAstVisitor) setupBasicBuiltinFunctions() {
	e.builtins[BuiltinPrint] = &ObjBuiltin{
		Name:       BuiltinPrint,
		ArgTypes:   ArgTypes{"any"},
		ReturnType: TypeVoid,
		Fn: func(env *Environment, args []Object) (Object, error) {
			fmt.Println(args[0].Inspect())
			return &ObjVoid{}, nil
		},
	}
	e.builtins[BuiltinEmpty] = &ObjBuiltin{
		Name:       BuiltinEmpty,
		ArgTypes:   ArgTypes{"any"},
		ReturnType: TypeBool,
		Fn: func(env *Environment, args []Object) (Object, error) {
			switch arg := args[0].(type) {
			case *ObjStruct:
				return nativeBooleanToBoolean(arg.Empty), nil
			case *ObjInteger:
				return nativeBooleanToBoolean(arg.Empty), nil
			case *ObjFloat:
				return nativeBooleanToBoolean(arg.Empty), nil
			case *ObjArray:
				return nativeBooleanToBoolean(arg.Empty), nil
			default:
				return nil, BuiltinFuncError("Type '%T' doesn't support emptiness", arg)
			}
		},
	}
	e.builtins[BuiltinLength] = &ObjBuiltin{
		Name:       BuiltinLength,
		ArgTypes:   ArgTypes{"array"},
		ReturnType: TypeInt,
		Fn: func(env *Environment, args []Object) (Object, error) {
			array := args[0].(*ObjArray)
			length := len(array.Elements)
			return &ObjInteger{Value: int64(length)}, nil
		},
	}
	e.builtins[BuiltinAbsInt] = &ObjBuiltin{
		Name:       BuiltinAbsInt,
		ArgTypes:   ArgTypes{TypeInt},
		ReturnType: TypeInt,
		Fn: func(env *Environment, args []Object) (Object, error) {
			num := args[0].(*ObjInteger).Value
			return &ObjInteger{Value: AbsInt64(num)}, nil
		},
	}
	e.builtins[BuiltinAbsFloat] = &ObjBuiltin{
		Name:       BuiltinAbsFloat,
		ArgTypes:   ArgTypes{TypeFloat},
		ReturnType: TypeFloat,
		Fn: func(env *Environment, args []Object) (Object, error) {
			float := args[0].(*ObjFloat).Value
			return &ObjFloat{Value: math.Abs(float)}, nil
		},
	}
}

func (e *ExecAstVisitor) AddBuiltinFunctions(builtins map[string]*ObjBuiltin) {
	for k, v := range builtins {
		e.builtins[k] = v
	}
}

func (e *ExecAstVisitor) checkArgs(builtin *ObjBuiltin, args []Object) error {
	if builtin.ArgTypes == nil {
		return nil
	}
	if len(builtin.ArgTypes) != len(args) {
		return fmt.Errorf(
			"wrong number of arguments for '%s'. need %d, got %d",
			builtin.Name,
			len(builtin.ArgTypes),
			len(args),
		)
	}
	for i, argType := range builtin.ArgTypes {
		if argType == "any" {
			continue
		} else if argType == "array" {
			if _, ok := args[i].(*ObjArray); !ok {
				return fmt.Errorf(
					"wrong type of argument #%d for '%s'. need %s, got %T",
					i+1,
					builtin.Name,
					argType,
					args[i],
				)
			}
		} else if argType != string(args[i].Type()) {
			return fmt.Errorf(
				"wrong type of argument #%d for '%s'. need %s, got %s",
				i+1,
				builtin.Name,
				argType,
				args[i].Type(),
			)
		}
	}
	return nil
}

// todo line and col
func BuiltinFuncError(format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}
