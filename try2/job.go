package reservoir

import (
	"log"
	"time"
)

const (
	JOB_PENDING = iota
	JOB_DISPATCHING
	JOB_RUNNING
	JOB_FINISHED
)

type ReservoirJob struct {
	JobId        uint64 `PK`
	JobScript    string
	JobStatus    int
	JobLog       string
	LastModified string
}

func Job_New(JobScript string) *ReservoirJob {
	job := &ReservoirJob{
		0,
		JobScript,
		JOB_PENDING,
		"",
		time.Now().Format("2006-01-02 15:04:05"),
	}

	if Job_Update(job) {
		return job
	}
	return nil
}

func Job_Get(id uint64) *ReservoirJob {
	model, err := RetrieveDB()
	if err != nil {
		log.Printf("Failed to retrieve DB model for job: %s", err)
		return nil
	}

	defer model.Db.Close()

	var job ReservoirJob
	err = model.Where("jobid=?", id).Find(&job)
	if err != nil {
		log.Printf("Failed to find job for ID %d: %s", id, err)
		return nil
	}

	return &job
}

func Job_Update(job *ReservoirJob) bool {
	model, err := RetrieveDB()
	if err != nil {
		log.Printf("Failed to retrieve DB model for job: %s", err)
		return false
	}

	defer model.Db.Close()

	job.LastModified = time.Now().Format("2006-01-02 15:04:05")

	err = model.Save(job)
	if err != nil {
		log.Printf("Failed to update job to DB model: %s", err)
		return false
	}

	return true
}

func (r *ReservoirJob) GetId() uint64 {
	return r.JobId
}

func (r *ReservoirJob) GetScript() string {
	return r.JobScript
}

func (r *ReservoirJob) GetStatus() int {
	return r.JobStatus
}

func (r *ReservoirJob) GetLog() string {
	return r.JobLog
}

func (r *ReservoirJob) GetLastModification() string {
	return r.LastModified
}
