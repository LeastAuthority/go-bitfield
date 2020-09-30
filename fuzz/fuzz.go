//+build gofuzz

package fuzz

import (
	"fmt"

	"github.com/filecoin-project/go-bitfield"
	rleplus "github.com/filecoin-project/go-bitfield/rle"
)

func FuzzNewFromBytes(data []byte) int {
	_, err := bitfield.NewFromBytes(data)
	if err != nil {
		return 0
	}

	return 1
}

func FuzzTwoNewAndUnion(data []byte) int {
	if len(data) > 4*1024*1024 {
		return -1
	}

	var (
		l         = len(data)
		dataLeft  = data[:l/2]
		dataRight = data[l/2:]
	)

	left, err := bitfield.NewFromBytes(dataLeft)
	if err != nil {
		return 0
	}

	right, err := bitfield.NewFromBytes(dataRight)
	if err != nil {
		return 0
	}

	union, err := bitfield.MergeBitFields(left, right)
	if err != nil {
		panic("merge should not return an error")
	}

	var (
		helper int64 = -1
		max          = uint64(helper) // use underflow for MAXINT
	)

	uall, err := union.AllMap(max)
	if err != nil {
		panic("allmap should not return an error (union)")
	}

	lall, err := left.AllMap(max)
	if err != nil {
		panic("allmap should not return an error (left)")
	}

	rall, err := right.AllMap(max)
	if err != nil {
		panic("allmap should not return an error (right)")
	}

	for v := range uall {
		_, lok := lall[v]
		_, rok := rall[v]

		if !lok && !rok {
			panic(fmt.Sprintf("element %v of union in neither left nor right", v))
		}
	}

	for v := range lall {
		if _, ok := uall[v]; !ok {
			panic(fmt.Sprintf("element %v of left is not in union", v))
		}
	}

	for v := range rall {
		if _, ok := uall[v]; !ok {
			panic(fmt.Sprintf("element %v of right is not in union", v))
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
