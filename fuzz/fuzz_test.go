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

func TestFuzz(t *testing.T) {
	var (
		data = []byte(" @\x140000000\x14-0\xbf\xef\x940000")
		args = make([]interface{}, 1, 9)
	)

	for line := 0; line < len(data); line += 16 {
		args = args[:1]
		args[0] = line

		lineData := data[line:]

		format := "%04d:"

		for i := 0; i < 8; i++ {
			if len(lineData) < 2 {
				break
			}

			format += " %x"
			args = append(args, lineData[:2])
			lineData = lineData[2:]
		}

		t.Logf(format, args...)
	}

	FuzzTwoNewAndUnion(data)
}
