package fin

import "time"

type Option func(engine *Engine)

func HandleNotFound(f HandlerFunc) Option {
	return func(engine *Engine) {
		engine.HandleNotFound = f
	}
}

func ReadTimeout(d time.Duration) Option {
	return func(engine *Engine) {
		engine.server.ReadTimeout = d
	}
}

func WriteTimeout(d time.Duration) Option {
	return func(engine *Engine) {
		engine.server.WriteTimeout = d
	}
}

func Name(name string) Option {
	return func(engine *Engine) {
		engine.server.Name = name
	}
}
