package main

import (
	"fmt"
	"io/ioutil"

	"github.com/corticph/g72x/goffmpeg"
)

// This small sample will decode a g723.1 audio file and output a raw PCM
// audio file. This file can be played with the following ffmpeg command:
// ffplay -f s16le -ar 8k -ac 1 outfile.wav
func main() {
	infile, err := ioutil.ReadFile("sample.wav")
	if err != nil {
		panic(err)
	}
	decoder := goffmpeg.NewG7231Decoder()
	defer decoder.Destroy()
	data := decoder.Decode(infile)

	if err := ioutil.WriteFile("outfile.wav", data, 0755); err != nil {
		panic(err)
	}

	ofile, err := ioutil.ReadFile("outfile.wav")
	if err != nil {
		panic(err)
	}
	fmt.Printf("outfile.wav was written to disk (%d bytes)\n", len(ofile))
}