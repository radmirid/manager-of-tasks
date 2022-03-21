package service

import (
	task "github.com/radmirid/manager-of-tasks"
	"github.com/radmirid/manager-of-tasks/pkg/repository"
)

type TaskListService struct {
	repo repository.TaskList
}

func NewTaskListService(repo repository.TaskList) *TaskListService {
	return &TaskListService{repo: repo}
}

func (s *TaskListService) Create(userID int, list task.TaskList) (int, error) {
	return s.repo.Create(userID, list)
}

func (s *TaskListService) GetAll(userID int) ([]task.TaskList, error) {
	return s.repo.GetAll(userID)
}

func (s *TaskListService) GetByID(userID, listID int) (task.TaskList, error) {
	return s.repo.GetByID(userID, listID)
}

func (s *TaskListService) Delete(userID, listID int) error {
	return s.repo.Delete(userID, listID)
}

func (s *TaskListService) Update(userID, listID int, input task.UpdateListInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	return s.repo.Update(userID, listID, input)
}
