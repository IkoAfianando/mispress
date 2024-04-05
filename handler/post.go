package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IkoAfianando/mispress/models"
	"github.com/gin-gonic/gin"
	"io"

	"net/http"
	"strings"
)

func (h *Handler) CreatePost(c *gin.Context) {
	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		h.Logger.Err(err).Msg("could not parse request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid request body: %s", err.Error())})
		return
	}
	err := h.DB.SavePost(&post)
	if err != nil {
		h.Logger.Err(err).Msg("could not save post")
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("could not save post: %s", err.Error())})
	} else {
		c.JSON(http.StatusCreated, gin.H{"post": post})
	}
}

func (h *Handler) GetPost(c *gin.Context) {
	postID := c.Param("post_id")
	post, err := h.DB.GetPost(postID)
	if err != nil {
		h.Logger.Err(err).Msg("could not fetch post")
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("could not fetch post: %s", err.Error())})
		return
	}
	c.JSON(http.StatusOK, gin.H{"post": post})
}

func (h *Handler) UpdatePost(c *gin.Context) {
	postID := c.Param("post_id")
	var updatedPost models.Post
	if err := c.ShouldBindJSON(&updatedPost); err != nil {
		h.Logger.Err(err).Msg("could not parse request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid request body: %s", err.Error())})
		return
	}
	err := h.DB.UpdatePost(postID, updatedPost)
	if err != nil {
		h.Logger.Err(err).Msg("could not update post")
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("could not update post: %s", err.Error())})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "post updated successfully"})
	}
}

func (h *Handler) DeletePost(c *gin.Context) {
	postID := c.Param("id")
	err := h.DB.DeletePost(postID)
	if err != nil {
		h.Logger.Err(err).Msg("could not delete post")
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("could not delete post: %s", err.Error())})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "post deleted successfully"})
	}
}

func (h *Handler) GetPosts(c *gin.Context) {
	posts, err := h.DB.GetPosts()
	if err != nil {
		h.Logger.Err(err).Msg("could not fetch posts")
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("could not fetch posts: %s", err.Error())})
		return
	}
	c.JSON(http.StatusOK, gin.H{"posts": posts})
}

func (h *Handler) SearchPosts(c *gin.Context) {
	var query string
	if query, _ = c.GetQuery("q"); query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no search query present"})
		return
	}

	body := fmt.Sprintf(
		`{"query": {"multi_match": {"query": "%s", "fields": ["title", "body"]}}}`,
		query)
	fmt.Println(body, "ini termasuk body")
	res, err := h.ESClient.Search(
		h.ESClient.Search.WithContext(context.Background()),
		h.ESClient.Search.WithIndex("posts"),
		h.ESClient.Search.WithBody(strings.NewReader(body)),
		h.ESClient.Search.WithPretty(),
	)
	fmt.Println(res, "ini res")
	if err != nil {
		h.Logger.Err(err).Msg("elasticsearch error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			h.Logger.Err(err).Msg("could not close the response body")
		}
	}(res.Body)
	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			h.Logger.Err(err).Msg("error parsing the response body")
		} else {
			h.Logger.Err(fmt.Errorf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)).Msg("failed to search query")
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": e["error"].(map[string]interface{})["reason"]})
		return
	}

	h.Logger.Info().Interface("res", res.Status())

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		h.Logger.Err(err).Msg("elasticsearch error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": r["hits"]})
}
