package controllers

import (
	interfaces2 "gitlab.com/aoterocom/changelog-guardian/controller/interfaces"
	middleware "gitlab.com/aoterocom/changelog-guardian/controller/pipes/release"
	middleware2 "gitlab.com/aoterocom/changelog-guardian/controller/pipes/tasks"
	infra "gitlab.com/aoterocom/changelog-guardian/infrastructure/models"
	"reflect"
	"testing"
)

func TestChangelogGuardianController_throughReleasePipes(t *testing.T) {
	type fields struct {
		releasePipes []interfaces2.ReleasePipe
	}
	type args struct {
		releases []infra.Release
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []infra.Release
	}{
		{
			name: "Test release pipeing",
			fields: fields{releasePipes: []interfaces2.ReleasePipe{
				interfaces2.ReleasePipe(&middleware.SemverReleasePipe{}),
			}},
			args: args{
				releases: []infra.Release{
					{Name: "NO-semver"},
					{Name: "v1.2.3"},
					{Name: "1.1.2-prerelease+meta"},
				},
			},
			want: []infra.Release{
				{Name: "v1.2.3"},
				{Name: "1.1.2-prerelease+meta"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cgc := &ChangelogGuardianController{
				releasePipes: tt.fields.releasePipes,
			}
			if got := cgc.throughReleasePipes(tt.args.releases); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("throughReleasePipes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChangelogGuardianController_throughTasksPipes(t *testing.T) {
	type fields struct {
		tasksPipes []interfaces2.TasksPipe
	}
	type args struct {
		tasks []infra.Task
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []infra.Task
	}{
		{
			name: "Test task pipeing",
			fields: fields{tasksPipes: []interfaces2.TasksPipe{
				// Need to go in reverse order!!
				interfaces2.TasksPipe(&middleware2.NaturalLanguageTasksPipe{}),
				interfaces2.TasksPipe(&middleware2.GitlabResolverTasksPipe{}),
			}},
			args: args{
				tasks: []infra.Task{
					{Title: "Resolve \"Add new feature\""},
					{Title: "Add new feature"},
					{Title: "Resolve \"This new feature\""},
					{Title: "No changes"},
				},
			},
			want: []infra.Task{
				{Title: "Added new feature"},
				{Title: "Added new feature"},
				{Title: "This new feature"},
				{Title: "No changes"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cgc := &ChangelogGuardianController{
				tasksPipes: tt.fields.tasksPipes,
			}
			if got := cgc.throughTasksPipes(tt.args.tasks); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("throughTasksPipes() = %v, want %v", got, tt.want)
			}
		})
	}
}
