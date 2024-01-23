package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides all functions to execute db queries and transactions
// 通过Store结构体，我们可以执行所有的数据库操作

type Store struct {
	*Queries
	db *sql.DB
}

// NewStore creates a new Store
// NewStore函数用于创建一个新的Store结构体
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transaction
// execTx函数用于执行一个数据库事务
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	// 开启一个事务
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	// 通过Queries结构体执行数据库操作
	q := New(tx)
	// 执行fn函数
	err = fn(q)
	// 判断fn函数是否执行成功
	if err != nil {
		// 如果fn函数执行失败，则回滚事务
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		// 如果fn函数执行失败，则回滚事务
		return err
	}
	// 如果fn函数执行成功，则提交事务
	return tx.Commit()
}

// TransferTxParams contains the input parameters of the transfer transaction
// TransferTxParams结构体用于存储转账事务的参数
type TransferTxParams struct {
	// 转出的账户
	FromAccountID int64 `json:"from_account_id"`
	// 接收转账的账户
	ToAccountID int64 `json:"to_account_id"`
	// 转账的金额
	Amount int64 `json:"amount"`
}

// TransferTxResult is the result of the transfer transaction
// TransferTxResult结构体用于存储转账事务的结果
type TransferTxResult struct {
	// 转账记录
	Transfer Transfer `json:"transfer"`
	// 转出账户的详情
	FromAccount Account `json:"from_account"`
	// 接收账户的详情
	ToAccount Account `json:"to_account"`
	// 转出记录
	FromEntry Entry `json:"from_entry"`
	// 转入记录
	ToEntry Entry `json:"to_entry"`
}

// TransferTx 转账事务
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	// 定义一个TransferTxResult结构体
	var result TransferTxResult
	// 执行execTx函数
	err := store.execTx(ctx, func(q *Queries) error {
		// 执行CreateTransfer函数
		var err error
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			// 转账金额必须为正数
			Amount: arg.Amount,
		})
		if err != nil {
			return err
		}

		// 创建转账记录
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			// 转出账户的金额必须为负数
			Amount: -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			// 接收转账的账户金额必须为正数
			Amount: arg.Amount,
		})
		if err != nil {
			return err
		}

		// 为账户添加金额
		// 如果转出账户的ID小于转入账户的ID，则先为转出账户添加金额，再为转入账户添加金额
		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
		}

		return err
	})
	// 返回结果
	return result, err
}

func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 int64,
	accountID2 int64,
	amount2 int64,
) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})
	if err != nil {
		return
	}

	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})
	return
}
