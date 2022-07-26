package literal

import (
	"fmt"
	"github.com/gojinja/literal_eval/ast"
	"github.com/gojinja/literal_eval/parser"
	"github.com/gojinja/literal_eval/py"
	"math/big"
	"reflect"
	"strings"
)

// Ellipsis type is used for indication that python type was Ellipsis.
type Ellipsis struct{}

// Eval evaluates python code into Golang type. It supports same types that Python's `ast.literal_eval` does.
// Those types are: strings, bytes, numbers, tuples, lists, dicts, sets, booleans, None and Ellipsis
func Eval(pythonCode string) (interface{}, error) {
	Ast, err := parser.ParseString(pythonCode, py.EvalMode)
	if err != nil {
		return nil, err
	}

	if Ast == nil {
		return nil, fmt.Errorf("got nil")
	}

	if Ast.Type().Name != "Expression" {
		return nil, fmt.Errorf("expected Expression, got: %s", Ast.Type().Name)
	}

	astValue := reflect.Indirect(reflect.ValueOf(Ast))
	astType := astValue.Type()
	for i := 0; i < astType.NumField(); i++ {
		fieldType := astType.Field(i)
		fieldValue := astValue.Field(i)
		if strings.ToLower(fieldType.Name) == "body" {
			if expr, ok := fieldValue.Interface().(ast.Expr); ok {
				return parseValue(expr)
			}
			return nil, fmt.Errorf("expected '*ast.Expr'")
		}
	}

	return nil, fmt.Errorf("couldn't find body")
}

func parseValue(astExpr ast.Expr) (interface{}, error) {
	switch v := interface{}(astExpr).(type) {
	case *ast.Dict:
		return parseDict(v)
	case *ast.List:
		return parseList(v)
	case *ast.Tuple:
		return parseTuple(v)
	case *ast.Set:
		return parseSet(v)
	case *ast.Str:
		return parseString(v)
	case *ast.Ellipsis:
		return parseEllipsis(v)
	case *ast.NameConstant:
		return parseNameConstant(v)
	case *ast.Bytes:
		return parseBytes(v)
	case *ast.Num:
		return parseNum(v)
	default:
		return nil, fmt.Errorf("unsupported type type")
	}
}

func parseNum(num *ast.Num) (interface{}, error) {
	switch v := num.N.(type) {
	case py.Int:
		return int64(v), nil
	case *py.BigInt:
		return big.Int(*v), nil
	case py.Float:
		return float64(v), nil
	case py.Complex:
		return complex128(v), nil
	default:
		return nil, fmt.Errorf("got unexpected number obj")
	}
}

func parseBytes(v *ast.Bytes) ([]byte, error) {
	return v.S, nil
}

func parseNameConstant(nameConstant *ast.NameConstant) (interface{}, error) {
	switch v := nameConstant.Value.(type) {
	case py.NoneType:
		return nil, nil
	case py.Bool:
		return bool(v), nil
	default:
		return nil, fmt.Errorf("unsupported name constant")
	}
}

func parseEllipsis(_ *ast.Ellipsis) (Ellipsis, error) {
	return Ellipsis{}, nil
}

func parseString(s *ast.Str) (string, error) {
	return string(s.S), nil
}

func parseSet(s *ast.Set) (Set, error) {
	elts, err := parseElts(s.Elts)
	if err != nil {
		return Set{}, err
	}
	return newSet(elts), nil
}

func parseElts(elts []ast.Expr) ([]interface{}, error) {
	res := make([]interface{}, 0, len(elts))
	var err error
	for i, el := range elts {
		res[i], err = parseValue(el)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func parseTuple(v *ast.Tuple) ([]interface{}, error) {
	return parseElts(v.Elts)
}

func parseList(v *ast.List) ([]interface{}, error) {
	return parseElts(v.Elts)
}

func parseDict(dict *ast.Dict) (res Dict, err error) {
	keys, err := parseElts(dict.Keys)
	if err != nil {
		return
	}
	values, err := parseElts(dict.Values)
	if err != nil {
		return Dict{}, err
	}
	return newDict(keys, values), nil
}
