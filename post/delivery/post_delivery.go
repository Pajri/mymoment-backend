package delivery

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/pajri/personal-backend/adapter/cerror"
	"github.com/pajri/personal-backend/domain"
	"github.com/pajri/personal-backend/global"
)

/* #region type helper */
type InsertPostResponse struct {
	Message []string    `json:"message"`
	Post    domain.Post `json:"post,omitempty"`
}

type InsertPostRequest struct {
	Content string `json:"content" binding:"required"`
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
	var (
		request   InsertPostRequest
		response  InsertPostResponse
		accountID string = c.GetString("account_id")
	)

	err := c.ShouldBind(&request)
	if err != nil {
		cerr := cerror.NewAndPrintWithTag("IPH00", err, global.FRIENDLY_MESSAGE)

		valError := err.(validator.ValidationErrors)
		if valError != nil {
			for _, elem := range valError {
				fieldName := elem.Field()
				field, _ := reflect.TypeOf(&request).Elem().FieldByName(fieldName)
				jsonField, _ := field.Tag.Lookup("json")

				switch elem.Tag() {
				case "required":
					msg := fmt.Sprintf(global.ERR_REQUIRED_FORMATTER, jsonField)
					response.Message = append(response.Message, msg)
					break
				}
			}

			c.JSON(http.StatusBadRequest, response)
			return
		}

		response.Message = []string{cerr.FriendlyMessageWithTag()}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	var post domain.Post
	post.Content = request.Content
	post.AccountID = accountID

	var storedPost *domain.Post
	storedPost, err = ph.useCase.InsertPost(post)
	if err != nil {
		response.Message = []string{err.(cerror.Error).FriendlyMessageWithTag()}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response = InsertPostResponse{nil, *storedPost}
	c.JSON(http.StatusCreated, response)
	return
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
