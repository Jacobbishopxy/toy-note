package service

import (
	"context"
	"io"
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
		logger: logger.NewSugar("ToyNoteService"),
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

var _ ToyNoteRepo = (*ToyNoteService)(nil)

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

func (s *ToyNoteService) UploadAffiliate(reader io.Reader, filename string) (string, error) {
	return s.mongo.UploadFile(reader, filename)
}

func (s *ToyNoteService) DownloadAffiliate(id uint) (entity.FileObject, error) {
	affiliate, err := s.pg.GetAffiliate(id)
	if err != nil {
		return entity.FileObject{}, err
	}
	return s.mongo.DownloadFile(affiliate.Filename, affiliate.ObjectId)
}

func (s *ToyNoteService) GetUnownedAffiliates(pagination entity.Pagination) ([]entity.Affiliate, error) {
	affiliates, err := s.pg.GetUnownedAffiliates(pagination)
	if err != nil {
		return nil, err
	}

	return affiliates, nil
}

func (s *ToyNoteService) RebindAffiliate(postId, affiliateId uint) error {
	// get affiliate, if not found, return error
	affiliate, err := s.pg.GetAffiliate(affiliateId)
	if err != nil {
		return err
	}

	// get post, if not found, return error
	post, err := s.pg.GetPost(postId)
	if err != nil {
		return err
	}

	// update post
	post.Affiliates = append(
		post.Affiliates,
		affiliate,
	)
	post, err = s.pg.UpdatePost(post)
	if err != nil {
		return err
	}

	return nil
}

func (s *ToyNoteService) DeleteUnownedAffiliates(ids []uint) error {
	// get all unowned affiliates from PG
	oa, err := s.pg.GetUnownedAffiliatesByIds(ids)
	if err != nil {
		return err
	}

	// extract object ids from unowned affiliates
	oids := make([]string, len(oa))
	for a := range oa {
		oids[a] = oa[a].ObjectId
	}

	// delete all unowned affiliates from Mongo
	err = s.mongo.DeleteFiles(oids)
	if err != nil {
		return err
	}

	// delete all unowned affiliates from PG
	err = s.pg.DeleteUnownedAffiliates(ids)
	if err != nil {
		return err
	}

	return nil
}

func (s *ToyNoteService) SearchPostsByTags(tagIds []uint, pagination entity.Pagination) ([]entity.Post, error) {
	return s.pg.GetPostsByTags(tagIds, pagination)
}

func (s *ToyNoteService) SearchPostsByTitle(title string, pagination entity.Pagination) ([]entity.Post, error) {
	return s.pg.GetPostsByTitle(title, pagination)
}
