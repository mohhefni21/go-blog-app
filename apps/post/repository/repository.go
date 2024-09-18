package repository

import (
	"context"
	"database/sql"
	"errors"
	"mohhefni/go-blog-app/apps/post/entity"
	"mohhefni/go-blog-app/infra/errorpkg"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Repository interface {
	VerifyAvailableTitle(ctx context.Context, title string) (err error)
	AddPost(ctx context.Context, model entity.PostEntity) (idPost int, err error)
	UpdateCover(ctx context.Context, cover string, idPost int) (err error)
	GetDataPosts(ctx context.Context, model entity.PostsPaginationEntity) (posts []entity.GetListPostsEntity, err error)
	GetDetailPostBySLug(ctx context.Context, slug string) (postDetail entity.GetDetailPostResponseEntity, err error)
	GetDetailPostBySLugAndInteraction(ctx context.Context, slug string, publicId uuid.UUID) (postDetail entity.GetDetailPostResponseEntity, err error)
	GetPostById(ctx context.Context, idPost int) (postDetail entity.PostEntity, err error)
	VerifyAvailableUsername(ctx context.Context, username string) (err error)
	GetDataPostsByUsername(ctx context.Context, model entity.PostsPaginationEntity, username string) (posts []entity.GetListPostsEntity, err error)
	GetDataPostsByUserLogin(ctx context.Context, publicId uuid.UUID) (posts []entity.GetListPostsByUserLoginEntity, err error)
	DeletePostById(ctx context.Context, idPost int) (err error)
	UpdatePostById(ctx context.Context, model entity.PostEntity) (err error)
	UploadImageContent(ctx context.Context, model entity.ContentImage) (filename string, err error)
	GetContentImageByPostId(ctx context.Context, postId int) (contentImage []entity.ContentImage, err error)
	GetUserByPublicId(ctx context.Context, publicId uuid.UUID) (model entity.UserEntity, err error)
	GetCommentsByPostId(ctx context.Context, postId int) ([]entity.Comment, error)
	AddOrGetTags(ctx context.Context, tags []string) (tagsId []int, err error)
	AddPostTags(ctx context.Context, postId int, tagId int) (err error)
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) VerifyAvailableTitle(ctx context.Context, title string) (err error) {
	query := `
		SELECT
			1
		FROM posts
		WHERE title=$1
	`
	var exits int8
	err = r.db.QueryRowContext(ctx, query, title).Scan(&exits)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}

		return
	}

	return errorpkg.ErrTitleAlreadyUsed
}

func (r *repository) AddPost(ctx context.Context, model entity.PostEntity) (idPost int, err error) {
	query := `
		INSERT INTO posts (
			user_id, title, slug, excerpt, content, published_at, status, created_at, updated_at
		) VALUES (
			:user_id, :title, :slug, :excerpt, :content, :published_at, :status, :created_at, :updated_at
		) RETURNING post_id
	`

	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return
	}

	err = stmt.GetContext(ctx, &idPost, &model)
	if err != nil {
		return
	}

	return
}

func (r *repository) UpdateCover(ctx context.Context, cover string, idPost int) (err error) {
	query := `
		UPDATE posts
		SET cover=$1
		WHERE post_id=$2
	`

	_, err = r.db.ExecContext(ctx, query, cover, idPost)
	if err != nil {
		return
	}

	return
}

