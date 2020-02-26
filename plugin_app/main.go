package main

import "plugin"

type goffmpegNewFuncSignature = func(string) (interface{}, error)

type Decoder interface {
	Decode([]byte) ([]byte, error)
	Destroy()
	GetRTPPayloadType() int
}

func loadPlugin(path string) Decoder {

	plug := openPlugin(path)
	symbol := getSymbol(plug)
	initFunc := getInitFunc(symbol)
	return initDecoder(initFunc)
}

func openPlugin(path string) *plugin.Plugin {

	plug, err := plugin.Open(path)
	if err != nil {
		panic(err)
	}

	return plug

}

func getSymbol(plug *plugin.Plugin) plugin.Symbol {
	symbol, err := plug.Lookup("NewFFMPEGDecoder")
	if err != nil {
		panic(err)
	}
	return symbol
}

func getInitFunc(symbol plugin.Symbol) goffmpegNewFuncSignature {

	var ok bool
	newFn, ok := symbol.(func(string) (interface{}, error))
	if !ok {
		panic("could not cast")
	}

	return newFn

}

func initDecoder(initFunc goffmpegNewFuncSignature) Decoder {

	iface, err := initFunc("G729")
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
