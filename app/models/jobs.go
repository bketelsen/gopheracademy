package models

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/robfig/revel"
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
	created                string
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

	rows, e := txn.Query("select id,title, location, jobtype, description, companyname, companywebsite, companylogourl, applyurl, applyemail, additionalinstructions, purchaseremail  from jobs where id=?", id)
	if e != nil {
		revel.ERROR.Println(e)
		return job, e
	}

	for rows.Next() {

		var id, jobtype int
		var title, location, description, companyname, companywebsite, companylogourl, applyurl, applyemail, additionalinstructions, purchaseremail string
		if e := rows.Scan(&id, &title, &location, &jobtype, &description, &companyname, &companywebsite, &companylogourl, &applyurl, &applyemail, &additionalinstructions, &purchaseremail); e != nil {
			revel.ERROR.Println(e)
		}

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
		}

	}
	if e := rows.Err(); e != nil {
		revel.ERROR.Println(e)
	}

	return job, e
}

func GetJobs(txn *sql.Tx, page, size int) (jobs []*Job, e error) {

	rows, e := txn.Query("select id,title, location, jobtype, description, companyname, companywebsite, companylogourl, applyurl, applyemail, additionalinstructions, purchaseremail  from jobs where approved=1 limit ?, ?", (page-1)*size, size)
	if e != nil {
		revel.ERROR.Println(e)
		return jobs, e
	}

	for rows.Next() {

		var id, jobtype int
		var title, location, description, companyname, companywebsite, companylogourl, applyurl, applyemail, additionalinstructions, purchaseremail string
		if e := rows.Scan(&id, &title, &location, &jobtype, &description, &companyname, &companywebsite, &companylogourl, &applyurl, &applyemail, &additionalinstructions, &purchaseremail); e != nil {
			revel.ERROR.Println(e)
		}

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
		}

		jobs = append(jobs, job)

	}
	if e := rows.Err(); e != nil {
		revel.ERROR.Println(e)
	}

	return jobs, e
}
