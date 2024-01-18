package zip

import (
	"archive/zip"
	"github.com/injoyai/goutil/oss"
	"github.com/injoyai/logs"
	"io"
	"os"
	"path/filepath"
)

// Encode 压缩文件
// @filePath,文件路径
// @zipPath,储存路劲
func Encode(filePath, zipPath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	zipFile, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	return compareZip(file, zipWriter, "", true)
}

// 压缩文件
func compareZip(file *os.File, zipWriter *zip.Writer, prefix string, top bool) error {
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}
	if fileInfo.IsDir() {
		if !top {
			prefix += "/" + fileInfo.Name()
		}
		fileInfoChilds, err := file.Readdir(-1)
		if err != nil {
			return err
		}
		if len(fileInfoChilds) == 0 {
			header, err := zip.FileInfoHeader(fileInfo)
			if err != nil {
				return err
			}
			header.Name = prefix + "/"
			_, err = zipWriter.CreateHeader(header)
			return err
		}
		for _, fileInfoChild := range fileInfoChilds {
			fileChild, err := os.Open(file.Name() + "/" + fileInfoChild.Name())
			if err != nil {
				return err
			}
			if err := compareZip(fileChild, zipWriter, prefix, false); err != nil {
				return err
			}
		}
		return nil
	}
	header, err := zip.FileInfoHeader(fileInfo)
	if err != nil {
		return err
	}
	header.Name = prefix + "/" + header.Name
	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, file)
	return err

}

// Decode 解压zip
func Decode(zipPath, filePath string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()
	for _, k := range r.Reader.File {
		var err error
		if k.FileInfo().IsDir() {
			logs.Debug("创建文件夹", filepath.Join(filePath, k.Name))
			if err := oss.New(filepath.Join(filePath, k.Name)); err != nil {
				return err
			}
		} else {
			r, err := k.Open()
			if err == nil {
				logs.Debug("创建文件", filepath.Join(filePath, k.Name))
				err = oss.New(filepath.Join(filePath, k.Name), r)
			}
		}
		if err != nil {
			return err
		}
	}
	return nil
}
