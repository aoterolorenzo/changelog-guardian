package pipes

import (
	settings "gitlab.com/aoterocom/changelog-guardian/config"
	infra "gitlab.com/aoterocom/changelog-guardian/infrastructure/models"
	"reflect"
	"testing"
)

func TestInclusionsExclusionsTasksPipe_Filter(t *testing.T) {

	type args struct {
		task            *infra.Task
		labelInclusions []string
		labelExclusions []string
		pathInclusions  []string
		pathExclusions  []string
	}
	tests := []struct {
		name    string
		args    args
		want    *infra.Task
		want1   bool
		wantErr bool
	}{
		{
			name: "All paths accepted",
			args: args{
				task: &infra.Task{
					Title:  "fix: feature cc message",
					Files:  []string{"/path/tso/file.txt"},
					Labels: []string{"kind::feature"},
				},
				labelInclusions: []string{"*all"},
				labelExclusions: []string{"internal"},
				pathInclusions:  []string{"*all"},
				pathExclusions:  []string{"/path/to/"},
			},
			want: &infra.Task{
				Title:  "fix: feature cc message",
				Files:  []string{"/path/tso/file.txt"},
				Labels: []string{"kind::feature"},
			},
			want1:   false,
			wantErr: false,
		},
		{
			name: "Path exclusion",
			args: args{
				task: &infra.Task{
					Title:  "fix: feature cc message",
					Files:  []string{"/path/tso/file.txt", "/path/to/file.txt"},
					Labels: []string{"kind::feature"},
				},
				labelInclusions: []string{"*all"},
				labelExclusions: []string{"internal"},
				pathInclusions:  []string{"*all"},
				pathExclusions:  []string{"/path/to/"},
			},
			want:    nil,
			want1:   true,
			wantErr: false,
		},
		{
			name: "Path included",
			args: args{
				task: &infra.Task{
					Title:  "fix: feature cc message",
					Files:  []string{"/path/tso/file.txt", "/otherpath/file.txt"},
					Labels: []string{"kind::feature"},
				},
				labelInclusions: []string{"*all"},
				labelExclusions: []string{"internal"},
				pathInclusions:  []string{"/otherpath/"},
				pathExclusions:  []string{"/path/to/"},
			},
			want: &infra.Task{
				Title:  "fix: feature cc message",
				Files:  []string{"/path/tso/file.txt", "/otherpath/file.txt"},
				Labels: []string{"kind::feature"},
			},
			want1:   false,
			wantErr: false,
		},
		{
			name: "Path not included",
			args: args{
				task: &infra.Task{
					Title:  "fix: feature cc message",
					Files:  []string{"/path/tso/file.txt", "/otherpath/file.txt"},
					Labels: []string{"kind::feature"},
				},
				labelInclusions: []string{"*all"},
				labelExclusions: []string{"internal"},
				pathInclusions:  []string{"/nop/"},
				pathExclusions:  []string{"/path/to/"},
			},
			want:    nil,
			want1:   true,
			wantErr: false,
		},
		{
			name: "All allowed but excluded path",
			args: args{
				task: &infra.Task{
					Title:  "fix: feature cc message",
					Files:  []string{"/path/to/"},
					Labels: []string{"kind::feature"},
				},
				labelInclusions: []string{"*all"},
				labelExclusions: []string{"internal"},
				pathInclusions:  []string{"*all"},
				pathExclusions:  []string{"/path/to/"},
			},
			want:    nil,
			want1:   true,
			wantErr: false,
		},
		{
			name: "All labels accepted",
			args: args{
				task: &infra.Task{
					Title:  "fix: feature cc message",
					Files:  []string{"/path/tso/file.txt"},
					Labels: []string{"kind::feature"},
				},
				labelInclusions: []string{"*all"},
				labelExclusions: []string{"internal"},
				pathInclusions:  []string{"*all"},
				pathExclusions:  []string{""},
			},
			want: &infra.Task{
				Title:  "fix: feature cc message",
				Files:  []string{"/path/tso/file.txt"},
				Labels: []string{"kind::feature"},
			},
			want1:   false,
			wantErr: false,
		},
		{
			name: "Label exclusion",
			args: args{
				task: &infra.Task{
					Title:  "fix: feature cc message",
					Files:  []string{"/path/tso/file.txt", "/path/to/file.txt"},
					Labels: []string{"kind::feature", "internal"},
				},
				labelInclusions: []string{"*all"},
				labelExclusions: []string{"internal"},
				pathInclusions:  []string{"*all"},
				pathExclusions:  []string{""},
			},
			want:    nil,
			want1:   true,
			wantErr: false,
		},
		{
			name: "Label included",
			args: args{
				task: &infra.Task{
					Title:  "fix: feature cc message",
					Files:  []string{"/path/tso/file.txt", "/otherpath/file.txt"},
					Labels: []string{"kind::feature", "hey-label"},
				},
				labelInclusions: []string{"hey-label"},
				labelExclusions: []string{"internal"},
				pathInclusions:  []string{"/otherpath/"},
				pathExclusions:  []string{"/path/to/"},
			},
			want: &infra.Task{
				Title:  "fix: feature cc message",
				Files:  []string{"/path/tso/file.txt", "/otherpath/file.txt"},
				Labels: []string{"kind::feature", "hey-label"},
			},
			want1:   false,
			wantErr: false,
		},
		{
			name: "Label not included",
			args: args{
				task: &infra.Task{
					Title:  "fix: feature cc message",
					Files:  []string{"/path/tso/file.txt", "/otherpath/file.txt"},
					Labels: []string{"kind::feature"},
				},
				labelInclusions: []string{"other-label"},
				labelExclusions: []string{"internal"},
				pathInclusions:  []string{"*all"},
				pathExclusions:  []string{""},
			},
			want:    nil,
			want1:   true,
			wantErr: false,
		},
		{
			name: "All allowed but excluded label",
			args: args{
				task: &infra.Task{
					Title:  "fix: feature cc message",
					Files:  []string{"/path/to/"},
					Labels: []string{"kind::feature", "other", "internal"},
				},
				labelInclusions: []string{"*all"},
				labelExclusions: []string{"internal"},
				pathInclusions:  []string{"*all"},
				pathExclusions:  []string{""},
			},
			want:    nil,
			want1:   true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tf := &InclusionsExclusionsTasksPipe{}
			settings.Settings.TasksPipesCfg.InclusionsExclusions.Paths.Inclusions = tt.args.pathInclusions
			settings.Settings.TasksPipesCfg.InclusionsExclusions.Paths.Exclusions = tt.args.pathExclusions
			settings.Settings.TasksPipesCfg.InclusionsExclusions.Labels.Inclusions = tt.args.labelInclusions
			settings.Settings.TasksPipesCfg.InclusionsExclusions.Labels.Exclusions = tt.args.labelExclusions

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
