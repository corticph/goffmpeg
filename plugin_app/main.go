package main

import "plugin"

type Decoder interface {
	Decode([]byte) ([]byte, error)
	Destroy()
	GetRTPPayloadType() int
}

func loadPlugin(path string) Decoder {
	plug, err := plugin.Open(path)
	if err != nil {
		panic(err)
	}
	symbol, err := plug.Lookup("NewFFMPEGDecoder")
	if err != nil {
		panic(err)
	}
	var ok bool
	newFn, ok := symbol.(func(string) (interface{}, error))
	if !ok {
		panic("could not cast")
	}
	iface, err := newFn("G729")
	if err != nil {
		panic("could not call init func")
	}
	dec, ok := iface.(Decoder)
	if !ok {
		panic("could not cast newFn to proper decoder type")
	}
	return dec
}

func main() {
	loadPlugin("../plugin/plugin.so")

}
