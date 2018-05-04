package models

import (
	"fmt"
)

// JobType is the type of job
type JobType int

// Job Type enum
const (
	JobTypeBuild JobType = iota
)

func (jt JobType) String() string {
	switch jt {
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
	JobStatusQueued JobStatus = iota
)

func (js JobStatus) String() string {
	switch js {
	case JobStatusQueued:
		return "queued"
	default:
		panic(fmt.Sprintf("Unknown job status %d", js))
	}
}

// Job is a builds a test, etc.
type Job struct {
	ID       uint
	Type     JobType
	Status   JobStatus
	UploadID uint
}
