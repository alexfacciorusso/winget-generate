package githuburl

import "testing"

func TestDestructureRepoURL(t *testing.T) {
	type args struct {
		repoURL string
	}
	tests := []struct {
		name         string
		args         args
		wantValid    bool
		wantOwner    string
		wantRepoName string
	}{
		{
			name:         "Valid url",
			args:         args{"https://github.com/alexfacciorusso/winget-generate"},
			wantValid:    true,
			wantOwner:    "alexfacciorusso",
			wantRepoName: "winget-generate",
		},
		{
			name:         "Valid url with extra slash",
			args:         args{"https://github.com/alexfacciorusso/winget-generate/other"},
			wantValid:    true,
			wantOwner:    "alexfacciorusso",
			wantRepoName: "winget-generate",
		},
		{
			name:         "Valid url with extra hash etc",
			args:         args{"https://github.com/alexfacciorusso/winget-generate#something"},
			wantValid:    true,
			wantOwner:    "alexfacciorusso",
			wantRepoName: "winget-generate",
		},
		{
			name:         "Valid url - no repo",
			args:         args{"https://github.com/alexfacciorusso/"},
			wantValid:    true,
			wantOwner:    "alexfacciorusso",
			wantRepoName: "",
		},
		{
			name:         "Invalid url - not github",
			args:         args{"https://test.com/alexfacciorusso/"},
			wantValid:    false,
			wantOwner:    "",
			wantRepoName: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValid, gotUsername, gotRepoName := DestructureRepoURL(tt.args.repoURL)
			if gotValid != tt.wantValid {
				t.Errorf("DestructureRepoURL() gotValid = %v, want %v", gotValid, tt.wantValid)
			}
			if gotUsername != tt.wantOwner {
				t.Errorf("DestructureRepoURL() gotUsername = %v, want %v", gotUsername, tt.wantOwner)
			}
			if gotRepoName != tt.wantRepoName {
				t.Errorf("DestructureRepoURL() gotRepoName = %v, want %v", gotRepoName, tt.wantRepoName)
			}
		})
	}
}
