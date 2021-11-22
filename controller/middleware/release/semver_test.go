package middleware

import (
	infra "gitlab.com/aoterocom/changelog-guardian/infrastructure/models"
	"reflect"
	"testing"
)

func TestSemverReleaseFilter_Filter(t *testing.T) {
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
			name:    "Test Natural Language filter with a non accepted task",
			args:    args{release: &infra.Release{Name: "1.2.0"}},
			want:    &infra.Release{Name: "1.2.0"},
			want1:   false,
			wantErr: false,
		}, {
			name:    "Test Natural Language filter with a non accepted task",
			args:    args{release: &infra.Release{Name: "non-semver"}},
			want:    nil,
			want1:   true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nlm := &SemverReleaseFilter{}
			got, got1, err := nlm.Filter(tt.args.release)
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
