package pipes

import (
	"gitlab.com/aoterocom/changelog-guardian/application/models"
	infra "gitlab.com/aoterocom/changelog-guardian/infrastructure/models"
	"reflect"
	"testing"
)

func TestConventionalCommitsTasksPipe_Filter(t *testing.T) {
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
			name:    "Conventional commit feat (Added)",
			args:    args{task: &infra.Task{Title: "feat: feature cc message"}},
			want:    &infra.Task{Title: "feat: feature cc message", Category: models.ADDED},
			want1:   true,
			wantErr: false,
		},
		{
			name:    "Conventional commit feat (Fixed)",
			args:    args{task: &infra.Task{Title: "fix: feature cc message"}},
			want:    &infra.Task{Title: "fix: feature cc message", Category: models.FIXED},
			want1:   true,
			wantErr: false,
		},
		{
			name:    "Not Conventional commit",
			args:    args{task: &infra.Task{Title: "Non cc feature"}},
			want:    nil,
			want1:   true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tf := &ConventionalCommitsTasksPipe{}
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
