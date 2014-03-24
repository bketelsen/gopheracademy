package controllers

import (
	"github.com/revel/revel"
)

type Jobs struct {
	*revel.Controller
}

func (c Jobs) Index() revel.Result {
	return c.Render()
}

func (c Jobs) Find(size, page int) revel.Result {
	return c.Render()
}

func (c Jobs) Show(id int) revel.Result {
	return c.Render()
}
