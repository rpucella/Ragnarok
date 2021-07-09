package ragnarok

import (
	"rpucella.net/ragnarok/internal/evaluator"
)

// A context contains anything interesting to the execution

type Context struct {
	HomeModule        string
	CurrentModule     string
	NextCurrentModule string // to switch modules, set nextCurrentModule != nil
	Ecosystem         Ecosystem
	CurrentEnv        *evaluator.Env
	Report            func(string)
	Bail              func()
	ReadAll           func(string, *Context) error
}
