// Package rule provides functions to evaluate boolean expression rules.
package rule

import (
	"errors"
	"time"

	"github.com/antonmedv/expr"
	"github.com/thingspect/api/go/common"
)

var ErrNotBool = errors.New("not a boolean expression")

// Eval evaluates a boolean expression using the Expr language.
// https://github.com/antonmedv/expr/blob/master/docs/Language-Definition.md
func Eval(point *common.DataPoint, ruleExpr string) (bool, error) {
	env := map[string]interface{}{
		"point":   point,
		"pointTS": point.Ts.Seconds,
		"currTS":  time.Now().Unix(),
	}

	// Populate point value for convenience. []byte is not supported.
	switch v := point.ValOneof.(type) {
	case *common.DataPoint_IntVal:
		env["pointVal"] = v.IntVal
	case *common.DataPoint_Fl64Val:
		env["pointVal"] = v.Fl64Val
	case *common.DataPoint_StrVal:
		env["pointVal"] = v.StrVal
	case *common.DataPoint_BoolVal:
		env["pointVal"] = v.BoolVal
	}

	res, err := expr.Eval(ruleExpr, env)
	if err != nil {
		return false, err
	}

	b, ok := res.(bool)
	if !ok {
		return false, ErrNotBool
	}

	return b, nil
}
