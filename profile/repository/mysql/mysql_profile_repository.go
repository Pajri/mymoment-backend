package mysql

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/pajri/personal-backend/adapter/cerror"
	"github.com/pajri/personal-backend/domain"
	"github.com/pajri/personal-backend/global"
	"github.com/pajri/personal-backend/util"
)

type MySqlProfileRepository struct {
	Db *sql.DB
}

func NewMySqlProfileRepository(db *sql.DB) domain.IProfileRepository {
	return &MySqlProfileRepository{
		Db: db,
	}
}

func (pr MySqlProfileRepository) InsertProfile(profile domain.Profile) error {
	/*start create sql*/
	query := sq.Insert("profile").
		Columns("profile_id, full_name, account_id").
		Values(util.GenerateUUID(), profile.FullName, profile.AccountID)

	sql, args, err := query.ToSql()
	if err != nil {
		return cerror.NewAndPrintWithTag("IPR00", err, global.FRIENDLY_MESSAGE)
	}
	/*end create sql*/

	/*start insert data*/
	tx, err := pr.Db.Begin()
	if err != nil {
		return cerror.NewAndPrintWithTag("IPR01", err, global.FRIENDLY_MESSAGE)
	}

	stmt, err := tx.Prepare(sql)
	if err != nil {
		tx.Rollback()
		return cerror.NewAndPrintWithTag("IPR02", err, global.FRIENDLY_MESSAGE)
	}
	defer stmt.Close()

	_, err = tx.Exec(sql, args...)
	if err != nil {
		tx.Rollback()
		return cerror.NewAndPrintWithTag("IPR03", err, global.FRIENDLY_MESSAGE)
	}

	err = tx.Commit()
	if err != nil {
		return cerror.NewAndPrintWithTag("IA04", err, global.FRIENDLY_MESSAGE)
	}
	/*end insert data*/

	return nil
}
