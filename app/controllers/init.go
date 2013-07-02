package controllers

import (
	"github.com/robfig/revel"

)

func init() {
	revel.OnAppStart(InitDB)
}
