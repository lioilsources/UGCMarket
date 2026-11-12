// Package storage keeps asset files on the local disk (a bind mount on JODA).
// Keys are relative paths like "products/12/model.fbx"; files are never
// handed out directly, only through market-api (media / download tokens).
package storage

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type FS struct{ root string }

func NewFS(root string) (*FS, error) {
	if err := os.MkdirAll(root, 0o755); err != nil {
		return nil, err
	}
	return &FS{root: root}, nil
}

func (f *FS) path(key string) (string, error) {
	clean := filepath.Clean("/" + key)
	if strings.Contains(key, "..") || clean == "/" {
		return "", fmt.Errorf("neplatný storage klíč %q", key)
	}
	return filepath.Join(f.root, clean), nil
}

// Put writes the content atomically (tmp + rename) and returns the size and sha256.
func (f *FS) Put(key string, r io.Reader) (int64, string, error) {
	dst, err := f.path(key)
	if err != nil {
		return 0, "", err
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return 0, "", err
	}
	tmp, err := os.CreateTemp(filepath.Dir(dst), ".put-*")
	if err != nil {
		return 0, "", err
	}
	defer os.Remove(tmp.Name())
	h := sha256.New()
	n, err := io.Copy(io.MultiWriter(tmp, h), r)
	if cerr := tmp.Close(); err == nil {
		err = cerr
	}
	if err != nil {
		return 0, "", err
	}
	if err := os.Rename(tmp.Name(), dst); err != nil {
		return 0, "", err
	}
	return n, hex.EncodeToString(h.Sum(nil)), nil
}

func (f *FS) Open(key string) (*os.File, error) {
	p, err := f.path(key)
	if err != nil {
		return nil, err
	}
	return os.Open(p)
}
