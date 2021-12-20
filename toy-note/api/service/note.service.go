package service

import (
	"context"
	"time"
	"toy-note/api/entity"
	"toy-note/api/persistence"
	"toy-note/logger"

	"go.uber.org/zap"
)

type ToyNoteService struct {
	logger *zap.SugaredLogger
	pg     *persistence.PgRepository
	mongo  *persistence.MongoRepository
}

func NewToyNoteService(
	logger *logger.ToyNoteLogger,
	pgConn persistence.PgConn,
	mongoConn persistence.MongoConn,
) (*ToyNoteService, error) {
	slog := logger.NewSugar("ToyNoteService")

	pg, err := persistence.NewPgRepository(logger, pgConn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mongo, err := persistence.NewMongoRepository(ctx, logger, mongoConn)
	if err != nil {
		return nil, err
	}

	return &ToyNoteService{
		logger: slog,
		pg:     &pg,
		mongo:  &mongo,
	}, nil
}

func (s *ToyNoteService) Init() error {
	s.logger.Debug("Initializing ToyNoteService")

	// make sure tables are updated to the latest schema
	if err := s.pg.AutoMigrate(); err != nil {
		return err
	}

	return nil
}

var _ ToyNoteRepo = &ToyNoteService{}

func (s *ToyNoteService) GetTags() ([]entity.Tag, error) {
	return s.pg.GetTags()
}

func (s *ToyNoteService) SaveTag(tag entity.Tag) (entity.Tag, error) {
	if tag.Id == 0 {
		return s.pg.CreateTag(tag)
	} else {
		return s.pg.UpdateTag(tag)
	}
}

func (s *ToyNoteService) DeleteTag(id uint) error {
	return s.pg.DeleteTag(id)
}

func (s *ToyNoteService) GetPosts(pagination entity.Pagination) ([]entity.Post, error) {
	return s.pg.GetPosts(pagination)
}

func (s *ToyNoteService) SavePost(post entity.Post) (entity.Post, error) {
	if post.Id == 0 {
		return s.pg.CreatePost(post)
	} else {
		return s.pg.UpdatePost(post)
	}
}

func (s *ToyNoteService) DeletePost(id uint) error {
	return s.pg.DeletePost(id)
}

func (s *ToyNoteService) DownloadAffiliate(id uint) ([]byte, error) {
	affiliate, err := s.pg.GetAffiliate(id)
	if err != nil {
		return nil, err
	}
	return s.mongo.DownloadFile(affiliate.ObjectId)
}
