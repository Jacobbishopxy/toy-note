package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"toy-note/api/entity"
	"toy-note/api/service"
	"toy-note/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// NoteController
// Work for `Gin.Router`
//
// filed service accepts a `ToyNoteRepo` interface
type NoteController struct {
	logger  *zap.SugaredLogger
	service service.ToyNoteRepo
}

func NewNoteController(
	logger *logger.ToyNoteLogger,
	service service.ToyNoteRepo,
) *NoteController {
	return &NoteController{
		logger:  logger.NewSugar("NoteController"),
		service: service,
	}
}

// ============================================================================
// Tag
// ============================================================================

func (c *NoteController) GetTags(ctx *gin.Context) {
	tags, err := c.service.GetTags()
	if err != nil {
		c.logger.Error(err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, tags)
}

func (c *NoteController) SaveTag(ctx *gin.Context) {
	var tag entity.Tag
	if err := ctx.ShouldBindJSON(&tag); err != nil {
		c.logger.Error(err)
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	tag, err := c.service.SaveTag(tag)
	if err != nil {
		c.logger.Error(err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, tag)
}

func (c *NoteController) DeleteTag(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err := c.service.DeleteTag(uint(id)); err != nil {
		c.logger.Error(err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.Status(http.StatusOK)
}

// ============================================================================
// Post
// ============================================================================

func (c *NoteController) GetPosts(ctx *gin.Context) {

	// get pagination's page from query string
	pageQuery := ctx.Query("page")
	page, err := strconv.ParseInt(pageQuery, 10, 64)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// get pagination's size from query string
	sizeQuery := ctx.Query("size")
	size, err := strconv.ParseInt(sizeQuery, 10, 64)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	pagination := entity.Pagination{
		Page: int(page),
		Size: int(size),
	}

	// get post from service
	posts, err := c.service.GetPosts(pagination)
	if err != nil {
		c.logger.Error(err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, posts)
}

/*
SavePost

The core function of this controller.
Save post can be used to create a new post or update an existing post.
*/
func (c *NoteController) SavePost(ctx *gin.Context) {
	// get multipart form
	form, err := ctx.MultipartForm()
	if err != nil {
		c.logger.Error(err)
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	files := form.File["files"]

	// get post from request body
	var post entity.Post
	if err := ctx.ShouldBindJSON(&post); err != nil {
		c.logger.Error(err)
		ctx.AbortWithError(http.StatusBadRequest, err)
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
		ctx.AbortWithError(http.StatusBadRequest, err)
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
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}
		defer file.Close()

		// upload file to MongoDB and get returned ObjectId
		oid, err := c.service.UploadAffiliate(file, filename)
		if err != nil {
			c.logger.Error(err)
			ctx.AbortWithError(http.StatusInternalServerError, err)
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
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, post)
}

func (c *NoteController) DeletePost(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if err := c.service.DeletePost(uint(id)); err != nil {
		c.logger.Error(err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.Status(http.StatusOK)
}

func (c *NoteController) DownloadAffiliate(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	fo, err := c.service.DownloadAffiliate(uint(id))
	if err != nil {
		c.logger.Error(err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.Header("Content-Disposition", "attachment; filename="+fo.Filename)
	ctx.Data(http.StatusOK, "application/octet-stream", fo.Content)
	ctx.JSON(http.StatusOK, gin.H{
		"message":  "success",
		"filename": fo.Filename,
		"size":     fo.Size,
	})
}
