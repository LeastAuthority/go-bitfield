//+build gofuzz

package fuzz

import (
	"testing"

	rleplus "github.com/filecoin-project/go-bitfield/rle"
	//fleece "github.com/leastauthority/fleece/fuzzing"
	"github.com/stretchr/testify/require"
)

func TestFuzzTwoNewAndUnion(t *testing.T) {
	left, err := rleplus.FromBuf([]byte{})
	require.NoError(t, err, "left new from bytes")

	_, err = left.RunIterator()
	require.NoError(t, err, "left validate")

	right, err := rleplus.FromBuf([]byte{0})
	require.NoError(t, err, "right new from bytes")

	_, err = right.RunIterator()
	require.EqualError(t, err, "decoding RLE: not minimally encoded: invalid encoding for RLE+ version 0", "right validate")
}
