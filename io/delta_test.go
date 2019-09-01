package io

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	"testing"
)

// func TestQtSize(t *testing.T) {
// 	q1 := QuadTree{size: 1}
// 	q2 := QuadTree{size: 2}

// 	_, err := q1.Compare(&q2)

// 	if err == nil {
// 		t.Error("No error for size misalignment: " + err.Error())
// 	}
// }

func TestBuildAndCompare(t *testing.T) {

	localSums := []uint64{3, 4, 345, 7547, 2345, 12, 0, 9543, 546543754, 923423, 348834693463463, 5,
		9, 43, 85238, 94934538543, 943, 99, 2, 33333, 4534535435435, 98534, 882565262465, 98453284325,
		77423, 54327, 88534, 8235834858326, 924832868, 7753720, 5466443643643, 84377547354, 88458358345, 7540, 7, 9,
		3, 4, 345, 7547, 35345435, 12, 0, 9543, 45432442342, 923423, 348834693463463, 45436436,
		9, 43, 85238, 8877, 943743, 99, 2, 33333, 4534543535353, 98534, 882565262465, 98453284325,
		77423, 3223232, 322323, 8235834858326, 43433, 7753720, 2000000, 84377547354, 88458358345, 7540, 7, 9,
		1, 2, 3, 4, 5}

	remoteSums := []uint64{3, 4, 345, 7547, 2345, 12, 0, 9543, 2, 923423, 348834693463463, 5,
		9, 43, 85238, 94934538543, 943, 99, 2, 33333, 4534535435435, 98534, 882565262465, 98453284325,
		77423, 5, 88534, 8235834858326, 924832868, 7753720, 5466443643643, 84377547354, 88458358345, 7540, 7, 9,
		3, 4, 345, 7547, 35345435, 12, 0, 3, 45432442342, 923423, 348834693463463, 45436436,
		9, 43, 85238, 8877, 943743, 99, 2, 33333, 4534543535353, 98534, 882565262465, 98453284325,
		77423, 3223232, 322323, 8, 43433, 7753720, 2000000, 84377547354, 88458358345, 0, 7, 9,
		1, 2, 3, 66545654, 5, 1}

	local := Build(localSums, 0)
	remote := Build(remoteSums, 0)

	b, overs, overe, err := local.Compare(&remote)

	if err != nil {
		t.Error("Unexpected error on Compare")
	}

	expected := []int{8, 25, 43, 63, 69, 75}

	if !cmp.Equal(b, expected) {
		fmt.Println(b)
		fmt.Println(expected)
		fmt.Println(cmp.Diff(b, expected))
		t.Error("Incorrect block delta")
	}

	if overe != 77 || overs != 77 {
		fmt.Println(overe)
		fmt.Println(overs)
		t.Error("Incorrect truncation")
	}
}
