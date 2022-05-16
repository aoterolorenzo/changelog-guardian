package mock

import (
	application "gitlab.com/aoterocom/changelog-guardian/application/models"
	infrastructure "gitlab.com/aoterocom/changelog-guardian/infrastructure/models"
	"time"
)

type MockProvider struct {
	getTaskResponse     GetTaskResponse
	getTasksResponse    GetTasksResponse
	getReleasesResponse GetReleasesResponse
}

type GetTaskResponse struct {
	Task  *infrastructure.Task
	Error error
}

type GetTasksResponse struct {
	Tasks *[]infrastructure.Task
	Error error
}

type GetReleasesResponse struct {
	Release *[]infrastructure.Release
	Error   error
}

func NewMockProvider(getTaskResponse GetTaskResponse, getTasksResponse GetTasksResponse,
	getReleasesResponse GetReleasesResponse) *MockProvider {
	return &MockProvider{
		getTaskResponse:     getTaskResponse,
		getTasksResponse:    getTasksResponse,
		getReleasesResponse: getReleasesResponse,
	}
}

func (gc *MockProvider) GetReleases(from *time.Time, to *time.Time) (*[]infrastructure.Release, error) {
	return gc.getReleasesResponse.Release, gc.getReleasesResponse.Error
}

func (gc *MockProvider) GetTasks(from *time.Time, to *time.Time, targetBranch string) (*[]infrastructure.Task, error) {
	return gc.getTasksResponse.Tasks, gc.getTasksResponse.Error
}

func (gc *MockProvider) DefineCategory(task infrastructure.Task) application.Category {
	return application.ADDED
}

func (gc *MockProvider) GetTask(taskId string) (*infrastructure.Task, error) {
	return gc.getTaskResponse.Task, gc.getTaskResponse.Error
}

func (gc *MockProvider) ReleaseURL(from *string, to string) (*string, error) {
	url := ""
	return &url, nil
}
