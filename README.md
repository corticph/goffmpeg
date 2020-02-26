# goffmpeg

This project aims at porting the C code from the [ffmpeg](https://www.ffmpeg.org/) library into GO. At this point in time, only decoding is supported.

## Installation

The standard library files only contain the minimal compiled libraries, which isn't good enough for the purpose of this Go library. Therefore, ensure that all shared libraries, headers and binaries are downloaded from a release in github.

As all the `.so` files have been compiled to run on a unix 64-bit system, it is recommeded to run this code with the following docker container (ran from the root of this repo):
>docker run --rm -it -e LD_LIBRARY_PATH='/go/src/github.com/corticph/goffmpeg/cffmpeg' -v $(pwd):/go/src/github.com/corticph/goffmpeg  golang

The environment variable `LD_LIBRARY_PATH` is essential for compilation and finding the necessary `.so:<n>` files.

## Usage

You can use the app in order to decode raw bitestreams to pcm files. For instance, do:
> go run app/main.go -i testfiles/G729.raw -o test.wav -c G729

To use the goffmpeg library in other applications it is reccomended to compile it as a plugin, just go to the `plugin` folder and run `make`.
