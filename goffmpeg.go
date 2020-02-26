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
	"errors"
	"fmt"
	"log"
	"unsafe"
)

const G729RTPPayloadType = 18
const G723RTPPayloadType = 4

// Codec wraps some properties of a given codec
type Codec struct {
	codecID        C.enum_AVCodecID
	rtpPayloadType int
}

var (
	G729 = Codec{
		codecID:        C.AV_CODEC_ID_G729,
		rtpPayloadType: G729RTPPayloadType,
	}

	G723 = Codec{
		codecID:        C.AV_CODEC_ID_G723_1,
		rtpPayloadType: G723RTPPayloadType,
	}

	DecoderDestroyedError = errors.New("Cannot decode frame after destroy has been called")
)

var (
	codecs = map[string]Codec{
		"G729": G729,
		"G723": G723,
	}
)

// Decoder is an interface borrowed from the `cart` project
type Decoder interface {
	Decode([]byte) ([]byte, error)
	Destroy()
	GetRTPPayloadType() int
}

var _ Decoder = &FFMPEGDecoder{}

// FFMPEGDecoder is a struct used for decoding audio in various
// ffmpeg supported protocols
type FFMPEGDecoder struct {
	freed        bool
	payloadType  int
	pkt          *C.struct_AVPacket
	codec        *C.struct_AVCodec
	parser       *C.struct_AVCodecParserContext
	context      *C.struct_AVCodecContext
	decodedFrame *C.struct_AVFrame
}

// NewFFMPEGDecoder will return a new FFMPEGDecoder
func NewFFMPEGDecoder(codecName string) (interface{}, error) {

	codec := codecs[codecName]
	ffmpegCodec := getCodec(codec.codecID)

	return &FFMPEGDecoder{
		freed:        false,
		payloadType:  codec.rtpPayloadType,
		pkt:          C.av_packet_alloc(),
		codec:        ffmpegCodec,
		parser:       getParser(C.int(ffmpegCodec.id)),
		context:      getContext(ffmpegCodec),
		decodedFrame: C.av_frame_alloc(),
	}, nil
}

func getCodec(codecType C.enum_AVCodecID) *C.struct_AVCodec {

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

// GetRTPPayloadType will return the given decoder payload type
// specified in the RTP payload type RFC 3550)
// https://en.wikipedia.org/wiki/RTP_payload_formats
func (decoder *FFMPEGDecoder) GetRTPPayloadType() int {
	return decoder.payloadType
}

// Decode will decode all of the input packets
func (decoder *FFMPEGDecoder) Decode(input []byte) ([]byte, error) {

	if decoder.freed {
		return []byte(""), DecoderDestroyedError
	}

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
	if !decoder.freed {
		decoder.destroy()
		decoder.freed = true
	}
}

func (decoder *FFMPEGDecoder) destroy() {
	C.av_parser_close(decoder.parser)
	C.avcodec_free_context(&decoder.context)
	C.av_packet_free(&decoder.pkt)
}
