package delivery

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pajri/personal-backend/adapter/cerror"
	"github.com/pajri/personal-backend/domain"
	"github.com/pajri/personal-backend/global"
)

/* #region type helper */
type InsertPostResponse struct {
	Message string
	Post    domain.Post
}

type DeletePostResponse struct {
	Message string
}

/* #endregion */

type PostHandler struct {
	useCase domain.IPostUsecase
}

func NewPostHandler(router *gin.Engine, postUsecase domain.IPostUsecase) {
	handler := &PostHandler{
		useCase: postUsecase,
	}

	router.POST("/api/post", handler.InsertPost)
	router.POST("/api/post/delete", handler.DeletePost)

}

func (ph PostHandler) InsertPost(c *gin.Context) {
	var post domain.Post
	err := c.ShouldBind(&post)
	if err != nil {
		cerr := cerror.NewAndPrintWithTag("IPH00", err, global.FRIENDLY_MESSAGE)
		response := InsertPostResponse{
			Message: cerr.FriendlyMessageWithTag(),
		}

		c.JSON(http.StatusInternalServerError, response)
		return
	}

	//temp
	post.AccountID = "20d4888b-11ae-11eb-a028-ac9e1790b6a2"

	var storedPost *domain.Post
	storedPost, err = ph.useCase.InsertPost(post)
	if err != nil {
		response := InsertPostResponse{
			Message: err.(cerror.Error).FriendlyMessageWithTag(),
		}

		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := InsertPostResponse{"", *storedPost}
	c.JSON(http.StatusCreated, response)
}

func (ph PostHandler) DeletePost(c *gin.Context) {
	var post domain.Post
	err := c.ShouldBind(&post)
	if err != nil {
		cerr := cerror.NewAndPrintWithTag("DPH00", err, global.FRIENDLY_MESSAGE)
		response := DeletePostResponse{
			Message: cerr.FriendlyMessageWithTag(),
		}

		c.JSON(http.StatusInternalServerError, response)
	}

	err = ph.useCase.DeletePost(post.PostID)
	if err != nil {
		response := DeletePostResponse{
			Message: err.(cerror.Error).FriendlyMessageWithTag(),
		}

		c.JSON(http.StatusInternalServerError, response)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
