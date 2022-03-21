package service

import (
	task "github.com/radmirid/manager-of-tasks"
	"github.com/radmirid/manager-of-tasks/pkg/repository"
)

type Authorization interface {
	CreateUser(user task.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

type TaskList interface {
	Create(userID int, list task.TaskList) (int, error)
	GetAll(userID int) ([]task.TaskList, error)
	GetByID(userID, listID int) (task.TaskList, error)
	Delete(userID, listID int) error
	Update(userID, listID int, input task.UpdateListInput) error
}

type TaskItem interface {
	Create(userID, listID int, item task.TaskItem) (int, error)
	GetAll(userID, listID int) ([]task.TaskItem, error)
	GetByID(userID, itemID int) (task.TaskItem, error)
	Delete(userID, itemID int) error
	Update(userID, itemID int, input task.UpdateItemInput) error
}

type Service struct {
	Authorization
	TaskList
	TaskItem
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		TaskList:      NewTaskListService(repos.TaskList),
		TaskItem:      NewTaskItemService(repos.TaskItem, repos.TaskList),
	}
}
