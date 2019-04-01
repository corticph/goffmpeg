
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "libavutil/frame.h"
#include "libavutil/mem.h"

#include "libavcodec/avcodec.h"

#define AUDIO_INBUF_SIZE 20480
#define AUDIO_REFILL_THRESH 4096

uint8_t * decode_single_channel(AVCodecContext *dec_ctx, AVPacket *pkt, AVFrame *frame,
                   FILE *outfile, int *result_size)
{
    int i, ch;
    int ret, data_size;

    /* send the packet with the compressed data to the decoder */
    ret = avcodec_send_packet(dec_ctx, pkt);
    if (ret < 0) {
        fprintf(stderr, "Error submitting the packet to the decoder\n");
        exit(1);
    }

    /* read all the output frames (in general there may be any number of them */
    while (ret >= 0) {
        ret = avcodec_receive_frame(dec_ctx, frame);
        if (ret == AVERROR(EAGAIN) || ret == AVERROR_EOF)
            return NULL;
        else if (ret < 0) {
            fprintf(stderr, "Error during decoding\n");
            exit(1);
        }
        data_size = av_get_bytes_per_sample(dec_ctx->sample_fmt);
        if (data_size < 0) {
            /* This should not occur, checking just for paranoia */
            fprintf(stderr, "Failed to calculate data size\n");
            exit(1);
        }

        *result_size = data_size*frame->nb_samples;

        uint8_t *decoded = (uint8_t *) malloc(data_size*frame->nb_samples);
        memcpy(decoded, frame->data[0], data_size*frame->nb_samples);
        // fwrite(decoded, sizeof(uint8_t), data_size*frame->nb_samples, outfile);

        return decoded;
    }
}


int _decode_frame(
    AVPacket *pkt,
    AVCodec *codec, AVCodecParserContext *parser,
    AVCodecContext *c, AVFrame *decoded_frame,
    FILE *outfile, uint8_t *data, size_t data_size,
    uint8_t **result, int *result_size)
{
    int len, ret;

    if (!decoded_frame) { // paranoia
        if (!(decoded_frame = av_frame_alloc())) {
            fprintf(stderr, "Could not allocate audio frame\n");
            return -1;
        }
    }

    ret = av_parser_parse2(parser, c, &pkt->data, &pkt->size,
                            data, data_size,
                            AV_NOPTS_VALUE, AV_NOPTS_VALUE, 0);
    if (ret < 0) {
        fprintf(stderr, "Error while parsing\n");
        return ret;
    }

    if (pkt->size)
        *result = decode_single_channel(c, pkt, decoded_frame, outfile, result_size);

    return pkt->size;
}

int main(int argc, char **argv)
{
    char *outfilename, *filename;
    AVCodec *codec;
    AVFrame *decoded_frame = NULL;
    AVCodecContext *c= NULL;
    AVCodecParserContext *parser = NULL;
    FILE *outfile, *f;
    AVPacket *pkt;
    uint8_t inbuf[AUDIO_INBUF_SIZE + AV_INPUT_BUFFER_PADDING_SIZE];
    uint8_t *data;
    size_t data_size;
    int len, ret;

    if (argc <= 2) {
        fprintf(stderr, "Usage: %s <input file> <output file>\n", argv[0]);
        exit(0);
    }
    filename    = argv[1];
    outfilename = argv[2];

    pkt = av_packet_alloc();

    /* find the MPEG audio decoder */
    codec = avcodec_find_decoder(AV_CODEC_ID_G723_1);
    if (!codec) {
        fprintf(stderr, "Codec not found\n");
        exit(1);
    }

    parser = av_parser_init(codec->id);
    if (!parser) {
        fprintf(stderr, "Parser not found\n");
        exit(1);
    }

    c = avcodec_alloc_context3(codec);
    if (!c) {
        fprintf(stderr, "Could not allocate audio codec context\n");
        exit(1);
    }

    c->channels = 1;

    /* open it */
    if (avcodec_open2(c, codec, NULL) < 0) {
        fprintf(stderr, "Could not open codec\n");
        exit(1);
    }

    outfile = fopen(outfilename, "wb");
    if (!outfile) {
        av_free(c);
        exit(1);
    }

    f = fopen(filename, "rb");
    if (!f) {
        av_free(f);
        exit(1);
    }

    FILE *cout;
    cout = fopen("cout2.wav", "wb");
    if (!cout) {
        av_free(cout);
        exit(1);
    }

    data = inbuf;
    data_size = fread(inbuf, 1, AUDIO_INBUF_SIZE, f);


    while (data_size > 0) {
        printf("Data size: %d\n", data_size);
        uint8_t *result;
        int result_size;
        ret = _decode_frame(pkt, codec, parser, c, decoded_frame, outfile, data, data_size, &result, &result_size);
        if (ret < 0) {
            fprintf(stderr, "Could not decode file\n");
            exit(1);
        }

        fwrite(result, sizeof(uint8_t), result_size, outfile);

        data      += ret;
        data_size -= ret;

        if (data_size < AUDIO_REFILL_THRESH) {
            memmove(inbuf, data, data_size);
            data = inbuf;
            len = fread(data + data_size, 1,
                        AUDIO_INBUF_SIZE - data_size, f);
            if (len > 0)
                data_size += len;
        }
    }

    fclose(outfile);

    avcodec_free_context(&c);
    av_parser_close(parser);
    av_packet_free(&pkt);

    return 0;
}