Not really a readme, sorry for getting you exciting...

Anyway. Since the compilation of the library is very sensitive to the platform architecture, it's recommended to use a docker container with a debian based platform. To use an interactive docker container, use the following docker command:

> docker run --rm -it -v $(pwd):/go/src/app golang

In this interative docker container, you can run the go application inside the `/go/src/app` directory:

> go run main.go

> NOTE: you might have to run a go get "github.com/davecgh/go-spew/spew" if I have been lazy enough to not update the files with a mod file.

If all goes well, expect to see something like the following output: 

```
([]uint8) (len=9 cap=9) {
 00000000  ea ea 02 7b 20 7d 42 4d  c6                       |...{ }BM.|
}
(main._Ctype_int) 24
(main._Ctype_short) -30644
([]uint8) (len=9 cap=9) {
 00000000  f4 da 0d 9d 08 04 00 9e  a2                       |.........|
}
(main._Ctype_void) {
}
(main._Ctype_short) 0
([]uint8) (len=9 cap=9) {
 00000000  f4 da 0d 9d 08 04 00 9e  a2                       |.........|
}
```

As for the code, the code is a static library which is compiled from the g723 directory of the ITU-T recommendation, used in the onionphone repository: https://github.com/gegel/onionphone

All the files in the `libcodecs/g723` directory has been compiled and then put into a static library.

> gcc -o file1.o file1.c

And then...

> ar -rsc libg723.a file1.o file2.o file3.o

Then all the `.h` files from `libcodecs/g723`, as well as `common/inc/ophtools.h`, have been copied into the `go723` directory of this repository, as well as the static library. 