package io

import (
	"io"
)

//Asset wraps a ReadWriter and encapsulates a Blocks struct
type Asset struct {
	rw     io2.ReadWriter
	blocks Blocks
}

//Blocks is a wrapper around the comm channels used to asynchronously communicate new checksums + receive requests for blocks
//as well as send byte ranges
type Blocks struct {
	Sums  <-chan SendRange
	Block chan<- int
	Bytes <-chan []byte
}

//SendRange wraps an xxhash checksum + block location
type SendRange struct {
	Sum   uint64
	Block int
}

//ReceiveRange wraps a byte slice and location
type ReceiveRange struct {
	Data  []byte
	Block int
}

//GenSums creates a set of sums for an already tracked file - for comparison rather than streaming something new (sync)
//Largely this package contains code should be agnostic of use, but this meets a specific use case - client daemon was down
//and is assessing if a modified file's state
func (asset *Asset) GenSums(size int) {

}

//Read returns an Asset struct that is asynchronously populated - can be used to create a Delta that also async populates
func Read(string locator) (asset *Asset, err error) {
	return
}

//Write asynchronously writes to a file - receives a channel that it streams data from
//Largely should be called in the context of a Delta, as that manages state of an Asset - applies to all of these funcs
func (asset *Asset) Write(chan<- Range) (err error) {

}
