package controllers

import (
	"github.com/revel/revel"
)

type GoJobs struct {
	*revel.Controller
}

func (c GoJobs) Index() revel.Result {
	return c.Render()
}

func (c GoJobs) Find(size, page int) revel.Result {
	revel.INFO.Println(len(Joblist))
	return c.Render(Joblist)
}

func (c GoJobs) Show(id int) revel.Result {
	job := Joblist[id]
	return c.Render(job)
}
