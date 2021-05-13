package fdalang

import (
	"fmt"
)

func execScalarBinOperation(left, right Object, operator TokenID) (Object, error) {
	if left.Type() == TypeInt {
		left, _ := left.(*ObjInteger)
		right, _ := right.(*ObjInteger)
		result, err := integerBinOperation(left, right, operator)
		return result, err
	} else if left.Type() == TypeFloat {
		left, _ := left.(*ObjFloat)
		right, _ := right.(*ObjFloat)
		result, err := floatBinOperation(left, right, operator)
		return result, err
	} else if left.Type() == TypeBool {
		left, _ := left.(*ObjBoolean)
		right, _ := right.(*ObjBoolean)
		result, err := booleanBinOperation(left, right, operator)
		return result, err
	}
	if _, ok := left.(*ObjEnum); ok {
		if operator != TokenEq {
			return nil, fmt.Errorf("unsupported operator '%s' for type: '%s'", operator, left.Type())
		}
		left := left.(*ObjEnum).Value
		right := right.(*ObjEnum).Value
		return nativeBooleanToBoolean(left == right), nil
	}
	return nil, fmt.Errorf("unsupported operator '%s' for type: '%s'", operator, left.Type())
}

func integerBinOperation(left, right *ObjInteger, operator TokenID) (Object, error) {
	switch operator {
	case TokenPlus:
		return &ObjInteger{Value: left.Value + right.Value}, nil
	case TokenMinus:
		return &ObjInteger{Value: left.Value - right.Value}, nil
	case TokenSlash:
		return &ObjInteger{Value: left.Value / right.Value}, nil
	case TokenAsterisk:
		return &ObjInteger{Value: left.Value * right.Value}, nil
	case TokenLt:
		return nativeBooleanToBoolean(left.Value < right.Value), nil
	case TokenGt:
		return nativeBooleanToBoolean(left.Value > right.Value), nil
	case TokenEq:
		return nativeBooleanToBoolean(left.Value == right.Value), nil
	case TokenNotEq:
		return nativeBooleanToBoolean(left.Value != right.Value), nil
	default:
		return nil, fmt.Errorf("unsupported operator for types: %s %s %s", left.Type(), operator, right.Type())
	}
}

func nativeBooleanToBoolean(value bool) *ObjBoolean {
	if value == true {
		return ReservedObjTrue
	}
	return ReservedObjFalse
}

func floatBinOperation(left, right *ObjFloat, operator TokenID) (Object, error) {
	switch operator {
	case TokenPlus:
		return &ObjFloat{Value: left.Value + right.Value}, nil
	case TokenMinus:
		return &ObjFloat{Value: left.Value - right.Value}, nil
	case TokenSlash:
		return &ObjFloat{Value: left.Value / right.Value}, nil
	case TokenAsterisk:
		return &ObjFloat{Value: left.Value * right.Value}, nil
	case TokenLt:
		return nativeBooleanToBoolean(left.Value < right.Value), nil
	case TokenGt:
		return nativeBooleanToBoolean(left.Value > right.Value), nil
	case TokenEq:
		return nativeBooleanToBoolean(left.Value == right.Value), nil
	case TokenNotEq:
		return nativeBooleanToBoolean(left.Value != right.Value), nil
	default:
		return nil, fmt.Errorf("unsupported operator for types: %s %s %s", left.Type(), operator, right.Type())
	}
}

func booleanBinOperation(left, right *ObjBoolean, operator TokenID) (Object, error) {
	switch operator {
	case TokenEq:
		return nativeBooleanToBoolean(left.Value == right.Value), nil
	case TokenNotEq:
		return nativeBooleanToBoolean(left.Value != right.Value), nil
	case TokenAnd:
		return nativeBooleanToBoolean(left.Value && right.Value), nil
	case TokenOr:
		return nativeBooleanToBoolean(left.Value || right.Value), nil
	default:
		return nil, fmt.Errorf("unsupported operator for types: %s %s %s", left.Type(), operator, right.Type())
	}
}
