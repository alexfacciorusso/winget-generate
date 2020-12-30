package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/AlecAivazis/survey/v2"
	"github.com/alexfacciorusso/winget-generate/debug"
	"github.com/alexfacciorusso/winget-generate/style"
	"github.com/alexfacciorusso/winget-generate/suggestion"
	"github.com/fatih/color"
	"gopkg.in/yaml.v2"
)

type manifestData struct {
	Name       string `yaml:"Name"`
	Publisher  string `yaml:"Publisher"`
	ID         string `yaml:"Id"`
	License    string `yaml:"License"`
	LicenseURL string `yaml:"LicenseUrl"`
	Version    string `yaml:"Version"`
}

func main() {
	verbose := flag.Bool("-v", false, "")
	flag.Parse()

	if *verbose {
		log.SetOutput(debug.DebugWriter)
	} else {
		log.SetOutput(ioutil.Discard)
	}

	var manifest manifestData

	fmt.Println("Hello! We are going to generate a manifest for winget.")

	// region GitHub
	var githubURL string
	survey.AskOne(&survey.Input{
		Message: fmt.Sprintf("If your app is on GitHub, insert its %s now, or leave empty otherwise:", style.QuestionElement("URL")),
	}, &githubURL)

	suggestions := suggestion.GetSuggestionsForRepo(githubURL)
	// endregion

	// region Name and Publisher
	prompt := []*survey.Question{
		{
			Name: "name",
			Prompt: &survey.Input{
				Message: fmt.Sprintf("What's the %s of the app? [E.g. Telegram]", style.QuestionElement("name")),
				Default: suggestions.Name,
			},
		},
		{
			Name: "publisher",
			Prompt: &survey.Input{
				Message: fmt.Sprintf("Who is the %s of the app? E.g. Telegram Messenger Inc.", style.QuestionElement("publisher")),
				Default: suggestions.Publisher,
			},
		},
	}

	err := survey.Ask(prompt, &manifest, getIconsOpt(), survey.WithValidator(survey.Required))
	// endregion

	// region ID
	err = survey.AskOne(&survey.Input{
		Message: fmt.Sprintf("What's the %s of the package?", style.QuestionElement("ID")),
		Default: fmt.Sprintf("%s.%s", manifest.Publisher, manifest.Name),
	}, &manifest.ID, getIconsOpt(), survey.WithValidator(survey.Required))

	if err != nil {
		log.Fatal(err)
	}
	// endregion

	debug.PrintJSON("All licenses", suggestions.LicenseList)
	debug.PrintJSON("Repo license", suggestions.License)

	// region License
	var licenseResponseIndex int
	survey.AskOne(&survey.Select{
		Message: fmt.Sprintf("Which %s does your project use?", style.QuestionElement("license")),
		Options: suggestions.GetLicenseNames(),
		Default: 0,
	}, &licenseResponseIndex, getIconsOpt(), survey.WithValidator(survey.Required))

	selectedLicense := suggestions.LicenseList[licenseResponseIndex]
	manifest.License = selectedLicense.GetName()
	manifest.LicenseURL = selectedLicense.GetURL()
	// endregion

	// region Version
	survey.AskOne(&survey.Input{
		Message: fmt.Sprintf("What's the %s of your project the manifest shall point to?", style.QuestionElement("version")),
		Default: suggestions.Version,
	}, &manifest.Version, getIconsOpt(), survey.WithValidator(survey.Required))
	// endregion

	writeManifest(&manifest)
}

func writeManifest(manifest *manifestData) {
	log.Printf("Data: %+v\n", &manifest)

	yamlContent, err := yaml.Marshal(&manifest)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("\nMarshaled yaml", string(yamlContent))
	dir := fmt.Sprintf("%s/%s", manifest.Publisher, manifest.Name)
	filename := fmt.Sprintf("%s-%s.yaml", manifest.Name, manifest.Version)
	fullPath := path.Join(dir, filename)

	os.MkdirAll(dir, os.ModePerm)
	f, err := os.Create(fullPath)

	if err != nil {
		log.Fatal("Error in opening file for writing: ", err)
	}

	_, err = f.Write(yamlContent)
	if err != nil {
		log.Fatal("Error in writing file: ", err)
	}

	fmt.Fprintf(color.Output, "The manifest file has been saved to %s", color.GreenString(fullPath))
}

func getIconsOpt() survey.AskOpt {
	return survey.WithIcons(func(is *survey.IconSet) {
		is.Question.Text = "[?]"
	})
}
