package repository

import (
	"github.com/jmoiron/sqlx"
	task "github.com/radmirid/manager-of-tasks"
)

type Authorization interface {
	CreateUser(user task.User) (int, error)
	GetUser(username, password string) (task.User, error)
}

type TaskList interface {
	Create(userID int, list task.TaskList) (int, error)
	GetAll(userID int) ([]task.TaskList, error)
	GetByID(userID, listID int) (task.TaskList, error)
	Delete(userID, listID int) error
	Update(userID, listID int, input task.UpdateListInput) error
}

type TaskItem interface {
	Create(listID int, item task.TaskItem) (int, error)
	GetAll(userID, listID int) ([]task.TaskItem, error)
	GetByID(userID, itemID int) (task.TaskItem, error)
	Delete(userID, itemID int) error
	Update(userID, itemID int, input task.UpdateItemInput) error
}

type Repository struct {
	Authorization
	TaskList
	TaskItem
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		TaskList:      NewTaskListPostgres(db),
		TaskItem:      NewTaskItemPostgres(db),
	}
}
