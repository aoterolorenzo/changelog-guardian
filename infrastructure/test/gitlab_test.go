package test

import (
	"fmt"
	"gitlab.com/aoterocom/changelog-guardian/infrastructure/providers"
	"reflect"
	"testing"
	"time"
)

func TestGitlabProvider_GetReleases(t *testing.T) {
	type args struct {
		repo *string
	}

	repoString1 := "https://gitlab.com/maxigaz/gitlab-dark.git"
	repoString2 := "https://gitlab.com/maxXXXXXXigaz/gitlab-dark.git"

	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name:    "Get Releases from maxigaz/gitlab-dark (archived repo)",
			args:    args{repo: &repoString1},
			want:    27,
			wantErr: false,
		},
		{
			name:    "Get Releases from non existent repo",
			args:    args{repo: &repoString2},
			want:    27,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			glc := providers.NewGitlabProvider()
			got, err := glc.GetReleases(tt.args.repo)
			if (err != nil) && tt.wantErr {
				return
			}
			if !reflect.DeepEqual(len(*got), tt.want) {
				t.Errorf("GetReleases() got = %v, want %v", len(*got), tt.want)
			}
		})
	}
}

func TestGitlabProvider_GetTasks(t *testing.T) {
	type args struct {
		from         *time.Time
		to           *time.Time
		repo         *string
		targetBranch string
	}

	layout := "2006-01-02"
	str := "2019-04-15"
	from1, _ := time.Parse(layout, str)
	fmt.Println(from1.Unix())

	str = "2019-04-23"
	to1, _ := time.Parse(layout, str)
	fmt.Println(to1.Unix())

	repoString1 := "https://gitlab.com/maxigaz/gitlab-dark.git"
	//repoString2 := "https://gitlab.com/maxXXXXXXigaz/gitlab-dark.git"

	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name:    "Get Tasks from non maxigaz/gitlab-dark (archived repo)",
			args:    args{repo: &repoString1, targetBranch: "master"},
			want:    11,
			wantErr: true,
		},
		{
			name:    "Get Tasks from non maxigaz/gitlab-dark (archived repo) with boudaries",
			args:    args{repo: &repoString1, targetBranch: "master", from: &from1, to: &to1},
			want:    2,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			glc := providers.NewGitlabProvider()
			got, err := glc.GetTasks(tt.args.from, tt.args.to, tt.args.repo, tt.args.targetBranch)
			if (err != nil) && tt.wantErr {
				return
			}
			if !reflect.DeepEqual(len(*got), tt.want) {
				t.Errorf("GetTasks() got = %v, want %v", len(*got), tt.want)
			}
		})
	}
}
