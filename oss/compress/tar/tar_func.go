package tar

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Decode 解压 tar.gz 文件到指定目录
func Decode(filename, dir string) error {
	// 打开源文件
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	// 创建 gzip reader
	gz, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gz.Close()

	// 创建 tar reader
	tr := tar.NewReader(gz)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break // 读完了
		}
		if err != nil {
			return err
		}

		target := filepath.Join(dir, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			// 创建目录
			if err := os.MkdirAll(target, os.FileMode(header.Mode)); err != nil {
				return err
			}
		case tar.TypeReg:
			// 创建父目录
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return err
			}
			// 创建文件
			outFile, err := os.Create(target)
			if err != nil {
				return err
			}
			if _, err := io.Copy(outFile, tr); err != nil {
				outFile.Close()
				return err
			}
			outFile.Close()
			// 设置权限
			if err := os.Chmod(target, os.FileMode(header.Mode)); err != nil {
				return err
			}
		default:
			fmt.Printf("跳过不支持的类型: %c in %s\n", header.Typeflag, header.Name)
		}
	}
	return nil
}
