package goporting

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

const (
	tmpDirPrefix = "gotestout"
	tmpFileName  = "tmpfile"
	G729FilePath = "testfiles/G729.wav"
)

func TestDecodeG729(t *testing.T) {

	tmpDir := createTmpDir(t, tmpDirPrefix)
	defer os.RemoveAll(tmpDir)
	tmpFilePath := filepath.Join(tmpDir, tmpFileName)

	decoder := getDecoder(t)
	defer decoder.Destroy()

	byteStream := readFile(t, G729FilePath)
	data := decodeData(t, decoder, byteStream)

	writeFile(t, tmpFilePath, data)
	outputFile := readFile(t, tmpFilePath)
}

func readFile(t *testing.T, path string) []byte {

	t.Helper()

	f, err := ioutil.ReadFile(path)
	if err != nil {
		t.Error(err)
	}

	return f

}

func createTmpDir(t *testing.T, prefix string) string {

	t.Helper()

	dir, err := ioutil.TempDir("", prefix)
	if err != nil {
		t.Error(err)
	}

	return dir

}

func getDecoder(t *testing.T) Decoder {

	d, err := New()
	if err != nil {
		t.Error(err)
	}

	decoder, ok := d.(Decoder)

	if !ok {
		t.Errorf("interface type not of Decoder")
	}

	return decoder

}

func decodeData(t *testing.T, decoder Decoder, in []byte) []byte {

	data, err := decoder.Decode(in)
	if err != nil {
		t.Error(err)
	}
	return data
}

func writeFile(t *testing.T, path string, data []byte) {
	if err := ioutil.WriteFile(path, data, 0755); err != nil {
		t.Error(err)
	}
}
