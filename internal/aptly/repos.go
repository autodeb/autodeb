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
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("unexpected status: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(repositories); err != nil {
		return nil, err
	}

	return repositories, err
}

// GetRepository will return the repository with the corresponding name
func (client *APIClient) GetRepository(name string) (*Repository, error) {
	var repository *Repository

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
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("unexpected status: %s", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(repository); err != nil {
		return nil, err
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
		return nil, err
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
		return nil, err
	}

	return createdRepository, nil
}

// AddPackageToRepository adds a package to a repositoryA
func (client *APIClient) AddPackageToRepository(pkg, dir, repository string) error {
	resp, err := client.do(
		http.MethodPost,
		fmt.Sprintf("/repos/%s/file/%s/%s", repository, dir, pkg),
		"",
		nil,
	)
	if err != nil {
		return err
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
