package postgresql

import (
	"database/sql"
	"errors"
	"se03/models"
)

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires)
VALUES($1, $2, current_timestamp,
current_timestamp + make_interval(days => $3) )  RETURNING id`
	var id int
	err := m.DB.QueryRow(stmt, title, content, expires ).Scan(&id)
	//result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		//return 0, err
		return 0, nil
	}
	//id, err := result.LastInsertId()

	return id, nil
}

func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
WHERE expires > current_timestamp AND id = $1`
	row := m.DB.QueryRow(stmt, id)
	s := &models.Snippet{}

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}
func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
WHERE expires > current_timestamp ORDER BY created DESC LIMIT 10`
	// our query.
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	snippets := []*models.Snippet{}

	for rows.Next() {
		s := &models.Snippet{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return snippets, nil
}
