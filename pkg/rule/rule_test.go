//go:build !integration

package rule

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/thingspect/proto/go/common"
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
			"invalid operation: int + string",
		},
		{&common.DataPoint{}, `"aaa"`, false, ErrNotBool.Error()},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Can evaluate %+v", test), func(t *testing.T) {
			t.Parallel()

			test.inpPoint.Ts = timestamppb.New(time.Now().Add(-time.Second))

			res, err := Eval(test.inpPoint, test.inpRuleExpr)
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
