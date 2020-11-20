package mysql

import (
	"database/sql"
	"os"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/pajri/personal-backend/adapter/cerror"
	"github.com/pajri/personal-backend/domain"
	"github.com/pajri/personal-backend/global"
)

type MySqlImageRepository struct {
	Db *sql.DB
}

func NewMySqlImageRepository(db *sql.DB) domain.IImageRepository {
	return &MySqlImageRepository{
		Db: db,
	}
}

func (im MySqlImageRepository) SaveImage(image domain.Image) error {
	if image.ImageID == "" {
		image.ImageID = uuid.New().String()
	}

	/*start create query*/
	query := sq.Insert("image").
		Columns("image_id, image_url").
		Values(image.ImageID, image.ImageURL)

	sql, args, err := query.ToSql()
	if err != nil {
		return cerror.NewAndPrintWithTag("IMR00", err, global.FRIENDLY_MESSAGE)
	}
	/*end create query*/

	/*start insert data*/
	tx, err := im.Db.Begin()
	if err != nil {
		return cerror.NewAndPrintWithTag("IMR01", err, global.FRIENDLY_MESSAGE)
	}

	stmt, err := tx.Prepare(sql)
	if err != nil {
		tx.Rollback()
		return cerror.NewAndPrintWithTag("IMR02", err, global.FRIENDLY_MESSAGE)
	}
	defer stmt.Close()

	_, err = tx.Exec(sql, args...)
	if err != nil {
		tx.Rollback()
		return cerror.NewAndPrintWithTag("IMR03", err, global.FRIENDLY_MESSAGE)
	}

	err = tx.Commit()
	if err != nil {
		return cerror.NewAndPrintWithTag("IMR04", err, global.FRIENDLY_MESSAGE)
	}

	return nil
	/*end insert data*/
}

func (im MySqlImageRepository) DeleteImage(image domain.Image, deleteFile bool) error {
	filter := domain.ImageFilter{ImageURL: image.ImageURL}
	tx, err := im.deleteFromDb(filter)
	if err != nil {
		return err
	}

	if deleteFile {
		err = im.deleteFile(image.ImageURL)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return cerror.NewAndPrintWithTag("DIG00", err, global.FRIENDLY_MESSAGE)
	}

	return nil
}

func (im MySqlImageRepository) deleteFromDb(filter domain.ImageFilter) (*sql.Tx, error) {
	query := sq.Delete("image")

	if filter.ImageURL != "" {
		query = query.Where(sq.Eq{"image_url": filter.ImageURL})
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, cerror.NewAndPrintWithTag("DIM00", err, global.FRIENDLY_MESSAGE)
	}

	tx, err := im.Db.Begin()
	if err != nil {
		return nil, cerror.NewAndPrintWithTag("DIM01", err, global.FRIENDLY_MESSAGE)
	}

	stmt, err := tx.Prepare(sql)
	if err != nil {
		tx.Rollback()
		return nil, cerror.NewAndPrintWithTag("DIM02", err, global.FRIENDLY_MESSAGE)
	}
	defer stmt.Close()

	_, err = tx.Exec(sql, args...)
	if err != nil {
		tx.Rollback()
		return nil, cerror.NewAndPrintWithTag("DIM03", err, global.FRIENDLY_MESSAGE)
	}

	return tx, nil
}

func (im MySqlImageRepository) deleteFile(path string) error {
	path = strings.ReplaceAll(path, "/", string(os.PathSeparator))
	fullPath := global.WD + path
	err := os.Remove(fullPath)
	if err != nil {
		return cerror.NewAndPrintWithTag("DIF00", err, global.FRIENDLY_MESSAGE)
	}
	return nil
}
