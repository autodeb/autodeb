package models

// AddBuildToRepositoryInput is the json input for
// the AddBuildToRepository job type
type AddBuildToRepositoryInput struct {
	RepositoryName string
	Distribution   string
}
