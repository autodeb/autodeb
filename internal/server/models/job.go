package models

// JobType is the type of job
type JobType int

// Job Type enum
const (
	JobTypeUnknown JobType = iota
	JobTypeBuildUpload
	JobTypeAutopkgtest
	JobTypeForwardUpload
	JobTypeSetupArchiveUpgrade
	JobTypePackageUpgrade
	JobTypeAddBuildToRepository
	JobTypeSetupArchiveBackport
	JobTypeBackport
)

func (jt JobType) String() string {
	switch jt {
	case JobTypeBuildUpload:
		return "build-upload"
	case JobTypeAutopkgtest:
		return "autopkgtest"
	case JobTypeForwardUpload:
		return "forward-upload"
	case JobTypeSetupArchiveUpgrade:
		return "setup-archive-upgrade"
	case JobTypePackageUpgrade:
		return "package-upgrade"
	case JobTypeAddBuildToRepository:
		return "add-build-to-repository"
	case JobTypeSetupArchiveBackport:
		return "setup-archive-backport"
	case JobTypeBackport:
		return "backport"
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
	JobParentTypeArchiveBackport
)

func (parentType JobParentType) String() string {
	switch parentType {
	case JobParentTypeUpload:
		return "upload"
	case JobParentTypeArchiveUpgrade:
		return "archive-upgrade"
	case JobParentTypeArchiveBackport:
		return "archive-backport"
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

	// Jobs can be created in the context of:
	//  - An Upload
	//  - An Archive Rebuild
	// The id of the job's ultimate parent is propagated
	// to all child jobs.
	ParentID   uint          `json:"parent_id"`
	ParentType JobParentType `json:"parent_type"`

	// Some jobs are associated to a build job. For example, an Autopkgtest job
	// will test the result of a build job
	BuildJobID uint `json:"build_job_id"`

	// == JOB INPUTS ==

	// Some jobs take an input.
	// For example, the input of an Autopkgtest job is an artifact id
	// that points to the .deb to test.
	Input string `json:"input_string"`
}
