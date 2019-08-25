package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Promo struct {
	Label       string `yaml:"label"`
	Description string `yaml:"description"`
	Path        string `yaml:"path"`
	Icon        string `yaml:"icon"`
}

type Linkbox struct {
	Name     string `yaml:"name"`
	Contents []struct {
		Label string `yaml:"label"`
		Path  string `yaml:"path"`
	} `yaml:"contents"`
}

type Footer struct {
	Footer []struct {
		Promos    []Promo   `yaml:"promos,omitempty"`
		Linkboxes []Linkbox `yaml:"linkboxes,omitempty"`
	} `yaml:"footer"`
}

func ParseFooter(filepath string) ([]Promo, []Linkbox, error) {
	footerContent, err := ioutil.ReadFile("src/content/" + filepath)
	if err != nil {
		return nil, nil, err
	}
	footer := Footer{}
	promos := []Promo{}
	linkboxes := []Linkbox{}
	err = yaml.Unmarshal(footerContent, &footer)
	for _, f := range footer.Footer {
		promos = append(promos, f.Promos...)
		linkboxes = append(linkboxes, f.Linkboxes...)
	}
	return promos, linkboxes, nil
}
