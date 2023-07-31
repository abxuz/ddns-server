package middleware

import (
	"path"

	"github.com/gin-gonic/gin"
	"github.com/xbugio/b-tools/bset"
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
