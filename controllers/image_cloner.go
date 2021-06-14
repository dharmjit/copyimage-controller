package controllers

import (
	"fmt"
	"os"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	v1 "k8s.io/api/core/v1"
)

// var (
// 	registry   = "index.docker.io"
// 	repository = "dharmjit"
// 	username   = "dharmjit"
// 	password   = "main22atthaan"
// )

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

func cloneImage(podTemplateSpec *v1.PodTemplateSpec) (*v1.PodTemplateSpec, error) {
	//TODO code can be refactored
	for i, image := range podTemplateSpec.Spec.InitContainers {
		OldRef, err := name.ParseReference(image.Image)
		if err != nil {
			return podTemplateSpec, err
		}
		img, err := remote.Image(OldRef)
		if err != nil {
			return podTemplateSpec, err
		}
		newImage := fmt.Sprintf("%s/%s/%s", registry, repository, image.Image)
		newRef, err := name.ParseReference(newImage)
		if err != nil {
			return podTemplateSpec, err
		}
		//TODO include more public docker registries
		if OldRef.Context().RegistryStr() == "index.docker.io" {
			//TODO check if image already exists
			err := remote.Write(newRef, img, remote.WithAuth(&authn.Basic{Username: username, Password: password}))
			if err != nil {
				return podTemplateSpec, err
			}
			podTemplateSpec.Spec.InitContainers[i].Image = newImage
		}
	}
	for i, image := range podTemplateSpec.Spec.Containers {
		OldRef, err := name.ParseReference(image.Image)
		if err != nil {
			return podTemplateSpec, err
		}
		img, err := remote.Image(OldRef)
		if err != nil {
			return podTemplateSpec, err
		}
		newImage := fmt.Sprintf("%s/%s/%s", "index.docker.io", "dharmjit", image.Image)
		newRef, err := name.ParseReference(newImage)
		if err != nil {
			return podTemplateSpec, err
		}
		if OldRef.Context().RegistryStr() == "index.docker.io" {
			//TODO check if image already exists
			err := remote.Write(newRef, img, remote.WithAuth(&authn.Basic{Username: "dharmjit", Password: "main22atthaan"}))
			if err != nil {
				return podTemplateSpec, err
			}
			podTemplateSpec.Spec.Containers[i].Image = newImage
		}
	}
	return podTemplateSpec, nil
}
