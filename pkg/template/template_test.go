//go:build !integration

package template

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/thingspect/proto/go/api"
	"github.com/thingspect/proto/go/common"
)

func TestGenerate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		inpPoint *common.DataPoint
		inpRule  *api.Rule
		inpDev   *api.Device
		inpTempl string
		res      string
		err      string
	}{
		{
			&common.DataPoint{}, nil, nil, `test`, "test", "",
		},
		{
			&common.DataPoint{ValOneof: &common.DataPoint_IntVal{IntVal: 40}},
			&api.Rule{Name: "test rule"}, nil, `point value is an ` +
				`integer: {{.pointVal}}, rule name is: {{.rule.Name}}`,
			"point value is an integer: 40, rule name is: test rule", "",
		},
		{
			&common.DataPoint{ValOneof: &common.DataPoint_Fl64Val{
				Fl64Val: 37.7,
			}}, nil, &api.Device{Status: api.Status_ACTIVE},
			`point value is a float: {{.pointVal}}, device status is: ` +
				`{{.device.Status}}`,
			"point value is a float: 37.7, device status is: ACTIVE", "",
		},
		{
			&common.DataPoint{ValOneof: &common.DataPoint_StrVal{
				StrVal: "line",
			}}, nil, nil, `point value is a string: {{.pointVal}}`,
			"point value is a string: line", "",
		},
		{
			&common.DataPoint{ValOneof: &common.DataPoint_BoolVal{
				BoolVal: true,
			}}, nil, nil, `point value is a bool: {{.pointVal}}`,
			"point value is a bool: true", "",
		},
		{
			&common.DataPoint{ValOneof: &common.DataPoint_BytesVal{
				BytesVal: []byte{0x00, 0x01},
			}}, nil, nil, `point value is a byte slice: {{.pointVal}}`,
			"point value is a byte slice: [0 1]", "",
		},
		{
			&common.DataPoint{}, nil, nil, `{{if`, "", "unclosed action",
		},
		{
			&common.DataPoint{}, nil, nil, `{{template "aaa"}}`, "",
			"no such template",
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Can generate %+v", test), func(t *testing.T) {
			t.Parallel()

			res, err := Generate(test.inpPoint, test.inpRule, test.inpDev,
				test.inpTempl)
			t.Logf("res, err: %v, %#v", res, err)
			require.Equal(t, test.res, res)
			if test.err == "" {
				require.NoError(t, err)
			} else {
				require.Contains(t, err.Error(), test.err)
			}
		})
	}
}
