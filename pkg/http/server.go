package http

import (
	"embed"
	healthcheck "github.com/RaMin0/gin-health-check"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/vsychov/go-rating-stars/pkg/config"
	"github.com/vsychov/go-rating-stars/pkg/html"
	"github.com/vsychov/go-rating-stars/pkg/voter"
	"net/http"
	"strings"
)

// ServeHttp start gin
func ServeHttp(distAssets embed.FS, voter voter.Voter, drawer html.Drawer, config config.Config) (err error) {
	controller := controller{
		Voter:    voter,
		Drawer:   drawer,
		Validate: validator.New(),
	}

	r := gin.Default()
	r.TrustedPlatform = config.ClientIpHeader
	r.Use(healthcheck.Default())
	r.Use(headersByRequestURI())

	r.StaticFS("/assets", http.FS(distAssets))

	r.GET("/:resource_id/results/", controller.results)
	r.POST("/:resource_id/vote/", controller.vote)

	err = r.Run()
	if err != nil {
		return
	}

	return
}

func headersByRequestURI() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.HasPrefix(c.Request.RequestURI, "/assets/") {
			c.Header("Cache-Control", "public, immutable, max-age=31536000")
		}
	}
}
