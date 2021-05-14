// Code generated by sqlc. DO NOT EDIT.

package db

import (
	"context"
	"database/sql"
	"fmt"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

func Prepare(ctx context.Context, db DBTX) (*Queries, error) {
	q := Queries{db: db}
	var err error
	if q.activateUserStmt, err = db.PrepareContext(ctx, activateUser); err != nil {
		return nil, fmt.Errorf("error preparing query ActivateUser: %w", err)
	}
	if q.countByEmailStmt, err = db.PrepareContext(ctx, countByEmail); err != nil {
		return nil, fmt.Errorf("error preparing query CountByEmail: %w", err)
	}
	if q.createAccountStmt, err = db.PrepareContext(ctx, createAccount); err != nil {
		return nil, fmt.Errorf("error preparing query CreateAccount: %w", err)
	}
	if q.createTransactionStmt, err = db.PrepareContext(ctx, createTransaction); err != nil {
		return nil, fmt.Errorf("error preparing query CreateTransaction: %w", err)
	}
	if q.createUserStmt, err = db.PrepareContext(ctx, createUser); err != nil {
		return nil, fmt.Errorf("error preparing query CreateUser: %w", err)
	}
	if q.createUserRegistrationStmt, err = db.PrepareContext(ctx, createUserRegistration); err != nil {
		return nil, fmt.Errorf("error preparing query CreateUserRegistration: %w", err)
	}
	if q.deleteAccountStmt, err = db.PrepareContext(ctx, deleteAccount); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteAccount: %w", err)
	}
	if q.deleteTransactionStmt, err = db.PrepareContext(ctx, deleteTransaction); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteTransaction: %w", err)
	}
	if q.deleteUserStmt, err = db.PrepareContext(ctx, deleteUser); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteUser: %w", err)
	}
	if q.existsByEmailStmt, err = db.PrepareContext(ctx, existsByEmail); err != nil {
		return nil, fmt.Errorf("error preparing query ExistsByEmail: %w", err)
	}
	if q.findByEmailStmt, err = db.PrepareContext(ctx, findByEmail); err != nil {
		return nil, fmt.Errorf("error preparing query FindByEmail: %w", err)
	}
	if q.getAccountStmt, err = db.PrepareContext(ctx, getAccount); err != nil {
		return nil, fmt.Errorf("error preparing query GetAccount: %w", err)
	}
	if q.getBalanceStmt, err = db.PrepareContext(ctx, getBalance); err != nil {
		return nil, fmt.Errorf("error preparing query GetBalance: %w", err)
	}
	if q.getTransactionStmt, err = db.PrepareContext(ctx, getTransaction); err != nil {
		return nil, fmt.Errorf("error preparing query GetTransaction: %w", err)
	}
	if q.getUserStmt, err = db.PrepareContext(ctx, getUser); err != nil {
		return nil, fmt.Errorf("error preparing query GetUser: %w", err)
	}
	if q.getUserRegistrationStmt, err = db.PrepareContext(ctx, getUserRegistration); err != nil {
		return nil, fmt.Errorf("error preparing query GetUserRegistration: %w", err)
	}
	if q.getUserRegistrationByUserIdStmt, err = db.PrepareContext(ctx, getUserRegistrationByUserId); err != nil {
		return nil, fmt.Errorf("error preparing query GetUserRegistrationByUserId: %w", err)
	}
	if q.isUserActivatedStmt, err = db.PrepareContext(ctx, isUserActivated); err != nil {
		return nil, fmt.Errorf("error preparing query IsUserActivated: %w", err)
	}
	if q.listAccountsStmt, err = db.PrepareContext(ctx, listAccounts); err != nil {
		return nil, fmt.Errorf("error preparing query ListAccounts: %w", err)
	}
	if q.listTransactionsStmt, err = db.PrepareContext(ctx, listTransactions); err != nil {
		return nil, fmt.Errorf("error preparing query ListTransactions: %w", err)
	}
	if q.listUsersStmt, err = db.PrepareContext(ctx, listUsers); err != nil {
		return nil, fmt.Errorf("error preparing query ListUsers: %w", err)
	}
	if q.updateAccountStmt, err = db.PrepareContext(ctx, updateAccount); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateAccount: %w", err)
	}
	return &q, nil
}

