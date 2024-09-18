package entity

type TagsList struct {
	Name       string `db:"name"`
	CountUsing int    `db:"post_count"`
}
