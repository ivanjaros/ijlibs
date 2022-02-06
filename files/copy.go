package files

import (
	"io"
	"os"
	"path/filepath"
)

func Copy(from, to string) error {
	if err := os.MkdirAll(filepath.Dir(to), 0666); err != nil {
		return err
	}

	src, err := os.Open(from)
	if err != nil {
		return err
	}
	defer src.Close()

	trg, err := os.OpenFile(to, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer src.Close()

	_, err = io.Copy(trg, src)
	return err
}
