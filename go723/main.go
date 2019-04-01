package main

// #cgo CFLAGS: -g -Icg723
// #cgo LDFLAGS: -L cg723 -lg723
// #include <stdlib.h>
// #include "lbccodec.h"
// #include "ophtools.h"
// #include "g723_const.h"
import "C"
import (
	"bytes"
	"fmt"
	"io/ioutil"
	"unsafe"

	"github.com/davecgh/go-spew/spew"
)

const (
	FRAME_SIZE          = 24
	DECODED_BUFFER_SIZE = 480
	WAV_HEADER          = 100
)

func init() {
	C.g723_i(0, 0)
}

func main() {
	audiobytes, err := ioutil.ReadFile("g723_1_example.wav")
	if err != nil {
		panic(err)
	}
	buf := audiobytes[WAV_HEADER:]

	dop := &bytes.Buffer{}
	// dop.Write(audiobytes[:WAV_HEADER])

	for i := 0; i < len(buf); i += FRAME_SIZE {
		if len(buf) < i+FRAME_SIZE {
			// spew.Dump(buf[:])
			dop.Write(decodeFrame(buf[:]))
			continue
		}
		fmt.Printf("Iteration %d - %d\n", i, i/FRAME_SIZE)
		// spew.Dump(buf[i : i+FRAME_SIZE])
		dop.Write(decodeFrame(buf[i : i+FRAME_SIZE]))
	}

	fffile, err := ioutil.ReadFile("right_decoded.wav")
	if err != nil {
		panic(err)
	}

	writeToFileDummy(dop.Bytes())
	fmt.Printf("Original File: (len)%d\n", len(audiobytes))
	fmt.Printf("Decoded File: (len)%d\n", len(dop.Bytes()))
	fmt.Printf("FFMpeg File: (len)%d\n", len(fffile))
}

func decodeFrame(frame []byte) []byte {
	buffer := make([]byte, DECODED_BUFFER_SIZE)
	frame_ptr := unsafe.Pointer(&frame[0])
	buffer_ptr := unsafe.Pointer(&buffer[0])

	C.g723_d((*C.uchar)(frame_ptr), (*C.short)(buffer_ptr)) // is void

	//fmt.Println("frame_ptr")
	//spew.Dump(C.GoBytes(frame_ptr, C.int(len(frame))))
	//
	//fmt.Println("buffer_ptr")
	//spew.Dump(C.GoBytes(buffer_ptr, C.int(len(buffer))))

	return C.GoBytes(buffer_ptr, C.int(len(buffer)))
}

// not even sure if this is encode frame or encode all audio >.<
func encodeFrame(frame []byte) []byte {
	buffer := make([]byte, DECODED_BUFFER_SIZE)
	frame_ptr := unsafe.Pointer(&frame[0])
	buffer_ptr := unsafe.Pointer(&buffer[0])

	spew.Dump(C.g723_e((*C.short)(buffer_ptr), (*C.uchar)(frame_ptr)))

	fmt.Println("frame_ptr")
	spew.Dump(C.GoBytes(frame_ptr, C.int(len(frame))))

	fmt.Println("buffer_ptr")
	spew.Dump(C.GoBytes(buffer_ptr, C.int(len(buffer))))

	return C.GoBytes(buffer_ptr, C.int(len(buffer)))
}

func writeToFileDummy(buf []byte) {
	if err := ioutil.WriteFile("output.wav", buf, 0700); err != nil {
		panic(err)
	}
}
