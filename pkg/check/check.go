package check

import (
	"github.com/jchesterpivotal/knative-service-resource/pkg"
	"github.com/jchesterpivotal/knative-service-resource/pkg/concourse"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"strconv"
	"errors"
	"fmt"
)

type Checker interface {
	Check() (Output, error)
}

type Input struct {
	Source  concourse.Source  `json:"source"`
	Version concourse.Version `json:"version,omitempty"`
}

type Output []concourse.Version

type checker struct {
	clients *clients.Clients

	source  *concourse.Source
	version *concourse.Version
}

func (c *checker) isFirstCheck() bool {
	return c.version == nil
}

func (c *checker) compareConcourseVersionTo(version string) (string, error) {
	cv, err := strconv.Atoi(c.version.ConfigurationGeneration)
	if err != nil {
		return "", nil
	}

	kv, err := strconv.Atoi(version)
	if err != nil {
		return "", nil
	}

	if kv > cv {
		return "KnativeVersionHigher", nil
	}

	if cv > kv {
		return "ConcourseVersionHigher", nil
	}

	return "VersionsEqual", nil
}

func (c *checker) latestGenerationInKnative() (string, error) {
	serviceName := c.source.Name
	service, err := c.clients.Service.Get(serviceName, v1.GetOptions{IncludeUninitialized: false})
	if err != nil {
		return "", err
	}

	observedGeneration := service.Status.ObservedGeneration
	return strconv.Itoa(int(observedGeneration)), nil
}

func (c *checker) versionsInKnativeSince(version string) ([]concourse.Version, error) {
	sel := fmt.Sprintf("serving.knative.dev/configuration=%s", c.source.Name)
	revs, err := c.clients.Revision.List(v1.ListOptions{LabelSelector: sel})
	if err != nil {
		return nil, err
	}

	versions := make([]concourse.Version, 0)
	for _, r := range revs.Items {
		gen := strconv.Itoa(int(r.GetGeneration()))
		versions = append(versions, concourse.Version{ConfigurationGeneration: gen})
	}

	return versions, nil
}

func (c *checker) Check() (Output, error) {
	latestInKnative, err := c.latestGenerationInKnative()
	if err != nil {
		return nil, fmt.Errorf("could not find Knative service '%s' in Kubernetes: %s", c.source.Name, err)
	}

	if c.isFirstCheck() {
		return []concourse.Version{
			{ConfigurationGeneration: latestInKnative},
		}, nil
	}

	compared, err := c.compareConcourseVersionTo(latestInKnative)
	if err != nil {
		return nil, err
	}

	switch compared {
	case "VersionsEqual":
		return []concourse.Version{*c.version}, nil
	case "KnativeVersionHigher":
		return c.versionsInKnativeSince(latestInKnative)
	case "ConcourseVersionHigher":
		return nil, fmt.Errorf(
			"version known to Concourse (%s) was ahead of version known to Kubernetes (%s)",
			c.version.ConfigurationGeneration,
			latestInKnative,
		)
	default:
		return nil, errors.New("'impossible' error occurred while comparing Knative Service versions in Concourse and Kubernetes")
	}

	return nil, nil
}

func NewChecker(clients *clients.Clients, source *concourse.Source, version *concourse.Version) Checker {
	return &checker{
		clients: clients,
		source:  source,
		version: version,
	}
}
