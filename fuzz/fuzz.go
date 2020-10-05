//+build gofuzz

package fuzz

import (
	"fmt"

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

func FuzzRLETwoNewAndIntersect(data []byte) int {
	if len(data) > 4*1024*1024 {
		return -1
	}

	var (
		dataLen   = len(data)
		dataLeft  = data[:dataLen/2]
		dataRight = data[dataLen/2:]
	)

	left, err := rleplus.FromBuf(dataLeft)
	if err != nil {
		return 0
	}

	right, err := rleplus.FromBuf(dataRight)
	if err != nil {
		return 0
	}

	leftRI, err := left.RunIterator()
	if err != nil {
		return 0
	}

	rightRI, err := right.RunIterator()
	if err != nil {
		return 0
	}

	isectRI, err := rleplus.And(leftRI, rightRI)
	if err != nil {
		return 0
	}

	uiter, err := rleplus.BitsFromRuns(isectRI)
	if err != nil {
		panic("bititerator should not return an error (isect)")
	}

	liter, err := rleplus.BitsFromRuns(leftRI)
	if err != nil {
		panic("bititerator should not return an error (left)")
	}

	riter, err := rleplus.BitsFromRuns(rightRI)
	if err != nil {
		panic("bititerator should not return an error (right)")
	}

	var (
		u, l, r      uint64
		lDone, rDone bool
		found        bool
	)

	if liter.HasNext() {
		l, err = liter.Next()
		if err != nil {
			panic("liter next should not return an error")
		}
	} else {
		lDone = true
	}

	if riter.HasNext() {
		r, err = riter.Next()
		if err != nil {
			panic("riter next should not return an error")
		}
	} else {
		rDone = true
	}

	for uiter.HasNext() {
		u, err = uiter.Next()
		if err != nil {
			panic("uiter next should not return an error")
		}

		found = false

		if !lDone {
			if l < u {
				panic(fmt.Sprintf("found element %d that is in left but not isect", l))
			} else if l == u {
				found = true
				if liter.HasNext() {
					l, err = liter.Next()
					if err != nil {
						panic("liter next should not return an error")
					}
				} else {
					lDone = true
				}
			}
		}

		if !rDone {
			if r < u {
				panic(fmt.Sprintf("found element %d that is in right but not isect", r))
			} else if r == u {
				found = true
				if riter.HasNext() {
					r, err = riter.Next()
					if err != nil {
						panic("riter next should not return an error")
					}
				} else {
					rDone = true
				}
			}
		}

		if !found {
			panic(fmt.Sprintf("element %d in isect found in neither left nor right", u))
		}
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
