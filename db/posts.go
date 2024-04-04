package db

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"gitlab.com/repo/mispress/models"
)

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

	post.ID = postId
	logQuery := "INSERT INTO post_logs(post_id, operation) VALUES ($1, $2)"
	_, err = db.Conn.Exec(logQuery, post.ID, updateOp)
	if err != nil {
		db.Logger.Err(err).Msg("could not log operation for logstash")
	}
	return nil
}

func (db Database) DeletePost(postId int) error {
	query := "DELETE FROM Posts WHERE id=$1"
	_, err := db.Conn.Exec(query, postId)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrNoRecord
		}
		return err
	}

	logQuery := "INSERT INTO post_logs(post_id, operation) VALUES ($1, $2)"
	_, err = db.Conn.Exec(logQuery, postId, deleteOp)
	if err != nil {
		db.Logger.Err(err).Msg("could not log operation for logstash")
	}
	return nil
}
