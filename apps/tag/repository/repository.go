package repository

import (
	"context"
	"database/sql"
	"errors"
	"mohhefni/go-blog-app/apps/tag/entity"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	GetTagByName(ctx context.Context, tagSearch string) (nameTag []entity.TagsList, err error)
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetTagByName(ctx context.Context, tagSearch string) (nameTag []entity.TagsList, err error) {
	query := `
		SELECT 
			tags.name, 
			COUNT(posts_tags.post_id) AS post_count
		FROM 
			tags
		LEFT JOIN 
			posts_tags ON tags.tag_id = posts_tags.tag_id
		WHERE
			tags.name LIKE $1
		GROUP BY 
			tags.tag_id
		ORDER BY 
			post_count DESC
		LIMIT 15;
	`

	searchParam := "%" + tagSearch + "%"
	err = r.db.SelectContext(ctx, &nameTag, query, searchParam)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return
	}

	return
}
