package fdalang

import (
	"errors"
	"fmt"
)

var (
	ReservedObjTrue  = &ObjBoolean{Value: true}
	ReservedObjFalse = &ObjBoolean{Value: false}
)

func registerStructDefinition(node *AstStructDefinition, env *Environment) error {
	s := &ObjStructDefinition{
		Name:   node.Name,
		Fields: CreateVarDefinitionsFromVarType(node.Fields),
	}
	if err := env.RegisterStructDefinition(s); err != nil {
		return err
	}
	return nil
}

func registerEnumDefinition(node *AstEnumDefinition, env *Environment) error {
	ed := &ObjEnumDefinition{
		Name:     node.Name,
		Elements: node.Elements,
	}
	if err := env.RegisterEnumDefinition(ed); err != nil {
		return err
	}
	return nil
}

func structTypeAndVarsChecks(n *AstAssignment, definition *ObjStructDefinition, result Object) error {
	fieldType, ok := definition.Fields[n.Left.Value]
	if !ok {
		return runtimeError(
			n, "Struct '%s' doesn't have the field '%s' in the definition", definition.Name, n.Left.Value)
	}
	if fieldType != string(result.Type()) {
		return runtimeError(
			n,
			"Field '%s' defined as '%s' but '%s' given",
			n.Left.Value,
			fieldType,
			result.Type())
	}
	return nil
}

func arrayElementsTypeCheck(node *AstArray, t string, es []Object) error {
	for i, el := range es {
		if string(el.Type()) != t {
			return runtimeError(node, "Array element #%d should be type '%s' but '%s' given", i+1, t, el.Type())
		}
	}
	return nil
}

func functionReturnTypeCheck(node *AstFunctionCall, result Object, functionReturnType string) error {
	if result.Type() != ObjectType(functionReturnType) {
		return runtimeError(node,
			"Return type mismatch: function declared as '%s' but in fact return '%s'",
			functionReturnType, result.Type())
	}
	return nil
}

func functionCallArgumentsCheck(node *AstFunctionCall, declaredArgs []*AstVarAndType, actualArgValues []Object) error {
	if len(declaredArgs) != len(actualArgValues) {
		return runtimeError(node, "Function call arguments count mismatch: declared %d, but called %d",
			len(declaredArgs), len(actualArgValues))
	}

	if len(actualArgValues) > 0 {
		for i, arg := range declaredArgs {
			if actualArgValues[i].Type() != ObjectType(arg.VarType) {
				return runtimeError(arg, "argument #%d type mismatch: expected '%s' by func declaration but called '%s'",
					i+1, arg.VarType, actualArgValues[i].Type())
			}
		}
	}

	return nil
}

func transferArgsToNewEnv(fn *ObjFunction, args []Object) *Environment {
	env := NewEnclosedEnvironment(fn.Env)

	for i, arg := range fn.Arguments {
		env.Set(arg.Var.Value, args[i])
	}

	return env
}

func runtimeError(node AstNode, format string, args ...interface{}) error {
	msg := fmt.Sprintf(format, args...)
	t := node.GetToken()
	return errors.New(fmt.Sprintf("%s\nline:%d, pos %d", msg, t.Line, t.Col))
}
