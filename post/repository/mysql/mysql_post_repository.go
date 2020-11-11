package mysql

import (
	"database/sql"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/pajri/personal-backend/adapter/cerror"
	"github.com/pajri/personal-backend/domain"
	"github.com/pajri/personal-backend/global"
	"github.com/pajri/personal-backend/util"
)

func NewMySqlPostRepository(db *sql.DB) domain.IPostRepository {
	return MySqlPostRepository{
		Db: db,
	}
}

type MySqlPostRepository struct {
	Db *sql.DB
}

func (ur MySqlPostRepository) InsertPost(post domain.Post) error {
	if post.PostID == "" {
		post.PostID = util.GenerateUUID()
	}
	/*start create query*/
	query := sq.Insert("post").
		Columns("post_id", "content", "image_url", "date", "last_updated", "account_id").
		Values(post.PostID, post.Content, post.ImageURL, post.Date, time.Now(), post.AccountID)

	sql, args, err := query.ToSql()
	if err != nil {
		return cerror.NewAndPrintWithTag("IP00", err, global.FRIENDLY_MESSAGE)
	}
	/*end create query*/

	/*start insert execution*/
	tx, err := ur.Db.Begin()
	if err != nil {
		return cerror.NewAndPrintWithTag("IP01", err, global.FRIENDLY_MESSAGE)
	}

	stmt, err := tx.Prepare(sql)
	if err != nil {
		tx.Rollback()
		return cerror.NewAndPrintWithTag("IP02", err, global.FRIENDLY_MESSAGE)
	}
	defer stmt.Close()

	_, err = tx.Exec(sql, args...)
	if err != nil {
		tx.Rollback()
		return cerror.NewAndPrintWithTag("IP03", err, global.FRIENDLY_MESSAGE)
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return cerror.NewAndPrintWithTag("IP04", err, global.FRIENDLY_MESSAGE)
	}

	return nil
	/*end insert execution*/
}

func (ur MySqlPostRepository) PostList(filter domain.PostFilter) ([]domain.Post, error) {
	query := sq.Select("post_id, content, image_url, date").
		From("post").
		OrderBy("date DESC")

	if filter.AccountID != "" {
		query = query.Where(sq.Eq{"account_id": filter.AccountID})
	}

	if filter.Limit != 0 {
		query = query.Limit(filter.Limit)
	}

	var zeroTime time.Time
	if filter.Date != zeroTime {
		query = query.Where(sq.Lt{"date": filter.Date})
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, cerror.NewAndPrintWithTag("PLI00", err, global.FRIENDLY_MESSAGE)
	}

	rows, err := ur.Db.Query(sql, args...)
	if rows != nil {
		defer rows.Close()
	}

	if err != nil {
		return nil, cerror.NewAndPrintWithTag("PLI01", err, global.FRIENDLY_MESSAGE)
	}

	var postList []domain.Post
	for rows.Next() {
		var post domain.Post
		err = rows.Scan(&post.PostID, &post.Content, &post.ImageURL, &post.Date)
		if err != nil {
			return nil, cerror.NewAndPrintWithTag("PLI02", err, global.FRIENDLY_MESSAGE)
		}

		postList = append(postList, post)
	}

	if err = rows.Close(); err != nil {
		// but what should we do if there's an error?
		log.Println(err)
	}

	return postList, nil
}

func (ur MySqlPostRepository) DeletePost(postID, accountID string) error {
	/*start create query*/
	query := sq.Delete("post").
		Where(sq.Eq{
			"post_id":    postID,
			"account_id": accountID,
		})

	sql, args, err := query.ToSql()
	if err != nil {
		return cerror.NewAndPrintWithTag("DP00", err, global.FRIENDLY_MESSAGE)
	}
	/*end create query*/

	tx, err := ur.Db.Begin()
	if err != nil {
		return cerror.NewAndPrintWithTag("DP01", err, global.FRIENDLY_MESSAGE)
	}

	stmt, err := tx.Prepare(sql)
	if err != nil {
		tx.Rollback()
		return cerror.NewAndPrintWithTag("DP02", err, global.FRIENDLY_MESSAGE)
	}
	defer stmt.Close()

	_, err = tx.Exec(sql, args...)
	if err != nil {
		tx.Rollback()
		return cerror.NewAndPrintWithTag("DP03", err, global.FRIENDLY_MESSAGE)
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return cerror.NewAndPrintWithTag("DP04", err, global.FRIENDLY_MESSAGE)
	}

	return tx.Commit()
}
