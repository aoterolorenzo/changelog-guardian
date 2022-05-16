package services

import (
	"gitlab.com/aoterocom/changelog-guardian/application/models"
	"reflect"
	"testing"
)

func TestChangelogMixer_MergeChangelogs(t *testing.T) {
	type args struct {
		changelog1 models.Changelog
		changelog2 models.Changelog
	}
	tests := []struct {
		name string
		args args
		want models.Changelog
	}{
		{
			name: "merging non-empty changelogs",
			args: args{
				changelog1: models.Changelog{
					Releases: []models.Release{
						{
							Date:     "",
							Version:  "Unreleased",
							Sections: nil,
						},
						{
							Date:    "01-01-2022",
							Version: "4.0",
							Sections: map[models.Category][]models.Task{
								models.ADDED: {
									models.Task{
										ID:         "#1",
										Name:       "Task 1",
										Href:       "",
										Title:      "Task 1",
										Author:     "aoterocom",
										AuthorHref: "/aoterocom",
										Category:   "ADDED",
									},
								},

								models.FIXED: {
									models.Task{
										ID:         "#3",
										Name:       "Task 3",
										Href:       "",
										Title:      "Task 3",
										Author:     "aoterocom",
										AuthorHref: "/aoterocom",
										Category:   "FIXED",
									},
									models.Task{
										ID:         "#4",
										Name:       "Task 4",
										Href:       "",
										Title:      "Task 4",
										Author:     "aoterocom",
										AuthorHref: "/aoterocom",
										Category:   "FIXED",
									},
								},
							},
						},
						{
							Date:    "01-01-2021",
							Version: "2.0",
							Sections: map[models.Category][]models.Task{
								models.ADDED: {
									models.Task{
										ID:         "#1",
										Name:       "Task 1",
										Href:       "",
										Title:      "Task 1",
										Author:     "aoterocom",
										AuthorHref: "/aoterocom",
										Category:   "ADDED",
									},
								},

								models.FIXED: {
									models.Task{
										ID:         "#3",
										Name:       "Task 3",
										Href:       "",
										Title:      "Task 3",
										Author:     "aoterocom",
										AuthorHref: "/aoterocom",
										Category:   "FIXED",
									},
									models.Task{
										ID:         "#4",
										Name:       "Task 4",
										Href:       "",
										Title:      "Task 4",
										Author:     "aoterocom",
										AuthorHref: "/aoterocom",
										Category:   "FIXED",
									},
								},
							},
						},
						{
							Date:    "01-09-2020",
							Version: "1.0",
							Sections: map[models.Category][]models.Task{
								models.ADDED: {
									models.Task{
										ID:         "#5",
										Name:       "Task 5",
										Href:       "",
										Title:      "Task 5",
										Author:     "aoterocom",
										AuthorHref: "/aoterocom",
										Category:   "ADDED",
									},
								},

								models.FIXED: {
									models.Task{
										ID:         "#6",
										Name:       "Task 6",
										Href:       "",
										Title:      "Task 6",
										Author:     "aoterocom",
										AuthorHref: "/aoterocom",
										Category:   "FIXED",
									},
									models.Task{
										ID:         "#7",
										Name:       "Task 7",
										Href:       "",
										Title:      "Task 7",
										Author:     "aoterocom",
										AuthorHref: "/aoterocom",
										Category:   "FIXED",
									},
								},
							},
						},
					},
				},
				changelog2: models.Changelog{
					Releases: []models.Release{
						{
							Date:    "01-03-2021",
							Version: "3.0",
							Sections: map[models.Category][]models.Task{
								models.ADDED: {
									models.Task{
										ID:         "#5",
										Name:       "Task 5",
										Href:       "",
										Title:      "Task 5",
										Author:     "aoterocom",
										AuthorHref: "/aoterocom",
										Category:   "ADDED",
									},
								},
								models.FIXED: {
									models.Task{
										ID:         "#6",
										Name:       "Task 6",
										Href:       "",
										Title:      "Task 6",
										Author:     "aoterocom",
										AuthorHref: "/aoterocom",
										Category:   "FIXED",
									},
									models.Task{
										ID:         "#7",
										Name:       "Task 7",
										Href:       "",
										Title:      "Task 7",
										Author:     "aoterocom",
										AuthorHref: "/aoterocom",
										Category:   "FIXED",
									},
								},
							},
						},
						{
							Date:    "01-01-2021",
							Version: "2.0",
							Sections: map[models.Category][]models.Task{
								models.ADDED: {
									models.Task{
										ID:         "#1",
										Name:       "Task 1",
										Href:       "",
										Title:      "Task 1",
										Author:     "aoterocom",
										AuthorHref: "/aoterocom",
										Category:   "ADDED",
									},
								},

								models.FIXED: {
									models.Task{
										ID:         "#3",
										Name:       "Task 3",
										Href:       "",
										Title:      "Task 3",
										Author:     "aoterocom",
										AuthorHref: "/aoterocom",
										Category:   "FIXED",
									},
									models.Task{
										ID:         "#44",
										Name:       "Task 44",
										Href:       "",
										Title:      "Task 44",
										Author:     "aoterocom",
										AuthorHref: "/aoterocom",
										Category:   "FIXED",
									},
								},
							},
						},
						{
							Date:    "01-09-2020",
							Version: "1.0",
							Sections: map[models.Category][]models.Task{
								models.ADDED: {
									models.Task{
										ID:         "#5",
										Name:       "Task 5",
										Href:       "",
										Title:      "Task 5",
										Author:     "aoterocom",
										AuthorHref: "/aoterocom",
										Category:   "ADDED",
									},
								},

								models.FIXED: {
									models.Task{
										ID:         "#6",
										Name:       "Task 6",
										Href:       "",
										Title:      "Task 6",
										Author:     "aoterocom",
										AuthorHref: "/aoterocom",
										Category:   "FIXED",
									},
									models.Task{
										ID:         "#7",
										Name:       "Task 7",
										Href:       "",
										Title:      "Task 7",
										Author:     "aoterocom",
										AuthorHref: "/aoterocom",
										Category:   "FIXED",
									},
								},
							},
						},
					},
				},
			},
			want: models.Changelog{
				Releases: []models.Release{
					{
						Date:     "",
						Version:  "Unreleased",
						Sections: nil,
					},

					{
						Date:    "01-01-2022",
						Version: "4.0",
						Sections: map[models.Category][]models.Task{
							models.ADDED: {
								models.Task{
									ID:         "#1",
									Name:       "Task 1",
									Href:       "",
									Title:      "Task 1",
									Author:     "aoterocom",
									AuthorHref: "/aoterocom",
									Category:   "ADDED",
								},
							},

							models.FIXED: {
								models.Task{
									ID:         "#3",
									Name:       "Task 3",
									Href:       "",
									Title:      "Task 3",
									Author:     "aoterocom",
									AuthorHref: "/aoterocom",
									Category:   "FIXED",
								},
								models.Task{
									ID:         "#4",
									Name:       "Task 4",
									Href:       "",
									Title:      "Task 4",
									Author:     "aoterocom",
									AuthorHref: "/aoterocom",
									Category:   "FIXED",
								},
							},
						},
					},
					{
						Date:    "01-03-2021",
						Version: "3.0",
						Sections: map[models.Category][]models.Task{
							models.ADDED: {
								models.Task{
									ID:         "#5",
									Name:       "Task 5",
									Href:       "",
									Title:      "Task 5",
									Author:     "aoterocom",
									AuthorHref: "/aoterocom",
									Category:   "ADDED",
								},
							},

							models.FIXED: {
								models.Task{
									ID:         "#6",
									Name:       "Task 6",
									Href:       "",
									Title:      "Task 6",
									Author:     "aoterocom",
									AuthorHref: "/aoterocom",
									Category:   "FIXED",
								},
								models.Task{
									ID:         "#7",
									Name:       "Task 7",
									Href:       "",
									Title:      "Task 7",
									Author:     "aoterocom",
									AuthorHref: "/aoterocom",
									Category:   "FIXED",
								},
							},
						},
					},
					{
						Date:    "01-01-2021",
						Version: "2.0",
						Sections: map[models.Category][]models.Task{
							models.ADDED: {
								models.Task{
									ID:         "#1",
									Name:       "Task 1",
									Href:       "",
									Title:      "Task 1",
									Author:     "aoterocom",
									AuthorHref: "/aoterocom",
									Category:   "ADDED",
								},
							},

							models.FIXED: {
								models.Task{
									ID:         "#4",
									Name:       "Task 4",
									Href:       "",
									Title:      "Task 4",
									Author:     "aoterocom",
									AuthorHref: "/aoterocom",
									Category:   "FIXED",
								},
								models.Task{
									ID:         "#3",
									Name:       "Task 3",
									Href:       "",
									Title:      "Task 3",
									Author:     "aoterocom",
									AuthorHref: "/aoterocom",
									Category:   "FIXED",
								},
								models.Task{
									ID:         "#44",
									Name:       "Task 44",
									Href:       "",
									Title:      "Task 44",
									Author:     "aoterocom",
									AuthorHref: "/aoterocom",
									Category:   "FIXED",
								},
							},
						},
					},

					{
						Date:    "01-09-2020",
						Version: "1.0",
						Sections: map[models.Category][]models.Task{
							models.ADDED: {
								models.Task{
									ID:         "#5",
									Name:       "Task 5",
									Href:       "",
									Title:      "Task 5",
									Author:     "aoterocom",
									AuthorHref: "/aoterocom",
									Category:   "ADDED",
								},
							},

							models.FIXED: {
								models.Task{
									ID:         "#6",
									Name:       "Task 6",
									Href:       "",
									Title:      "Task 6",
									Author:     "aoterocom",
									AuthorHref: "/aoterocom",
									Category:   "FIXED",
								},
								models.Task{
									ID:         "#7",
									Name:       "Task 7",
									Href:       "",
									Title:      "Task 7",
									Author:     "aoterocom",
									AuthorHref: "/aoterocom",
									Category:   "FIXED",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "changelog vs empty changelog",
			args: args{
				changelog1: models.Changelog{
					Releases: []models.Release{
						{
							Date:     "",
							Version:  "Unreleased",
							Sections: nil,
						},
						{
							Date:    "01-01-2021",
							Version: "2.0",
							Sections: map[models.Category][]models.Task{
								models.ADDED: {
									models.Task{
										ID:         "#1",
										Name:       "Task 1",
										Href:       "",
										Title:      "Task 1",
										Author:     "aoterocom",
										AuthorHref: "/aoterocom",
										Category:   "ADDED",
									},
								},

								models.FIXED: {
									models.Task{
										ID:         "#3",
										Name:       "Task 3",
										Href:       "",
										Title:      "Task 3",
										Author:     "aoterocom",
										AuthorHref: "/aoterocom",
										Category:   "FIXED",
									},
									models.Task{
										ID:         "#4",
										Name:       "Task 4",
										Href:       "",
										Title:      "Task 4",
										Author:     "aoterocom",
										AuthorHref: "/aoterocom",
										Category:   "FIXED",
									},
								},
							},
						},
						{
							Date:    "01-09-2020",
							Version: "1.0",
							Sections: map[models.Category][]models.Task{
								models.ADDED: {
									models.Task{
										ID:         "#5",
										Name:       "Task 5",
										Href:       "",
										Title:      "Task 5",
										Author:     "aoterocom",
										AuthorHref: "/aoterocom",
										Category:   "ADDED",
									},
								},

								models.FIXED: {
									models.Task{
										ID:         "#6",
										Name:       "Task 6",
										Href:       "",
										Title:      "Task 6",
										Author:     "aoterocom",
										AuthorHref: "/aoterocom",
										Category:   "FIXED",
									},
									models.Task{
										ID:         "#7",
										Name:       "Task 7",
										Href:       "",
										Title:      "Task 7",
										Author:     "aoterocom",
										AuthorHref: "/aoterocom",
										Category:   "FIXED",
									},
								},
							},
						},
					},
				},
				changelog2: models.Changelog{
					Releases: []models.Release{},
				},
			},
			want: models.Changelog{
				Releases: []models.Release{
					{
						Date:     "",
						Version:  "Unreleased",
						Sections: nil,
					},
					{
						Date:    "01-01-2021",
						Version: "2.0",
						Sections: map[models.Category][]models.Task{
							models.ADDED: {
								models.Task{
									ID:         "#1",
									Name:       "Task 1",
									Href:       "",
									Title:      "Task 1",
									Author:     "aoterocom",
									AuthorHref: "/aoterocom",
									Category:   "ADDED",
								},
							},

							models.FIXED: {
								models.Task{
									ID:         "#3",
									Name:       "Task 3",
									Href:       "",
									Title:      "Task 3",
									Author:     "aoterocom",
									AuthorHref: "/aoterocom",
									Category:   "FIXED",
								},
								models.Task{
									ID:         "#4",
									Name:       "Task 4",
									Href:       "",
									Title:      "Task 4",
									Author:     "aoterocom",
									AuthorHref: "/aoterocom",
									Category:   "FIXED",
								},
							},
						},
					},
					{
						Date:    "01-09-2020",
						Version: "1.0",
						Sections: map[models.Category][]models.Task{
							models.ADDED: {
								models.Task{
									ID:         "#5",
									Name:       "Task 5",
									Href:       "",
									Title:      "Task 5",
									Author:     "aoterocom",
									AuthorHref: "/aoterocom",
									Category:   "ADDED",
								},
							},

							models.FIXED: {
								models.Task{
									ID:         "#6",
									Name:       "Task 6",
									Href:       "",
									Title:      "Task 6",
									Author:     "aoterocom",
									AuthorHref: "/aoterocom",
									Category:   "FIXED",
								},
								models.Task{
									ID:         "#7",
									Name:       "Task 7",
									Href:       "",
									Title:      "Task 7",
									Author:     "aoterocom",
									AuthorHref: "/aoterocom",
									Category:   "FIXED",
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := NewChangelogMixer()
			if got := cm.MergeChangelogs(tt.args.changelog1, tt.args.changelog2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MergeChangelogs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChangelogMixer_ChangelogContainsTask(t *testing.T) {
	type args struct {
		changelog models.Changelog
		task      models.Task
	}
	addedCat := models.ADDED
	tests := []struct {
		name  string
		args  args
		want  *models.Category
		want1 *models.Task
		want2 bool
	}{
		{
			name: "empty changelog",
			args: args{
				changelog: models.Changelog{},
				task: models.Task{
					ID: "1",
				},
			},
			want:  nil,
			want2: false,
		},
		{
			name: "remove task case",
			args: args{
				changelog: models.Changelog{},
				task: models.Task{
					ID:       "1",
					Category: models.REMOVED,
				},
			},
			want:  nil,
			want2: false,
		},

		{
			name: "existing task",
			args: args{
				changelog: models.Changelog{
					Releases: []models.Release{
						{
							Date:    "01-01-2021",
							Version: "2.0",
							Sections: map[models.Category][]models.Task{
								models.ADDED: {
									models.Task{
										ID: "1",
									},
								},
							},
						},
					},
				},
				task: models.Task{
					ID:       "1",
					Category: addedCat,
				},
			},
			want: &addedCat,
			want1: &models.Task{
				ID:       "1",
				Category: addedCat,
			},
			want2: true,
		},
		{
			name: "existing task but removed category",
			args: args{
				changelog: models.Changelog{
					Releases: []models.Release{
						{
							Date:    "01-01-2021",
							Version: "2.0",
							Sections: map[models.Category][]models.Task{
								models.REMOVED: {
									models.Task{
										ID:       "1",
										Category: models.REMOVED,
									},
								},
							},
						},
					},
				},
				task: models.Task{
					ID:       "1",
					Category: addedCat,
				},
			},
			want:  nil,
			want1: nil,
			want2: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &ChangelogMixer{}
			got, got1, got2 := cm.ChangelogContainsTask(tt.args.changelog, tt.args.task)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChangelogContainsTask() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ChangelogContainsTask() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("ChangelogContainsTask() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func TestChangelogMixer_MergeReleases(t *testing.T) {
	type args struct {
		release1 models.Release
		release2 models.Release
	}
	tests := []struct {
		name string
		args args
		want *models.Release
	}{
		{
			name: "normal merge",
			args: args{
				release1: models.Release{
					Date:    "01-09-2020",
					Version: "1.0",
					Sections: map[models.Category][]models.Task{
						models.ADDED: {
							models.Task{
								ID:         "#55",
								Name:       "Task 55",
								Href:       "",
								Title:      "Task 55",
								Author:     "aoterocom",
								AuthorHref: "/aoterocom",
								Category:   "ADDED",
							},
						},

						models.FIXED: {
							models.Task{
								ID:         "#6",
								Name:       "Task 6",
								Href:       "",
								Title:      "Task 6",
								Author:     "aoterocom",
								AuthorHref: "/aoterocom",
								Category:   "FIXED",
							},
							models.Task{
								ID:         "#7",
								Name:       "Task 7",
								Href:       "",
								Title:      "Task 7",
								Author:     "aoterocom",
								AuthorHref: "/aoterocom",
								Category:   "FIXED",
							},
						},
					},
				},
				release2: models.Release{
					Date:    "01-09-2020",
					Version: "1.0",
					Sections: map[models.Category][]models.Task{
						models.ADDED: {
							models.Task{
								ID:         "#5",
								Name:       "Task 5",
								Href:       "",
								Title:      "Task 5",
								Author:     "aoterocom",
								AuthorHref: "/aoterocom",
								Category:   "ADDED",
							},
						},

						models.FIXED: {
							models.Task{
								ID:         "#6",
								Name:       "Task 6",
								Href:       "",
								Title:      "Task 6",
								Author:     "aoterocom",
								AuthorHref: "/aoterocom",
								Category:   "FIXED",
							},
							models.Task{
								ID:         "#8",
								Name:       "Task 8",
								Href:       "",
								Title:      "Task 8",
								Author:     "aoterocom",
								AuthorHref: "/aoterocom",
								Category:   "FIXED",
							},
						},
					},
				},
			},
			want: &models.Release{
				Date:    "01-09-2020",
				Version: "1.0",
				Sections: map[models.Category][]models.Task{
					models.ADDED: {
						models.Task{
							ID:         "#55",
							Name:       "Task 55",
							Href:       "",
							Title:      "Task 55",
							Author:     "aoterocom",
							AuthorHref: "/aoterocom",
							Category:   "ADDED",
						},
						models.Task{
							ID:         "#5",
							Name:       "Task 5",
							Href:       "",
							Title:      "Task 5",
							Author:     "aoterocom",
							AuthorHref: "/aoterocom",
							Category:   "ADDED",
						},
					},

					models.FIXED: {
						models.Task{
							ID:         "#7",
							Name:       "Task 7",
							Href:       "",
							Title:      "Task 7",
							Author:     "aoterocom",
							AuthorHref: "/aoterocom",
							Category:   "FIXED",
						},
						models.Task{
							ID:         "#6",
							Name:       "Task 6",
							Href:       "",
							Title:      "Task 6",
							Author:     "aoterocom",
							AuthorHref: "/aoterocom",
							Category:   "FIXED",
						},
						models.Task{
							ID:         "#8",
							Name:       "Task 8",
							Href:       "",
							Title:      "Task 8",
							Author:     "aoterocom",
							AuthorHref: "/aoterocom",
							Category:   "FIXED",
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &ChangelogMixer{}
			if got := cm.MergeReleases(tt.args.release1, tt.args.release2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MergeReleases() = %v, want %v", got, tt.want)
			}
		})
	}
}
