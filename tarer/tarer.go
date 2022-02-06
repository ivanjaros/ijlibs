package tarer

// Un/Packs an entire directory into a single tarball, optionally with gzip compression.
// Inspired by: https://golangdocs.com/tar-gzip-in-golang

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Tar(source string, compress bool) (string, error) {
	info, err := os.Stat(source)
	if err != nil {
		return "", err
	}

	if info.IsDir() == false {
		return "", errors.New("source must be a directory")
	}

	fileName := source + ".tar"

	if compress {
		fileName += ".gz"
	}

	target, err := os.Create(fileName)
	if err != nil {
		return "", err
	}
	defer target.Close()

	var w *tar.Writer

	if compress {
		gW, err := gzip.NewWriterLevel(target, gzip.BestCompression)
		if err != nil {
			return "", err
		}
		defer gW.Close()
		w = tar.NewWriter(gW)
	} else {
		w = tar.NewWriter(target)
	}

	defer w.Close()

	return fileName, filepath.Walk(source,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			header, err := tar.FileInfoHeader(info, info.Name())
			if err != nil {
				return err
			}

			fp, _ := filepath.Rel(source, filepath.Join(filepath.Dir(path), info.Name()))
			header.Name = fp

			if err := w.WriteHeader(header); err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			file, err := os.Open(path)
			if err != nil {
				return err
			}

			_, err = io.Copy(w, file)

			_ = file.Close()

			return err
		})
}

// if the target has ".gz" suffix, it will automatically decompress it
func Untar(target string) error {
	if strings.HasSuffix(target, ".tar") == false && strings.HasSuffix(target, ".tar.gz") == false {
		return errors.New("file must be a tar archive")
	}

	src, err := os.Open(target)
	if err != nil {
		return err
	}
	defer src.Close()

	cut := ".tar"

	var reader io.Reader
	if strings.HasSuffix(target, ".gz") {
		gzR, err := gzip.NewReader(src)
		if err != nil {
			return err
		}
		defer gzR.Close()
		reader = gzR
		cut += ".gz"
	} else {
		reader = src
	}

	dir := target[:len(target)-len(cut)]
	tarReader := tar.NewReader(reader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		path := filepath.Join(dir, header.Name)
		info := header.FileInfo()
		if info.IsDir() {
			if err = os.MkdirAll(path, info.Mode()); err != nil {
				return err
			}
			continue
		}

		file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			return err
		}

		_, err = io.Copy(file, tarReader)

		_ = file.Close()

		if err != nil {
			return err
		}
	}

	return nil
}
