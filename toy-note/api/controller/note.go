package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"toy-note/api/entity"
	"toy-note/api/service"
	"toy-note/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ToyNoteController
// Work for `Gin.Router`
//
// filed service accepts a `ToyNoteRepo` interface
type ToyNoteController struct {
	logger  *zap.SugaredLogger
	service service.ToyNoteRepo
}

func NewToyNoteController(
	logger *logger.ToyNoteLogger,
	service service.ToyNoteRepo,
) *ToyNoteController {
	return &ToyNoteController{
		logger:  logger.NewSugar("NoteController"),
		service: service,
	}
}

// ============================================================================
// Tag
// ============================================================================

// @Summary      get all tags
// @Description  get all tags
// @Tags         tag
// @Produce      json
// @Success      200  {array}  entity.Tag
// @Router       /get-tags [get]
func (c *ToyNoteController) GetTags(ctx *gin.Context) {
	tags, err := c.service.GetTags()
	if err != nil {
		c.logger.Error(err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, tags)
}

// @Summary      create/update a tag
// @Description  create a new tag or update an existing tag, based on whether the tag ID is provided
// @Tags         tag
// @Produce      json
// @Param        data  body      entity.Tag  true  "tag data"
// @Success      200   {object}  entity.Tag
// @Failure      400    {object}  errorMessage
// @Router       /save-tag [post]
func (c *ToyNoteController) SaveTag(ctx *gin.Context) {
	var tag entity.Tag
	if err := ctx.ShouldBindJSON(&tag); err != nil {
		c.logger.Error(err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	tag, err := c.service.SaveTag(tag)
	if err != nil {
		c.logger.Error(err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, tag)
}

// @Summary      delete a tag by ID
// @Description  delete a tag by ID
// @Tags         tag
// @Produce      json
// @Param        id   path      int  true  "tag ID"
// @Success      200  {object}  successMessage
// @Failure      500  {object}  errorMessage
// @Router       /delete-tag/{id} [delete]
func (c *ToyNoteController) DeleteTag(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if err := c.service.DeleteTag(uint(id)); err != nil {
		c.logger.Error(err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, successResponse(id))
}

// ============================================================================
// Post
// ============================================================================

// @Summary      get all posts
// @Description  get all posts
// @Tags         post
// @Param        page   query  int     true  "page number"
// @Param        size   query  int     true  "page size"
// @Produce      json
// @Success      200  {array}  entity.Post
// @Router       /get-posts [get]
func (c *ToyNoteController) GetPosts(ctx *gin.Context) {
	pagination, err := getPaginationFromQuery(ctx)
	if err != nil {
		c.logger.Error(err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// get post from service
	posts, err := c.service.GetPosts(pagination)
	if err != nil {
		c.logger.Error(err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, posts)
}

// @Summary      create/update a post
// @Description  Save post can be used to create a new post or update an existing post.
// @Description  If id is not provided, it will create a new post; Otherwise, it will update an existing post.
// @Description  Form-data should also be provided if the post has any new affiliate.
// @Tags         post
// @Accept       multipart/form-data
// @Produce      json
// @Param        data   formData  string  true   "post data"
// @Param        files  formData  file    false  "affiliate files"
// @Success      200    {object}  entity.Post
// @Failure      400   {object}  errorMessage
// @Router       /save-post [post]
func (c *ToyNoteController) SavePost(ctx *gin.Context) {
	// get multipart form
	form, err := ctx.MultipartForm()
	if err != nil {
		c.logger.Error(err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	files := form.File["files"]
	data := form.Value["data"]
	if len(data) < 1 {
		err := errors.New("filed: data is missing")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// get post from request body
	var post entity.Post

	err = json.Unmarshal([]byte(data[0]), &post)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}

	if err := ctx.ShouldBindJSON(&post); err != nil {
		c.logger.Error(err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// check if new affiliates has the same length as multipart form's length
	newAffiliatesLen := 0
	filesLen := len(files)
	for _, a := range post.Affiliates {
		if a.Id == 0 {
			newAffiliatesLen++
		}
	}
	if filesLen != newAffiliatesLen {
		c.logger.Error(fmt.Sprintf(
			"new affiliates length %d not match files length %d",
			newAffiliatesLen,
			filesLen,
		))
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// oids used for storing object ids returned from mongo service,
	// later we need it to be recorded in post's affiliate.
	var oids []string
	for _, file := range files {
		filename := file.Filename

		// open file
		file, err := file.Open()
		if err != nil {
			c.logger.Error(err)
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		defer file.Close()

		// upload file to MongoDB and get returned ObjectId
		oid, err := c.service.UploadAffiliate(file, filename)
		if err != nil {
			c.logger.Error(err)
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		oids = append(oids, oid)
	}

	// rebind ObjectIds to post's new affiliates
	oidsIdx := 0
	for idx, a := range post.Affiliates {
		if a.Id == 0 {
			post.Affiliates[idx].ObjectId = oids[oidsIdx]
			oidsIdx++
		}
	}

	// save post to PG
	post, err = c.service.SavePost(post)
	if err != nil {
		c.logger.Error(err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, post)
}

// @Summary      delete a post by ID
// @Description  delete a post by ID
// @Tags         post
// @Produce      json
// @Param        id   path      string  true  "post ID"
// @Success      200  {object}  successMessage
// @Failure      500  {object}  errorMessage
// @Router       /delete-post/{id} [delete]
func (c *ToyNoteController) DeletePost(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := c.service.DeletePost(uint(id)); err != nil {
		c.logger.Error(err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, successResponse(id))
}

// @Summary      download an affiliate by ID
// @Description  download an affiliate by ID
// @Tags         affiliate
// @Produce      json
// @Param        id   path      int  true  "affiliate ID"
// @Success      200  {object}  downloadSuccess
// @Failure      500  {object}  errorMessage
// @Router       /download-file [get]
func (c *ToyNoteController) DownloadAffiliate(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	fo, err := c.service.DownloadAffiliate(uint(id))
	if err != nil {
		c.logger.Error(err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.Header("Content-Disposition", "attachment; filename="+fo.Filename)
	ctx.Data(http.StatusOK, "application/octet-stream", fo.Content)
	ctx.JSON(http.StatusOK, downloadSuccess{
		Filename: fo.Filename,
		Size:     fo.Size,
	})
}

// @Summary      get posts by tags
// @Description  get posts by tags
// @Tags         post
// @Param        page  query  int     true  "page number"
// @Param        size  query  int     true  "page size"
// @Param        ids   query  string  true  "tag ids"
// @Produce      json
// @Success      200  {array}  entity.Post
// @Router       /search-posts-by-tags [get]
func (c *ToyNoteController) SearchPostsByTags(ctx *gin.Context) {
	pagination, err := getPaginationFromQuery(ctx)
	if err != nil {
		c.logger.Error(err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	idsQuery := ctx.Query("ids")
	ids := strings.Split(idsQuery, ",")
	var idsUint []uint
	for _, id := range ids {
		idUint, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			c.logger.Error(err)
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		idsUint = append(idsUint, uint(idUint))
	}

	posts, err := c.service.SearchPostsByTags(idsUint, pagination)
	if err != nil {
		c.logger.Error(err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, posts)
}

// @Summary      get posts by title
// @Description  get posts by title
// @Tags         post
// @Param        page  query  int  true  "page number"
// @Param        size  query  int  true  "page size"
// @Param        title  query  string  true  "post title"
// @Produce      json
// @Success      200  {array}  entity.Post
// @Router       /search-posts-by-title [get]
func (c *ToyNoteController) SearchPostsByTitle(ctx *gin.Context) {
	pagination, err := getPaginationFromQuery(ctx)
	if err != nil {
		c.logger.Error(err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	titleQuery := ctx.Query("title")

	posts, err := c.service.SearchPostsByTitle(titleQuery, pagination)
	if err != nil {
		c.logger.Error(err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, posts)
}