func (r *repository) GetDataPosts(ctx context.Context, model entity.PostsPaginationEntity) (posts []entity.GetListPostsEntity, err error) {
	var query string
	if model.Search == "" {
		query = `
        SELECT
            posts.post_id, posts.cover, posts.title, posts.slug, posts.excerpt, 
            posts.published_at, users.fullname, users.username, users.picture
        FROM 
            posts
        INNER JOIN
            users ON posts.user_id = users.user_id
        WHERE 
            posts.post_id > $1
        ORDER BY 
            posts.post_id DESC
        LIMIT $2
        `
		err := r.db.SelectContext(ctx, &posts, query, model.Cursor, model.Limit)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}
			return nil, err
		}
	} else {
		query = `
        SELECT
            posts.post_id, posts.cover, posts.title, posts.slug, posts.excerpt, 
            posts.published_at, users.fullname, users.username, users.picture
        FROM 
            posts
        INNER JOIN
            users ON posts.user_id = users.user_id
        WHERE 
            posts.post_id > $1
        AND
            posts.title ILIKE $2
        ORDER BY 
            posts.post_id DESC
        LIMIT $3
        `
		searchParam := "%" + model.Search + "%"
		err := r.db.SelectContext(ctx, &posts, query, model.Cursor, searchParam, model.Limit)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}
			return nil, err
		}
	}

	return posts, nil
}

func (r *repository) GetDetailPostBySLug(ctx context.Context, slug string) (postDetail entity.GetDetailPostResponseEntity, err error) {
	query := `
        SELECT
            posts.post_id, posts.cover, posts.title, posts.content, posts.published_at, 
			users.fullname AS "author.fullname", 
			users.username AS "author.username", 
			users.picture AS "author.picture"
        FROM 
            posts
        INNER JOIN
            users ON posts.user_id = users.user_id
        WHERE 
            posts.slug=$1
        `

	err = r.db.GetContext(ctx, &postDetail, query, slug)
	if err != nil {
		if err == sql.ErrNoRows {
			err = errorpkg.ErrNotFound
			return
		}
		return
	}

	return
}

func (r *repository) GetDetailPostBySLugAndInteraction(ctx context.Context, slug string, publicId uuid.UUID) (postDetail entity.GetDetailPostResponseEntity, err error) {
	query := `
        SELECT
            posts.post_id, posts.cover, posts.title, posts.content, posts.published_at, 
			users.fullname AS "author.fullname", 
			users.username AS "author.username", 
			users.picture AS "author.picture",
			CASE 
				WHEN likes.user_id IS NOT NULL THEN true 
				ELSE false 
			END AS "interaction.liked",
			CASE 
				WHEN shares.user_id IS NOT NULL THEN true 
				ELSE false 
			END AS "interaction.shared",
			CASE 
				WHEN bookmarks.user_id IS NOT NULL THEN true 
				ELSE false 
			END AS "interaction.bookmarked"
        FROM 
            posts
        INNER JOIN
            users ON posts.user_id = users.user_id
		LEFT JOIN 
			interactions AS likes ON posts.post_id = likes.post_id AND users.public_id = $2 AND likes.type = 'like'
		LEFT JOIN 
			interactions AS shares ON posts.post_id = shares.post_id AND users.public_id = $2 AND shares.type = 'share'
		LEFT JOIN 
			interactions AS bookmarks ON posts.post_id = bookmarks.post_id AND users.public_id = $2 AND bookmarks.type = 'bookmark'
        WHERE 
            posts.slug=$1
        `

	err = r.db.GetContext(ctx, &postDetail, query, slug, publicId)
	if err != nil {
		if err == sql.ErrNoRows {
			err = errorpkg.ErrNotFound
			return
		}
		return
	}

	return
}

func (r *repository) VerifyAvailableUsername(ctx context.Context, username string) (err error) {
	query := `
		SELECT 
			1
		FROM users
		WHERE username=$1
	`

	var exists int
	err = r.db.QueryRowContext(ctx, query, username).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errorpkg.ErrNotFound
		}
		return err
	}

	return nil
}

