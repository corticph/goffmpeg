
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "libavutil/frame.h"
#include "libavutil/mem.h"

#include "libavcodec/avcodec.h"

#define AUDIO_INBUF_SIZE 20480
#define AUDIO_REFILL_THRESH 4096

uint8_t * decode_single_channel(AVCodecContext *dec_ctx, AVPacket *pkt, AVFrame *frame,
                   int *result_size)
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

        return decoded;
    }
    return NULL;
}

int decode_frame(
    AVPacket *pkt,
    AVCodec *codec, AVCodecParserContext *parser,
    AVCodecContext *c, AVFrame *decoded_frame,
    uint8_t *data, size_t data_size,
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
        *result = decode_single_channel(c, pkt, decoded_frame, result_size);

    return pkt->size;
}
