package app

import (
	"beyimtech-test/configs"
	"beyimtech-test/internal/repositories"
	"beyimtech-test/internal/services/images"
	"beyimtech-test/internal/storage"
	"beyimtech-test/internal/transport/http"
	"beyimtech-test/internal/util/converter"
	"context"
)

type Application struct {
	cfg    *configs.Config
	server *http.Server

	Repository repositories.ImageRepository
	Service    images.Service
	Converter  converter.Converter
	Storage    storage.Storage
}

func NewApplication(_ context.Context) (*Application, error) {
	cfg, err := configs.Parse()
	if err != nil {
		return nil, err
	}

	return &Application{
		cfg: cfg,
	}, nil
}

func InitApp(ctx context.Context) (*Application, error) {
	app, err := NewApplication(ctx)
	if err != nil {
		return nil, err
	}

	for _, init := range []func(context.Context) error{
		app.initRepository,
		app.initClients,
		app.initServices,
		app.initServer,
	} {
		if err := init(ctx); err != nil {
			return nil, err
		}
	}

	return app, nil
}

func (a *Application) initRepository(ctx context.Context) error {
	var err error
	a.Repository, err = repositories.NewRepository(ctx, a.cfg.Env.Get("postgres_dsn"))
	if err != nil {
		return err
	}

	return nil
}

func (a *Application) initClients(_ context.Context) error {
	var err error
	a.Storage, err = storage.NewStorage(a.cfg)
	if err != nil {
		return err
	}

	a.Converter = converter.NewConverter(90)
	return nil
}

func (a *Application) initServices(_ context.Context) error {
	a.Service = images.NewService(
		a.Storage,
		a.Converter,
		a.Repository,
	)

	return nil
}

func (a *Application) initServer(_ context.Context) error {
	a.server = http.NewServer(
		a.Service,
	)

	return nil
}

func (a *Application) Start(ctx context.Context) error {
	go func() {
		a.Service.ReceiveFiles(ctx)
	}()

	return a.server.Start()
}