func (r *repository) GetDataPostsByUsername(ctx context.Context, model entity.PostsPaginationEntity, username string) (posts []entity.GetListPostsEntity, err error) {
	var query string
	if model.Search == "" {
		query = `
			SELECT
				posts.post_id, posts.cover, posts.title, posts.slug, posts.excerpt, 
				posts.published_at, users.fullname, users.username, users.picture
			FROM 
				posts
			INNER JOIN
				users ON posts.user_id = users.user_id
			WHERE
				 users.username=$1
			AND
				posts.post_id > $2
			ORDER BY 
				posts.post_id DESC
			LIMIT $3
			`
		err := r.db.SelectContext(ctx, &posts, query, username, model.Cursor, model.Limit)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}
			return nil, err
		}
	} else {
		query = `
			SELECT
				posts.post_id, posts.cover, posts.title, posts.slug, posts.excerpt, 
				posts.published_at, users.fullname, users.username, users.picture
			FROM 
				posts
			INNER JOIN
				users ON posts.user_id = users.user_id
			WHERE
				 users.username=$1
			AND
				posts.post_id > $2
			ORDER BY 
				posts.post_id DESC
			LIMIT $3
			`
		searchParam := "%" + model.Search + "%"
		err := r.db.SelectContext(ctx, &posts, query, username, model.Cursor, searchParam, model.Limit)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}
			return nil, err
		}
	}

	return posts, nil
}

func (r *repository) GetDataPostsByUserLogin(ctx context.Context, publicId uuid.UUID) (posts []entity.GetListPostsByUserLoginEntity, err error) {
	query := `
			SELECT
				posts.post_id, posts.cover, posts.title, posts.slug, posts.status, 
				posts.published_at, posts.created_at
			FROM 
				posts
			INNER JOIN
				users ON posts.user_id = users.user_id
			WHERE
				 users.public_id=$1
			`
	err = r.db.SelectContext(ctx, &posts, query, publicId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return posts, nil
}

func (r *repository) DeletePostById(ctx context.Context, idPost int) (err error) {
	query := `
			DELETE 
			FROM 
				posts
			WHERE
				post_id=$1
			`

	_, err = r.db.ExecContext(ctx, query, idPost)
	if err != nil {
		if err == sql.ErrNoRows {
			return errorpkg.ErrorNotFound
		}
		return
	}

	return
}

func (r *repository) GetPostById(ctx context.Context, idPost int) (postDetail entity.PostEntity, err error) {
	query := `
        SELECT
            post_id, user_id, cover, title, slug, excerpt, content, published_at, status, created_at, updated_at
        FROM 
            posts
        WHERE 
            post_id=$1
        `

	err = r.db.GetContext(ctx, &postDetail, query, idPost)
	if err != nil {
		if err == sql.ErrNoRows {
			err = errorpkg.ErrNotFound
			return
		}
		return
	}

	return
}

func (r *repository) UpdatePostById(ctx context.Context, model entity.PostEntity) (err error) {
	query := `
        UPDATE
			posts
		SET
            title=:title, excerpt=:excerpt, slug=:slug, content=:content, status=:status, published_at=:published_at, updated_at=:updated_at
        WHERE 
            post_id=:post_id
        `

	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return
	}

	_, err = stmt.ExecContext(ctx, model)
	if err != nil {
		if err == sql.ErrNoRows {
			return errorpkg.ErrorNotFound
		}
		return
	}

	return
}

func (r *repository) UploadImageContent(ctx context.Context, model entity.ContentImage) (filename string, err error) {
	query := `
		INSERT INTO content_image (
			id_post, filename, created_at, updated_at
		) VALUES (
			:id_post, :filename, :created_at, :updated_at
		) RETURNING filename
	`

	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return
	}

	err = stmt.GetContext(ctx, &filename, &model)
	if err != nil {
		return
	}

	return
}

func (r *repository) GetContentImageByPostId(ctx context.Context, postId int) (contentImage []entity.ContentImage, err error) {
	query := `
			SELECT
				id_post, filename
			FROM 
				content_image
			WHERE
				 id_post=$1
			`
	err = r.db.SelectContext(ctx, &contentImage, query, postId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return contentImage, nil
}

func (r *repository) GetUserByPublicId(ctx context.Context, publicId uuid.UUID) (model entity.UserEntity, err error) {
	query := `
		SELECT
			user_id, public_id, username, email
		FROM users
		WHERE public_id=$1
	`

	err = r.db.GetContext(ctx, &model, query, publicId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.UserEntity{}, errorpkg.ErrorNotFound
		}
		return
	}

	return
}