func (q *Queries) Close() error {
	var err error
	if q.activateUserStmt != nil {
		if cerr := q.activateUserStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing activateUserStmt: %w", cerr)
		}
	}
	if q.countByEmailStmt != nil {
		if cerr := q.countByEmailStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing countByEmailStmt: %w", cerr)
		}
	}
	if q.createAccountStmt != nil {
		if cerr := q.createAccountStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createAccountStmt: %w", cerr)
		}
	}
	if q.createTransactionStmt != nil {
		if cerr := q.createTransactionStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createTransactionStmt: %w", cerr)
		}
	}
	if q.createUserStmt != nil {
		if cerr := q.createUserStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createUserStmt: %w", cerr)
		}
	}
	if q.createUserRegistrationStmt != nil {
		if cerr := q.createUserRegistrationStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createUserRegistrationStmt: %w", cerr)
		}
	}
	if q.deleteAccountStmt != nil {
		if cerr := q.deleteAccountStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteAccountStmt: %w", cerr)
		}
	}
	if q.deleteTransactionStmt != nil {
		if cerr := q.deleteTransactionStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteTransactionStmt: %w", cerr)
		}
	}
	if q.deleteUserStmt != nil {
		if cerr := q.deleteUserStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteUserStmt: %w", cerr)
		}
	}
	if q.existsByEmailStmt != nil {
		if cerr := q.existsByEmailStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing existsByEmailStmt: %w", cerr)
		}
	}
	if q.findByEmailStmt != nil {
		if cerr := q.findByEmailStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing findByEmailStmt: %w", cerr)
		}
	}
	if q.getAccountStmt != nil {
		if cerr := q.getAccountStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getAccountStmt: %w", cerr)
		}
	}
	if q.getBalanceStmt != nil {
		if cerr := q.getBalanceStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getBalanceStmt: %w", cerr)
		}
	}
	if q.getTransactionStmt != nil {
		if cerr := q.getTransactionStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getTransactionStmt: %w", cerr)
		}
	}
	if q.getUserStmt != nil {
		if cerr := q.getUserStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getUserStmt: %w", cerr)
		}
	}
	if q.getUserRegistrationStmt != nil {
		if cerr := q.getUserRegistrationStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getUserRegistrationStmt: %w", cerr)
		}
	}
	if q.getUserRegistrationByUserIdStmt != nil {
		if cerr := q.getUserRegistrationByUserIdStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getUserRegistrationByUserIdStmt: %w", cerr)
		}
	}
	if q.isUserActivatedStmt != nil {
		if cerr := q.isUserActivatedStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing isUserActivatedStmt: %w", cerr)
		}
	}
	if q.listAccountsStmt != nil {
		if cerr := q.listAccountsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing listAccountsStmt: %w", cerr)
		}
	}
	if q.listTransactionsStmt != nil {
		if cerr := q.listTransactionsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing listTransactionsStmt: %w", cerr)
		}
	}
	if q.listUsersStmt != nil {
		if cerr := q.listUsersStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing listUsersStmt: %w", cerr)
		}
	}
	if q.updateAccountStmt != nil {
		if cerr := q.updateAccountStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateAccountStmt: %w", cerr)
		}
	}
	return err
}

func (q *Queries) exec(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (sql.Result, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).ExecContext(ctx, args...)
	case stmt != nil:
		return stmt.ExecContext(ctx, args...)
	default:
		return q.db.ExecContext(ctx, query, args...)
	}
}

func (q *Queries) query(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (*sql.Rows, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryContext(ctx, args...)
	default:
		return q.db.QueryContext(ctx, query, args...)
	}
}

func (q *Queries) queryRow(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) *sql.Row {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryRowContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryRowContext(ctx, args...)
	default:
		return q.db.QueryRowContext(ctx, query, args...)
	}
}

type Queries struct {
	db                              DBTX
	tx                              *sql.Tx
	activateUserStmt                *sql.Stmt
	countByEmailStmt                *sql.Stmt
	createAccountStmt               *sql.Stmt
	createTransactionStmt           *sql.Stmt
	createUserStmt                  *sql.Stmt
	createUserRegistrationStmt      *sql.Stmt
	deleteAccountStmt               *sql.Stmt
	deleteTransactionStmt           *sql.Stmt
	deleteUserStmt                  *sql.Stmt
	existsByEmailStmt               *sql.Stmt
	findByEmailStmt                 *sql.Stmt
	getAccountStmt                  *sql.Stmt
	getBalanceStmt                  *sql.Stmt
	getTransactionStmt              *sql.Stmt
	getUserStmt                     *sql.Stmt
	getUserRegistrationStmt         *sql.Stmt
	getUserRegistrationByUserIdStmt *sql.Stmt
	isUserActivatedStmt             *sql.Stmt
	listAccountsStmt                *sql.Stmt
	listTransactionsStmt            *sql.Stmt
	listUsersStmt                   *sql.Stmt
	updateAccountStmt               *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:                              tx,
		tx:                              tx,
		activateUserStmt:                q.activateUserStmt,
		countByEmailStmt:                q.countByEmailStmt,
		createAccountStmt:               q.createAccountStmt,
		createTransactionStmt:           q.createTransactionStmt,
		createUserStmt:                  q.createUserStmt,
		createUserRegistrationStmt:      q.createUserRegistrationStmt,
		deleteAccountStmt:               q.deleteAccountStmt,
		deleteTransactionStmt:           q.deleteTransactionStmt,
		deleteUserStmt:                  q.deleteUserStmt,
		existsByEmailStmt:               q.existsByEmailStmt,
		findByEmailStmt:                 q.findByEmailStmt,
		getAccountStmt:                  q.getAccountStmt,
		getBalanceStmt:                  q.getBalanceStmt,
		getTransactionStmt:              q.getTransactionStmt,
		getUserStmt:                     q.getUserStmt,
		getUserRegistrationStmt:         q.getUserRegistrationStmt,
		getUserRegistrationByUserIdStmt: q.getUserRegistrationByUserIdStmt,
		isUserActivatedStmt:             q.isUserActivatedStmt,
		listAccountsStmt:                q.listAccountsStmt,
		listTransactionsStmt:            q.listTransactionsStmt,
		listUsersStmt:                   q.listUsersStmt,
		updateAccountStmt:               q.updateAccountStmt,
	}
}
