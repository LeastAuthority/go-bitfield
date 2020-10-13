//+build gofuzz

package fuzz

import (
	"github.com/filecoin-project/go-bitfield"
	rleplus "github.com/filecoin-project/go-bitfield/rle"
	//"github.com/google/gofuzz"
)

func FuzzNewFromBytes(data []byte) int {
	_, err := bitfield.NewFromBytes(data)
	if err != nil {
		return 0
	}

	return 1
}

func FuzzRLE_FromBuf(data []byte) int {
	_, err := rleplus.FromBuf(data)
	if err != nil {
		return 0
	}

	return 1
}

/*
func FuzzRLE_BitIterator(data []byte) int {
	var (
		runs []rleplus.Run
	)
	fuzz.NewFromGoFuzz(data).NilChance(0).NumElements(1, 20).Fuzz(&run)
	runit := rleplus.RunSliceIterator{Runs: runs}
	bitit, err := rleplus.BitsFromRuns(runit)

	require.NoError(t, err, "bitsfromruns")
	t.Log(bitit)

}
*/
