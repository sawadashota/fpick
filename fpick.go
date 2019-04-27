package fpick

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Client of directory walker
type Client struct {
	src string
	dst string
}

// New client
// src is starting point of seeking files
// dst is directory for output
// if dst doesn't exists, make directory recursively
func New(src, dst string) (*Client, error) {
	if _, err := os.Stat(src); err != nil {
		return nil, err
	}

	return &Client{
		src: src,
		dst: dst,
	}, nil
}

// FileMatcher return matcher function
type FileMatcher func(string) bool

// FilenameExtractMatch .
func FilenameExtractMatch(filename string) FileMatcher {
	return func(s string) bool {
		return s == filename
	}
}

// FilenameRegexMatch .
func FilenameRegexMatch(regex string) (FileMatcher, error) {
	r, err := regexp.Compile(regex)
	if err != nil {
		return nil, err
	}

	return func(s string) bool {
		return r.MatchString(s)
	}, nil
}

// File .
type File struct {
	Path string
	Perm os.FileMode
}

func (f *File) filename() string {
	return filepath.Base(f.Path)
}

func (f *File) dir() string {
	return filepath.Dir(f.Path)
}

// Read the file body
func (f *File) Read() ([]byte, error) {
	body, err := ioutil.ReadFile(f.Path)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// Copy to given dir
func (f *File) Copy(dst string) error {
	if err := mkdirUnlessExists(dst); err != nil {
		return err
	}

	sourceFileStat, err := os.Stat(f.Path)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", f.Path)
	}

	source, err := os.Open(f.Path)
	if err != nil {
		return err
	}
	defer source.Close()

	dstFilePath := filepath.Join(dst, f.filename())

	destination, err := os.Create(dstFilePath)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}

func mkdirUnlessExists(dir string) error {
	if _, err := os.Stat(dir); err != nil {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

// FileList .
func (c *Client) FileList(match FileMatcher) ([]*File, error) {
	var fs []*File
	err := filepath.Walk(c.src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		if !match(info.Name()) {
			return nil
		}

		f := &File{
			Path: path,
			Perm: info.Mode(),
		}

		fs = append(fs, f)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return fs, nil
}

type OutputOption func(string) string

func OutputFlatDirOption(path string) string {
	return strings.ReplaceAll(path, "/", "__")
}

// Pick files and copy to another directory
func (c *Client) Pick(match FileMatcher, opts ...OutputOption) error {
	fs, err := c.FileList(match)
	if err != nil {
		return err
	}

	for _, f := range fs {
		dir := strings.TrimPrefix(f.dir(), c.src)
		dir = strings.TrimPrefix(dir, "/")

		for _, opt := range opts {
			dir = opt(dir)
		}

		dst := filepath.Join(c.dst, dir)
		if err := f.Copy(dst); err != nil {
			return err
		}
	}
	return nil
}
