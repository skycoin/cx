package controllers

import (
	"github.com/revel/revel"
	. "github.com/skycoin/cx/src/base"
	. "github.com/skycoin/cx/src/cxgo"
)

type App struct {
	*revel.Controller
}

func (c App) Tutorial () revel.Result {
	return c.Render()
}

func (c App) Examples () revel.Result {
	return c.Render()
}

func (c App) Index (code string) revel.Result {
	var evalResult string
	if code != "" {
		evalResult = Eval(code)
	}
	return c.Render(code, evalResult)
}