func (r *repository) GetCommentsByPostId(ctx context.Context, postId int) ([]entity.Comment, error) {
	query := `
        WITH RECURSIVE comment_tree AS (
            SELECT 
                comment_id, 
                post_id, 
                user_id, 
                parent_id, 
                content, 
                created_at, 
                updated_at, 
                0 AS level
            FROM comments
            WHERE post_id = $1 AND parent_id=0

            UNION ALL

            SELECT 
                c.comment_id, 
                c.post_id, 
                c.user_id, 
                c.parent_id, 
                c.content, 
                c.created_at, 
                c.updated_at, 
                ct.level + 1 AS level
            FROM comments c
            INNER JOIN comment_tree ct ON c.parent_id = ct.comment_id
        )
        SELECT * FROM comment_tree
        ORDER BY level, created_at;
    `

	rows, err := r.db.QueryContext(ctx, query, postId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []entity.Comment
	for rows.Next() {
		var comment entity.Comment
		err := rows.Scan(
			&comment.CommentId,
			&comment.PostId,
			&comment.UserId,
			&comment.ParentId,
			&comment.Content,
			&comment.CreatedAt,
			&comment.UpdatedAt,
			&comment.Level,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

func (r *repository) AddOrGetTags(ctx context.Context, tags []string) (tagsId []int, err error) {

	// mengambil tags yang sudah ada
	query := `
		SELECT
			tag_id, name
		FROM
			tags
		WHERE
			name=ANY($1)
	`

	var existsTags []struct {
		TagId int    `db:"tag_id"`
		Name  string `db:"name"`
	}

	err = r.db.SelectContext(ctx, &existsTags, query, pq.Array(tags))
	if err != nil {
		return nil, err
	}

	mapExistsTags := make(map[string]int)
	for _, tag := range existsTags {
		// Simpan tag yang sudah ada kedalam map
		mapExistsTags[tag.Name] = tag.TagId
		// Untuk tag yang sudah ada bisa langsung dikembalikan
		tagsId = append(tagsId, tag.TagId)
	}

	// Untuk mengecek tag apa saja yang tidak ada di parameter dari hasil database
	var newTags []string
	for _, tag := range tags {
		_, found := mapExistsTags[tag]
		// jika tidak ada didalam database maka insert ke dalam newTags []
		if !found {
			newTags = append(newTags, tag)
		}
	}

	if len(newTags) > 0 {
		queryInsert := `
			INSERT INTO
				tags (name)
			VALUES
				(:name)
			RETURNING
				tag_id
		`

		// Buat slice untuk batch insert
		newTagsRecord := make([]map[string]interface{}, len(newTags))
		for i, tagName := range newTags {
			newTagsRecord[i] = map[string]interface{}{"name": tagName}
		}

		// Batch insert menggunakan NamedExec atau NamedQuery untuk banyak record
		rows, err := r.db.NamedQueryContext(ctx, queryInsert, newTagsRecord)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		// Ambil tag_id yang baru diinsert
		for rows.Next() {
			var tagId int
			err = rows.Scan(&tagId)
			if err != nil {
				return nil, err
			}

			tagsId = append(tagsId, tagId)
		}
	}

	return tagsId, nil

}

func (r *repository) AddPostTags(ctx context.Context, postId int, tagId int) (err error) {
	query := `
		INSERT INTO posts_tags (
			tag_id, post_id, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4
		)
	`

	_, err = r.db.ExecContext(ctx, query, tagId, postId, time.Now(), time.Now())
	if err != nil {
		return
	}

	return
}
