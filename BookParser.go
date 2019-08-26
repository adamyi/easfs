package main

import (
	"fmt"
	"html"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

type LowerTabContent struct {
	Title          string            `yaml:"title,omitempty"`
	Path           string            `yaml:"path,omitempty"`
	Include        string            `yaml:"include,omitempty"`
	Status         string            `yaml:"status,omitempty"`
	Heading        string            `yaml:"heading,omitempty"`
	Style          string            `yaml:"style,omitempty"`
	AlternatePaths []string          `yaml:"alternate_paths,omitempty"`
	Section        []LowerTabContent `yaml:"section,omitempty"`
}

type LowerTab struct {
	Name     string            `yaml:"name"`
	Path     string            `yaml:"path,omitempty"`
	Selected bool              `yaml:"selected,omitempty"`
	Contents []LowerTabContent `yaml:"contents"`
}

type UpperTab struct {
	Include    string `yaml:"include,omitempty"`
	Name       string `yaml:"name,omitempty"`
	Heading    string `yaml:"heading,omitempty"`
	Path       string `yaml:"path,omitempty"`
	Attributes []struct {
		Name  string `yaml:"name"`
		Value string `yaml:"value"`
	} `yaml:"attributes,omitempty"`
	LowerTabs struct {
		Other  []LowerTab        `yaml:"other"`
		Guides []LowerTabContent `yaml:"guides"`
	} `yaml:"lower_tabs,omitempty"`
}

type ToCLowerTabYAML struct {
	ToC []LowerTabContent `yaml:"toc"`
}

type ToCUpperTabYAML struct {
	ToC []UpperTab `yaml:"toc"`
}

type Book struct {
	UpperTabs []UpperTab `yaml:"upper_tabs"`
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func ExpandUpperTab(tab UpperTab) UpperTab {
	ret := UpperTab{}
	ret.Name = tab.Name
	ret.Path = tab.Path
	ret.Heading = tab.Heading
	ret.Attributes = tab.Attributes
	ret.LowerTabs.Guides = tab.LowerTabs.Guides
	for _, tab := range tab.LowerTabs.Other {
		ltab := LowerTab{}
		ltab.Name = tab.Name
		for _, tabcontent := range tab.Contents {
			ltabc := tabcontent
			if tabcontent.Include != "" {
				tocContent, err := ioutil.ReadFile(flagSitePath + tabcontent.Include)
				if err != nil {
					fmt.Println(err)
				}
				// fmt.Println(tocContent)
				tocYaml := ToCLowerTabYAML{}
				err = yaml.Unmarshal(tocContent, &tocYaml)
				if err != nil {
					fmt.Println(err)
				}
				for _, ltabcc := range tocYaml.ToC {
					ltab.Contents = append(ltab.Contents, ltabcc)
				}
			} else {
				ltab.Contents = append(ltab.Contents, ltabc)
			}
		}
		ret.LowerTabs.Other = append(ret.LowerTabs.Other, ltab)
	}
	return ret
}

func ExpandBook(book Book) Book {
	ret := Book{}
	for _, tab := range book.UpperTabs {
		if tab.Include != "" {
			tocContent, err := ioutil.ReadFile(flagSitePath + tab.Include)
			if err != nil {
				fmt.Println(err)
			}
			newTab := ToCUpperTabYAML{}
			err = yaml.Unmarshal(tocContent, &newTab)
			// fmt.Println(newTab)
			if err != nil {
				fmt.Println(err)
			}
			for _, utab := range newTab.ToC {
				ret.UpperTabs = append(ret.UpperTabs, ExpandUpperTab(utab))
			}
		} else {
			ret.UpperTabs = append(ret.UpperTabs, ExpandUpperTab(tab))
		}
	}
	// a, _ := yaml.Marshal(ret)
	// fmt.Println(string(a))
	return ret
}

func ParseBook(filepath string) (Book, error) {
	bookContent, err := ioutil.ReadFile(flagSitePath + filepath)
	if err != nil {
		return Book{}, err
	}
	bookYAML := Book{}
	err = yaml.Unmarshal(bookContent, &bookYAML)
	if err != nil {
		return Book{}, err
	}
	return ExpandBook(bookYAML), nil
}

func GetLowerTabs(requestPath string, book Book) []LowerTab {
	ret := []LowerTab{}
	for _, tab := range book.UpperTabs {
		for _, ltab := range tab.LowerTabs.Other {
			lt := LowerTab{}
			lt.Name = ltab.Name
			if len(ltab.Contents) > 0 {
				lt.Path = GetFirstTabPath(ltab.Contents)
				if lt.Path != "" {
					if LowerTabMatchesPath(requestPath, ltab.Contents) {
						lt.Selected = true
					}
					ret = append(ret, lt)
				}
			}
		}
	}
	return ret
}

func GetFirstTabPath(tabContent []LowerTabContent) string {
	for _, tc := range tabContent {
		if tc.Path != "" {
			return tc.Path
		}
		if len(tc.Section) > 0 {
			return GetFirstTabPath(tc.Section)
		}
	}
	return ""
}

func LowerTabMatchesPath(requestPath string, tabContent []LowerTabContent) bool {
	for _, tc := range tabContent {
		if tc.Path == requestPath {
			return true
		}
		if len(tc.Section) > 0 {
			if LowerTabMatchesPath(requestPath, tc.Section) {
				return true
			}
		}
	}
	return false
}

func GetLeftNav(requestPath string, book Book) string {
	var currentUpperTab UpperTab
	for _, upperTab := range book.UpperTabs {
		if upperTab.Path != "" && strings.HasPrefix(requestPath, upperTab.Path) {
			currentUpperTab = upperTab
			break
		}
	}
	if currentUpperTab.Path == "" {
		return ""
	}
	for _, lt := range currentUpperTab.LowerTabs.Other {
		if len(lt.Contents) == 1 && len(lt.Contents[0].Section) == 0 {
			continue
		}
		if LowerTabMatchesPath(requestPath, lt.Contents) {
			result := "<ul class=\"devsite-nav-list devsite-nav-expandable\">\n"
			result += BuildLeftNav(requestPath, lt.Contents)
			result += "</ul>\n"
			return result
		}

	}
	return ""
}

func BuildLeftNav(requestPath string, ltc []LowerTabContent) string {
	result := ""
	for _, item := range ltc {
		if item.Path != "" {
			itemClass := "devsite-nav-item"
			if item.Status != "" {
				itemClass += " devsite-nav-has-status devsite-nav-" + item.Status
			}
			if strings.HasPrefix(requestPath, item.Path) && strings.Count(requestPath, "/") == strings.Count(item.Path, "/") {
				itemClass += " devsite-nav-active"
			}
			result += "<li class=\"" + itemClass + "\">\n"
			result += "<a href=\"" + item.Path + "\" class=\"devsite-nav-title\">\n"
			result += "<span class=\"devsite-nav-text\">"
			result += "<span>" + html.EscapeString(item.Title) + "</span>\n"
			result += "</span>"
			if item.Status != "" {
				result += "<span class=\"devsite-nav-icon-wrapper\">"
				result += "<span class=\"devsite-nav-icon material-icons\"></span>"
				result += "</span>"
			}
			result += "</a>\n"
			result += "</li>\n"
		} else if item.Heading != "" {
			result += "<li class=\"devsite-nav-item devsite-nav-item-heading\">\n"
			result += "<span class=\"devsite-nav-title devsite-nav-title-no-path\">\n"
			result += "</span>\n</li>\n"
		} else if len(item.Section) > 0 {
			itemClass := "devsite-nav-item devsite-nav-item-section-expandable x"
			if item.Style != "" {
				itemClass += " devsite-nav-accordion"
			}
			if item.Status != "" {
				itemClass += " devsite-nav-has-status devsite-nav-" + item.Status
			}
			result += "<li class=\"" + itemClass + "\">\n"
			result += "<span class=\"devsite-nav-title devsite-nav-title-no-path\">\n"
			result += "<span>" + html.EscapeString(item.Title) + "</span>\n"
			if item.Status != "" {
				result += "<span class=\"devsite-nav-icon-wrapper\">"
				result += "<span class=\"devsite-nav-icon material-icons\"></span>"
				result += "</span>"
			}
			result += "</span>"
			result += "<a class=\"devsite-nav-toggle devsite-nav-toggle-collapsed material-icons\">\n"
			result += "</a>"
			result += "<ul class=\"devsite-nav-section devsite-nav-section-collapsed\">\n"
			result += BuildLeftNav(requestPath, item.Section)
			result += "</ul>\n"
		}
	}
	return result
}
