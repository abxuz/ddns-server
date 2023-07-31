package service

import (
	"ddns-server/internal/api"
	"ddns-server/internal/config"
	"ddns-server/internal/dao"
	"ddns-server/internal/logic"
	"ddns-server/internal/middleware"
	"ddns-server/resource"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/xbugio/b-tools/bhttp"
)

var Web = sWeb{
	lock: &sync.RWMutex{},
}

type sWeb struct {
	lock   *sync.RWMutex
	config *config.Web
	server *http.Server
}

func init() {
	logic.Web.RegisterService(&Web)
}

func (s *sWeb) ReloadService() error {
	s.lock.Lock()
	defer s.lock.Unlock()

	config, err := dao.Web.Get()
	if err != nil {
		return err
	}

	if s.config == nil || s.config.Listen == "" {
		// off -> off
		if config == nil || config.Listen == "" {
			return nil
		}

		// off -> on
		handler, err := s.createHttpHandler(config)
		if err != nil {
			return err
		}
		s.server = &http.Server{
			Addr:     config.Listen,
			Handler:  handler,
			ErrorLog: log.New(io.Discard, "", log.LstdFlags),
		}
		go s.server.ListenAndServe()
		s.config = config
		return nil
	}

	// on -> off
	if config == nil || config.Listen == "" {
		s.server.Close()
		s.server = nil
		s.config = config
		return nil
	}

	// on -> on
	middleware.Auth.SetAuth(config.Auth)
	s.config.Auth = config.Auth

	// port not changed
	if config.Listen == s.config.Listen {
		return nil
	}

	// port changed
	handler, err := s.createHttpHandler(config)
	if err != nil {
		return err
	}
	s.server.Close()
	s.server = &http.Server{
		Addr:     config.Listen,
		Handler:  handler,
		ErrorLog: log.New(io.Discard, "", log.LstdFlags),
	}
	go s.server.ListenAndServe()
	s.config = config
	return nil
}

func (s *sWeb) CloseService() {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.server != nil {
		s.server.Close()
		s.server = nil
		s.config = nil
	}
}

func (s *sWeb) createHttpHandler(config *config.Web) (http.Handler, error) {
	gin.SetMode(gin.ReleaseMode)
	h := gin.New()

	middleware.Auth.SetAuth(config.Auth)
	h.Use(middleware.Auth.Auth)
	h.Use(middleware.Cache.Cache)

	apiGroup := h.Group("/api/")
	{
		recordApiGroup := apiGroup.Group("/record/")
		{
			recordApiGroup.GET("/", api.Record.List)
			recordApiGroup.GET("/:domain", api.Record.Get)
			recordApiGroup.PUT("/", api.Record.Set)
			recordApiGroup.DELETE("/:domain", api.Record.Delete)
		}

		systemApiGroup := apiGroup.Group("/system/")
		{
			systemApiGroup.GET("/", api.System.Get)
			systemApiGroup.PUT("/", api.System.Set)
		}
	}

	fs, err := resource.HtmlFs()
	if err != nil {
		return nil, err
	}
	filesystem := &bhttp.NoAutoIndexFileSystem{
		FileSystem: http.FS(fs),
	}
	fileserver := http.FileServer(filesystem)
	h.NoRoute(gin.WrapH(fileserver))

	return h, nil
}
