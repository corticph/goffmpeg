GO111MODULE := on
VERSION := v1.0
NAME := goffmpeg
BUILDSTRING := $(shell git log --pretty=format:'%h' -n 1)
VERSIONSTRING := $(NAME) version $(VERSION)+$(BUILDSTRING)

LDFLAGS := "-extldflags \"-static\" -X \"main.VERSION=$(VERSIONSTRING)\" -X \"main.BUILDDATE=$(BUILDDATE)\""

default: build_plugin

build_static_plugin:
	go build -buildmode=plugin -o plugin.so -ldflags=$(LDFLAGS)

build_plugin:
	go build -buildmode=plugin -o plugin.so
