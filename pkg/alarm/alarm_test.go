// +build !integration

package alarm

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/thingspect/api/go/common"
)

func TestGenerate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		inpPoint *common.DataPoint
		inpRule  *common.Rule
		inpDev   *common.Device
		inpTempl string
		res      string
		err      string
	}{
		{&common.DataPoint{}, nil, nil, `test`, "test", ""},
		{&common.DataPoint{ValOneof: &common.DataPoint_IntVal{IntVal: 40}},
			&common.Rule{Name: "test rule"}, nil, `point value is an ` +
				`integer: {{.pointVal}}, rule name is: {{.rule.Name}}`,
			"point value is an integer: 40, rule name is: test rule", ""},
		{&common.DataPoint{ValOneof: &common.DataPoint_Fl64Val{Fl64Val: 37.7}},
			nil, &common.Device{Status: common.Status_ACTIVE}, `point value ` +
				`is a float: {{.pointVal}}, device status is: ` +
				`{{.device.Status}}`, "point value is a float: 37.7, device " +
				"status is: ACTIVE", ""},
		{&common.DataPoint{ValOneof: &common.DataPoint_StrVal{StrVal: "batt"}},
			nil, nil, `point value is a string: {{.pointVal}}`,
			"point value is a string: batt", ""},
		{&common.DataPoint{ValOneof: &common.DataPoint_BoolVal{BoolVal: true}},
			nil, nil, `point value is a bool: {{.pointVal}}`,
			"point value is a bool: true", ""},
		{&common.DataPoint{ValOneof: &common.DataPoint_BytesVal{
			BytesVal: []byte{0x00, 0x01}}}, nil, nil, `point value is a byte ` +
			`slice: {{.pointVal}}`, "point value is a byte slice: [0 1]", ""},
		{&common.DataPoint{}, nil, nil, `{{if`, "", "unclosed action"},
		{&common.DataPoint{}, nil, nil, `{{template "aaa"}}`, "",
			"no such template"},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Can generate %+v", lTest), func(t *testing.T) {
			t.Parallel()

			res, err := Generate(lTest.inpPoint, lTest.inpRule, lTest.inpDev,
				lTest.inpTempl)
			t.Logf("res, err: %v, %#v", res, err)
			require.Equal(t, lTest.res, res)
			if lTest.err == "" {
				require.NoError(t, err)
			} else {
				require.Contains(t, err.Error(), lTest.err)
			}
		})
	}
}
