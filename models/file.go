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

func (folder *Folder) Delete() error {
	return db.Delete(folder).Error
}

func (folder *Folder) GetChildFolders() ([]Folder, error) {
	var folders []Folder
	result := db.Where("parent_id = ?", folder.ID).Find(&folders)
	return folders, result.Error
}

func (folder *Folder) GetChildFiles() ([]File, error) {
	var files []File
	result := db.Where("folder_id = ?", folder.ID).Find(&files)
	return files, result.Error
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

func GetFileByID(userID uint, fileID uint) (*File, error) {
	file := &File{}
	result := db.Where("user_id = ?", userID).First(file, fileID)
	return file, result.Error
}

func GetFolderByID(userID uint, folderID uint) (*Folder, error) {
	folder := &Folder{}
	result := db.Where("user_id = ?", userID).First(folder, folderID)
	return folder, result.Error
}

func GetObjectByID(userID uint, objectID uint) (*File, *Folder, error) {
	if file, err := GetFileByID(userID, objectID); err == nil {
		if folder, err := GetFolderByID(userID, file.FolderID); err == nil {
			return file, folder, nil
		} else {
			return nil, nil, err
		}
	} else if err != gorm.ErrRecordNotFound {
		return nil, nil, err
	}

	if folder, err := GetFolderByID(userID, objectID); err == nil {
		return nil, folder, nil
	} else if err != gorm.ErrRecordNotFound {
		return nil, nil, err
	}

	return nil, nil, gorm.ErrRecordNotFound
}

func DeleteChildFiles(userID uint, folderID uint) error {
	err := db.Where("user_id = ? and folder_id = ?", userID, folderID).Delete(&File{}).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	} else {
		return nil
	}
}

func DeleteFile(userID uint, fileID uint) error {
	err := db.Where("user_id = ?", userID).Delete(&File{}, fileID).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	} else {
		return nil
	}
}
