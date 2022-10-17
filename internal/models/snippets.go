package models

import (
	"database/sql"
	"errors"
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

	query := `SELECT * FROM snippets WHERE expires > current_date AND id = $1`

	statement, err := model.DB.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("error preparing query %v", err)
	}
	defer statement.Close()

	snippet := &Snippet{}

	row := statement.QueryRow(id)

	err = row.Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.Created, &snippet.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return snippet, nil
}

func (model *SnippetModel) Latest() ([]*Snippet, error) {
	query := `SELECT * FROM snippets WHERE expires > current_date ORDER BY id DESC LIMIT 10`

	rows, err := model.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	snippets := []*Snippet{}

	for rows.Next() {
		snippet := &Snippet{}
		err = rows.Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.Created, &snippet.Expires)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, snippet)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
