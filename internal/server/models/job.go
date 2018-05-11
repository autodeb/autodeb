package models

import (
	"fmt"
)

// JobType is the type of job
type JobType int

// Job Type enum
const (
	JobTypeUnknown JobType = iota
	JobTypeBuild
)

func (jt JobType) String() string {
	switch jt {
	case JobTypeUnknown:
		return "unknown"
	case JobTypeBuild:
		return "build"
	default:
		panic(fmt.Sprintf("Unknown job type %d", jt))
	}
}

// JobStatus is the status of the job
type JobStatus int

// Job Status enum
const (
	JobStatusUnknown JobStatus = iota
	JobStatusQueued
	JobStatusAssigned
	JobStatusSuccess
	JobStatusFailed
)

func (js JobStatus) String() string {
	switch js {
	case JobStatusUnknown:
		return "unknown"
	case JobStatusQueued:
		return "queued"
	case JobStatusAssigned:
		return "assigned"
	case JobStatusSuccess:
		return "success"
	case JobStatusFailed:
		return "failed"
	default:
		panic(fmt.Sprintf("Unknown job status %d", js))
	}
}

// Job is a builds a test, etc.
type Job struct {
	ID       uint      `json:"id"`
	Type     JobType   `json:"type"`
	Status   JobStatus `json:"status"`
	UploadID uint      `json:"upload_id"`
}
