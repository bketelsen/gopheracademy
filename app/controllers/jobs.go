package controllers

import "github.com/robfig/revel"
import "gopheracademy/app/models"

type Jobs struct {
	*revel.Controller
}

func (c Jobs) Index() revel.Result {
	return c.Render()
}

func (c Jobs) Find(size, page int) revel.Result {
	if page == 0 {
		page = 1
	}
	if size == 0 {
		size = 5
	}
	nextPage := page + 1

	jobs, err := models.GetJobs(c.Txn, page, size)
	if err != nil {
		revel.ERROR.Println(err)
	}

	return c.Render(jobs, size, page, nextPage)
}

func (c Jobs) Post() revel.Result {
	return c.Render()
}

func (c Jobs) Show(id int) revel.Result {
	job, err := models.GetJob(c.Txn, id)
	if err != nil {
		revel.ERROR.Println(err)
	}
	revel.ERROR.Println(job)
	return c.Render(job)
}

func (c Jobs) Confirm(id int) revel.Result {
	models.ApproveJob(c.Txn, id)
	return c.Render(id)
}

func (c Jobs) HandlePostSubmit(job *models.Job) revel.Result {
	job.Validate(c.Validation)

	// Handle errors
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(Jobs.Post)
	}

	err := job.Save(c.Txn)
	if err != nil {
		revel.ERROR.Println(err)
	}
	// Ok, display the created job
	return c.Render(job)
}
