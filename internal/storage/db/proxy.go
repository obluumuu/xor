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

func (s *DbStorage) CreateProxy(ctx context.Context, proxy *models.Proxy) error {
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

	createProxyParams := sqlc.CreateProxyParams{
		Id:          proxy.Id,
		Name:        proxy.Name,
		Description: proxy.Description,
		Schema:      proxy.Schema,
		Host:        proxy.Host,
		Port:        proxy.Port,
		Username:    null.StringFromPtr(proxy.Username),
		Password:    null.StringFromPtr(proxy.Password),
	}

	if err := q.CreateProxy(ctx, createProxyParams); err != nil {
		logrus.Errorf("failed to insert proxy to db: %v", err)
		return fmt.Errorf("insert proxy to db: %w", err)
	}

	proxy.Tags, err = upsertTags(ctx, q, proxy.Tags)
	if err != nil {
		logrus.Errorf("failed to upsert tags: %v", err)
		return fmt.Errorf("upsert tags: %w", err)
	}

	createProxyTagParams := make([]sqlc.CreateProxyTagParams, 0, len(proxy.Tags))
	for _, tag := range proxy.Tags {
		createProxyTagParams = append(createProxyTagParams, sqlc.CreateProxyTagParams{
			ProxyId: proxy.Id,
			TagId:   tag.Id,
		})
	}

	if _, err := q.CreateProxyTag(ctx, createProxyTagParams); err != nil {
		logrus.Errorf("failed to insert proxytags to db: %v", err)
		return fmt.Errorf("insert proxy_tags to db: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		logrus.Errorf("failed to commit transaction: %v", err)
		return fmt.Errorf("commit transaction: %w", err)
	}
	return nil
}

func (s *DbStorage) GetProxy(ctx context.Context, id uuid.UUID) (*models.Proxy, error) {
	rows, err := s.query.GetProxyWithTags(ctx, id)
	if err != nil {
		logrus.Errorf("failed to select proxy with tags: %v", err)
		return nil, fmt.Errorf("select proxy with tags: %w", err)
	}

	if len(rows) == 0 {
		return nil, storage.ErrProxyNotFound
	}

	proxy := models.NewFromSqlcProxy(&rows[0].Proxy)
	proxy.Tags = make([]models.Tag, 0, len(rows))
	for _, row := range rows {
		if row.Id == nil || row.Name.IsZero() {
			break
		}
		proxy.Tags = append(proxy.Tags, models.Tag{
			Id:    *row.Id,
			Name:  row.Name.String,
			Color: row.Color.Ptr(),
		})
	}

	return proxy, nil
}

func (s *DbStorage) UpdateProxy(ctx context.Context, proxy *models.Proxy, fieldMask []string) error {
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

	_, err = q.SelectProxyForNoKeyUpdate(ctx, proxy.Id)
	if errors.Is(err, sql.ErrNoRows) {
		return storage.ErrProxyNotFound
	}
	if err != nil {
		logrus.Errorf("failed to select proxy for no key update: %v", err)
		return fmt.Errorf("select proxy for no key update: %w", err)
	}

	updateProxyParams := sqlc.UpdateProxyParams{Id: proxy.Id}

	if slices.Contains(fieldMask, "name") {
		updateProxyParams.Name = null.StringFrom(proxy.Name)
	}
	if slices.Contains(fieldMask, "description") {
		updateProxyParams.Description = null.StringFrom(proxy.Description)
	}
	if slices.Contains(fieldMask, "schema") {
		updateProxyParams.Schema = null.StringFrom(proxy.Schema)
	}
	if slices.Contains(fieldMask, "host") {
		updateProxyParams.Host = null.StringFrom(proxy.Host)
	}
	if slices.Contains(fieldMask, "port") {
		updateProxyParams.Port = null.Int32From(proxy.Port)
	}
	if slices.Contains(fieldMask, "username") {
		updateProxyParams.Username = null.StringFromPtr(proxy.Username)
	}
	if slices.Contains(fieldMask, "password") {
		updateProxyParams.Password = null.StringFromPtr(proxy.Password)
	}

	if err := q.UpdateProxy(ctx, updateProxyParams); err != nil {
		logrus.Errorf("failed to update proxy fields: %v", err)
		return fmt.Errorf("update proxy fields: %w", err)
	}

	if slices.Contains(fieldMask, "tags") {
		if err := q.DeleteProxyTags(ctx, proxy.Id); err != nil {
			logrus.Errorf("failed to delete proxy tags: %v", err)
			return fmt.Errorf("delete proxy tags: %w", err)
		}

		proxy.Tags, err = upsertTags(ctx, q, proxy.Tags)
		if err != nil {
			logrus.Errorf("failed to upsert tags: %v", err)
			return fmt.Errorf("upsert tags: %w", err)
		}

		createProxyTagParams := make([]sqlc.CreateProxyTagParams, 0, len(proxy.Tags))
		for _, tag := range proxy.Tags {
			createProxyTagParams = append(createProxyTagParams, sqlc.CreateProxyTagParams{
				ProxyId: proxy.Id,
				TagId:   tag.Id,
			})
		}

		if _, err := q.CreateProxyTag(ctx, createProxyTagParams); err != nil {
			logrus.Errorf("failed to insert proxytags to db: %v", err)
			return fmt.Errorf("insert proxy_tags: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		logrus.Errorf("failed to commit transaction: %v", err)
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func (s *DbStorage) DeleteProxy(ctx context.Context, id uuid.UUID) error {
	_, err := s.query.DeleteProxy(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return storage.ErrProxyNotFound
	}
	if err != nil {
		logrus.Errorf("failed to delete proxy from db: %v", err)
		return fmt.Errorf("delete proxy: %w", err)
	}

	return nil
}
