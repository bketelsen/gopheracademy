package models

import "github.com/robfig/revel"

type Job struct {
	Id                     int
	Title                  string
	Location               string
	JobType                int
	Description            string
	CompanyName            string
	CompanyWebsite         string
	CompanyLogoURL         string
	ApplyURL               string
	ApplyEmail             string
	AdditionalInstructions string
	PurchaserEmail         string
}

func (job *Job) Validate(v *revel.Validation) {
	v.Required(job.Title)
	v.Required(job.Location)
	v.Required(job.Description)
	v.Required(job.CompanyName)
	v.Required(job.CompanyWebsite)
	v.Required(job.CompanyLogoURL)
	v.Required(job.AdditionalInstructions)
	v.Required(job.PurchaserEmail)

}
