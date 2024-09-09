package utility

import "github.com/gosimple/slug"

func GenerateSlug(title string) string {
	return slug.Make(title)
}
