package db

import (
	"database/sql"
	"fmt"
	"github.com/IkoAfianando/mispress/models"
	"github.com/google/uuid"
)

type Post struct {
	PostId string `json:"post_id"`
	PostName string `json:"post_name"`
}

var (
	ErrNoRecord = fmt.Errorf("no matching record found")
	insertOp    = "insert"
	deleteOp    = "delete"
	updateOp    = "update"
)

func (db Database) SavePost(post *models.Post) error {
	var postId = uuid.NewString()
	query := `INSERT INTO posts(post_id, title, body) VALUES ($1, $2, $3)`
	err := db.Conn.QueryRow(query, postId, post.Title, post.Body).Err()
	if err != nil {
		return err
	}
	logQuery := `INSERT INTO post_logs(post_log_id, post_id, operation) VALUES ($1, $2, $3)`

	var postLogId = uuid.NewString()
	_, err = db.Conn.Exec(logQuery, postLogId, postId, insertOp)
	if err != nil {
		db.Logger.Err(err).Msg("could not log operation for logstash")
	}
	return nil
}

func (db Database) UpdatePost(postId string, post models.Post) error {
	query := "UPDATE posts SET title=$1, body=$2 WHERE post_id=$3"
	_, err := db.Conn.Exec(query, post.Title, post.Body, postId)
	if err != nil {
		return err
	}

	var postLogId = uuid.NewString()
	logQuery := "INSERT INTO post_logs(post_log_id, post_id, operation) VALUES ($1, $2, $3)"
	_, err = db.Conn.Exec(logQuery, postLogId, postId, updateOp)
	if err != nil {
		db.Logger.Err(err).Msg("could not log operation for logstash")
	}
	return nil
}

func (db Database) DeletePost(postId string) error {
	query := "DELETE FROM Posts WHERE post_id=$1"
	_, err := db.Conn.Exec(query, postId)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrNoRecord
		}
		return err
	}

	var postLogId = uuid.NewString()
	logQuery := "INSERT INTO post_logs(post_log_id, post_id, operation) VALUES ($1, $2, $3)"
	_, err = db.Conn.Exec(logQuery, postLogId, postId, deleteOp)
	if err != nil {
		db.Logger.Err(err).Msg("could not log operation for logstash")
	}
	return nil
}

func (db Database) GetPosts() ([]models.Post, error) {
	rows, err := db.Conn.Query("SELECT post_id, title, body FROM posts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.PostId, &post.Title, &post.Body)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (db Database) GetPost(postId string) (*models.Post, error) {
	var post models.Post
	query := "SELECT post_id, title, body FROM posts WHERE post_id=$1"
	err := db.Conn.QueryRow(query, postId).Scan(&post.PostId, &post.Title, &post.Body)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNoRecord
		}
		return nil, err
	}
	return &post, nil
}
