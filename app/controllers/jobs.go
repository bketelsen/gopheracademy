package controllers

import "github.com/robfig/revel"

type Jobs struct {
	*revel.Controller
}

func (c Jobs) Index() revel.Result {
	return c.Render()
}

func (c Jobs) Find() revel.Result {
	return c.Render()
}

func (c Jobs) Post() revel.Result {
	return c.Render()
}
