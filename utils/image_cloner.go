package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-logr/logr"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	v1 "k8s.io/api/core/v1"
)

var (
	registry   string
	repository string
	username   string
	password   string
)

func init() {
	registry = os.Getenv("PRIV_OCI_REGISTRY")
	if registry == "" {
		panic("PRIV_OCI_REGISTRY not set")
	}
	repository = os.Getenv("PRIV_OCI_REPOSITORY")
	if repository == "" {
		panic("PRIV_OCI_REPOSITORY not set")
	}
	username = os.Getenv("PRIV_OCI_REGISTRY_USERNAME")
	if username == "" {
		panic("PRIV_OCI_REGISTRY_USERNAME not set")
	}
	password = os.Getenv("PRIV_OCI_REGISTRY_PASSWORD")
	if password == "" {
		panic("PRIV_OCI_REGISTRY_PASSWORD not set")
	}
}

func CloneImage(logr logr.Logger, podTemplateSpec *v1.PodTemplateSpec) (*v1.PodTemplateSpec, error) {
	//TODO code can be refactored
	for i, container := range podTemplateSpec.Spec.InitContainers {
		OldRef, err := name.ParseReference(container.Image)
		if err != nil {
			continue
		}
		img, err := remote.Image(OldRef)
		if err != nil {
			continue
		}

		//TODO externalize the hardcoded registry name
		if OldRef.Context().RegistryStr() == "index.docker.io" && !strings.HasPrefix(OldRef.Context().RepositoryStr(), repository) {
			//TODO check if image already exists
			index := strings.IndexAny(OldRef.Context().RepositoryStr(), "/")
			if index != -1 {
				container.Image = container.Image[index+1:]
			}
			newImage := fmt.Sprintf("%s/%s/%s", registry, repository, container.Image)
			newRef, err := name.ParseReference(newImage)
			if err != nil {
				continue
			}
			//TODO - requires authentication each time.
			err = remote.Write(newRef, img, remote.WithAuth(&authn.Basic{Username: username, Password: password}))
			if err != nil {
				return podTemplateSpec, err
			}
			podTemplateSpec.Spec.InitContainers[i].Image = newImage
		}
	}
	for i, container := range podTemplateSpec.Spec.Containers {
		logr.Info("Looping over containers", "Image", container.Image)
		OldRef, err := name.ParseReference(container.Image)
		if err != nil {
			continue
		}
		img, err := remote.Image(OldRef)
		if err != nil {
			continue
		}
		//TODO externalize the hardcoded registry name
		if OldRef.Context().RegistryStr() == "index.docker.io" && !strings.HasPrefix(OldRef.Context().RepositoryStr(), repository) {
			logr.Info("Matches Business criteria of being a public Image", "Repo", OldRef.Context().RepositoryStr())
			//TODO check if image already exists
			index := strings.IndexAny(container.Image, "/")
			if index != -1 {
				container.Image = container.Image[index+1:]
			}
			newImage := fmt.Sprintf("%s/%s/%s", registry, repository, container.Image)
			logr.Info("Container Image Spec will be updated", "New Image", newImage)
			newRef, err := name.ParseReference(newImage)
			if err != nil {
				continue
			}
			err = remote.Write(newRef, img, remote.WithAuth(&authn.Basic{Username: username, Password: password}))
			if err != nil {
				return podTemplateSpec, err
			}
			podTemplateSpec.Spec.Containers[i].Image = newImage
		}
	}
	return podTemplateSpec, nil
}
