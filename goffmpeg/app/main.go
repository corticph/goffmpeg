package main

import (
	"fmt"
	"github/corticph/g72x/goffmpeg/goporting"
	"io/ioutil"
)

// This small sample will decode a g723.1 audio file and output a raw PCM
// audio file. This file can be played with the following ffmpeg command:
// ffplay -f s16le -ar 8k -ac 1 outfile.wav
func main() {
	infile, err := ioutil.ReadFile("../goporting/testfiles/G729.raw")
	if err != nil {
		panic(err)
	}
	d, err := goporting.New()
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

	if err := ioutil.WriteFile("outfile.wav", data, 0755); err != nil {
		panic(err)
	}

	ofile, err := ioutil.ReadFile("outfile.wav")
	if err != nil {
		panic(err)
	}
	fmt.Printf("outfile.wav was written to disk (%d bytes)\n", len(ofile))
}
