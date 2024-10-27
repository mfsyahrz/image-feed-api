package file

import "mime/multipart"

// FileHeader abstracts the functionality needed to handle file uploads
type IFileHeader interface {
	Open() (multipart.File, error)
	GetFilename() string
	GetSize() int64 // in byte
}

func NewFileHeader(fh *multipart.FileHeader) IFileHeader {
	return &customFileHeader{fh}
}

// customFileHeader wraps multipart.FileHeader to satisfy FileHeader interface
type customFileHeader struct {
	*multipart.FileHeader
}

// Open overrides the Open method to return the associated file
func (c *customFileHeader) Open() (multipart.File, error) {
	return c.FileHeader.Open()
}

// GetFilename returns the filename
func (c *customFileHeader) GetFilename() string {
	return c.FileHeader.Filename
}

// GetSize returns the size of the file
func (c *customFileHeader) GetSize() int64 {
	return c.FileHeader.Size
}
