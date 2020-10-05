//+build gofuzz

package fuzz

import (
	"fmt"

	rleplus "github.com/filecoin-project/go-bitfield/rle"
	//"github.com/davecgh/go-spew/spew"
	//"github.com/google/gofuzz"
)

func FuzzRLETwoNewAndUnion(data []byte) int {
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

	//fmt.Println("right:")
	//spew.Dump(rightRI)

	unionRI, err := rleplus.Or(leftRI, rightRI)
	if err != nil {
		return 0
	}

	//fmt.Println("union:")
	//spew.Dump(unionRI)

	leftRI, err = left.RunIterator()
	if err != nil {
		return 0
	}

	rightRI, err = right.RunIterator()
	if err != nil {
		return 0
	}

	uiter, err := rleplus.BitsFromRuns(unionRI)
	if err != nil {
		panic("bititerator should not return an error (union)")
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
		u, l, r      uint64 // offset of set bit
		lDone, rDone bool   // has the iterator returned false for HasNext
		found        bool
	)

	if liter.HasNext() {
		l, err = liter.Next()
		if err != nil {
			panic("liter next should not return an error")
		}
	} else {
		fmt.Println("l done from the start")
		lDone = true
	}

	if riter.HasNext() {
		r, err = riter.Next()
		if err != nil {
			panic("riter next should not return an error")
		}
	} else {
		fmt.Println("r done from the start")
		rDone = true
	}

	for uiter.HasNext() {
		u, err = uiter.Next()
		if err != nil {
			panic("uiter next should not return an error")
		}

		//fmt.Println("u", u)
		//fmt.Println("l", l, lDone)
		//fmt.Println("r", r, rDone)

		found = false

		if !lDone {
			if l < u {
				panic(fmt.Sprintf("found element %d that is in left but not union", l))
			} else if l == u {
				found = true
				if liter.HasNext() {
					l, err = liter.Next()
					if err != nil {
						panic("liter next should not return an error")
					}
				} else {
					fmt.Println("l done after", l)
					lDone = true
				}
			}
		}

		if !rDone {
			if r < u {
				panic(fmt.Sprintf("found element %d that is in right but not union", r))
			} else if r == u {
				found = true
				if riter.HasNext() {
					r, err = riter.Next()
					if err != nil {
						panic("riter next should not return an error")
					}
				} else {
					fmt.Println("r done after", r)
					rDone = true
				}
			}
		}

		if !found {
			panic(fmt.Sprintf("element %d in union found in neither left nor right", u))
		}
	}

	return 1
}