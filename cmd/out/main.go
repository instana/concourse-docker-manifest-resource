package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/mbialon/concourse-docker-manifest-resource/pkg/docker"
	"github.com/mbialon/concourse-docker-manifest-resource/pkg/docker/manifest"
)

type Request struct {
	Source *Source `json:"source"`
	Params *Params `json:"params"`
}

type Source struct {
	Repository string `json:"repository"`
	Tag        string `json:"tag"`
	Username   string `json:"username"`
	Password   string `json:"password"`
}

type Params struct {
	Manifests []Manifest `json:"manifests"`
}

type Manifest struct {
	Arch    string `json:"arch"`
	OS      string `json:"os"`
	TagFile string `json:"tag_file"`
}

func main() {
	if err := os.Chdir(os.Args[1]); err != nil {
		log.Fatalf("cannot change dir: %v", err)
	}
	var request Request
	if err := json.NewDecoder(os.Stdin).Decode(&request); err != nil {
		log.Fatalf("cannot decode input: %v", err)
	}
	if err := docker.Login(request.Source.Username, request.Source.Password); err != nil {
		log.Fatalf("cannot login to docker hub: %v", err)
	}
	manifestList := request.Source.Repository + ":" + request.Source.Tag
	var manifests []string
	var annotations []manifest.Annotation
	for _, m := range request.Params.Manifests {
		tag, err := readTag(m.TagFile)
		if err != nil {
			log.Fatalf("cannot read tag: %v", err)
		}
		ref := request.Source.Repository + ":" + tag
		manifests = append(manifests, ref)
		annotations = append(annotations, manifest.Annotation{
			Manifest:     ref,
			Architecture: m.Arch,
			OS:           m.OS,
		})
	}
	if err := manifest.Create(manifestList, manifests); err != nil {
		log.Fatalf("cannot create manifest: %v", err)
	}
	if err := manifest.Annotate(manifestList, annotations); err != nil {
		log.Fatalf("cannot annotate manifest: %v", err)
	}
	digest, err := manifest.Push(manifestList)
	if err != nil {
		log.Fatalf("cannot push manifest: %v", err)
	}
	output := map[string]interface{}{
		"version": map[string]string{
			"digest": digest,
		},
	}
	if err := json.NewEncoder(os.Stdout).Encode(output); err != nil {
		log.Fatalf("cannot encode output: %v", err)
	}
}

func readTag(filename string) (string, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("cannot read tag file: %v", err)
	}
	return strings.TrimSpace(string(b)), nil
}
