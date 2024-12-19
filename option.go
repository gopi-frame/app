package app

import (
	"github.com/gopi-frame/contract"
	"github.com/gopi-frame/contract/repository"
)

type Option = contract.Option[*App]

type OptionFunc func(app *App) error

func (f OptionFunc) Apply(app *App) error {
	return f(app)
}

var noneOption = OptionFunc(func(app *App) error {
	return nil
})

func WithName(name string) OptionFunc {
	return func(app *App) error {
		app.name = name
		return nil
	}
}

func WithVersion(version string) OptionFunc {
	return func(app *App) error {
		app.version = version
		return nil
	}
}

func WithDebug(debug bool) OptionFunc {
	return func(app *App) error {
		app.debug = debug
		return nil
	}
}

func WithStoragePath(path string) OptionFunc {
	return func(app *App) error {
		app.storagePath = path
		return nil
	}
}

func WithResourcePath(path string) OptionFunc {
	return func(app *App) error {
		app.resourcePath = path
		return nil
	}
}

func WithConfigPath(path string) OptionFunc {
	return func(app *App) error {
		app.configPath = path
		return nil
	}
}

func WithConfigParser(parser repository.Parser) OptionFunc {
	return func(app *App) error {
		app.configParser = parser
		return nil
	}
}
