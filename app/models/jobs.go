package models

import (
	"strconv"

	"strings"
)

type Job struct {
	Id             int
	Title          string
	Location       string
	JobType        int
	Description    string
	Requirements   string
	CompanyName    string
	CompanyWebsite string
	MoreInfoURL    string
}

func GetJob(id int) (job *Job, e error) {

	return job, e
}

func GetJobs(page, size int) (jobs []*Job, e error) {

	return jobs, e
}

func (j *Job) FromRecord(record []string) {

	j.Id, _ = strconv.Atoi(record[0])
	j.Title = record[1]
	j.CompanyName = record[7]
	j.Location = record[4] + ", " + record[5] + ", " + record[6]
	j.Description = strings.Replace(record[2], "\\n", "\n", -1)
	j.Requirements = strings.Replace(record[3], "\\n", "\n", -1)
	j.CompanyWebsite = record[8]
	j.MoreInfoURL = record[10]
}
