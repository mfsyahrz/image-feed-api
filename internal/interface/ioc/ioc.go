package ioc

import (
	"github.com/mfsyahrz/image_feed_api/internal/config"
	"github.com/mfsyahrz/image_feed_api/internal/db"
	"github.com/mfsyahrz/image_feed_api/internal/infrastructure/filestore"
	postgresComment "github.com/mfsyahrz/image_feed_api/internal/infrastructure/postgres/comment"
	postgresPost "github.com/mfsyahrz/image_feed_api/internal/infrastructure/postgres/post"
	commentService "github.com/mfsyahrz/image_feed_api/internal/service/comment"
	postService "github.com/mfsyahrz/image_feed_api/internal/service/post"
)

type IOC struct {
	Config         *config.Config
	PostService    postService.PostService
	CommentService commentService.CommentService
}

func Setup() *IOC {
	cfg, err := config.New(".env")
	if err != nil {
		panic("failed to setup config. " + err.Error())
	}

	postgresDB, err := db.NewPostgres(cfg.Postgres)
	if err != nil {
		panic("failed to setup db. " + err.Error())
	}

	// construct infrastructure
	postRepo := postgresPost.NewPostRepo(postgresDB.Conn)
	commentRepo := postgresComment.NewCommentRepo(postgresDB.Conn)
	filestore, err := filestore.NewFileStore(&cfg.FileStorage)
	if err != nil {
		panic("failed to setup fileStore. " + err.Error())
	}

	// construct services
	postSvc := postService.NewPostService(postRepo, commentRepo, filestore)
	commentSvc := commentService.NewCommentService(commentRepo, postRepo)

	return &IOC{
		Config:         cfg,
		PostService:    postSvc,
		CommentService: commentSvc,
	}
}
