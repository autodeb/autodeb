package aptly

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"salsa.debian.org/autodeb-team/autodeb/internal/errors"
)

// Source is a local repository or a snapshot name
type Source struct {
	Component string `json:"Component"`
	Name      string `json:"Name"`
}

// SigningOptions contain gpg options
type SigningOptions struct {
	Skip           bool   `json:"Skip"`
	Batch          bool   `json:"Batch"`
	GpgKey         string `json:"GpgKey"`
	Keyring        string `json:"Keyring"`
	SecretKeyring  string `json:"SecretKeyring"`
	Passphrase     string `json:"Passphrase"`
	PassphraseFile string `json:"PassphraseFile"`
}

// PublishParameters are used to publish a snapshot or local repository
type PublishParameters struct {
	SourceKind           string          `json:"SourceKind"`
	Sources              []*Source       `json:"Sources"`
	Distribution         string          `json:"Distribution"`
	Label                string          `json:"Label"`
	Origin               string          `json:"Origin"`
	ForceOverwrite       bool            `json:"ForceOverwrite"`
	Architectures        []string        `json:"Architectures"`
	Signing              *SigningOptions `json:"Signing"`
	NotAutomatic         string          `json:"NotAutomatic"`
	ButAutomaticUpgrades string          `json:"ButAutomaticUpgrades"`
	SkipCleanup          string          `json:"SkipCleanup"`
	AcquireByHash        bool            `json:"AcquireByHash"`
}

// Publish will publish a snapshot or a local repo
func (client *APIClient) Publish(prefix string, publishParameters *PublishParameters) error {
	paramsBytes, err := json.Marshal(publishParameters)
	if err != nil {
		return err
	}

	resp, err := client.do(
		http.MethodPost,
		fmt.Sprintf("/publish/%s", prefix),
		"application/json",
		bytes.NewReader(paramsBytes),
	)
	if err != nil {
		return errors.WithMessage(err, "request failed")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return errors.Errorf("unexpected response status: %d", resp.StatusCode)
	}

	createdRepository := &Repository{}

	if err := json.NewDecoder(resp.Body).Decode(createdRepository); err != nil {
		return errors.WithMessage(err, "could not decode repository")
	}

	return nil
}

// PublishDefaults will publish a snapshot or a local repo with sensible defaults
func (client *APIClient) PublishDefaults(repositoryName string) error {
	publishParameters := &PublishParameters{
		SourceKind: "local",
		Sources: []*Source{
			&Source{
				Name: repositoryName,
			},
		},
		ForceOverwrite: true,
		Architectures: []string{
			"amd64",
			"arm64",
			"armel",
			"armhf",
			"i386",
			"mips",
			"mips64el",
			"ppc64el",
			"s390x",
			"all",
		},
		Signing: &SigningOptions{
			Skip: true,
		},
	}

	return client.Publish(repositoryName, publishParameters)
}

// PublishedRepositoryUpdateParameters holds parameters required to update a published repository
type PublishedRepositoryUpdateParameters struct {
	Snapshots      []*Source       `json:"Snapshots"`
	ForceOverwrite bool            `json:"ForceOverwrite"`
	Signing        *SigningOptions `json:"Signing"`
	AcquireByHash  bool            `json:"AcquireByHash"`
}

// UpdatePublishedRepository updates a published repository
func (client *APIClient) UpdatePublishedRepository(prefix, distribution string, params *PublishedRepositoryUpdateParameters) error {
	paramsBytes, err := json.Marshal(params)
	if err != nil {
		return err
	}

	resp, err := client.do(
		http.MethodPut,
		fmt.Sprintf("/publish/%s/%s", prefix, distribution),
		"application/json",
		bytes.NewReader(paramsBytes),
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

// UpdatePublishedRepositoryDefaults updates a published repository with sensible defaults
func (client *APIClient) UpdatePublishedRepositoryDefaults(repositoryName, distribution string) error {
	params := &PublishedRepositoryUpdateParameters{
		ForceOverwrite: true,
		Signing: &SigningOptions{
			Skip: true,
		},
		AcquireByHash: true,
	}

	return client.UpdatePublishedRepository(repositoryName, distribution, params)
}
