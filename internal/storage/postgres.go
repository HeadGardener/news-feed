package storage

import (
	"fmt"
	"github.com/HeadGardener/news-feed/internal/configs"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

const (
	sourcesTable   = "sources"
	articlesTable  = "articles"
	usersTable     = "users"
	favoritesTable = "favorites"
)

// source table queries
var (
	saveSourceQuery = fmt.Sprintf(`INSERT INTO %s (name, feed_url, created_at) VALUES ($1,$2,$3) RETURNING id`,
		sourcesTable)
	getSourcesQuery    = fmt.Sprintf(`SELECT * FROM %s`, sourcesTable)
	getSourceByIDQuery = fmt.Sprintf(`SELECT * FROM %s WHERE id=$1`, sourcesTable)
)

// article table queries
var (
	saveArticleQuery = fmt.Sprintf(`INSERT INTO %s (source_id, title, link, summary, published_at)
											VALUES ($1,$2,$3,$4,$5)`, articlesTable)
	getArticleByIDQuery = fmt.Sprintf(`SELECT * FROM %s WHERE id=$1`, articlesTable)
)

func getArticlesQuery(v1, v2 int) string {
	return fmt.Sprintf(`SELECT * FROM %s WHERE id BETWEEN %d AND %d`, articlesTable, v1, v2)
}

// user table queries
var (
	createUserQuery = fmt.Sprintf(`INSERT INTO %s (username, email, password_hash) VALUES ($1,$2,$3) RETURNING id`,
		usersTable)
	getUsersForSendQuery = fmt.Sprintf(`SELECT * FROM %s WHERE send_flag=1`, usersTable)
	getUserWithInput     = fmt.Sprintf(`SELECT * FROM %s WHERE username=$1 AND email=$2 AND password_hash=$3`,
		usersTable)
)

// favorites table queries
var (
	addFavoriteQuery    = fmt.Sprintf(`INSERT INTO %s (user_id, article_id) VALUES ($1,$2)`, favoritesTable)
	getArticlesByUserID = fmt.Sprintf(`SELECT * FROM %s WHERE user_id=$1`, favoritesTable)
	deleteFavoriteQuery = fmt.Sprintf(`DELETE FROM %s WHERE user_id=$1 AND article_id=$2`)
)

func NewDB(conf configs.DBConfig) (*sqlx.DB, error) {
	db, err := sqlx.Open("pgx",
		fmt.Sprintf("host=%s dbname=%s sslmode=%s", conf.Host, conf.DBName, conf.SSLMode))

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
