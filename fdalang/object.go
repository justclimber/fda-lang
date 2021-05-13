package fdalang

import (
	"bytes"
	"fmt"
	"strings"
)

type ObjectType string

const (
	TypeInt         = "int"
	TypeFloat       = "float"
	TypeBool        = "bool"
	TypeReturnValue = "return_value"
	TypeFunction    = "function_obj"
	TypeBuiltinFn   = "builtin_fn_obj"
	TypeVoid        = "void"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Emptier struct {
	Empty bool
}

func (e *Emptier) IsEmpty() bool { return e.Empty }

type ObjStructDefinition struct {
	Name   string
	Fields map[string]string
}

type ObjEnumDefinition struct {
	Name     string
	Elements []string
}

func CreateVarDefinitionsFromVarType(varTypes map[string]*VarAndType) map[string]string {
	varDefinitions := make(map[string]string)
	for k, v := range varTypes {
		varDefinitions[k] = v.VarType
	}
	return varDefinitions
}

type IIdentifier interface{}
type IStatements interface{}

type ObjInteger struct {
	Emptier
	Value int64
}

func (i *ObjInteger) Type() ObjectType { return TypeInt }
func (i *ObjInteger) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

type ObjFloat struct {
	Emptier
	Value float64
}

func (f *ObjFloat) Type() ObjectType { return TypeFloat }
func (f *ObjFloat) Inspect() string  { return fmt.Sprintf("%.2f", f.Value) }

type ObjBoolean struct {
	Value bool
}

func (b *ObjBoolean) Type() ObjectType { return TypeBool }
func (b *ObjBoolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

type ObjEnum struct {
	Definition *ObjEnumDefinition
	Value      int8
}

func (e *ObjEnum) Type() ObjectType { return ObjectType(e.Definition.Name) }
func (e *ObjEnum) Inspect() string {
	return fmt.Sprintf("%s", e.Definition.Elements[e.Value])
}

type ObjArray struct {
	Emptier
	ElementsType string
	Elements     []Object
}

func (a *ObjArray) Type() ObjectType {
	varType := fmt.Sprintf("[]%s", a.ElementsType)
	return ObjectType(varType)
}
func (a *ObjArray) Inspect() string {
	var elements []string
	for _, e := range a.Elements {
		elements = append(elements, e.Inspect())
	}

	return fmt.Sprintf("[]%s{%s}", a.ElementsType, strings.Join(elements, ", "))
}

type ObjReturnValue struct {
	Value Object
}

func (rv *ObjReturnValue) Type() ObjectType { return TypeReturnValue }
func (rv *ObjReturnValue) Inspect() string  { return rv.Value.Inspect() }

type ObjFunction struct {
	Arguments  []*VarAndType
	Statements *StatementsBlock
	ReturnType string
	Env        *Environment
}

func (f *ObjFunction) Type() ObjectType { return TypeFunction }
func (f *ObjFunction) Inspect() string {
	return "function"
}

type ObjStruct struct {
	Emptier
	Definition *ObjStructDefinition
	Fields     map[string]Object
}

func (s *ObjStruct) Type() ObjectType { return ObjectType(s.Definition.Name) }
func (s *ObjStruct) Inspect() string {
	var out bytes.Buffer

	var elements []string
	for k, v := range s.Fields {
		elements = append(elements, fmt.Sprintf("%s: %s", k, v.Inspect()))
	}

	out.WriteString(s.Definition.Name)
	out.WriteString("{")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("}")

	return out.String()
}

type BuiltinFunction func(env *Environment, args []Object) (Object, error)

type ArgTypes []string

type ObjBuiltin struct {
	Name       string
	ArgTypes   ArgTypes
	Fn         BuiltinFunction
	ReturnType string
}

func (b *ObjBuiltin) Type() ObjectType { return TypeBuiltinFn }
func (b *ObjBuiltin) Inspect() string  { return "builtin function" }

type ObjVoid struct{}

func (v *ObjVoid) Type() ObjectType { return TypeVoid }
func (v *ObjVoid) Inspect() string  { return "void" }
