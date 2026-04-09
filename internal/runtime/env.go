package runtime

import (
	"github.com/mattn/anko/core"
	"github.com/mattn/anko/env"
)

func newEnv() *env.Env {
	e := env.NewEnv()

	core.Import(e)

	return e
}
