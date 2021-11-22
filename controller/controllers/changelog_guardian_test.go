package controllers

import (
	interfaces2 "gitlab.com/aoterocom/changelog-guardian/controller/interfaces"
	middleware "gitlab.com/aoterocom/changelog-guardian/controller/middleware/release"
	middleware2 "gitlab.com/aoterocom/changelog-guardian/controller/middleware/tasks"
	infra "gitlab.com/aoterocom/changelog-guardian/infrastructure/models"
	"reflect"
	"testing"
)

func TestChangelogGuardianController_throughReleaseFilters(t *testing.T) {
	type fields struct {
		releaseFilters []interfaces2.ReleaseFilter
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
			name: "Test release filtering",
			fields: fields{releaseFilters: []interfaces2.ReleaseFilter{
				interfaces2.ReleaseFilter(&middleware.SemverReleaseFilter{}),
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
				releaseFilters: tt.fields.releaseFilters,
			}
			if got := cgc.throughReleaseFilters(tt.args.releases); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("throughReleaseFilters() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChangelogGuardianController_throughTaskFilters(t *testing.T) {
	type fields struct {
		taskFilters []interfaces2.TaskFilter
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
			name: "Test task filtering",
			fields: fields{taskFilters: []interfaces2.TaskFilter{
				// Need to go in reverse order!!
				interfaces2.TaskFilter(&middleware2.NaturalLanguageTaskFilter{}),
				interfaces2.TaskFilter(&middleware2.GitlabResolverTaskFilter{}),
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
				taskFilters: tt.fields.taskFilters,
			}
			if got := cgc.throughTaskFilters(tt.args.tasks); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("throughTaskFilters() = %v, want %v", got, tt.want)
			}
		})
	}
}
