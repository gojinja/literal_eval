package literal

import "github.com/gojinja/literal_eval/ast"

func isHashable(expr ast.Expr) bool {
	switch expr.(type) {
	case *ast.Dict, *ast.Set:
		return false
	default:
		return true
	}
}
