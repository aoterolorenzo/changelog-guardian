package services

import (
	"gitlab.com/aoterocom/changelog-guardian/application/models"
	"testing"
)

func TestSemVerService_CalculateNextVersion(t *testing.T) {
	type args struct {
		tasks         []models.Category
		versionToBump string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Calculate patch",
			args: args{
				tasks: []models.Category{
					models.FIXED,
					models.DOCUMENTATION,
					models.DEPENDENCIES,
				},
				versionToBump: "1.0.0-alpha.beta.1",
			},
			want: "1.0.1",
		},
		{
			name: "Calculate minor",
			args: args{
				tasks: []models.Category{
					models.FIXED,
					models.FIXED,
					models.DOCUMENTATION,
					models.ADDED,
					models.DEPENDENCIES,
				},
				versionToBump: "1.0.0-alpha.beta.1",
			},
			want: "1.1.0",
		},
		{
			name: "Calculate major",
			args: args{
				tasks: []models.Category{
					models.FIXED,
					models.DOCUMENTATION,
					models.DEPENDENCIES,
					models.BREAKING_CHANGE,
				},
				versionToBump: "1.0.0-alpha.beta.1",
			},
			want: "2.0.0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svs := &SemVerService{}
			if got := svs.CalculateNextVersion(tt.args.tasks, tt.args.versionToBump); got != tt.want {
				t.Errorf("CalculateNextVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}
