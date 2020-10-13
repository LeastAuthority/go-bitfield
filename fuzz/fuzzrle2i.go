//+build gofuzz

package fuzz

import (
	"fmt"

	rleplus "github.com/filecoin-project/go-bitfield/rle"
	//"github.com/google/gofuzz"
)

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

	leftRI, err = left.RunIterator()
	if err != nil {
		return 0
	}

	rightRI, err = right.RunIterator()
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

		found = true

		if !lDone {
			for l < u {
				if liter.HasNext() {
					l, err = liter.Next()
					if err != nil {
						panic("liter next should not return an error")
					}
				} else {
					lDone = true
					found = false
				}
			}

			if l > u {
				found = false
			}
		}

		if !rDone {
			for r < u {
				if riter.HasNext() {
					r, err = riter.Next()
					if err != nil {
						panic("riter next should not return an error")
					}
				} else {
					rDone = true
					found = false
				}
			}

			if r > u {
				found = false
			}
		}

		if !found {
			panic(fmt.Sprintf("element %d in isect not found in either left nor right", u))
		}
	}

	return 1
}
