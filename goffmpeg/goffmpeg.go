package goffmpeg

// #cgo CFLAGS: -I./include
// #cgo LDFLAGS: -L cffmpeg -lavcodec -lavdevice -lavfilter -lavformat -lavutil -lffmpeg -lswresample -lswscale
// #include <stdio.h>
// #include <stdlib.h>
// #include <string.h>
// #include "libavcodec/avcodec.h"
// #include "include/decoder.h"
import "C"
import (
	"bytes"
	"fmt"
	"log"
	"unsafe"
)

type Codec C.enum_AVCodecID

var (
	G729 Codec = C.AV_CODEC_ID_G729
	G723 Codec = C.AV_CODEC_ID_G723_1
)

var (
	codecs = map[string]Codec{
		"G729": G729,
		"G723": G723,
	}
)

const G729RTPPayloadType = 18

// Decoder is an interface borrowed from the `cart` project
type Decoder interface {
	Decode([]byte) ([]byte, error)
	Destroy()
	ConsumesPayloadType(int) bool
}

var _ Decoder = &FFMPEGDecoder{}

// FFMPEGDecoder is a struct used for decoding audio in various
// ffmpeg supported protocols
type FFMPEGDecoder struct {
	pkt          *C.struct_AVPacket
	codec        *C.struct_AVCodec
	parser       *C.struct_AVCodecParserContext
	context      *C.struct_AVCodecContext
	decodedFrame *C.struct_AVFrame
}

// New will return a new g723.1 decoder
func NewFFMPEGDecoder(codecName string) (interface{}, error) {
	codecType := codecs[codecName]
	pkt := C.av_packet_alloc()
	codec := getCodec(codecType)
	parser := getParser(C.int(codec.id))
	context := getContext(codec)
	decodedFrame := C.av_frame_alloc()

	return &FFMPEGDecoder{
		pkt:          pkt,
		codec:        codec,
		parser:       parser,
		context:      context,
		decodedFrame: decodedFrame,
	}, nil
}

func getCodec(codecType Codec) *C.struct_AVCodec {

	c := C.avcodec_find_decoder(uint32(codecType))
	if c == nil {
		log.Fatal("Codec not found")
	}
	return c

}

func getParser(id C.int) *C.struct_AVCodecParserContext {

	parser := C.av_parser_init(id)
	if parser == nil {
		log.Fatal("Parser not found")
	}

	return parser

}

func getContext(codec *C.struct_AVCodec) *C.struct_AVCodecContext {

	context := C.avcodec_alloc_context3(codec)
	if context == nil {
		log.Fatal("Could not allocate audio codec context")
	}

	context.channels = 1
	openContext(context, codec)
	return context

}

func openContext(context *C.struct_AVCodecContext, codec *C.struct_AVCodec) {
	if C.avcodec_open2(context, codec, nil) < 0 {
		log.Fatal("Could not open codec")
	}
}

// ConsumesPayloadType will return whether or not the given decoder
// consumes the payload type (specified in the RTP payload type RFC 3550)
// https://en.wikipedia.org/wiki/RTP_payload_formats
func (decoder *FFMPEGDecoder) ConsumesPayloadType(plt int) bool {
	return plt == G729RTPPayloadType
}

// Decode will decode all of the input packets
func (decoder *FFMPEGDecoder) Decode(input []byte) ([]byte, error) {

	firstIndex := 0
	lastIndex := len(input) - 1
	var result *C.uchar
	resultSize := C.int(0)
	buf := &bytes.Buffer{}

	for lastIndex > 0 {
		data := unsafe.Pointer(&input[firstIndex])
		bytesProcessed := C.decode_frame(
			decoder.pkt,
			decoder.codec,
			decoder.parser,
			decoder.context,
			decoder.decodedFrame,
			(*C.uchar)(data),
			(C.ulong)(lastIndex),
			&result,
			&resultSize,
		)

		if isFlushPackage(bytesProcessed) {
			break
		}

		if isDecodeError(bytesProcessed) {
			C.free(unsafe.Pointer(result))
			return []byte(""), fmt.Errorf("error while decoding frame in byte number %d (0 indexed)", firstIndex)

		}

		appendToBuffer(result, buf, resultSize)
		C.free(unsafe.Pointer(result))

		firstIndex += int(bytesProcessed)
		lastIndex -= int(bytesProcessed)
	}
	return buf.Bytes(), nil
}

func isFlushPackage(n C.int) bool {
	return n == 0
}

func isDecodeError(n C.int) bool {
	return n < 0
}

func appendToBuffer(data *C.uchar, buf *bytes.Buffer, size C.int) {

	dataBytes := C.GoBytes(unsafe.Pointer(data), size)
	buf.Write(dataBytes)
}

// Destroy will free all of the memory allocated by the decoder
func (decoder *FFMPEGDecoder) Destroy() {
	C.avcodec_free_context(&decoder.context)
	C.av_parser_close(decoder.parser)
	C.av_packet_free(&decoder.pkt)
}
