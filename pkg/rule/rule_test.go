//go:build !integration

package rule

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/thingspect/api/go/common"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestEval(t *testing.T) {
	t.Parallel()

	tests := []struct {
		inpPoint    *common.DataPoint
		inpRuleExpr string
		res         bool
		err         string
	}{
		{&common.DataPoint{}, `true`, true, ""},
		{&common.DataPoint{}, `10 > 15`, false, ""},
		{&common.DataPoint{}, `point.Token == ""`, true, ""},
		{&common.DataPoint{}, `pointTS < currTS`, true, ""},
		{
			&common.DataPoint{ValOneof: &common.DataPoint_IntVal{IntVal: 40}},
			`point.GetIntVal() == 40`, true, "",
		},
		{
			&common.DataPoint{ValOneof: &common.DataPoint_IntVal{IntVal: 40}},
			`pointVal > 32`, true, "",
		},
		{
			&common.DataPoint{ValOneof: &common.DataPoint_Fl64Val{
				Fl64Val: 37.7,
			}}, `pointVal < 32`, false, "",
		},
		{
			&common.DataPoint{ValOneof: &common.DataPoint_StrVal{
				StrVal: "line",
			}}, `pointVal == battery`, false, "",
		},
		{
			&common.DataPoint{ValOneof: &common.DataPoint_BoolVal{
				BoolVal: true,
			}}, `pointVal`, true, "",
		},
		{
			&common.DataPoint{}, `1 + "aaa"`, false,
			"invalid operation: + (mismatched types int and string)",
		},
		{&common.DataPoint{}, `"aaa"`, false, ErrNotBool.Error()},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Can evaluate %+v", lTest), func(t *testing.T) {
			t.Parallel()

			lTest.inpPoint.Ts = timestamppb.New(time.Now().Add(-time.Second))

			res, err := Eval(lTest.inpPoint, lTest.inpRuleExpr)
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
