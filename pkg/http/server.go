package http

import (
	"embed"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/vsychov/go-rating-stars/pkg/html"
	"github.com/vsychov/go-rating-stars/pkg/voter"
	"net/http"
)

// ServeHttp start gin
func ServeHttp(distAssets embed.FS, voter voter.Voter, drawer html.Drawer) (err error) {
	controller := controller{
		Voter:    voter,
		Drawer:   drawer,
		Validate: validator.New(),
	}

	r := gin.Default()
	r.StaticFS("/assets", http.FS(distAssets))
	r.GET("/:resource_id/results/", controller.results)
	r.POST("/:resource_id/vote/", controller.vote)

	err = r.Run()
	if err != nil {
		return
	}

	return
}
