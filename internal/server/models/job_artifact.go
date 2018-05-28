package models

// JobArtifact is part of the result of a job
type JobArtifact struct {
	ID       uint   `json:"id"`
	JobID    uint   `json:"job_id"`
	Filename string `json:"filename"`
}
