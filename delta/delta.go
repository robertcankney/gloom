package delta

import (
	"encoding/binary"
	"errors"
	"github.com/cespare/xxhash"
)

//QuadTree is a Quadtree for storing data about byte ranges - premature optimization for sure
// type QuadTree struct {
// 	sum      uint64
// 	loc      uint32 //Always incremented by one compared to the actual index so that 0 values work for aggregated trees
// 	size     int32  //To ensure needed re-checksumming is performed
// 	children uint64
// 	first    *QuadTree
// 	second   *QuadTree
// 	third    *QuadTree
// 	fourth   *QuadTree
// }

// Delta contains a block size reference, aggregations and aggregated size, and the raw array of checksums
// First either contains direct hashes of blocks of sums of reference size, or aggregates second
// Regardless of number of aggregated sums, first will have a max of 10 elements, second a max of 100
type Delta struct {
	reference int32
	csum      uint64
	// size int32
	// first []uint64
	// second []uint64
	// aggregator aggregator
	sums []uint64
}

// Different encode block sizes
var ErrSizeMisalignment = errors.New("Block size differs between deltas")

//Calculates fields
func Build(sums []uint64, size int32) (delta Delta) {
	delta.reference = size
	delta.csum = encode(sums)
	delta.sums = sums
	return
}

//Local in this case would be the source - not necessarily a local copy
//Client -> cloud -> other client would have local be the cloud for the second compare
//overS and overE are start and end indices to delete - remote being truncated in this case
func (local *Delta) Compare(remote *Delta) (blocks []int, overS int, overE int, err error) {

	if local.reference != remote.reference {
		err = ErrSizeMisalignment
		return
	}

	if local.csum == remote.csum {
		return
	}

	if len(local.sums) >= len(remote.sums) {
		overE = -1
		overS = -1
	} else {
		l := len(local.sums)
		r := len(remote.sums)
		overS = l
		overE = r - 1
	}

	blocks = make([]int, 0, len(local.sums))

	for i := 0; i < len(local.sums); i++ {
		if local.sums[i] != remote.sums[i] {
			blocks = append(blocks, i)
		}
	}
	return

}

// func (local *QuadTree) Compare(remote *QuadTree) (blocks []uint32, err error) {

// 	if local.size != remote.size {
// 		err = ErrSizeMisalignment
// 		return
// 	}
// 	if local.sum == 0 {

// 		if local.children == remote.children {
// 			return
// 		} else {
// 			fi, _ := local.first.Compare(remote.first)
// 			se, _ := local.second.Compare(remote.second)
// 			th, _ := local.third.Compare(remote.third)
// 			fo, _ := local.fourt h.Compare(remote.fourth)
// 			// blocks = make([] uint64, len(fi) + len(se) + len(th) + len(fo))
// 			blocks = append(blocks, fi...)
// 			blocks = append(blocks, se...)
// 			blocks = append(blocks, th...)
// 			blocks = append(blocks, fo...)
// 		}

// 	} else {
// 		if local.sum != remote.sum {
// 			blocks = append(blocks, local.loc)
// 		}
// 	}
// 	return
// }

//Need to finish this prior to comparing - want to ensure I am maintaining a hierarchical model
//No tests other than integration for this, since it doesn't care too much about inputs
// func QtBuild(sums []uint64, start uint32) (tree QuadTree) {

// 	length := len(sums)
// 	if length/4.0 < 1 {

// 		remainder := length % 4
// 		switch remainder {
// 		case 0:
// 			first := QuadTree{sum: sums[0], loc: start}
// 			second := QuadTree{sum: sums[1], loc: start + 1}
// 			third := QuadTree{sum: sums[2], loc: start + 2}
// 			fourth := QuadTree{sum: sums[3], loc: start + 3}
// 			tree = QuadTree{children: childEncode([]QuadTree{first, second, third, fourth}), first: &first,
// 				second: &second, third: &third, fourth: &fourth}

// 		case 1:
// 			first := QuadTree{sum: sums[0], loc: start}
// 			second := QuadTree{sum: sums[1], loc: start + 1}
// 			third := QuadTree{sum: sums[2], loc: start + 2}
// 			tree = QuadTree{children: childEncode([]QuadTree{first, second, third}), first: &first,
// 				second: &second, third: &third}

