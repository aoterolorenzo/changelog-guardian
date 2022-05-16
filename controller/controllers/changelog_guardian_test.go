package controllers

import (
	"github.com/pkg/errors"
	"gitlab.com/aoterocom/changelog-guardian/application/models"
	controllerInterfaces "gitlab.com/aoterocom/changelog-guardian/controller/interfaces"
	middleware "gitlab.com/aoterocom/changelog-guardian/controller/pipes/release"
	middleware2 "gitlab.com/aoterocom/changelog-guardian/controller/pipes/tasks"
	infraInterfaces "gitlab.com/aoterocom/changelog-guardian/infrastructure/interfaces"
	infra "gitlab.com/aoterocom/changelog-guardian/infrastructure/models"
	"gitlab.com/aoterocom/changelog-guardian/infrastructure/providers/mock"
	"reflect"
	"testing"
)

func TestChangelogGuardianController_throughReleasePipes(t *testing.T) {
	type fields struct {
		releasePipes []controllerInterfaces.ReleasePipe
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
			fields: fields{releasePipes: []controllerInterfaces.ReleasePipe{
				controllerInterfaces.ReleasePipe(&middleware.SemverReleasePipe{}),
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
		tasksPipes []controllerInterfaces.TasksPipe
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
			fields: fields{tasksPipes: []controllerInterfaces.TasksPipe{
				controllerInterfaces.TasksPipe(&middleware2.GitlabResolverTasksPipe{}),
				controllerInterfaces.TasksPipe(&middleware2.NaturalLanguageTasksPipe{}),
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

func TestChangelogGuardianController_GetTask(t *testing.T) {
	type fields struct {
		releaseProvider infraInterfaces.Provider
		taskProvider    infraInterfaces.Provider
		releasePipes    []controllerInterfaces.ReleasePipe
		tasksPipes      []controllerInterfaces.TasksPipe
	}
	type args struct {
		taskId string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.Task
		wantErr bool
	}{
		{
			name: "simple test getting 1 task",
			fields: fields{
				releaseProvider: nil,
				taskProvider: mock.NewMockProvider(
					struct {
						Task  *infra.Task
						Error error
					}{
						Task:  &infra.Task{ID: "#3"},
						Error: nil,
					},
					struct {
						Tasks *[]infra.Task
						Error error
					}{
						Tasks: nil,
						Error: nil,
					},
					struct {
						Release *[]infra.Release
						Error   error
					}{
						Release: nil,
						Error:   nil},
				),
				releasePipes: nil,
				tasksPipes:   nil,
			},
			args: args{taskId: "#3"},
			want: &models.Task{ID: "#3"},
		},
		{
			name: "simple test getting no task",
			fields: fields{
				releaseProvider: nil,
				taskProvider: mock.NewMockProvider(
					struct {
						Task  *infra.Task
						Error error
					}{
						Task:  nil,
						Error: errors.New("error"),
					},
					struct {
						Tasks *[]infra.Task
						Error error
					}{
						Tasks: nil,
						Error: nil,
					},
					struct {
						Release *[]infra.Release
						Error   error
					}{
						Release: nil,
						Error:   nil},
				),
				releasePipes: nil,
				tasksPipes:   nil,
			},
			args:    args{taskId: "#5"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cgc := &ChangelogGuardianController{
				releaseProvider: tt.fields.releaseProvider,
				taskProvider:    tt.fields.taskProvider,
				releasePipes:    tt.fields.releasePipes,
				tasksPipes:      tt.fields.tasksPipes,
			}
			got, err := cgc.GetTask(tt.args.taskId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTask() got = %v, want %v", got, tt.want)
			}
		})
	}
}
