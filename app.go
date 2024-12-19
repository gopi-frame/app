package app

import (
	"fmt"
	"path/filepath"

	"github.com/gopi-frame/collection/kv"

	"github.com/gopi-frame/config"
	"github.com/gopi-frame/container"
	"github.com/gopi-frame/contract/app"
	"github.com/gopi-frame/contract/repository"
	"github.com/gopi-frame/env"
)

type App struct {
	container.Container[any]
	kernel     app.Kernel
	config     repository.Repository
	components *kv.LinkedMap[string, app.Component]

	name         string
	version      string
	debug        bool
	root         string
	wd           string
	storagePath  string
	resourcePath string
	configPath   string
	booted       bool
}

func NewApp(kernel app.Kernel, opts ...Option) (*App, error) {
	debug, err := env.GetBoolOr("APP_DEBUG", false)
	if err != nil {
		return nil, fmt.Errorf("failed to get APP_DEBUG: %w", err)
	}
	app := &App{
		kernel:       kernel,
		config:       config.NewRepository(),
		components:   kv.NewLinkedMap[string, app.Component](),
		debug:        debug,
		root:         env.Get("APP_ROOT"),
		wd:           env.Get("APP_WD"),
		storagePath:  env.GetOr("APP_STORAGE_PATH", filepath.Join(env.Get("APP_ROOT"), "storage")),
		resourcePath: env.GetOr("APP_RESOURCE_PATH", filepath.Join(env.Get("APP_ROOT"), "resource")),
		configPath:   env.GetOr("APP_CONFIG_PATH", filepath.Join(env.Get("APP_ROOT"), "config")),
	}
	for _, opt := range opts {
		if err := opt.Apply(app); err != nil {
			return nil, err
		}
	}
	return app, nil
}

func (app *App) Kernel() app.Kernel {
	return app.kernel
}

func (app *App) Config() repository.Repository {
	return app.config
}

func (app *App) Name() string {
	return app.name
}

func (app *App) Version() string {
	return app.version
}

func (app *App) Debug() bool {
	return app.debug
}

func (app *App) Root() string {
	return app.root
}

func (app *App) WorkingDirectory() string {
	return app.wd
}

func (app *App) StoragePath() string {
	return app.storagePath
}

func (app *App) ResourcePath() string {
	return app.resourcePath
}

func (app *App) ConfigPath() string {
	return app.configPath
}

func (app *App) Register(component app.Component) error {
	app.components.Lock()
	defer app.components.Unlock()
	if app.components.ContainsKey(component.Name()) {
		return fmt.Errorf("component %s already registered", component.Name())
	}
	app.components.Set(component.Name(), component)
	if err := component.Register(app); err != nil {
		return err
	}
	if app.booted {
		if err := component.Boot(); err != nil {
			return err
		}
	}
	return nil
}

func (app *App) MustRegister(component app.Component) {
	if err := app.Register(component); err != nil {
		panic(err)
	}
}

func (app *App) Unregister(component app.Component) error {
	app.components.Lock()
	defer app.components.Unlock()
	if !app.components.ContainsKey(component.Name()) {
		return nil
	}
	app.components.Remove(component.Name())
	if err := component.Unregister(app); err != nil {
		return err
	}
	if app.booted {
		if err := component.Shutdown(); err != nil {
			return err
		}
	}
	return nil
}

func (app *App) Run() error {
	app.components.Lock()
	for _, component := range app.components.Values() {
		if err := component.Boot(); err != nil {
			return err
		}
	}
	app.booted = true
	app.components.Unlock()
	return app.kernel.Run()
}
