package main

import (
	"embed"
	"github.com/vsychov/go-rating-stars/pkg/config"
	"github.com/vsychov/go-rating-stars/pkg/html"
	"github.com/vsychov/go-rating-stars/pkg/http"
	"github.com/vsychov/go-rating-stars/pkg/storage"
	"github.com/vsychov/go-rating-stars/pkg/voter"
	"go.uber.org/fx"
	"os"
)

//go:embed assets/dist
var distAssets embed.FS

func main() {
	var target interface{}
	if len(os.Args) > 1 && os.Args[1] == "--cron" {
		target = storage.Cron
	} else {
		target = http.ServeHttp
	}

	fx.New(
		fx.Provide(config.CreateFromEnv),
		fx.Provide(storage.CreatePostgres),
		fx.Provide(html.Create),
		fx.Provide(voter.Create),
		fx.Provide(assetsFSProvider),
		fx.Invoke(target),
	).Run()
}

func assetsFSProvider() embed.FS {
	return distAssets
}
