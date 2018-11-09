package quotes

import (
	"context"
	"database/sql"
	"fmt"
	"os"
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

var (
	sqlUser = os.Getenv("MYSQL_USER")
	sqlPwd  = os.Getenv("MYSQL_PASS")
	sqlHost = os.Getenv("MYSQL_HOST")
	sqlPort = os.Getenv("MYSQL_PORT")
	sqlDb   = os.Getenv("MYSQL_DB")
)

type Service struct {
	db *sql.DB
}

func NewService() (Interactor, error) {
	db, err := sql.Open(
		"mysql",
		fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", sqlUser, sqlPwd, sqlHost, sqlPort, sqlDb),
	)
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

func (s *Service) Get(ctx context.Context, id string) (*Quote, error) {

	stmt, err := s.db.PrepareContext(ctx, "select id, author, quote, created from quotes where id=?")
	if err != nil {
		return nil, errors.Wrap(err, "failed to prepare query")
	}

	row := stmt.QueryRowContext(ctx, id)
	if row == nil {
		return nil, nil
	}

	quote := &Quote{}
	if err := row.Scan(&quote.ID, &quote.Author, &quote.Quote, &quote.Created); err != nil {
		return nil, errors.Wrap(err, "failed to read row from database")
	}

	return quote, nil
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

func (s *Service) Update(ctx context.Context, q *Quote) (*Quote, error) {
	upd, err := s.db.Prepare("update quotes set author=?, quote=? where id=?") // ? = placeholder
	if err != nil {
		return nil, errors.Wrap(err, "failed to prepare update statement")
	}
	defer upd.Close() // Close the statement when we leave main() / the program terminates

	_, err = upd.Exec(q.Author, q.Quote, q.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update quote")
	}

	return q, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	del, err := s.db.Prepare("delete from quotes where id=?") // ? = placeholder
	if err != nil {
		return errors.Wrap(err, "failed to prepare delete statement")
	}
	defer del.Close() // Close the statement when we leave main() / the program terminates

	_, err = del.Exec(id)
	if err != nil {
		return errors.Wrap(err, "failed to delete quote")
	}

	return nil
}