// 		case 2:
// 			first := QuadTree{sum: sums[0], loc: start}
// 			second := QuadTree{sum: sums[1], loc: start + 1}
// 			tree = QuadTree{children: childEncode([]QuadTree{first, second}), first: &first,
// 				second: &second}

// 		case 3:
// 			first := QuadTree{sum: sums[0], loc: start}
// 			tree = QuadTree{children: childEncode([]QuadTree{first}), first: &first}

// 		}
// 		return

// 	} else {

// 		partition := int(math.Ceil(float64(length) / float64(4)))
// 		children := make([]QuadTree, 4)
// 		tmpLen := length

// 		for i := 1; tmpLen > partition; i++ {

// 			if tmpLen-partition > 0 {
// 				child := make([]uint64, partition)
// 				copy(sums[(length-tmpLen):(length-(tmpLen-partition))-1], child)
// 				children[i] = Build(child, uint32(length-tmpLen))
// 			} else {
// 				child := make([]uint64, tmpLen)
// 				copy(sums[(length-partition):(length-1)], child)
// 				children[i] = Build(child, uint32(length-partition))
// 			}
// 			tmpLen = tmpLen - partition
// 		}

// 		fi := children[0]
// 		se := children[1]
// 		th := children[2]
// 		fo := children[3]

// 		tree = QuadTree{children: childEncode(children), first: &fi, second: &se,
// 			third: &th, fourth: &fo}

// 		return
// 	}
// }

// func (delta *Delta) buildDep(sums []uint64) (err error) {
// 	length := len(sums)

// 	//If necessary, start by populating second
// 	if length > 10*delta.size {

// 		for i := 0; i*delta.size < length; i++ {
// 			if (i+1)*delta.size > length {
// 				second[i] = childEncode(sums[(i * delta.size):(length - 1)])
// 			} else {
// 				second[i] = childEncode(sums[(i * delta.size):((i + 1*delta.size) - 1)])
// 			}
// 		}
// 		//Aggregates second - breaks if we encounter a block of all zeroes, and strips out zeroes in semi-populated block from hashing
// 		var tmp []uint64
// 		var encode []uint64
// 		var zeroes int32

// 		for i := 0; i*delta.size < length; i++ {
// 			zeroes = 0
// 			encode = make([]uint64, 0, length)

// 			if (i+1)*delta.size > length {
// 				tmp = sums[(i * delta.size):(length - 1)]
// 				for i := range tmp {
// 					if tmp[i] == 0 {
// 						zeroes += 1
// 					} else {
// 						encode = append(encode, tmp[i])
// 					}
// 				}
// 				if len(tmp) == zeroes {
// 					break
// 				}
// 				first[i] = childEncode(tmp)
// 			} else {
// 				first[i] = childEncode(sums[(i * delta.size):((i + 1*delta.size) - 1)])
// 				for i := range tmp {
// 					if tmp[i] == 0 {
// 						zeroes += 1
// 					}
// 				}
// 				if len(tmp) == zeroes {
// 					break
// 				}
// 				first[i] = childEncode(tmp)
// 			}
// 		}

// 	} else {
// 		//Executes if we only need to aggregate to first

// 		for i := 1; i*delta.size < length; i++ {
// 			if (i+1)*delta.size > length {
// 				first[i] = childEncode(sums[(i * delta.size):(length - 1)])
// 			} else {
// 				first[i] = childEncode(sums[(i * delta.size):((i + 1*delta.size) - 1)])
// 			}
// 		}

// 		return
// 	}
// }

// func qtChildEncode(trees []QuadTree) (val uint64) {
// 	tmpBytes := make([]byte, binary.MaxVarintLen64)
// 	total := make([]byte, len(trees)*8)i -

// 	for i := range trees {
// 		binary.LittleEndian.PutUint64(tmpBytes, trees[i].sum)
// 		total = append(total, tmpBytes...)
// 	}
// 	val = xxhash.Sum64(total)
// 	return
// }

func encode(sums []uint64) (val uint64) {
	tmpBytes := make([]byte, binary.MaxVarintLen64)
	total := make([]byte, len(sums)*8)

	for i := range sums {
		binary.LittleEndian.PutUint64(tmpBytes, sums[i])
		total = append(total, tmpBytes...)
	}
	val = xxhash.Sum64(total)
	return
}
