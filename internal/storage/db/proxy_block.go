package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"

	"github.com/google/uuid"
	"github.com/guregu/null/v5"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"

	"github.com/obluumuu/xor/gen/sqlc"
	"github.com/obluumuu/xor/internal/models"
	"github.com/obluumuu/xor/internal/storage"
)

func (s *DbStorage) CreateProxyBlock(ctx context.Context, proxyBlock *models.ProxyBlock) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		logrus.Errorf("failed to begin transaction: %v", err)
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		err := tx.Rollback(ctx)
		if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			logrus.Errorf("rollback transaction: %v", err)
		}
	}()

	q := s.query.WithTx(tx)

	createProxyBlockParams := sqlc.CreateProxyBlockParams{
		Id:          proxyBlock.Id,
		Name:        proxyBlock.Name,
		Description: proxyBlock.Description,
	}
	if err := q.CreateProxyBlock(ctx, createProxyBlockParams); err != nil {
		logrus.Errorf("failed to insert proxy block to db: %v", err)
		return fmt.Errorf("insert proxy block: %w", err)
	}

	proxyBlock.Tags, err = upsertTags(ctx, q, proxyBlock.Tags)
	if err != nil {
		logrus.Errorf("failed to upsert tags: %v", err)
		return fmt.Errorf("upsert tags: %w", err)
	}

	createProxyBlockTagParams := make([]sqlc.CreateProxyBlockTagParams, 0, len(proxyBlock.Tags))
	for _, tag := range proxyBlock.Tags {
		createProxyBlockTagParams = append(createProxyBlockTagParams, sqlc.CreateProxyBlockTagParams{
			ProxyBlockId: proxyBlock.Id,
			TagId:        tag.Id,
		})
	}

	if _, err := q.CreateProxyBlockTag(ctx, createProxyBlockTagParams); err != nil {
		logrus.Errorf("failed to insert proxy block tags to db: %v", err)
		return fmt.Errorf("insert proxy block tags: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		logrus.Errorf("failed to commit transaction: %v", err)
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func (s *DbStorage) GetProxyBlock(ctx context.Context, id uuid.UUID) (*models.ProxyBlock, error) {
	rows, err := s.query.GetProxyBlockWithTags(ctx, id)
	if err != nil {
		logrus.Errorf("failed to select proxy_block with tags: %v", err)
		return nil, fmt.Errorf("select proxy_block with tags: %w", err)
	}

	if len(rows) == 0 {
		return nil, storage.ErrProxyBlockNotFound
	}

	proxyBlock := models.NewFromSqlcProxyBlock(&rows[0].ProxyBlock)
	proxyBlock.Tags = make([]models.Tag, 0, len(rows))
	for _, row := range rows {
		if row.Id == nil || row.Name.IsZero() {
			break
		}
		proxyBlock.Tags = append(proxyBlock.Tags, models.Tag{
			Id:    *row.Id,
			Name:  row.Name.String,
			Color: row.Color.Ptr(),
		})
	}

	return proxyBlock, nil
}

func (s *DbStorage) UpdateProxyBlock(ctx context.Context, proxyBlock *models.ProxyBlock, fieldMask []string) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		logrus.Errorf("failed to begin transaction: %v", err)
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		err := tx.Rollback(ctx)
		if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			logrus.Errorf("rollback transaction: %v", err)
		}
	}()

	q := s.query.WithTx(tx)

	_, err = q.SelectProxyBlockForNoKeyUpdate(ctx, proxyBlock.Id)
	if errors.Is(err, sql.ErrNoRows) {
		return storage.ErrProxyBlockNotFound
	}
	if err != nil {
		logrus.Errorf("failed to select proxy_block for no key update: %v", err)
		return fmt.Errorf("select proxy_block for no key update: %w", err)
	}

	updateProxyBlockParams := sqlc.UpdateProxyBlockParams{Id: proxyBlock.Id}

	if slices.Contains(fieldMask, "name") {
		updateProxyBlockParams.Name = null.StringFrom(proxyBlock.Name)
	}
	if slices.Contains(fieldMask, "description") {
		updateProxyBlockParams.Description = null.StringFrom(proxyBlock.Description)
	}

	if err := q.UpdateProxyBlock(ctx, updateProxyBlockParams); err != nil {
		logrus.Errorf("failed to update proxy_block fields: %v", err)
		return fmt.Errorf("update proxy_block: %w", err)
	}

	if slices.Contains(fieldMask, "tags") {
		if err := q.DeleteProxyBlockTags(ctx, proxyBlock.Id); err != nil {
			logrus.Errorf("failed to delete proxy_block tags: %v", err)
			return fmt.Errorf("delete proxy_block tags: %w", err)
		}

		proxyBlock.Tags, err = upsertTags(ctx, q, proxyBlock.Tags)
		if err != nil {
			logrus.Errorf("failed to upsert tags: %v", err)
			return fmt.Errorf("upsert tags: %w", err)
		}

		createProxyBlockTagParams := make([]sqlc.CreateProxyBlockTagParams, 0, len(proxyBlock.Tags))
		for _, tag := range proxyBlock.Tags {
			createProxyBlockTagParams = append(createProxyBlockTagParams, sqlc.CreateProxyBlockTagParams{
				ProxyBlockId: proxyBlock.Id,
				TagId:        tag.Id,
			})
		}

		if _, err := q.CreateProxyBlockTag(ctx, createProxyBlockTagParams); err != nil {
			logrus.Errorf("failed to insert proxy block tags to db: %v", err)
			return fmt.Errorf("insert proxy_block tags: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		logrus.Errorf("failed to commit transaction: %v", err)
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func (s *DbStorage) DeleteProxyBlock(ctx context.Context, id uuid.UUID) error {
	_, err := s.query.DeleteProxyBlock(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return storage.ErrProxyBlockNotFound
	}
	if err != nil {
		logrus.Errorf("failed to delete proxy_block from db: %v", err)
		return fmt.Errorf("delete proxy_block: %w", err)
	}

	return nil
}

func (s *DbStorage) GetProxiesByProxyBlockId(ctx context.Context, id uuid.UUID) ([]*models.Proxy, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		logrus.Errorf("failed to begin transaction: %v", err)
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		err := tx.Rollback(ctx)
		if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			logrus.Errorf("rollback transaction: %v", err)
		}
	}()

	q := s.query.WithTx(tx)

	_, err = q.SelectProxyBlockForNoKeyUpdate(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, storage.ErrProxyBlockNotFound
	}
	if err != nil {
		logrus.Errorf("failed to select proxy_block for no key update: %v", err)
		return nil, fmt.Errorf("select proxy_block for no key update: %w", err)
	}

	tagsCount, err := q.GetProxyBlockTagsCount(ctx, id)
	if err != nil {
		logrus.Errorf("failed to get proxy_block tags count: %v", err)
		return nil, fmt.Errorf("get proxy_block tags count: %w", err)
	}

	proxies, err := q.GetProxiesByProxyBlockId(ctx, sqlc.GetProxiesByProxyBlockIdParams{ProxyBlockId: id, Column2: int32(tagsCount)})
	if err != nil {
		logrus.Errorf("failed to get proxies by proxy_block id %s: %v", id, err)
		return nil, fmt.Errorf("get proxies by proxy_block id: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		logrus.Errorf("failed to commit transaction: %v", err)
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	res := make([]*models.Proxy, 0, len(proxies))
	for _, proxy := range proxies {
		res = append(res, models.NewFromSqlcProxy(&proxy))
	}

	return res, nil
}
