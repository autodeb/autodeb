package aptly

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"salsa.debian.org/autodeb-team/autodeb/internal/errors"
)

// Repository represents a repository
type Repository struct {
	Name                string `json:"Name"`
	Comment             string `json:"Comment"`
	DefaultDistribution string `json:"DefaultDistribution"`
	DefaultComponent    string `json:"DefaultComponent"`
}

// GetRepositories returns the list of existing repositories
func (client *APIClient) GetRepositories() ([]*Repository, error) {
	var repositories []*Repository

	resp, err := client.do(
		http.MethodGet,
		"/repos",
		"application/json",
		nil,
	)
	if err != nil {
		return nil, errors.WithMessage(err, "request failed")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("unexpected status: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(repositories); err != nil {
		return nil, errors.WithMessage(err, "could not unmarshal repositories")
	}

	return repositories, err
}

// GetRepository will return the repository with the corresponding name
func (client *APIClient) GetRepository(name string) (*Repository, error) {
	repository := &Repository{}

	resp, err := client.do(
		http.MethodGet,
		fmt.Sprintf(
			"/repos/%s",
			name,
		),
		"application/json",
		nil,
	)
	if err != nil {
		return nil, errors.WithMessage(err, "request failed")
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("unexpected status: %s", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(repository); err != nil {
		return nil, errors.WithMessage(err, "could not unmarshal repository")
	}

	return repository, err
}

// CreateRepository creates a new repository
func (client *APIClient) CreateRepository(name, comment, defaultDistribution, defaultComponent string) (*Repository, error) {
	postedRepository := Repository{
		Name:                name,
		Comment:             comment,
		DefaultDistribution: defaultDistribution,
		DefaultComponent:    defaultComponent,
	}

	postedRepoBytes, err := json.Marshal(postedRepository)
	if err != nil {
		return nil, errors.WithMessagef(err, "could not encode repository %+v", postedRepository)
	}

	resp, err := client.do(
		http.MethodPost,
		"/repos",
		"application/json",
		bytes.NewReader(postedRepoBytes),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, errors.Errorf("unexpected response status: %d", resp.StatusCode)
	}

	createdRepository := &Repository{}

	if err := json.NewDecoder(resp.Body).Decode(createdRepository); err != nil {
		return nil, errors.WithMessage(err, "could not decode repository")
	}

	return createdRepository, nil
}

// CreateRepositoryDefaults creates a new repository with sensible defaults
func (client *APIClient) CreateRepositoryDefaults(name, distribution string) (*Repository, error) {
	repo, err := client.CreateRepository(
		name,
		fmt.Sprintf("Packages for %s on autodeb", name),
		distribution,
		"main",
	)
	return repo, err
}

// GetOrCreateAndPublishRepository will get or create and publish a repository with sensible defaults
func (client *APIClient) GetOrCreateAndPublishRepository(name, distribution string) (*Repository, error) {

	// Get the repository, if it already exists
	if repo, err := client.GetRepository(name); err != nil {
		return nil, errors.WithMessagef(err, "could not retrieve repository %s", name)
	} else if repo != nil {
		return repo, nil
	}

	// Create the repository
	repo, err := client.CreateRepositoryDefaults(name, distribution)
	if err != nil {
		return nil, errors.WithMessagef(err, "could not create repository %s", name)
	}

	// Publish it for the first time
	if err := client.PublishDefaults(name); err != nil {
		return nil, errors.WithMessagef(err, "could not publish newly created repository %s", name)
	}

	return repo, nil
}

// AddPackageToRepository adds a package to a repository
func (client *APIClient) AddPackageToRepository(pkg, dir, repository string) error {
	resp, err := client.do(
		http.MethodPost,
		fmt.Sprintf("/repos/%s/file/%s/%s", repository, dir, pkg),
		"",
		nil,
	)
	if err != nil {
		return errors.WithMessage(err, "request failed")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("unexpected response status: %d", resp.StatusCode)
	}

	return nil
}

// UploadPackageAndAddToRepository will upload a package and add it to a repository
func (client *APIClient) UploadPackageAndAddToRepository(packageName string, packageContent io.Reader, repositoryName string) error {

	// Upload the file
	if err := client.UploadFileInDirectory(
		packageContent,
		packageName,
		repositoryName,
	); err != nil {
		return errors.WithMessagef(err, "could not upload %s to aptly", packageName)
	}

	// Add the package
	if err := client.AddPackageToRepository(
		packageName,
		repositoryName,
		repositoryName,
	); err != nil {
		return errors.WithMessagef(err, "could not add %s to the repository", packageName)
	}

	return nil
}
