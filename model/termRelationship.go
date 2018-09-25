package model

// TermRelationship 分类与文章信息表（wp_posts）、链接表(wp_links)的关联表。
type TermRelationship struct {
	ObjectID       uint64 `gorm:"primary_key"`
	TermTaxonomyID uint64
	TermOrder      int
}
