package quotes

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const (
	createTable = `
	CREATE TABLE IF NOT EXISTS quotes (
		id VARCHAR(255) NOT NULL,
		author VARCHAR(255) NOT NULL,
		quote TEXT,
		created INT NOT NULL,
		PRIMARY KEY (id)
	)  ENGINE=INNODB;
	`
)

type Service struct {
	db *sql.DB
}

func NewService() (Interactor, error) {
	db, err := sql.Open("mysql", "cvcc-user:cvcc-pass@tcp(127.0.0.1:3307)/sample")
	if err != nil {
		return nil, errors.Wrap(err, "failed to create SQL connection")
	}

	if err := db.Ping(); err != nil {
		return nil, errors.Wrap(err, "failed to ping SQL server")
	}

	_, err = db.ExecContext(context.Background(), createTable)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create table")
	}

	return &Service{db: db}, nil
}

func (s *Service) List(ctx context.Context) ([]*Quote, error) {

	rows, err := s.db.QueryContext(ctx, "select id, author, quote, created from quotes order by created desc")
	if err != nil {
		return nil, errors.Wrap(err, "failed to query database")
	}

	results := []*Quote{}
	for rows.Next() {
		quote := &Quote{}
		if err := rows.Scan(&quote.ID, &quote.Author, &quote.Quote, &quote.Created); err != nil {
			return nil, errors.Wrap(err, "failed to read row from database")
		}

		results = append(results, quote)
	}

	return results, nil
}

func (s *Service) Put(ctx context.Context, q *Quote) (*Quote, error) {
	ins, err := s.db.Prepare("INSERT INTO quotes VALUES( ?, ?, ?, ? )") // ? = placeholder
	if err != nil {
		return nil, errors.Wrap(err, "failed to prepare insert statement")
	}
	defer ins.Close() // Close the statement when we leave main() / the program terminates

	q.ID = uuid.New().String()
	q.Created = time.Now().UTC().Unix()
	_, err = ins.Exec(q.ID, q.Author, q.Quote, q.Created)
	if err != nil {
		return nil, errors.Wrap(err, "failed to insert quote")
	}

	return q, nil
}
