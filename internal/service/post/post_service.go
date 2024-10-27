package post

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"net/http"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/mfsyahrz/image_feed_api/internal/common/file"
	"github.com/mfsyahrz/image_feed_api/internal/common/logger"
	"github.com/mfsyahrz/image_feed_api/internal/common/pagination"
	"github.com/mfsyahrz/image_feed_api/internal/domain/entity"
	"github.com/mfsyahrz/image_feed_api/internal/domain/repository"
	"github.com/mfsyahrz/image_feed_api/internal/infrastructure/filestore"
	"github.com/nfnt/resize"
)

type PostService interface {
	GetPaginated(ctx context.Context, input repository.GetPostInput) (entity.Posts, *pagination.PostCursor, error)
	Save(ctx context.Context, input *CreatePostInput) (*entity.Post, error)
	SaveImage(ctx context.Context, fileHeader file.IFileHeader) (string, string, error)
}

type postService struct {
	postRepo    repository.PostRepository
	commentRepo repository.CommentRepository
	fileStore   filestore.FileStore
}

func NewPostService(postRepo repository.PostRepository, commentRepo repository.CommentRepository, fileStore filestore.FileStore) PostService {
	return &postService{
		postRepo:    postRepo,
		commentRepo: commentRepo,
		fileStore:   fileStore,
	}
}

func (s *postService) GetPaginated(ctx context.Context, input repository.GetPostInput) (entity.Posts, *pagination.PostCursor, error) {
	logger.FromCtx(ctx).Info("Getting Paginated posts...")

	posts, nextCursor, err := s.postRepo.FetchPaginated(ctx, input)
	if err != nil {
		return nil, nil, err
	}

	if posts.Len() == 0 {
		return nil, nil, nil
	}

	s.loadURL(posts)
	return posts, nextCursor, nil
}

func (s *postService) loadURL(posts entity.Posts) {
	url := s.fileStore.GetBaseURL()
	for _,p := range posts {
		p.SetDisplayImgURL(url)
	}
}
 
func (s *postService) Save(ctx context.Context, input *CreatePostInput) (*entity.Post, error) {
	logger.FromCtx(ctx).Info("Saving posts...")

	post := &entity.Post{
		Creator:    input.Creator,
		Caption:    input.Caption,
		SrcImg:     input.SrcImgPath,
		DisplayImg: input.DisplayImgPath,
	}

	if err := s.postRepo.Save(ctx, post); err != nil {
		return nil, err
	}

	return post, nil
}

func (s *postService) SaveImage(ctx context.Context, fh file.IFileHeader) (string, string, error) {
	if err := s.validateImage(fh); err != nil {
		return "", "", err
	}

	displayImg, err := s.resizeAndConvert(fh)
	if err != nil {
		return "", "", fmt.Errorf("unable to resize file %s. +%v", fh.GetFilename(), err)
	}

	srcImgName := s.buildImgName(srcPrefix, filepath.Ext(fh.GetFilename()))
	displayImgName := s.buildImgName(displayPrefix, extJpeg)

	srcImg, err := fh.Open()
	if err != nil {
		return "", "", fmt.Errorf("unable to open file %s. +%v", fh.GetFilename(), err)
	}
	defer srcImg.Close()

	if err := s.fileStore.Save(ctx, filestore.File{
		Name:   srcImgName,
		Dir:    srcImgDir,
		Object: srcImg,
	}); err != nil {
		return "", "", fmt.Errorf("unable to store source image %s. +%v", srcImgName, err)
	}

	if err := s.fileStore.Save(ctx, filestore.File{
		Name:   displayImgName,
		Dir:    displayImgDir,
		Object: bytes.NewReader(displayImg),
	}); err != nil {
		return "", "", fmt.Errorf("unable to store display image %s. +%v", displayImgName, err)
	}

	srcImgPath := filepath.Join(srcImgDir, srcImgName)
	displayImgPath := filepath.Join(displayImgDir, displayImgName)
	return srcImgPath, displayImgPath, nil
}

func (s *postService) buildImgName(prefix string, ext string) string {
	uid := uuid.New().String()
	return fmt.Sprintf(imgNameFmt, prefix, uid, ext)
}

// validateImage ensures the uploaded image is of a valid format and size.
func (s *postService) validateImage(fh file.IFileHeader) error {
	if fh.GetSize() > maxFileSize {
		return fmt.Errorf("file size exceeds the limit of %d MB", maxFileSize)
	}

	f, err := fh.Open()
	if err != nil {
		return fmt.Errorf("unable to open file %s. +%v", fh.GetFilename(), err)
	}
	defer f.Close()

	buf := make([]byte, 512)
	if _, err := f.Read(buf); err != nil {
		return fmt.Errorf("unable to read file %s. +%v", fh.GetFilename(), err)
	}

	fileType := http.DetectContentType(buf)
	if !allowedFormats[fileType] {
		return fmt.Errorf("unsupported file format: %s", fileType)
	}

	return nil
}

// resizeAndConvert resizes the image to expected width and converts it to acceptable format.
func (s *postService) resizeAndConvert(fh file.IFileHeader) ([]byte, error) {
	f, err := fh.Open()
	if err != nil {
		return nil, fmt.Errorf("unable to open file %s. +%v", fh.GetFilename(), err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("unable to decode image: %v", err)
	}

	resizedImage := resize.Resize(defaultImageRatio, 0, img, resize.Lanczos3)

	buf := new(bytes.Buffer)
	if err := jpeg.Encode(buf, resizedImage, nil); err != nil {
		return nil, fmt.Errorf("unable to convert image to JPEG: %v", err)
	}

	return buf.Bytes(), nil
}
