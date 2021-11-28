package pipes

import (
	infra "gitlab.com/aoterocom/changelog-guardian/infrastructure/models"
	"reflect"
	"testing"
)

func TestNaturalLanguageTasksPipe_Filter(t *testing.T) {
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
			name:    "Test Natural Language pipe with an accepted task",
			args:    args{task: &infra.Task{Title: "Add new feature"}},
			want:    &infra.Task{Title: "Added new feature"},
			want1:   true,
			wantErr: false,
		},
		{
			name:    "Test Natural Language pipe with a non accepted task",
			args:    args{task: &infra.Task{Title: "Rare verb new feature"}},
			want:    &infra.Task{Title: "Rare verb new feature"},
			want1:   false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nlm := &NaturalLanguageTasksPipe{}
			got, got1, err := nlm.Filter(tt.args.task)
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
