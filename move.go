package main

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"syscall"
)

func moveEntry(src, dst string) error {
	if err := os.Rename(src, dst); err == nil {
		return nil
	} else {
		var le *os.LinkError
		if errors.As(err, &le) && errors.Is(le.Err, syscall.EXDEV) {
			return copyThenDelete(src, dst)
		}
		return err
	}
}

func copyThenDelete(src, dst string) error {
	info, err := os.Stat(src)
	if err != nil {
		return err
	}
	if info.IsDir() {
		if err := copyDir(src, dst); err != nil {
			os.RemoveAll(dst) //nolint — best-effort cleanup on partial copy
			return err
		}
		return os.RemoveAll(src)
	}
	if err := copyFile(src, dst, info.Mode()); err != nil {
		return err
	}
	return os.Remove(src)
}

func copyDir(src, dst string) error {
	info, err := os.Stat(src)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dst, info.Mode()); err != nil {
		return err
	}
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	for _, e := range entries {
		s := filepath.Join(src, e.Name())
		d := filepath.Join(dst, e.Name())
		if e.Type()&os.ModeSymlink != 0 {
			target, err := os.Readlink(s)
			if err != nil {
				return err
			}
			if err := os.Symlink(target, d); err != nil {
				return err
			}
		} else if e.IsDir() {
			if err := copyDir(s, d); err != nil {
				return err
			}
		} else {
			fi, err := e.Info()
			if err != nil {
				return err
			}
			if err := copyFile(s, d, fi.Mode()); err != nil {
				return err
			}
		}
	}
	return nil
}

func copyFile(src, dst string, mode os.FileMode) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}
