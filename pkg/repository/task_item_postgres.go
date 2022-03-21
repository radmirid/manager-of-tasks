package repository

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	task "github.com/radmirid/manager-of-tasks"
)

type TaskItemPostgres struct {
	db *sqlx.DB
}

func NewTaskItemPostgres(db *sqlx.DB) *TaskItemPostgres {
	return &TaskItemPostgres{db: db}
}

func (r *TaskItemPostgres) Create(listID int, item task.TaskItem) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var itemID int
	createItemQuery := fmt.Sprintf("INSERT INTO %s (title, description) values ($1, $2) RETURNING id", taskItemsTable)

	row := tx.QueryRow(createItemQuery, item.Title, item.Description)
	err = row.Scan(&itemID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	createListItemsQuery := fmt.Sprintf("INSERT INTO %s (list_id, item_id) values ($1, $2)", listsItemsTable)
	_, err = tx.Exec(createListItemsQuery, listID, itemID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return itemID, tx.Commit()
}

func (r *TaskItemPostgres) GetAll(userID, listID int) ([]task.TaskItem, error) {
	var items []task.TaskItem
	query := fmt.Sprintf(`SELECT ti.id, ti.title, ti.description, ti.done FROM %s ti INNER JOIN %s li on li.item_id = ti.id
									INNER JOIN %s ul on ul.list_id = li.list_id WHERE li.list_id = $1 AND ul.user_id = $2`,
		taskItemsTable, listsItemsTable, usersListsTable)
	if err := r.db.Select(&items, query, listID, userID); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *TaskItemPostgres) GetByID(userID, itemID int) (task.TaskItem, error) {
	var item task.TaskItem
	query := fmt.Sprintf(`SELECT ti.id, ti.title, ti.description, ti.done FROM %s ti INNER JOIN %s li on li.item_id = ti.id
									INNER JOIN %s ul on ul.list_id = li.list_id WHERE ti.id = $1 AND ul.user_id = $2`,
		taskItemsTable, listsItemsTable, usersListsTable)
	if err := r.db.Get(&item, query, itemID, userID); err != nil {
		return item, err
	}

	return item, nil
}

func (r *TaskItemPostgres) Delete(userID, itemID int) error {
	query := fmt.Sprintf(`DELETE FROM %s ti USING %s li, %s ul 
									WHERE ti.id = li.item_id AND li.list_id = ul.list_id AND ul.user_id = $1 AND ti.id = $2`,
		taskItemsTable, listsItemsTable, usersListsTable)
	_, err := r.db.Exec(query, userID, itemID)
	return err
}

func (r *TaskItemPostgres) Update(userID, itemID int, input task.UpdateItemInput) error {
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

	if input.Done != nil {
		setValues = append(setValues, fmt.Sprintf("done=$%d", argID))
		args = append(args, *input.Done)
		argID++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf(`UPDATE %s ti SET %s FROM %s li, %s ul
									WHERE ti.id = li.item_id AND li.list_id = ul.list_id AND ul.user_id = $%d AND ti.id = $%d`,
		taskItemsTable, setQuery, listsItemsTable, usersListsTable, argID, argID+1)
	args = append(args, userID, itemID)

	_, err := r.db.Exec(query, args...)
	return err
}
