package pipes

import (
	infra "gitlab.com/aoterocom/changelog-guardian/infrastructure/models"
	"reflect"
	"testing"
)

func TestGitlabResolverTasksPipe_Filter(t *testing.T) {
	type args struct {
		task *infra.Task
	}
	tests := []struct {
		name    string
		args    args
		want    *infra.Task
		want1   bool
		wantErr bool
	}{
		{
			name:    "Test Gitlab Resolver pipe with an accepted task",
			args:    args{task: &infra.Task{Title: "Resolve \"Add new feature\""}},
			want:    &infra.Task{Title: "Add new feature"},
			want1:   true,
			wantErr: false,
		},
		{
			name:    "Test Gitlab Resolver pipe with a non accepted task",
			args:    args{task: &infra.Task{Title: "Rare verb new feature"}},
			want:    &infra.Task{Title: "Rare verb new feature"},
			want1:   false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tf := &GitlabResolverTasksPipe{}
			got, got1, err := tf.Filter(tt.args.task)
			if (err != nil) != tt.wantErr {
				t.Errorf("Filter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Filter() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Filter() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
