package post

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/gommon/log"
	"github.com/mfsyahrz/image_feed_api/internal/common/file"
	"github.com/mfsyahrz/image_feed_api/internal/config"
	"github.com/mfsyahrz/image_feed_api/internal/domain/entity"
	"github.com/mfsyahrz/image_feed_api/internal/domain/repository"
	"github.com/mfsyahrz/image_feed_api/internal/infrastructure/filestore"
	repoMock "github.com/mfsyahrz/image_feed_api/mocks/domain/repository"
	fileStoreMock "github.com/mfsyahrz/image_feed_api/mocks/infrastructure/filestore"
)

var (
	testLog = log.New("test")
)

func TestNewPost(t *testing.T) {

	t.Run("setup", func(t *testing.T) {
		_ = NewPostService(nil, nil, nil)
	})
}

func TestPostService_GetPaginated(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostRepo := repoMock.NewMockPostRepository(ctrl)
	mockCommentRepo := repoMock.NewMockCommentRepository(ctrl)
	mockFileStore := fileStoreMock.NewMockFileStore(ctrl)

	postService := NewPostService(mockPostRepo, mockCommentRepo, mockFileStore)

	expectedPosts := entity.Posts{
		&entity.Post{ID: 1, Comments: entity.Comments{{PostID: 1}}},
	}

	tests := []struct {
		name          string
		input         repository.GetPostInput
		mockPosts     entity.Posts
		mockComments  entity.Comments
		setupMock     func()
		expectedPosts entity.Posts
		wantErr       bool
	}{
		{
			name: "successful retrieval of paginated posts with comments",
			input: repository.GetPostInput{
				Limit: 5,
			},
			mockPosts: entity.Posts{
				&entity.Post{ID: 1, Comments: nil},
			},
			mockComments: entity.Comments{
				&entity.Comment{PostID: 1},
			},
			setupMock: func() {
				mockPostRepo.EXPECT().FetchPaginated(gomock.Any(), gomock.Any()).Return(expectedPosts, nil, nil)
				mockFileStore.EXPECT().GetBaseURL().Return("url")
			},
			expectedPosts: entity.Posts{
				&entity.Post{ID: 1, DisplayImgURL: "url", Comments: entity.Comments{{PostID: 1}}},
			},

			wantErr: false,
		},
		{
			name: "error retrieving posts",
			input: repository.GetPostInput{
				Limit: 5,
			},
			mockPosts:    nil,
			mockComments: nil,
			setupMock: func() {
				mockPostRepo.EXPECT().FetchPaginated(gomock.Any(), gomock.Any()).Return(nil, nil, errors.New("unexpected error"))
			},
			expectedPosts: nil,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			posts, _, err := postService.GetPaginated(context.Background(), tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("expected error: %v, got: %v", tt.wantErr, err)
			}

			if !reflect.DeepEqual(posts, tt.expectedPosts) {
				t.Errorf("expected posts %v, got %v", tt.expectedPosts, posts)
			}
		})
	}
}

func TestPostService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostRepo := repoMock.NewMockPostRepository(ctrl)
	mockCommentRepo := repoMock.NewMockCommentRepository(ctrl)
	mockFileStore := fileStoreMock.NewMockFileStore(ctrl)

	postService := NewPostService(mockPostRepo, mockCommentRepo, mockFileStore)

	tests := []struct {
		name          string
		input         *CreatePostInput
		expectedPost  *entity.Post
		expectedError error
	}{
		{
			name: "successful creation of post",
			input: &CreatePostInput{
				Creator:        "user1",
				Caption:        "A new post",
				SrcImgPath:     "source.jpg",
				DisplayImgPath: "display.jpg",
			},
			expectedPost: &entity.Post{
				Creator:    "user1",
				Caption:    "A new post",
				SrcImg:     "source.jpg",
				DisplayImg: "display.jpg",
			},
			expectedError: nil,
		},
		{
			name: "error saving post",
			input: &CreatePostInput{
				Creator: "user1",
				Caption: "A new post",
			},
			expectedPost:  nil,
			expectedError: errors.New("error saving post"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectedError != nil {
				mockPostRepo.EXPECT().
					Save(gomock.Any(), gomock.Any()).
					Return(tt.expectedError)
			} else {
				mockPostRepo.EXPECT().
					Save(gomock.Any(), gomock.Any()).
					Return(nil)
			}

			post, err := postService.Save(context.Background(), tt.input)
			if err != nil && err.Error() != tt.expectedError.Error() {
				t.Fatalf("expected error %v, got %v", tt.expectedError, err)
			}

			if post != nil && fmt.Sprint(entity.Posts{post}) == fmt.Sprint(entity.Posts{tt.expectedPost}) {
				t.Errorf("expected post %v, got %v", tt.expectedPost, post)
			}

			if !reflect.DeepEqual(post, tt.expectedPost) {
				t.Errorf("expected post %v, got %v", tt.expectedPost, post)
			}
		})
	}
}

func TestPostService_SaveImage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostRepo := repoMock.NewMockPostRepository(ctrl)
	mockCommentRepo := repoMock.NewMockCommentRepository(ctrl)

	basePath := "storage/"
	mockFileStore, _ := filestore.NewFileStore(&config.FileStorage{
		BaseURL:  "localhost:8080",
		BasePath: basePath,
	})

	postService := NewPostService(mockPostRepo, mockCommentRepo, mockFileStore)

	imgName := "test_img.jpeg"
	img, err := os.ReadFile(imgName)
	if err != nil {
		t.Fatalf("Failed to open sample image file: %v", err)
	}

	tests := []struct {
		name               string
		fileHeader         file.IFileHeader
		wantSrcImgPath     bool
		wantDisplayImgPath bool
		setupMock          func()
		wantErr            bool
	}{
		{
			name: "successful image save",
			fileHeader: &mockFileHeader{
				Filename: imgName,
				Size:     int64(len(img)),
				Content:  img,
			},
			wantSrcImgPath:     true,
			wantDisplayImgPath: true,
			wantErr:            false,
		},
		{
			name: "error validating image",
			fileHeader: &mockFileHeader{
				Filename: imgName,
				Size:     100000,
			},
			wantSrcImgPath:     false,
			wantDisplayImgPath: false,
			setupMock:          nil,
			wantErr:            true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupMock != nil {
				tt.setupMock()
			}

			srcImgPath, displayImgPath, err := postService.SaveImage(context.Background(), tt.fileHeader)
			if (err != nil) != tt.wantErr {
				t.Fatalf("expected error: %v, got: %v", tt.wantErr, err)
			}

			if (srcImgPath != "") != tt.wantSrcImgPath {
				t.Errorf("expected srcImgPath: %v, got %v", tt.wantSrcImgPath, srcImgPath)
			}

			if (displayImgPath != "") != tt.wantDisplayImgPath {
				t.Errorf("expected displayImgPath: %v, got %v", tt.wantDisplayImgPath, displayImgPath)
			}

		})
	}

	_ = os.RemoveAll(basePath)
}

type mockFileHeader struct {
	Filename string
	Size     int64
	Content  []byte
}

func (m *mockFileHeader) Open() (multipart.File, error) {
	return &mockFile{bytes.NewReader(m.Content)}, nil
}

func (m *mockFileHeader) GetFilename() string {
	return m.Filename
}

func (m *mockFileHeader) GetSize() int64 {
	return m.Size
}

type mockFile struct {
	*bytes.Reader
}

func (m *mockFile) Close() error {
	return nil
}
