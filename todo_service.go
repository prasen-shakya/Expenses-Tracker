package main

import (
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3"
)

// Todo represents a single task.
type Todo struct {
	id        int
	Task      string
	Completed bool
}

// ID returns the immutable ID of the Todo.
func (t Todo) ID() int {
	return t.id
}

// Store manages Todo tasks using SQLite.
type Store struct {
	db *sql.DB
}

// NewStore initializes the SQLite database and returns a Store.
func NewStore() (*Store, error) {
	db, err := sql.Open("sqlite3", "./todo.db")
	if err != nil {
		return nil, err
	}

	sqlStmt := `
	CREATE TABLE IF NOT EXISTS todos (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		task TEXT NOT NULL,
		completed BOOLEAN DEFAULT 0
	);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return nil, err
	}

	return &Store{db: db}, nil
}

// CloseStore closes the database connection
func (s *Store) CloseStore() error {
	return s.db.Close()
}

// GetAllTasks retrieves all tasks from the database.
func (s *Store) GetAllTasks() ([]Todo, error) {
	rows, err := s.db.Query("SELECT id, task, completed FROM todos")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var allTasks []Todo
	for rows.Next() {
		var todoItem Todo
		var id int
		if err := rows.Scan(&id, &todoItem.Task, &todoItem.Completed); err != nil {
			return nil, err
		}
		todoItem.id = id
		allTasks = append(allTasks, todoItem)
	}

	return allTasks, rows.Err()
}

// CreateNewTask inserts a new task into the database and returns the created Todo.
func (s *Store) CreateNewTask(taskDescription string) (Todo, error) {
	if taskDescription == "" {
		return Todo{}, errors.New("taskDescription is empty")
	}

	result, err := s.db.Exec("INSERT INTO todos(task) VALUES(?)", taskDescription)
	if err != nil {
		return Todo{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return Todo{}, err
	}

	return Todo{id: int(id), Task: taskDescription, Completed: false}, nil
}

// RemoveTask deletes a task by its ID.
func (s *Store) RemoveTask(taskId int) error {
	if taskId < 0 {
		return errors.New("taskId cannot be negative")
	}

	result, err := s.db.Exec("DELETE FROM todos WHERE id = ?", taskId)
	if err != nil {
		return err
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return errors.New("task not found")
	}

	return nil
}

// UpdateTask updates Task and Completed fields for a given task ID.
func (s *Store) UpdateTask(taskId int, updatedTask Todo) error {
	result, err := s.db.Exec(
		"UPDATE todos SET task = ?, completed = ? WHERE id = ?",
		updatedTask.Task, updatedTask.Completed, taskId,
	)
	if err != nil {
		return err
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return errors.New("task not found")
	}

	return nil
}
