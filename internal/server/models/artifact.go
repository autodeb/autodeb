package models

// Artifact is part of the result of a job
type Artifact struct {
	ID       uint   `json:"id"`
	JobID    uint   `json:"job_id"`
	Filename string `json:"filename"`
}
