package repository

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	task "github.com/radmirid/manager-of-tasks"
	"github.com/sirupsen/logrus"
)

type TaskListPostgres struct {
	db *sqlx.DB
}

func NewTaskListPostgres(db *sqlx.DB) *TaskListPostgres {
	return &TaskListPostgres{db: db}
}

func (r *TaskListPostgres) Create(userID int, list task.TaskList) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	createListQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id", taskListsTable)
	row := tx.QueryRow(createListQuery, list.Title, list.Description)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	createUsersListQuery := fmt.Sprintf("INSERT INTO %s (user_id, list_id) VALUES ($1, $2)", usersListsTable)
	_, err = tx.Exec(createUsersListQuery, userID, id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}

func (r *TaskListPostgres) GetAll(userID int) ([]task.TaskList, error) {
	var lists []task.TaskList

	query := fmt.Sprintf("SELECT tl.id, tl.title, tl.description FROM %s tl INNER JOIN %s ul on tl.id = ul.list_id WHERE ul.user_id = $1",
		taskListsTable, usersListsTable)
	err := r.db.Select(&lists, query, userID)

	return lists, err
}

func (r *TaskListPostgres) GetByID(userID, listID int) (task.TaskList, error) {
	var list task.TaskList

	query := fmt.Sprintf(`SELECT tl.id, tl.title, tl.description FROM %s tl
								INNER JOIN %s ul on tl.id = ul.list_id WHERE ul.user_id = $1 AND ul.list_id = $2`,
		taskListsTable, usersListsTable)
	err := r.db.Get(&list, query, userID, listID)

	return list, err
}

func (r *TaskListPostgres) Delete(userID, listID int) error {
	query := fmt.Sprintf("DELETE FROM %s tl USING %s ul WHERE tl.id = ul.list_id AND ul.user_id=$1 AND ul.list_id=$2",
		taskListsTable, usersListsTable)
	_, err := r.db.Exec(query, userID, listID)

	return err
}

func (r *TaskListPostgres) Update(userID, listID int, input task.UpdateListInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argID := 1

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argID))
		args = append(args, *input.Title)
		argID++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argID))
		args = append(args, *input.Description)
		argID++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE %s tl SET %s FROM %s ul WHERE tl.id = ul.list_id AND ul.list_id=$%d AND ul.user_id=$%d",
		taskListsTable, setQuery, usersListsTable, argID, argID+1)
	args = append(args, listID, userID)

	logrus.Debugf("updateQuery: %s", query)
	logrus.Debugf("args: %s", args)

	_, err := r.db.Exec(query, args...)
	return err
}
