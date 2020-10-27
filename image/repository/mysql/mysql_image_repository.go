package mysql

import (
	"database/sql"

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
