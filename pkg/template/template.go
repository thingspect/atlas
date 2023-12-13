// Package template provides functions to generate HTML-safe output from
// templates.
package template

import (
	"html/template"
	"strings"

	"github.com/thingspect/proto/go/api"
	"github.com/thingspect/proto/go/common"
)

// Generate generates HTML-safe output from templates using the Go template
// engine: https://golang.org/pkg/html/template/
func Generate(
	point *common.DataPoint, rule *api.Rule, dev *api.Device, templ string,
) (string, error) {
	env := map[string]interface{}{
		"point":  point,
		"rule":   rule,
		"device": dev,
	}

	// Populate point value for convenience. If point doesn't validate, pointVal
	// remains unset.
	switch v := point.GetValOneof().(type) {
	case *common.DataPoint_IntVal:
		env["pointVal"] = v.IntVal
	case *common.DataPoint_Fl64Val:
		env["pointVal"] = v.Fl64Val
	case *common.DataPoint_StrVal:
		env["pointVal"] = v.StrVal
	case *common.DataPoint_BoolVal:
		env["pointVal"] = v.BoolVal
	case *common.DataPoint_BytesVal:
		env["pointVal"] = v.BytesVal
	}

	t, err := template.New("template").Parse(templ)
	if err != nil {
		return "", err
	}

	res := &strings.Builder{}
	err = t.Execute(res, env)
	if err != nil {
		return "", err
	}

	return res.String(), nil
}
