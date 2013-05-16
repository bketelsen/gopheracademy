package controllers

import "github.com/robfig/revel"
import "gopheracademy/app/models"

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

func (c Jobs) Preview() revel.Result {
	return c.Render()
}
func (c Jobs) Payment() revel.Result {
	return c.Render()
}
func (c Jobs) Confirm() revel.Result {
	return c.Render()
}

func (c Jobs) HandlePostSubmit(job *models.Job) revel.Result {
	job.Validate(c.Validation)

	// Handle errors
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(Jobs.Post)
	}

	// Ok, display the created user
	return c.Render(job)
}
