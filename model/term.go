package model

import "github.com/jinzhu/gorm"

// Term Information table for article classification, link classification, and label.
type Term struct {
	ID        uint64 `gorm:"primary_key"`
	Name      string `gorm:"not null default '' index VARCHAR(200)"`
	Slug      string `gorm:"not null default '' index VARCHAR(200)"`
	TermGroup uint64 `gorm:"not null default 0 BIGINT(10)"`
}

// NewTerm
func NewTerm(term *Term) (db *gorm.DB) {
	db = Database.Create(term)
	return
}

// AddTerm
func AddTerm(name, slug string, termGroup uint64) (db *gorm.DB) {
	db = Database.Create(&Term{Name: name, Slug: slug, TermGroup: termGroup})
	return
}

// GetTerm
func GetTerm(id uint64) (db *gorm.DB, Term Term) {
	db = Database.First(&Term, "id = ?", id)
	return
}

// UpdateTerm
func UpdateTerm(id uint64, name string) (db *gorm.DB, term Term) {
	term = Term{Name: name}
	db = Database.Model(&term).Update("id", id)
	return
}

// DeleteTerm
func DeleteTerm(id uint64) (db *gorm.DB) {
	db = Database.Delete(Term{}, "id = ?", id)
	return
}
