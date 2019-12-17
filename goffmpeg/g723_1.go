package goffmpeg

// #cgo CFLAGS: -I./include
// #cgo LDFLAGS: -L lib -lavcodec -lavdevice -lavfilter -lavformat -lavutil -lffmpeg -lswresample -lswscale
// #include <stdio.h>
// #include <stdlib.h>
// #include <string.h>
// #include "libavcodec/avcodec.h"
// #include "include/g723_1_decode.h"
import "C"
import (
	"bytes"
	"log"
	"unsafe"
)

// Decoder is an interface borrowed from the `cart` project
type Decoder interface {
	Decode([]byte) ([]byte, error)
	Destroy()
	ConsumesPayloadType(int) bool
}

var _ Decoder = &G7231Decoder{}

// G7231Decoder is a struct for decoding g723.1 packets
type G7231Decoder struct {
	pkt          *C.struct_AVPacket
	codec        *C.struct_AVCodec
	parser       *C.struct_AVCodecParserContext
	c            *C.struct_AVCodecContext
	decodedFrame *C.struct_AVFrame
}

// NewG7231Decoder will return a new g723.1 decoder
func NewG7231Decoder() *G7231Decoder {
	pkt := C.av_packet_alloc()
	codec := C.avcodec_find_decoder(C.AV_CODEC_ID_G723_1)
	if codec == nil {
		log.Fatal("Codec not found")
	}
	parser := C.av_parser_init(C.int(codec.id))
	if parser == nil {
		log.Fatal("Parser not found")
	}
	c := C.avcodec_alloc_context3(codec)
	if c == nil {
		log.Fatal("Could not allocate audio codec context")
	}
	c.channels = 1
	if C.avcodec_open2(c, codec, nil) < 0 {
		log.Fatal("Could not open codec")
	}
	decodedFrame := C.av_frame_alloc()

	return &G7231Decoder{
		pkt:          pkt,
		codec:        codec,
		parser:       parser,
		c:            c,
		decodedFrame: decodedFrame,
	}
}

// ConsumesPayloadType will return whether or not the given decoder
// consumes the payload type (specified in the RTP payload type RFC 3550)
// https://en.wikipedia.org/wiki/RTP_payload_formats
func (decoder *G7231Decoder) ConsumesPayloadType(plt int) bool {
	return plt == 4
}

// Decode will decode all of the input packets
func (decoder *G7231Decoder) Decode(input []byte) ([]byte, error) {
	data := unsafe.Pointer(&input[0])

	ptrindex := 0
	dataSize := len(input) - 1

	var result *C.uchar
	resultSize := C.int(0)
	buf := &bytes.Buffer{}

	for dataSize > 0 {
		ret := C.decode_frame(
			decoder.pkt,
			decoder.codec,
			decoder.parser,
			decoder.c,
			decoder.decodedFrame,
			(*C.uchar)(data),
			(C.ulong)(dataSize),
			&result,
			&resultSize,
		)

		if ret < 0 {
			break
		}
		buf.Write(C.GoBytes(unsafe.Pointer(result), resultSize))
		C.free(unsafe.Pointer(result))

		if ptrindex+int(ret) > len(input)-1 {
			break
		}

		ptrindex += int(ret)
		dataSize -= int(ret)
		data = unsafe.Pointer(&input[ptrindex])
	}
	return buf.Bytes(), nil
}

// Destroy will free all of the memory allocated by the decoder
func (decoder *G7231Decoder) Destroy() {
	C.avcodec_free_context(&decoder.c)
	C.av_parser_close(decoder.parser)
	C.av_packet_free(&decoder.pkt)
}
