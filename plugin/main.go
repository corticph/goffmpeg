package main

import "github.com/corticph/goffmpeg"

func NewFFMPEGDecoder(codecName string) (interface{}, error) {
	return goffmpeg.NewFFMPEGDecoder(codecName)
}
