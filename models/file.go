package models

import "gorm.io/gorm"

type File struct {
	gorm.Model
	Name     string `gorm:"unique_index:idx_only_one"`
	UserID   uint   `gorm:"index:user_id;unique_index:idx_only_one"`
	Size     uint64 `gorm:"not null"`
	MIMEType string `gorm:"not null"`
	FolderID string `gorm:"index:folder_id;unique_index:idx_only_one"`
}

type Folder struct {
	gorm.Model
	Name     string `gorm:"unique_index:idx_only_one_name"`
	FullPath string `gorm:"index:full_path"`
	ParentID *uint  `gorm:"index:parent_id;unique_index:idx_only_one_name"`
	UserID   uint   `gorm:"index:user_id"`
}
