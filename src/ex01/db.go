package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func connectToDB() (*sql.DB, error) {
	connStr := "postgres://postgres:postgres@localhost/blogdb?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	sqlStmt := `
	create table if not exists articles (id serial primary key, title text, content text);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func dbCreateArticle(article *Article, db *sql.DB) error {
	query, err := db.Prepare("insert into articles(title,content) values ($1,$2)")
	if err != nil {
		return err
	}
	defer query.Close()
	_, err = query.Exec(article.Title, article.Content)
	if err != nil {
		return err
	}
	return nil
}

func dbGetAllArticles(db *sql.DB) ([]*Article, error) {
	articles := make([]*Article, 0)
	query, err := db.Prepare("select * from articles")
	if err != nil {
		return nil, err
	}
	defer query.Close()
	result, err := query.Query()
	if err != nil {
		return nil, err
	}
	for result.Next() {
		data := new(Article)
		err := result.Scan(
			&data.ID,
			&data.Title,
			&data.Content,
		)
		if err != nil {
			return nil, err
		}
		articles = append(articles, data)
	}
	return articles, nil
}

func dbGetArticle(id string, db *sql.DB) (*Article, error) {
	query, err := db.Prepare("select * from articles where id = $1")
	if err != nil {
		return nil, err
	}
	result := query.QueryRow(id)
	data := new(Article)
	err = result.Scan(&data.ID, &data.Title, &data.Content)
	if err != nil {
		return nil, err
	}
	return data, nil
}
