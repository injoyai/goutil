package gzip

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"
)

/*
gzip 标识 31 139 8 0 0 0 0 0 0 255
*/

// EncodeBytes 压缩字节
func EncodeBytes(input []byte) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	gzipWriter := gzip.NewWriter(buf)
	_, err := gzipWriter.Write(input)
	gzipWriter.Close()
	return buf.Bytes(), err
}

// DecodeBytes 解压字节
func DecodeBytes(input []byte) ([]byte, error) {
	reader := bytes.NewReader(input)
	gzipReader, err := gzip.NewReader(reader)
	if err != nil {
		return nil, err
	}
	defer gzipReader.Close()
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(gzipReader)
	return buf.Bytes(), err
}

func EncodeFile(inPath, outPath string) error {
	inFile, err := os.Open(inPath)
	if err != nil {
		return err
	}
	inBytes, err := ioutil.ReadAll(inFile)
	if err != nil {
		return err
	}
	gzFile, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer gzFile.Close()
	gzipWriter := gzip.NewWriter(gzFile)
	defer gzipWriter.Close()
	gzipWriter.Name = inFile.Name()
	_, err = gzipWriter.Write(inBytes)
	return err
}

// DecodeFile 解压文件
func DecodeFile(input []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewBuffer(input))
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	var buff bytes.Buffer
	_, err = io.Copy(&buff, reader)
	return buff.Bytes(), err
}
