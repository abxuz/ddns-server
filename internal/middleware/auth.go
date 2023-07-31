package middleware

import (
	"ddns-server/internal/config"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

var Auth = mAuth{
	lock: &sync.RWMutex{},
}

type mAuth struct {
	lock *sync.RWMutex
	auth *config.Auth
}

func (m *mAuth) SetAuth(auth *config.Auth) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.auth = auth
}

func (m *mAuth) Auth(ctx *gin.Context) {
	if m.auth == nil || (m.auth.Username == "" && m.auth.Password == "") {
		ctx.Next()
		return
	}

	username, password, ok := ctx.Request.BasicAuth()
	if ok && username == m.auth.Username && password == m.auth.Password {
		ctx.Next()
		return
	}

	ctx.Header("WWW-Authenticate", "Basic realm=Authorization Required")
	ctx.AbortWithStatus(http.StatusUnauthorized)
}
