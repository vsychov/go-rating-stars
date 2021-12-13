package html

import (
	"bytes"
	"embed"
	"github.com/Masterminds/sprig"
	"github.com/vsychov/go-rating-stars/pkg/voter"
	"html/template"
)

// Drawer is used for draw html
type Drawer struct {
	distAssets embed.FS
}

// Create new object instance
func Create(distAssets embed.FS) Drawer {
	return Drawer{
		distAssets: distAssets,
	}
}

type renderData struct {
	Rating     float64
	TotalVotes int
	ReadOnly   bool
	ResourceId string
}

// RenderHtml generate html
func (drawer *Drawer) RenderHtml(resourceId string, results voter.VoteResults) (html []byte, err error) {
	tmpl := template.New("index.html").Funcs(sprig.FuncMap())
	t, err := tmpl.ParseFS(drawer.distAssets, "assets/dist/*")

	if err != nil {
		return
	}

	data := renderData{
		TotalVotes: results.TotalVotes,
		Rating:     results.Rating,
		ReadOnly:   !results.AllowedToVote,
		ResourceId: resourceId,
	}

	buf := new(bytes.Buffer)
	err = t.Execute(buf, data)

	if err != nil {
		return
	}

	html = buf.Bytes()

	return
}
