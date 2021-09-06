package main

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strings"
	"unicode"

	"gopkg.in/yaml.v3"
)

type (
	Metadata struct {
		Categories []string `yaml:",omitempty"`
		Tags       []string `yaml:",omitempty"`
		Date       string   `yaml:",omitempty"`
		Title      string
	}

	Post struct {
		Metadata
		Content []byte
	}
)

func ReadPost(r io.Reader) (*Post, error) {
	scanner := bufio.NewScanner(r)
	firstline := true
	read_metadata := false
	post := Post{}
	var mdbuf []byte

	for scanner.Scan() {
		line := scanner.Bytes()
		line = append(line, '\n')

		if firstline {
			firstline = false

			if bytes.Equal(line[:3], []byte("---")) {
				read_metadata = true
				continue
			}
		}

		if read_metadata {
			if bytes.Equal(line[:3], []byte("---")) {
				read_metadata = false
				if err := yaml.Unmarshal(mdbuf, &post.Metadata); err != nil {
					return nil, err
				}
				continue
			}

			mdbuf = append(mdbuf, line...)
		} else {
			post.Content = append(post.Content, line...)
		}
	}
	return &post, nil
}

func ReadPostFromFile(path string) (*Post, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return ReadPost(f)
}

func mogrifyTitle(r rune) rune {
	switch {
	case unicode.IsLetter(r):
		return r
	case r == ' ' || r == '_' || r == '-':
		return '-'
	default:
		return -1
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}

func (post *Post) Slug(maxlen int) string {
	slug := strings.Map(mogrifyTitle, strings.ToLower(post.Metadata.Title))
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}

	return slug[:min(len(slug), maxlen)]
}

func (post *Post) Write(w io.Writer) error {
	encoder := yaml.NewEncoder(w)
	encoder.SetIndent(2)

	if _, err := w.Write([]byte("---\n")); err != nil {
		return err
	}

	err := encoder.Encode(post.Metadata)
	if err != nil {
		return err
	}
	if _, err := w.Write([]byte("---\n")); err != nil {
		return err
	}

	if _, err := w.Write(post.Content); err != nil {
		return err
	}

	return nil
}

func (post *Post) WriteToFile(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return post.Write(f)
}
