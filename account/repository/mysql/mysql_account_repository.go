package mysql

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/go-sql-driver/mysql"
	"github.com/pajri/personal-backend/adapter/cerror"
	"github.com/pajri/personal-backend/domain"
	"github.com/pajri/personal-backend/global"
	"github.com/pajri/personal-backend/util"
)

func NewMySqlAccountRepository(db *sql.DB) domain.IAccountRepository {
	return &MySqlUserRepository{
		Db: db,
	}
}

type MySqlUserRepository struct {
	Db *sql.DB
}

func (ur MySqlUserRepository) GetAccount(filter domain.AccountFilter) (*domain.Account, error) {
	query := sq.Select("account_id, password, email").
		From("account")

	if filter.Email != "" {
		query = query.Where(sq.Eq{"email": filter.Email})
	}

	sqlString, args, err := query.ToSql()
	if err != nil {
		return nil, cerror.NewAndPrintWithTag("GA00", err, global.FRIENDLY_MESSAGE)
	}

	row := ur.Db.QueryRow(sqlString, args...)

	account := new(domain.Account)
	err = row.Scan(&account.AccountID, &account.Password, &account.Email)
	if err != nil {
		return nil, cerror.NewAndPrintWithTag("GA02", err, global.FRIENDLY_MESSAGE)
	}

	return account, nil
}

func (ur MySqlUserRepository) InsertAccount(account domain.Account) (*domain.Account, error) {
	account.AccountID = util.GenerateUUID()

	query := sq.Insert("account").
		Columns("account_id, email, password, salt").
		Values(account.AccountID, account.Email, account.Password, account.Salt)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, cerror.NewAndPrintWithTag("IA00", err, global.FRIENDLY_MESSAGE)
	}

	tx, err := ur.Db.Begin()
	if err != nil {
		return nil, cerror.NewAndPrintWithTag("IA01", err, global.FRIENDLY_MESSAGE)
	}

	stmt, err := tx.Prepare(sql)
	if err != nil {
		tx.Rollback()
		return nil, cerror.NewAndPrintWithTag("IA02", err, global.FRIENDLY_MESSAGE)
	}
	defer stmt.Close()

	_, err = tx.Exec(sql, args...)
	if err != nil {
		errMySQL, ok := err.(*mysql.MySQLError)
		if ok && errMySQL.Number == 1062 {
			return nil, cerror.NewAndPrintWithTag("IA03", err, global.FRIENDLY_DUPLICATE_EMAIL)
		}
		tx.Rollback()
		return nil, cerror.NewAndPrintWithTag("IA03", err, global.FRIENDLY_MESSAGE)
	}

	err = tx.Commit()
	if err != nil {
		return nil, cerror.NewAndPrintWithTag("IA04", err, global.FRIENDLY_MESSAGE)
	}

	return &account, nil
}
