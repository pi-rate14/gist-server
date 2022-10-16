package models

import (
	"database/sql"
	"fmt"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *sql.DB
}

func (model *SnippetModel) Insert(title, content string, expires int) (int, error) {

	expireDate := time.Now().AddDate(0, 0, expires)

	query := `INSERT INTO snippets ("title", "content", "created", "expires") VALUES($1, $2, current_date, $3) RETURNING id`

	statement, err := model.DB.Prepare(query)
	if err != nil {
		return 0, fmt.Errorf("error preparing query %v", err)
	}
	defer statement.Close()

	var lastInsertId int

	err = statement.QueryRow(title, content, expireDate).Scan(&lastInsertId)
	if err != nil {
		return 0, fmt.Errorf("error executing query %v", err)
	}

	return int(lastInsertId), nil
}

func (model *SnippetModel) Get(id int) (*Snippet, error) {
	return nil, nil
}

func (model *SnippetModel) Latest() ([]*Snippet, error) {
	return nil, nil
}
