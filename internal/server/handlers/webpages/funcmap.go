package webpages

import (
	"html/template"

	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

// FuncMap defines functions used in the templates
func FuncMap() template.FuncMap {
	funcs := template.FuncMap{
		"jobStatusTableClass": jobStatusTableClass,
	}
	return funcs
}

func jobStatusTableClass(jobStatus models.JobStatus) string {
	class := ""

	switch jobStatus {
	case models.JobStatusFailed:
		return "table-danger"
	case models.JobStatusSuccess:
		return "table-success"
	}

	return class
}
