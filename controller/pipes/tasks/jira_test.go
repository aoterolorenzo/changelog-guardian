package pipes

import (
	"gitlab.com/aoterocom/changelog-guardian/application/models"
	"gitlab.com/aoterocom/changelog-guardian/controller/controllers/providers/mock"
	infra "gitlab.com/aoterocom/changelog-guardian/infrastructure/models"
	"reflect"
	"testing"
)

func TestJiraTasksPipe_Filter(t *testing.T) {
	mockJiraController := services.MockJiraController{}
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
			name:    "Test Jira pipe with an accepted task",
			args:    args{task: &infra.Task{Title: "[TES-1] Done task"}},
			want:    &infra.Task{Title: "TES1", Category: models.ADDED},
			want1:   true,
			wantErr: false,
		},
		{
			name:    "Test Jira pipe with a non accepted task",
			args:    args{task: &infra.Task{Title: "Rare verb new feature"}},
			want:    nil,
			want1:   true,
			wantErr: false,
		},
		{
			name:    "Test Jira pipe with a non accepted task",
			args:    args{task: &infra.Task{Title: "[TES-2] Done task"}},
			want:    &infra.Task{Title: "TES2", Category: models.ADDED},
			want1:   true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tf := &JiraTasksPipe{providerController: mockJiraController}
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
