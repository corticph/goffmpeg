package main

import (
	"fmt"
	"github/corticph/g72x/goffmpeg/goporting"
	"io/ioutil"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	codecs = map[string]goporting.Codec{
		"G729": goporting.G729,
		"G723": goporting.G723,
	}
)

// This small sample will decode a g723.1 audio file and output a raw PCM
// audio file. This file can be played with the following ffmpeg command:
// ffplay -f s16le -ar 8k -ac 1 outfile.wav
func main() {

	pflag.StringP("codec", "c", "", "The name of the codec to use")
	pflag.StringP("input", "i", "", "The path of the file to decode")
	pflag.StringP("output", "o", "", "The path where to save the decoded file")

	pflag.Parse()

	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		panic(err)
	}

	infile, err := ioutil.ReadFile(viper.GetString("input"))
	if err != nil {
		panic(err)
	}
	d, err := goporting.NewFFMPEGDecoder(codecs[viper.GetString("codec")])
	if err != nil {
		panic(err)
	}

	decoder, ok := d.(goporting.Decoder)

	if !ok {
		panic("oh no")
	}

	defer decoder.Destroy()
	data, err := decoder.Decode(infile)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := ioutil.WriteFile(viper.GetString("output"), data, 0755); err != nil {
		panic(err)
	}

	ofile, err := ioutil.ReadFile(viper.GetString("output"))
	if err != nil {
		panic(err)
	}
	fmt.Printf("file was written to disk (%d bytes)\n", len(ofile))
}
