package db

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/guregu/null/v5"
	"github.com/sirupsen/logrus"

	"github.com/obluumuu/xor/gen/sqlc"
	"github.com/obluumuu/xor/internal/models"
)

func upsertTags(ctx context.Context, query *sqlc.Queries, tags []models.Tag) ([]models.Tag, error) {
	upsertTagsParams := make([]sqlc.UpsertTagsParams, 0, len(tags))
	for _, tag := range tags {
		upsertTagsParams = append(upsertTagsParams, sqlc.UpsertTagsParams{
			Id:    tag.Id,
			Name:  tag.Name,
			Color: null.StringFromPtr(tag.Color),
		})
	}

	var batchLastError error
	rowsCnt := 0
	batchRes := query.UpsertTags(ctx, upsertTagsParams)
	batchRes.QueryRow(func(_ int, u uuid.UUID, err error) {
		if err != nil {
			logrus.Errorf("failed to upsert tag: %v", err)
			batchLastError = err
		} else if rowsCnt < len(tags) {
			tags[rowsCnt].Id = u
		}
		rowsCnt += 1
	})

	if batchLastError != nil {
		return nil, fmt.Errorf("batch error: %w", batchLastError)
	}
	if rowsCnt != len(tags) {
		return nil, fmt.Errorf("query returned %d rows (expected %d)", rowsCnt, len(tags))
	}

	return tags, nil
}
