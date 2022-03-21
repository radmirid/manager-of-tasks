package service

import (
	task "github.com/radmirid/manager-of-tasks"
	"github.com/radmirid/manager-of-tasks/pkg/repository"
)

type TaskItemService struct {
	repo     repository.TaskItem
	listRepo repository.TaskList
}

func NewTaskItemService(repo repository.TaskItem, listRepo repository.TaskList) *TaskItemService {
	return &TaskItemService{repo: repo, listRepo: listRepo}
}

func (s *TaskItemService) Create(userID, listID int, item task.TaskItem) (int, error) {
	_, err := s.listRepo.GetByID(userID, listID)
	if err != nil {
		return 0, err
	}

	return s.repo.Create(listID, item)
}

func (s *TaskItemService) GetAll(userID, listID int) ([]task.TaskItem, error) {
	return s.repo.GetAll(userID, listID)
}

func (s *TaskItemService) GetByID(userID, itemID int) (task.TaskItem, error) {
	return s.repo.GetByID(userID, itemID)
}

func (s *TaskItemService) Delete(userID, itemID int) error {
	return s.repo.Delete(userID, itemID)
}

func (s *TaskItemService) Update(userID, itemID int, input task.UpdateItemInput) error {
	return s.repo.Update(userID, itemID, input)
}
