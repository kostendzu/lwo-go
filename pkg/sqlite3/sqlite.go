package sqlite3

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type Sqlite struct {
	conn *sql.DB
}

func ConnectorInit() (*Sqlite, error) {
	filepath := os.Getenv("FILEPATH")

	connector, err := newSqliteConnector(filepath)

	return connector, err
}

func newSqliteConnector(filepath string) (*Sqlite, error) {
	connector := &Sqlite{}
	dbConn, err := connector.setConn(filepath)
	if err != nil {
		return nil, err
	}
	return dbConn, nil
}

func (p *Sqlite) setConn(filepath string) (*Sqlite, error) {
	if p.conn != nil {
		return p, nil
	}

	dbConn, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return nil, err
	}
	p.conn = dbConn

	if err := p.conn.Ping(); err != nil {
		return nil, err
	}

	return p, nil
}

func (p *Sqlite) Query(query string, args ...any) (*sql.Rows, error) {
	return p.conn.Query(query, args...)
}

func (p *Sqlite) QueryRow(query string, args ...any) *sql.Row {
	return p.conn.QueryRow(query, args...)
}

func (p *Sqlite) Exec(query string, args ...any) (sql.Result, error) {
	return p.conn.Exec(query, args...)
}
