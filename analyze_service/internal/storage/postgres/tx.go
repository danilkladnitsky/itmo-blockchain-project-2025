package postgres

import (
	"memory-book/internal/lib/utils"
)

func (s *Storage) BeginTransaction() (*StorageTx, error) {

	ctx, _ := utils.GetContext()
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}

	return &StorageTx{tx: tx}, nil
}

func (st *StorageTx) Commit() error {
	st.mu.Lock()
	defer st.mu.Unlock()
	ctx, _ := utils.GetContext()

	return st.tx.Commit(ctx)
}

func (st *StorageTx) Rollback() error {
	st.mu.Lock()
	defer st.mu.Unlock()

	ctx, _ := utils.GetContext()

	return st.tx.Rollback(ctx)
}
