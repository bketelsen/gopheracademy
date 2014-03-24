package models

import ()

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
	Approved               bool
	Created                string
}

func GetJob(id int) (job *Job, e error) {

	return job, e
}

func GetJobs(page, size int) (jobs []*Job, e error) {

	return jobs, e
}
