package response

import "mohhefni/go-blog-app/apps/tag/entity"

type TagsList struct {
	Name       string `json:"name"`
	CountUsing int    `json:"post_count"`
}

func NewTagsListResponse(tags []entity.TagsList) []TagsList {
	var tagsList []TagsList

	for _, tag := range tags {
		tagsList = append(tagsList, TagsList{
			Name:       tag.Name,
			CountUsing: tag.CountUsing,
		})
	}

	return tagsList
}
