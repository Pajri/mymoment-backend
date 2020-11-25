package delivery

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/microcosm-cc/bluemonday"
	"github.com/pajri/personal-backend/adapter/cerror"
	"github.com/pajri/personal-backend/domain"
	"github.com/pajri/personal-backend/global"
)

/* #region type helper */
type InsertPostResponse struct {
	Message []string           `json:"message"`
	Post    PostListingElement `json:"post,omitempty"`
}

type InsertPostRequest struct {
	Content  string `form:"content" binding:"required"`
	ImageURL string `form:"image_url"`
}

type DeletePostRequest struct {
	PostID string `json:"post_id" binding:"required"`
}

type DeletePostResponse struct {
	Message []string `json:"message"`
}

type PostListingRequest struct {
	Date  time.Time `form:"date"`
	Limit uint64    `form:"limit"`
}

type PostListingResponse struct {
	Message  string               `json:"message"`
	PostList []PostListingElement `json:"post_list"`
}

type PostListingElement struct {
	PostID     string `json:"post_id"`
	Content    string `json:"content"`
	ImageURL   string `json:"image_url"`
	Date       string `json:"date"`
	HiddenDate string `json:"hidden_date"`
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
	router.GET("/api/post", handler.PostListing)
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
				jsonField, _ := field.Tag.Lookup("form")

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

	//sanitize input
	p := bluemonday.UGCPolicy()
	request.Content = p.Sanitize(request.Content)

	var post domain.Post
	post.Content = request.Content
	post.ImageURL = request.ImageURL
	post.AccountID = accountID

	var storedPost *domain.Post
	storedPost, err = ph.useCase.InsertPost(post)

	//the response will be used as first element of listing
	//so the post response uses PostListingElement type
	var postResponse PostListingElement
	postResponse = ph.creatPostListingElement(*storedPost)

	if err != nil {
		response.Message = []string{err.(cerror.Error).FriendlyMessageWithTag()}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response = InsertPostResponse{nil, postResponse}
	c.JSON(http.StatusCreated, response)
	return
}

func (ph PostHandler) PostListing(c *gin.Context) {
	var (
		request   PostListingRequest
		response  PostListingResponse
		accountID string = c.GetString("account_id")
	)

	//validate account id
	if accountID == "" {
		fmt.Println("account id is empty")
		cerr := cerror.NewAndPrintWithTag("PLD00", errors.New("Account id is empty"), global.FRIENDLY_MESSAGE)
		response.Message = cerr.FriendlyMessageWithTag()
		c.JSON(http.StatusBadRequest, response)
		return
	}

	//get request param
	err := c.ShouldBind(&request)
	if err != nil {
		cerr := cerror.NewAndPrintWithTag("PLD01", err, global.FRIENDLY_INVALID_PARAM)
		response.Message = cerr.FriendlyMessageWithTag()
		c.JSON(http.StatusBadRequest, response)
		return
	}

	postList, err := ph.useCase.PostListing(accountID, request.Limit, request.Date)
	if err != nil {
		response.Message = err.(cerror.Error).FriendlyMessageWithTag()
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	var postListElements []PostListingElement
	for _, post := range postList {
		var new PostListingElement
		new = ph.creatPostListingElement(post)

		postListElements = append(postListElements, new)
	}
	response.PostList = postListElements

	c.JSON(http.StatusOK, response)
	return
}

func (ph PostHandler) DeletePost(c *gin.Context) {
	var (
		request   DeletePostRequest
		response  DeletePostResponse
		accountID string = c.GetString("account_id")
	)

	err := c.ShouldBind(&request)
	if err != nil {
		cerr := cerror.NewAndPrintWithTag("DPH00", err, global.FRIENDLY_MESSAGE)

		valError := err.(validator.ValidationErrors)
		if valError != nil {
			for _, elem := range valError {
				fieldName := elem.Field()
				field, _ := reflect.TypeOf(&request).Elem().FieldByName(fieldName)
				jsonField, _ := field.Tag.Lookup("form")

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
	}

	err = ph.useCase.DeletePost(request.PostID, accountID)
	if err != nil {
		response.Message = []string{err.(cerror.Error).FriendlyMessageWithTag()}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (ph PostHandler) creatPostListingElement(post domain.Post) PostListingElement {
	var postListingElement PostListingElement
	postListingElement.PostID = post.PostID
	postListingElement.Content = post.Content
	postListingElement.ImageURL = post.ImageURL
	postListingElement.Date = post.Date.Format(global.TIME_FORMAT)
	postListingElement.HiddenDate = post.Date.Format(global.TIME_ISO8601)

	return postListingElement
}
