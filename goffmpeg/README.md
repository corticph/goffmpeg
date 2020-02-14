# goffmpeg

This project aims at porting the C code from the [ffmpeg](https://www.ffmpeg.org/) library into GO.

The standard library files only contain the minimal compiled libraries, which isn't good enough for the purpose of this Go library. Therefore, ensure that all shared libraries, headers and binaries are downloaded from here: https://github.com/corticph/g72x/releases/tag/g723_1_v1

As all the `.so` files have been compiled to run on a unix 64-bit system, it is recommeded to run this code with the following docker container:
> docker run --rm -it -e LD_LIBRARY_PATH='/go/src/github.com/corticph/g72x/goffmpeg/lib' -v (pwd):/go/src/github.com/corticph/g72x golang

The environment variable `LD_LIBRARY_PATH` is essential for compilation and finding the necessary `.so:<n>` files. To see an example run, go to the `goffmpeg/samples` directory and run the go application

> go run main.go

In this directory there is also a c program sample, which should be moved to the `goffmpeg` root before being compiled with the following command:

> gcc -o ffmpeg main.c -I./include -I. -L lib -lavcodec -lavdevice -lavfilter -lavformat -lavutil -lffmpeg -lswresample -lswscale; ./ffmpeg sample.wav cout.wav

