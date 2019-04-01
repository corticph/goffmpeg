package main

import (
	"fmt"
	"io/ioutil"

	"github.com/corticph/g72x/goffmpeg"
)

func main() {
	infile, err := ioutil.ReadFile("sample.wav")
	if err != nil {
		panic(err)
	}
	decoder := goffmpeg.NewG7231Decoder()
	defer decoder.Destroy()
	data, err := decoder.Decode(infile)
	if err != nil {
		panic(err)
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
