package middleware

import (
	infra "gitlab.com/aoterocom/changelog-guardian/infrastructure/models"
	"reflect"
	"testing"
)

func TestSemverReleasePipe_Pipe(t *testing.T) {
	type args struct {
		release *infra.Release
	}
	tests := []struct {
		name    string
		args    args
		want    *infra.Release
		want1   bool
		wantErr bool
	}{
		{
			name:    "Test Natural Language pipe with a non accepted task",
			args:    args{release: &infra.Release{Name: "1.2.0"}},
			want:    &infra.Release{Name: "1.2.0"},
			want1:   false,
			wantErr: false,
		}, {
			name:    "Test Natural Language pipe with a non accepted task",
			args:    args{release: &infra.Release{Name: "non-semver"}},
			want:    nil,
			want1:   true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nlm := &SemverReleasePipe{}
			got, got1, err := nlm.Pipe(tt.args.release)
			if (err != nil) != tt.wantErr {
				t.Errorf("Pipe() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Pipe() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Pipe() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
