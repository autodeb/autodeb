package aptly

import (
	"bytes"
	"encoding/json"
	"fmt"
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
