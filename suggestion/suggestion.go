package suggestion

import (
	"context"
	"log"
	"regexp"

	"github.com/alexfacciorusso/winget-generate/debug"
	"github.com/alexfacciorusso/winget-generate/githuburl"
	"github.com/fatih/color"
	"github.com/google/go-github/v33/github"
)

// RepoSuggestions contains the data retrieved from GitHub.
type RepoSuggestions struct {
	Name        string
	Publisher   string
	License     *github.License
	LicenseList []*github.License
	Version     string
}

// GetLicenseNames returns all the possible license names.
func (rs RepoSuggestions) GetLicenseNames() []string {
	return getLicenseNames(rs.LicenseList)
}

// GetSuggestionsForRepo gets the license associated with the given GitHub repository.
func GetSuggestionsForRepo(repoURL string) *RepoSuggestions {

	suggestions := &RepoSuggestions{}

	ctx := context.Background()
	githubClient := github.NewClient(nil)

	hasRepoInfo, owner, repoName := fillRepoInfo(ctx, repoURL, githubClient, suggestions)
	fillLicenses(ctx, githubClient, suggestions)

	if hasRepoInfo {
		fillVersion(ctx, githubClient, owner, repoName, suggestions)
	}

	return suggestions
}

func fillRepoInfo(ctx context.Context, repoURL string, githubClient *github.Client, suggestions *RepoSuggestions) (hasRepoInfo bool, owner string, repoName string) {
	hasRepoInfo, owner, repoName = githuburl.DestructureRepoURL(repoURL)

	if !hasRepoInfo {
		log.Printf("The inserted repo %s is not a valid github repo\n", repoURL)
		return
	}

	log.Println("Github user: ", owner, "and repo", repoName)

	repo, _, err := githubClient.Repositories.Get(ctx, owner, repoName)

	if err != nil {
		return
	}

	debug.PrintJSON("Repository", repo)

	suggestions.Name = repo.GetName()
	suggestions.Publisher = repo.GetOwner().GetLogin()
	suggestions.License = repo.GetLicense()

	return
}

func fillVersion(ctx context.Context, githubClient *github.Client, owner string, repoName string, suggestions *RepoSuggestions) {
	release, _, err := githubClient.Repositories.GetLatestRelease(ctx, owner, repoName)

	if err != nil {
		log.Println("Error in getting latest release: ", err)
		return
	}

	versionRegex := regexp.MustCompile(`([0-9]+)\.([0-9]+)\.([0-9]+)(?:-([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*))?(?:\+[0-9A-Za-z-]+)?`)
	version := versionRegex.FindString(*release.TagName)

	debug.PrintJSON("Latest release", release)
	log.Println("Matched version: ", color.CyanString(version))

	suggestions.Version = version
}
