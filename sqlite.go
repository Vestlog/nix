package nix

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

var (
	poolsize = 10
)

type SQLiteDatabase struct {
	db             *sql.DB
	connectionPool chan struct{}
}

func (db *SQLiteDatabase) Close() error {
	return db.db.Close()
}

func (db *SQLiteDatabase) SavePost(post Post) error {
	db.connectionPool <- struct{}{}
	defer func() {
		<-db.connectionPool
	}()
	_, err := db.db.Exec(
		"INSERT INTO posts (user_id, id, title, body) VALUES ($1, $2, $3, $4)",
		post.UserID, post.ID, post.Title, post.Body,
	)
	if err != nil {
		return err
	}
	return nil
}

func (db *SQLiteDatabase) SaveComment(comment Comment) error {
	db.connectionPool <- struct{}{}
	defer func() {
		<-db.connectionPool
	}()
	_, err := db.db.Exec(
		`INSERT INTO comments
		(post_id, id, name, email, body)
		VALUES ($1, $2, $3, $4, $5)`,
		comment.PostID, comment.ID, comment.Name, comment.Email, comment.Body,
	)
	if err != nil {
		return err
	}
	return nil
}

func (db *SQLiteDatabase) CreateTables() error {
	if err := db.CreatePostsTable(); err != nil {
		return fmt.Errorf("could not creating posts table: %w", err)
	}
	if err := db.CreateCommentsTable(); err != nil {
		return fmt.Errorf("could not creating comments table: %w", err)
	}
	return nil
}

func (db *SQLiteDatabase) CreatePostsTable() error {
	q := `CREATE TABLE IF NOT EXISTS posts (
		user_id INTEGER,
		id INTEGER,
		title TEXT,
		body TEXT
	)`
	if _, err := db.db.Exec(q); err != nil {
		return err
	}
	return nil
}

func (db *SQLiteDatabase) CreateCommentsTable() error {
	q := `CREATE TABLE IF NOT EXISTS comments (
		post_id INTEGER,
		id INTEGER,
		name TEXT,
		email TEXT,
		body TEXT
	)`
	if _, err := db.db.Exec(q); err != nil {
		return err
	}
	return nil
}

func CreateSQLiteDatabase(dsn string) (*SQLiteDatabase, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	pool := make(chan struct{}, poolsize)
	return &SQLiteDatabase{db, pool}, nil
}
