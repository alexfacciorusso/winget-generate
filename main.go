package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
)

type manifestData struct {
	Name      string
	Publisher string
	ID        string
}

func getIcons(is *survey.IconSet) {
	is.Question.Text = "[?]"
}

func main() {
	var manifest manifestData

	prompt := []*survey.Question{
		{
			Name: "name",
			Prompt: &survey.Input{
				Message: fmt.Sprintf("Hello! We are going to generate a manifest for winget. \nWhat's the %s of the app? [E.g. Telegram]", color.HiBlueString("name")),
			},
		},
		{
			Name: "publisher",
			Prompt: &survey.Input{
				Message: fmt.Sprintf("Who is the %s of the app? E.g. Telegram Messenger Inc.", "publisher"),
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

	fmt.Printf("Data: %+v", &manifest)
}
