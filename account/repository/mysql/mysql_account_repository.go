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
	query := sq.Select("account_id, password, email, salt, email_token, is_verified").
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
	err = row.Scan(
		&account.AccountID,
		&account.Password,
		&account.Email,
		&account.Salt,
		&account.EmailToken,
		&account.IsVerified,
	)
	if err != nil {
		return nil, cerror.NewAndPrintWithTag("GA02", err, global.FRIENDLY_MESSAGE)
	}

	return account, nil
}

func (ur MySqlUserRepository) InsertAccount(account domain.Account) (*domain.Account, error) {
	account.AccountID = util.GenerateUUID()

	query := sq.Insert("account").
		Columns(`
			account_id, 
			email, 
			password, 
			salt, 
			email_token, 
			is_verified`).
		Values(
			account.AccountID,
			account.Email,
			account.Password,
			account.Salt,
			account.EmailToken,
			account.IsVerified,
		)

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
		return nil, cerror.NewAndPrintWithTag("IA05", err, global.FRIENDLY_MESSAGE)
	}

	err = tx.Commit()
	if err != nil {
		return nil, cerror.NewAndPrintWithTag("IA04", err, global.FRIENDLY_MESSAGE)
	}

	return &account, nil
}

func (ur MySqlUserRepository) UpdateIsVerified(accountId string, isVerified bool) error {
	/*start create query*/
	query := sq.Update("account").
		Set("is_verified", isVerified).
		Where(sq.Eq{"account_id": accountId})

	sqlString, args, err := query.ToSql()
	if err != nil {
		return cerror.NewAndPrintWithTag("UIV00", err, global.FRIENDLY_MESSAGE)
	}
	/*start create query*/

	tx, err := ur.Db.Begin()
	if err != nil {
		return cerror.NewAndPrintWithTag("UIV01", err, global.FRIENDLY_MESSAGE)
	}

	stmt, err := tx.Prepare(sqlString)
	if err != nil {
		tx.Rollback()
		return cerror.NewAndPrintWithTag("UIV02", err, global.FRIENDLY_MESSAGE)
	}
	defer stmt.Close()

	_, err = tx.Exec(sqlString, args...)
	if err != nil {
		tx.Rollback()
		return cerror.NewAndPrintWithTag("UIV03", err, global.FRIENDLY_MESSAGE)
	}

	err = tx.Commit()
	if err != nil {
		return cerror.NewAndPrintWithTag("UIV04", err, global.FRIENDLY_MESSAGE)
	}
	return nil
}
