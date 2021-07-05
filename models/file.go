package models

import "gorm.io/gorm"

type File struct {
	gorm.Model
	Name     string `gorm:"not null;index:idx_only_one,unique"`
	UserID   uint   `gorm:"not null;index:idx_user_id;index:idx_only_one,unique"`
	Size     uint64 `gorm:"not null"`
	MIMEType string `gorm:"not null"`
	FolderID uint   `gorm:"not null;index:idx_folder_id;index:idx_only_one,unique"`
}

type Folder struct {
	gorm.Model
	Name     string `gorm:"not null;index:idx_only_one,unique"`
	FullPath string `gorm:"not null;index:idx_full_path"`
	ParentID *uint  `gorm:"index:parent_id;index:idx_only_one,unique"`
	UserID   uint   `gorm:"not null;index:idx_user_id"`
}

func GetFolderByFullPath(userID uint, fullPath string) (*Folder, error) {
	folder := &Folder{}
	result := db.Where("user_id = ? and full_path = ?", userID, fullPath).First(folder)
	return folder, result.Error
}

func GetRootFolder(userID uint) (*Folder, error) {
	folder := &Folder{}
	result := db.Where("user_id = ? and parent_id is null", userID).First(folder)
	return folder, result.Error
}

func (folder *Folder) Create() error {
	return db.Create(folder).Error
}

func CreateRootFolder(userID uint) (*Folder, error) {
	rootFolder := &Folder{
		Name:     "/",
		FullPath: "/",
		ParentID: nil,
		UserID:   userID,
	}
	err := rootFolder.Create()
	return rootFolder, err
}

func GetFileByNameAndFolderID(userID uint, name string, folderID uint) (*File, error) {
	file := &File{}
	result := db.Where("user_id = ? and name = ? and folder_id = ?", userID, name, folderID).First(file)
	return file, result.Error
}

func (file *File) Create() error {
	return db.Create(file).Error
}

func (file *File) Update() error {
	return db.Save(file).Error
}
