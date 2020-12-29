package suggestion

import (
	"context"
	"log"
	"strings"

	"github.com/alexfacciorusso/winget-generate/debug"
	"github.com/alexfacciorusso/winget-generate/githuburl"
	"github.com/alexfacciorusso/winget-generate/slices"
	"github.com/google/go-github/v33/github"
)

// RepoSuggestions contains the data retrieved from GitHub.
type RepoSuggestions struct {
	Name        string
	Publisher   string
	License     *github.License
	LicenseList []*github.License
}

// GetLicenseNames returns all the possible license names.
func (rs RepoSuggestions) GetLicenseNames() []string {
	return getLicenseNames(rs.LicenseList)
}

// GetSuggestionsForRepo gets the license associated with the given GitHub repository.
func GetSuggestionsForRepo(repoURL string) *RepoSuggestions {
	valid, user, repoName := githuburl.DestructureRepoURL(repoURL)

	suggestions := &RepoSuggestions{}

	if !valid {
		log.Println("The inserted repo is not a valid github repo")
		return suggestions
	}

	log.Println("Github user: ", user, "and repo", repoName)

	ctx := context.Background()
	githubClient := github.NewClient(nil)

	repo, _, err := githubClient.Repositories.Get(ctx, user, repoName)

	if err != nil {
		return suggestions
	}

	debug.PrintJSON("Repository", repo)

	suggestions.Name = strings.Title(*repo.Name)
	suggestions.Publisher = strings.Title(*repo.Owner.Login)
	suggestions.License = repo.GetLicense()

	githubLicenses, _, _ := githubClient.Licenses.List(ctx)

	if githubLicenses != nil {
		githubLicenses = orderLicenses(githubLicenses, suggestions.License)
	}

	suggestions.LicenseList = githubLicenses

	return suggestions
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
