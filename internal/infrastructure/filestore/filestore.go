package filestore

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/mfsyahrz/image_feed_api/internal/config"
)

type File struct {
	Name   string
	Dir    string
	Object io.Reader
}

type FileStore interface {
	Save(ctx context.Context, file File) error
	GetBaseURL() string
}

// localFileStore implements FileStore for local file system storage
type localFileStore struct {
	basePath string
	baseURL  string
}

func NewFileStore(cfg *config.FileStorage) (FileStore, error) {
	if cfg.BasePath == "" || cfg.BaseURL == "" {
		return nil, errors.New("invalid filestore config")
	}

	return &localFileStore{
		basePath: cfg.BasePath,
		baseURL:  cfg.BaseURL,
	}, nil
}

func (f *localFileStore) GetBaseURL() string {
	return f.baseURL + f.basePath
}

func (f *localFileStore) Save(ctx context.Context, file File) error {
	dir := filepath.Join(f.basePath, file.Dir)

	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create directories %s. error: +%v", file.Dir, err)
	}

	fileDir := filepath.Join(dir, file.Name)
	outFile, err := os.Create(fileDir)
	if err != nil {
		return fmt.Errorf("failed to create file %s. error: +%v. ", fileDir, err)
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, file.Object)
	if err != nil {
		return fmt.Errorf("failed to copy file. error: +%v ", err)
	}

	return nil
}
