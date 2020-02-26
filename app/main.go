package main

import (
	"fmt"
	"github.com/corticph/goffmpeg"
	"io/ioutil"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// This small sample will decode a g723.1 audio file and output a raw PCM
// audio file. This file can be played with the following ffmpeg command:
// ffplay -f s16le -ar 8k -ac 1 outfile.wav
func main() {

	pflag.StringP("codec", "c", "", "The name of the codec to use")
	pflag.StringP("input", "i", "", "The path of the file to decode")
	pflag.StringP("output", "o", "", "The path where to save the decoded file")
	pflag.Parse()
	bindViperFlags()

	input := readFile(viper.GetString("input"))
	decoder := getDecoder(viper.GetString("codec"))
	defer decoder.Destroy()

	result := decode(decoder, input)
	writeFile(result, viper.GetString("output"))
	fmt.Printf("%v was written to disk (%d bytes)\n", viper.GetString("output"), len(result))
}

func decode(decoder goffmpeg.Decoder, input []byte) []byte {

	result, err := decoder.Decode(input)
	if err != nil {
		panic(err)
	}

	return result
}

func bindViperFlags() {

	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		panic(err)
	}
}

func getDecoder(codec string) goffmpeg.Decoder {

	d, err := goffmpeg.NewFFMPEGDecoder(codec)
	if err != nil {
		panic(err)
	}

	decoder, ok := d.(goffmpeg.Decoder)

	if !ok {
		panic("could not cast decoder")
	}

	return decoder

}

func writeFile(data []byte, path string) {

	if err := ioutil.WriteFile(path, data, 0755); err != nil {
		panic(err)
	}
}

func readFile(path string) []byte {

	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return data
}
