package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/alexfacciorusso/winget-generate/debug"
	"github.com/alexfacciorusso/winget-generate/suggestion"
	"github.com/fatih/color"
	"gopkg.in/yaml.v2"
)

type manifestData struct {
	Name      string `yaml:"Name"`
	Publisher string `yaml:"Publisher"`
	ID        string `yaml:"Id"`
}

func getIcons(is *survey.IconSet) {
	is.Question.Text = "[?]"
}

func main() {
	log.SetOutput(debug.DebugWriter)

	var manifest manifestData

	fmt.Println("Hello! We are going to generate a manifest for winget.")

	var githubURL string
	survey.AskOne(&survey.Input{
		Message: fmt.Sprintf("If your app is on GitHub, insert its url now, or leave empty otherwise:"),
	}, &githubURL)

	suggestions := suggestion.GetSuggestionsForRepo(githubURL)

	prompt := []*survey.Question{
		{
			Name: "name",
			Prompt: &survey.Input{
				Message: fmt.Sprintf("What's the %s of the app? [E.g. Telegram]", color.HiBlueString("name")),
				Default: suggestions.Name,
			},
		},
		{
			Name: "publisher",
			Prompt: &survey.Input{
				Message: fmt.Sprintf("Who is the %s of the app? E.g. Telegram Messenger Inc.", "publisher"),
				Default: suggestions.Publisher,
			},
		},
	}

	err := survey.Ask(prompt, &manifest, survey.WithIcons(getIcons), survey.WithValidator(survey.Required))

	titleStrings := func(str string) string {
		return strings.ReplaceAll(strings.Title(str), " ", "")
	}

	err = survey.AskOne(&survey.Input{
		Message: fmt.Sprintf("What's the %s of the package?", "ID"),
		Default: fmt.Sprintf("%s.%s", titleStrings(manifest.Publisher), titleStrings(manifest.Name)),
	}, &manifest.ID, survey.WithIcons(getIcons), survey.WithValidator(survey.Required))

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Data: %+v\n", &manifest)

	it, err := yaml.Marshal(&manifest)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(it))
	ioutil.WriteFile("manifest.yaml", it, 0)
}
