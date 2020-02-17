package main

import "github/corticph/g72x/goffmpeg"

func NewFFMPEGDecoder(codecName string) (interface{}, error) {
	return goffmpeg.NewFFMPEGDecoder(codecName)
}
