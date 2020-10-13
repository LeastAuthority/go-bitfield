//+build gofuzz

package fuzz

import (
	//"encoding/hex"
	"os"
	"testing"

	"net/http"
	_ "net/http/pprof"

	rleplus "github.com/filecoin-project/go-bitfield/rle"

	"github.com/davecgh/go-spew/spew"
	"github.com/leastauthority/fleece/fuzzing"
	"github.com/stretchr/testify/require"
)

func init() {
	go http.ListenAndServe(":12345", nil)
}

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
		//data = []byte(" @\x140000000\x14-0\xbf\xef\x940000")
		data = []byte("tee0q0n00%00000")
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

	t.Run("toplevel", func(t *testing.T) {
		FuzzTwoNewAndUnion(data)
	})

	t.Run("rle", func(t *testing.T) {
		FuzzRLETwoNewAndUnion(data)
	})
}

func TestFleece(t *testing.T) {
	t.Log(os.Getwd())
	env := fuzzing.NewEnv("../fleece")
	ci := fuzzing.MustNewCrasherIterator(env, FuzzRLETwoNewAndUnion)
	_, panics, _ := ci.TestFailingLimit(t, 100000000)
	t.Log(panics)

}

func TestBitIterator(t *testing.T) {
	runit := &rleplus.RunSliceIterator{
		Runs: []rleplus.Run{
			{Len: 1},
			{Len: 550, Val: true},
			{Len: 1},
			{Len: 6, Val: true},
			{Len: 1},
			{Len: 6, Val: true},
			{Len: 1},
		},
	}

	spew.Dump(runit)

	bitit, err := rleplus.BitsFromRuns(runit)
	require.NoError(t, err, "bitfromruns")

	for bitit.HasNext() {
		t.Log(bitit.Next())
	}
}
