package goporting

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

const (
	tmpDirPrefix    = "gotestout"
	tmpFileName     = "tmpfile"
	G729InFilePath  = "testfiles/G729.raw"
	G729OutFilePath = "testfiles/G729.wav"
)

func TestDecodeG729(t *testing.T) {

	tmpDir := createTmpDir(t, tmpDirPrefix)
	defer os.RemoveAll(tmpDir)
	tmpFilePath := filepath.Join(tmpDir, tmpFileName)

	decoder := getDecoder(t)
	defer decoder.Destroy()

	byteStream := readFile(t, G729InFilePath)
	data := decodeData(t, decoder, byteStream)

	writeFile(t, tmpFilePath, data)
	writtenFile := readFile(t, tmpFilePath)

	expectedOutput := readFile(t, G729OutFilePath)

	assertFilesEqual(t, writtenFile, expectedOutput)
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

	t.Helper()

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

	t.Helper()

	data, err := decoder.Decode(in)
	if err != nil {
		t.Error(err)
	}
	return data
}

func writeFile(t *testing.T, path string, data []byte) {

	t.Helper()

	if err := ioutil.WriteFile(path, data, 0755); err != nil {
		t.Error(err)
	}
}

func assertFilesEqual(t *testing.T, b1, b2 []byte) {

	t.Helper()

	if !bytes.Equal(b1, b2) {
		t.Errorf("bytestream of files are different. Lengths: %d vs %d", len(b1), len(b2))
	}
}
