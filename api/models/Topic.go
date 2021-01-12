package models

// Topic is model for topics table in the db
type Topic struct {
	Name string `gorm:"size:255;primary_key;" json:"name"`
}
