package controllers

import (
	"github.com/revel/revel"
)

type Newsletter struct {
	*revel.Controller
}

func (c Newsletter) Index() revel.Result {
	return c.Render()
}
