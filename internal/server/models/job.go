package models

// JobType is the type of job
type JobType int

// Job Type enum
const (
	JobTypeUnknown JobType = iota
	JobTypeBuild
	JobTypeAutopkgtest
	JobTypeForward
	JobTypeSetupArchiveUpgrade
	JobTypePackageUpgrade
)

func (jt JobType) String() string {
	switch jt {
	case JobTypeBuild:
		return "build"
	case JobTypeAutopkgtest:
		return "autopkgtest"
	case JobTypeForward:
		return "forward"
	case JobTypeSetupArchiveUpgrade:
		return "setup-archive-upgrade"
	case JobTypePackageUpgrade:
		return "package-upgrade"
	default:
		return "unknown"
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
	case JobStatusQueued:
		return "queued"
	case JobStatusAssigned:
		return "assigned"
	case JobStatusSuccess:
		return "success"
	case JobStatusFailed:
		return "failed"
	default:
		return "unknown"
	}
}

// JobParentType represents the parent of the job
type JobParentType int

// JobParentType enum
const (
	JobParentTypeUnknown JobParentType = iota
	JobParentTypeUpload
	JobParentTypeArchiveUpgrade
)

func (parentType JobParentType) String() string {
	switch parentType {
	case JobParentTypeUpload:
		return "upload"
	case JobParentTypeArchiveUpgrade:
		return "archive-upgrade"
	default:
		return "unknown"
	}
}

// Job is a builds a test, etc.
type Job struct {
	ID     uint      `json:"id"`
	Type   JobType   `json:"type"`
	Status JobStatus `json:"status"`

	// == Job parent ==

	// The ID of the Job's parent.
	// It is propagated to all child jobs.
	ParentID   uint          `json:"parent_id"`
	ParentType JobParentType `json:"parent_type"`

	// == JOB INPUTS ==

	// Some jobs take an input.
	// For example, the input of an Autopkgtest job is an artifact id
	// that points to the .deb to test.
	Input string `json:"input_string"`
}
