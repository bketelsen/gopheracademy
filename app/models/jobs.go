package models

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/robfig/revel"
	"time"
)

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

func (job *Job) Save(txn *sql.Tx) (e error) {
	// save to the database
	_, e = txn.Exec(
		"INSERT INTO jobs (title, location, jobtype, description, companyname, companywebsite, companylogourl, applyurl, applyemail,additionalinstructions, purchaseremail, approved) VALUES (?, ?, ?, ?, ?,?,?,?,?,?,?,?)",
		job.Title,
		job.Location,
		job.JobType,
		job.Description,
		job.CompanyName,
		job.CompanyWebsite,
		job.CompanyLogoURL,
		job.ApplyURL,
		job.ApplyEmail,
		job.AdditionalInstructions,
		job.PurchaserEmail,
		0)

	if e != nil {
		revel.ERROR.Println(e)
		return e
	}
	var id int
	err := txn.QueryRow("SELECT LAST_INSERT_ID() ").Scan(&id)
	if err != nil {
		revel.ERROR.Println(err)
		return err
	}
	job.Id = id

	return nil
}

func ApproveJob(txn *sql.Tx, id int) (e error) {
	//update the approved flag to 1
	_, e = txn.Exec(
		"update jobs set approved = 1 where id = ?",
		id,
	)

	if e != nil {
		revel.ERROR.Println(e)
		return e
	}
	return nil
}
func GetJob(txn *sql.Tx, id int) (job *Job, e error) {

	rows, e := txn.Query("select id,title, location, jobtype, description, companyname, companywebsite, companylogourl, applyurl, applyemail, additionalinstructions, purchaseremail,created  from jobs where id=?", id)
	if e != nil {
		revel.ERROR.Println(e)
		return job, e
	}

	for rows.Next() {

		var id, jobtype int
		var title, location, description, companyname, companywebsite, companylogourl, applyurl, applyemail, additionalinstructions, purchaseremail, created string
		if e := rows.Scan(&id, &title, &location, &jobtype, &description, &companyname, &companywebsite, &companylogourl, &applyurl, &applyemail, &additionalinstructions, &purchaseremail, &created); e != nil {
			revel.ERROR.Println(e)
		}
		revel.ERROR.Println(created)

		var parsedDate time.Time
		parsedDate, err := time.Parse("2006-01-02 15:04:05", created)

		if err != nil {
			parsedDate, _ = time.Parse("2006-01-02 15:04:05 MST", "January 1, 2013")

		}

		formattedDate := parsedDate.Format("2006-01-02")

		job = &Job{
			Id:                     id,
			Title:                  title,
			Location:               location,
			JobType:                jobtype,
			Description:            description,
			CompanyName:            companyname,
			CompanyWebsite:         companywebsite,
			CompanyLogoURL:         companylogourl,
			ApplyURL:               applyurl,
			ApplyEmail:             applyemail,
			AdditionalInstructions: additionalinstructions,
			PurchaserEmail:         purchaseremail,
			Created:                formattedDate,
		}

	}
	if e := rows.Err(); e != nil {
		revel.ERROR.Println(e)
	}

	return job, e
}

func GetJobs(txn *sql.Tx, page, size int) (jobs []*Job, e error) {

	rows, e := txn.Query("select id,title, location, jobtype, description, companyname, companywebsite, companylogourl, applyurl, applyemail, additionalinstructions, purchaseremail,created  from jobs where approved=1 limit ?, ?", (page-1)*size, size)
	if e != nil {
		revel.ERROR.Println(e)
		return jobs, e
	}

	for rows.Next() {

		var id, jobtype int
		var title, location, description, companyname, companywebsite, companylogourl, applyurl, applyemail, additionalinstructions, purchaseremail, created string
		if e := rows.Scan(&id, &title, &location, &jobtype, &description, &companyname, &companywebsite, &companylogourl, &applyurl, &applyemail, &additionalinstructions, &purchaseremail, &created); e != nil {
			revel.ERROR.Println(e)
		}

		var parsedDate time.Time
		parsedDate, err := time.Parse("2006-01-02 15:04:05", created)

		if err != nil {
			parsedDate, _ = time.Parse("2006-01-02 15:04:05 MST", "January 1, 2013")

		}

		formattedDate := parsedDate.Format("2006-01-02")

		job := &Job{
			Id:                     id,
			Title:                  title,
			Location:               location,
			JobType:                jobtype,
			Description:            description,
			CompanyName:            companyname,
			CompanyWebsite:         companywebsite,
			CompanyLogoURL:         companylogourl,
			ApplyURL:               applyurl,
			ApplyEmail:             applyemail,
			AdditionalInstructions: additionalinstructions,
			PurchaserEmail:         purchaseremail,
			Created:                formattedDate,
		}

		jobs = append(jobs, job)

	}
	if e := rows.Err(); e != nil {
		revel.ERROR.Println(e)
	}

	return jobs, e
}
