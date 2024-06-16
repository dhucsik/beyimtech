package http

import (
	imagesS "beyimtech-test/internal/services/images"
	"beyimtech-test/internal/transport/http/handlers/images"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

type Server struct {
	srv      *fasthttp.Server
	basePath string
	started  chan struct{}
	errors   chan error
	router   *router.Router
}

type Controller interface {
	Init(*router.Router)
}

func NewServer(imagesService imagesS.Service) *Server {
	server := &Server{
		srv:      &fasthttp.Server{},
		basePath: "",
		started:  make(chan struct{}, 1),
		errors:   make(chan error, 1),
		router:   router.New(),
	}
	server.WithController(images.New(imagesService))

	return server
}

func (s *Server) Serve() {
	err := s.srv.ListenAndServe("0.0.0.0:8080")
	if err == nil || errors.Is(err, http.ErrServerClosed) {
		s.started <- struct{}{}
		return
	}

	s.errors <- err
}

func (s *Server) WaitsForStarted() error {
	select {
	case err := <-s.errors:
		return err

	case <-s.started:
		return nil
	}
}

func (s *Server) WithController(c Controller) {
	c.Init(s.router)
}

func (s *Server) Start() error {
	if s == nil {
		return nil
	}

	handler := s.router.Handler
	s.srv.Handler = handler

	go s.Serve()
	err := s.WaitsForStarted()
	if err != nil {
		log.Fatalln(err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	if err = s.Close(); err != nil {
		log.Println(err)
	}

	return nil
}

func (s *Server) Close() error {
	if s == nil || s.srv == nil {
		return nil
	}

	return s.srv.Shutdown()
}
