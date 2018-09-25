package model

import "github.com/jinzhu/gorm"

type Site struct {
	ID     uint64 `gorm:"primary_key"`
	Domain string `gorm:"not null default '' index(domain) VARCHAR(200)"`
	Path   string `gorm:"not null default '' index(domain) VARCHAR(100)"`
}

// NewSite
func NewSite(site *Site) (db *gorm.DB) {
	db = Database.Create(site)
	return
}

// AddSite
func AddSite(domain, path string) (db *gorm.DB) {
	db = Database.Create(&Site{Domain: domain, Path: path})
	return
}

// GetSite
func GetSite(id uint64) (db *gorm.DB, Site Site) {
	db = Database.First(&Site, "id = ?", id)
	return
}

// UpdateSite
func UpdateSite(id uint64, domain, path string) (db *gorm.DB, site Site) {
	site = Site{Domain: domain, Path: path}
	db = Database.Model(&site).Update("id", id)
	return
}

// DeleteSite
func DeleteSite(id uint64) (db *gorm.DB) {
	db = Database.Delete(Site{}, "id = ?", id)
	return
}
