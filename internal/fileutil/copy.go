// Package fileutil provides file operation utilities.
package fileutil

import (
	"io"
	"os"
)

// CopyFile copies the file at src to dst. If dst already exists it is
// overwritten. Partial files are cleaned up on error.
func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}

	_, copyErr := io.Copy(out, in)
	closeErr := out.Close()
	if copyErr != nil {
		os.Remove(dst) // clean up partial file
		return copyErr
	}
	if closeErr != nil {
		os.Remove(dst)
		return closeErr
	}
	return nil
}
