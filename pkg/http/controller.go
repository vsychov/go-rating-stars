package http

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/vsychov/go-rating-stars/pkg/html"
	"github.com/vsychov/go-rating-stars/pkg/voter"
	"log"
	"net/http"
)

// VoteRequest is user-input with vote
type VoteRequest struct {
	Vote float64 `json:"vote" binding:"required" validate:"max=5,min=1"`
}

// VoteResult response to user after vote
type VoteResult struct {
	Rating     float64 `json:"rating"`
	TotalVotes int     `json:"totalVotes"`
}

type controller struct {
	Voter    voter.Voter
	Drawer   html.Drawer
	Validate *validator.Validate
}

func (cont controller) results(c *gin.Context) {
	clientIp := c.ClientIP()
	results := cont.Voter.Results(c.Param("resource_id"), clientIp)
	cont.renderHtml(results, c)
}

func (cont controller) vote(c *gin.Context) {
	var voteRequest VoteRequest
	if err := c.ShouldBindJSON(&voteRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := cont.Validate.Struct(voteRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	clientIp := c.ClientIP()
	results := cont.Voter.Vote(c.Param("resource_id"), clientIp, voteRequest.Vote)
	answer := VoteResult{
		Rating:     results.Rating,
		TotalVotes: results.TotalVotes,
	}

	c.JSON(200, answer)
}

func (cont controller) renderHtml(results voter.VoteResults, c *gin.Context) {
	htmlDoc, err := cont.Drawer.RenderHtml(c.Param("resource_id"), results)

	if err != nil {
		log.Println(err)
		htmlDoc = []byte{}
	}

	c.Header("Content-Type", "text/html")
	c.String(200, string(htmlDoc))
}
