// +build !integration

package parse

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrFormat(t *testing.T) {
	t.Parallel()

	err := ErrFormat("function", "bad identifier", []byte{0x00, 0x01, 0x02})
	t.Logf("err: %v", err)

	require.EqualError(t, err, "function format bad identifier: 000102")
}
