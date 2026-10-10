package kernelproto

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
)

//go:embed proto/tolo
var embedded embed.FS

func Files() fs.FS {
	sub, err := fs.Sub(embedded, "proto")
	if err != nil {
		panic(err)
	}

	return sub
}

func WriteProto(dir string) (written []string, err error) {
	err = os.MkdirAll(dir, 0o750)
	if err != nil {
		return nil, fmt.Errorf("create %s: %w", dir, err)
	}

	root, err := os.OpenRoot(dir)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", dir, err)
	}

	defer func() {
		err = errors.Join(err, root.Close())
	}()

	files := Files()

	err = fs.WalkDir(files, ".", func(name string, entry fs.DirEntry, err error) error {
		if err != nil || entry.IsDir() || path.Ext(name) != ".proto" {
			return err
		}

		content, err := fs.ReadFile(files, name)
		if err != nil {
			return fmt.Errorf("read embedded %s: %w", name, err)
		}

		written = append(written, filepath.Join(dir, filepath.FromSlash(name)))

		current, err := root.ReadFile(name)
		if err == nil && bytes.Equal(current, content) {
			return nil
		}

		err = root.MkdirAll(path.Dir(name), 0o750)
		if err != nil {
			return fmt.Errorf("create directory for %s: %w", name, err)
		}

		err = root.WriteFile(name, content, 0o644)
		if err != nil {
			return fmt.Errorf("write %s: %w", name, err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return written, nil
}
