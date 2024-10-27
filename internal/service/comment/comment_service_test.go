package comment

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mfsyahrz/image_feed_api/internal/domain/entity"
	repoMock "github.com/mfsyahrz/image_feed_api/mocks/domain/repository"

	"github.com/stretchr/testify/assert"
)

func TestNewComment(t *testing.T) {

	t.Run("setup", func(t *testing.T) {
		_ = NewCommentService(nil, nil)
	})
}

func TestCommentService_Save(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repoMock.NewMockCommentRepository(ctrl)
	mockPostRepo := repoMock.NewMockPostRepository(ctrl)
	svc := NewCommentService(mockRepo, mockPostRepo)

	tests := []struct {
		name      string
		comment   *entity.Comment
		setupMock func()
		wantErr   bool
	}{
		{
			name: "Success",
			comment: &entity.Comment{
				ID:      1,
				Content: "Test comment",
			},
			setupMock: func() {
				mockPostRepo.EXPECT().FetchOne(gomock.Any(), gomock.Any()).Return(&entity.Post{ID: 1}, nil)
				mockRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Error invalid postID",
			comment: &entity.Comment{
				ID:      1,
				Content: "Test comment",
			},
			setupMock: func() {
				mockPostRepo.EXPECT().FetchOne(gomock.Any(), gomock.Any()).Return(nil, errors.New("invalid post"))
			},
			wantErr: true,
		},
		{
			name: "Error saving comment",
			comment: &entity.Comment{
				ID:      1,
				Content: "Test comment",
			},
			setupMock: func() {
				mockPostRepo.EXPECT().FetchOne(gomock.Any(), gomock.Any()).Return(&entity.Post{ID: 1}, nil)
				mockRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			err := svc.Save(context.Background(), tt.comment)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCommentService_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repoMock.NewMockCommentRepository(ctrl)
	mockPostRepo := repoMock.NewMockPostRepository(ctrl)
	svc := NewCommentService(mockRepo, mockPostRepo)

	tests := []struct {
		name      string
		comment   *entity.Comment
		setupMock func()
		wantErr   bool
	}{
		{
			name: "Success",
			comment:&entity.Comment{
				ID:     1,
				PostID: 1,
			},
			setupMock: func() {
				mockPostRepo.EXPECT().FetchOne(gomock.Any(), gomock.Any()).Return(&entity.Post{ID: 1}, nil)
				mockRepo.EXPECT().Delete(gomock.Any(), &entity.Comment{
					ID:     1,
					PostID: 1,
				}).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Error deleting comment",
			comment: &entity.Comment{
				ID:     1,
				PostID: 1,
			},
			setupMock: func() {
				mockPostRepo.EXPECT().FetchOne(gomock.Any(), gomock.Any()).Return(&entity.Post{ID: 1}, nil)
				mockRepo.EXPECT().Delete(gomock.Any(), &entity.Comment{
				ID:     1,
				PostID: 1,
			}).Return(errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			err := svc.Delete(context.Background(), tt.comment)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
