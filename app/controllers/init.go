package controllers

import (
	"encoding/csv"
	"github.com/bketelsen/gopheracademy/app/models"
	"github.com/revel/revel"
	"github.com/revel/revel/modules/jobs/app/jobs"
	"net/http"
	"os"
)

var Joblist map[int]*models.Job

func init() {

	Joblist = make(map[int]*models.Job)
	revel.OnAppStart(func() {
		jobs.Now(PopulateJobs{})
		jobs.Schedule("@every 24h", PopulateJobs{})
	})

}

type PopulateJobs struct{}

func (j PopulateJobs) Run() {

	gjpw := os.Getenv("GJPW")
	gjuser := os.Getenv("GJUSER")
	revel.INFO.Println(os.Getenv("GJPW"))
	// do the thing with the stuff
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://golangprojects.com/gatabfeed/active.txt", nil)
	req.SetBasicAuth(gjuser, gjpw)
	resp, err := client.Do(req)
	if err != nil {
		revel.ERROR.Printf("Error : %s", err)
	}

	defer resp.Body.Close()

	r := csv.NewReader(resp.Body)
	r.Comma = '\t'
	r.LazyQuotes = true

	records, err := r.ReadAll()
	if err != nil {
		revel.INFO.Println(err)
	}

	revel.INFO.Println(records)

	for i := 1; i < len(records); i++ {
		job := &models.Job{}
		revel.INFO.Printf("Processing job %d\n", i)
		job.FromRecord(records[i])
		revel.INFO.Println(job)
		Joblist[job.Id] = job
	}

}
