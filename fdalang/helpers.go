package fdalang

import (
	"encoding/json"
	"fmt"
	"log"
)

func (e *Environment) Print() {
	fmt.Println("Env content:")
	for k, v := range e.store {
		fmt.Printf("%s: %s\n", k, v.Inspect())
	}
}

func (e *Environment) GetVarsAsJson() ([]byte, error) {
	varMap := make(map[string]string)
	for k, v := range e.store {
		varMap[k] = v.Inspect()
	}
	return json.Marshal(varMap)
}

func (e *Environment) ToStrings() []string {
	result := make([]string, 0)
	for k, v := range e.store {
		result = append(result, fmt.Sprintf("%s: %s\n", k, v.Inspect()))
	}
	return result
}

func NewEmptyStruct(def *ObjStructDefinition) *ObjStruct {
	return &ObjStruct{
		Emptier:    Emptier{Empty: true},
		Definition: def,
		Fields:     make(map[string]Object),
	}
}

func (e *Environment) Keys() []string {
	keys := make([]string, len(e.store))

	i := 0
	for k := range e.store {
		keys[i] = k
		i++
	}
	return keys
}

func (e *Environment) LoadVarsInStruct(definition *ObjStructDefinition, s map[string]interface{}) *ObjStruct {
	fields := make(map[string]Object)
	for k, v := range s {
		fields[k] = getLangObject(v)
	}
	return &ObjStruct{
		Definition: definition,
		Fields:     fields,
	}
}

func getLangType(t interface{}) string {
	switch t.(type) {
	case float64:
		return TypeFloat
	case int:
		return TypeInt
	case int32:
		return TypeInt
	case bool:
		return TypeBool
	default:
		log.Fatalf("Unsupported type for struct creation: '%T'", t)
	}
	return ""
}

func getLangObject(t interface{}) Object {
	switch tt := t.(type) {
	case float64:
		return &ObjFloat{Value: tt}
	case int:
		return &ObjInteger{Value: int64(tt)}
	case int32:
		return &ObjInteger{Value: int64(tt)}
	case uint32:
		return &ObjInteger{Value: int64(tt)}
	case bool:
		return &ObjBoolean{Value: tt}
	default:
		log.Fatalf("Unsupported type for struct creation: '%T'", t)
	}
	return nil
}
