package middleware

import (
	"path"

	"github.com/abxuz/b-tools/bset"
	"github.com/gin-gonic/gin"
)

var Cache = mCache{
	cacheExts: bset.NewSetString(".js", ".css", ".ttf", ".woff", ".html"),
}

type mCache struct {
	cacheExts *bset.SetString
}

func (m *mCache) Cache(ctx *gin.Context) {
	ext := path.Ext(ctx.Request.URL.Path)
	if m.cacheExts.Has(ext) {
		ctx.Header("Cache-Control", "max-age=3600")
	}
	ctx.Next()
}
