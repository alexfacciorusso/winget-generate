package suggestion

import (
	"context"
	"log"
	"regexp"
	"strings"

	"github.com/alexfacciorusso/winget-generate/debug"
	"github.com/google/go-github/v33/github"
)

func getRepoRegex() *regexp.Regexp {
	return regexp.MustCompile(`^(?:https|git)(?:://|@)github\.com[/:]([^/:#]+)/([^/#]*)$`)
}

// destructureRepoURL returns the user and the repo name from a repository url
func destructureRepoURL(repoURL string) (bool, string, string) {
	sanitisedRepoURL := strings.TrimSuffix(repoURL, ".git")
	submatches := getRepoRegex().FindStringSubmatch(sanitisedRepoURL)[1:]
	if len(submatches) < 2 {
		return false, "", ""
	}
	return true, submatches[0], submatches[1]
}

// RepoSuggestions contains the data retrieved from GitHub.
type RepoSuggestions struct {
	Name        string
	Publisher   string
	LicenseName string
}

// GetSuggestionsForRepo gets the license associated with the given GitHub repository.
func GetSuggestionsForRepo(repoURL string) *RepoSuggestions {
	valid, user, repoName := destructureRepoURL(repoURL)

	suggestions := &RepoSuggestions{}

	log.Println("Github user: ", user, "and repo", repoName)

	if !valid {
		return nil
	}

	ctx := context.Background()
	githubClient := github.NewClient(nil)

	repo, _, err := githubClient.Repositories.Get(ctx, user, repoName)

	if err != nil {
		return nil
	}

	debug.PrintJSON("Repository", repo)

	suggestions.Name = strings.Title(*repo.Name)
	suggestions.Publisher = strings.Title(*repo.Owner.Login)

	license := repo.GetLicense()

	if license != nil {
		suggestions.LicenseName = *license.Name
	}

	return suggestions
}
