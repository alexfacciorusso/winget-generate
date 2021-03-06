package suggestion

import (
	"context"
	"log"

	"github.com/alexfacciorusso/winget-generate/debug"
	"github.com/alexfacciorusso/winget-generate/slices"
	"github.com/google/go-github/v33/github"
)

func fillLicenses(ctx context.Context, githubClient *github.Client, suggestions *RepoSuggestions) {
	githubLicenses, _, _ := githubClient.Licenses.List(ctx)

	if githubLicenses != nil {
		githubLicenses = orderLicenses(githubLicenses, suggestions.License)
	}

	suggestions.LicenseList = githubLicenses
}

func getLicenseNames(licenses []*github.License) []string {
	var licenseNames = make([]string, 0, len(licenses))
	for _, v := range licenses {
		licenseNames = append(licenseNames, *v.Name)
	}
	return licenseNames
}

func getLicenseKeys(licenses []*github.License) []string {
	var licenseKeys = make([]string, 0, len(licenses))
	for _, v := range licenses {
		licenseKeys = append(licenseKeys, *v.Key)
	}
	return licenseKeys
}

func orderLicenses(licenses []*github.License, projectLicense *github.License) []*github.License {
	debug.PrintJSON("Ordering Licenses", licenses)

	if projectLicense == nil {
		return licenses
	}

	licenseMap := make(map[string]*github.License, 0)
	for _, v := range licenses {
		licenseMap[*v.Key] = v
	}

	licenseKeys := getLicenseKeys(licenses)
	log.Printf("License keys: %v", licenseKeys)

	orderedKeys := slices.ElementToFirst(licenseKeys, *projectLicense.Key)
	log.Printf("Ordered keys %v", orderedKeys)

	orderedLicenses := make([]*github.License, 0, len(licenses))
	for _, v := range orderedKeys {
		orderedLicenses = append(orderedLicenses, licenseMap[v])
	}
	return orderedLicenses
}
