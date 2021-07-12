package models

import (
	"gorm.io/gorm"
	"path"
)

type ObjectType uint

const (
	FILE ObjectType = 0
	DIR  ObjectType = 1
)

type Object struct {
	gorm.Model
	DeletedAt gorm.DeletedAt `gorm:"index:idx_deleted_at;index:idx_only_one,unique"`
	Type      ObjectType     `gorm:"not null;index:idx_type;index:idx_only_one,unique"`
	Name      string         `gorm:"not null;index:idx_only_one,unique,priority:3"`
	UserID    uint           `gorm:"not null;index:idx_user_id;index:idx_only_one,unique"`
	ParentID  *uint          `gorm:"index:idx_parent_id;index:idx_only_one,unique"`
	FullPath  string         `gorm:"not null;index:idx_full_path"`
	Size      *uint64
	MIMEType  *string
}

func (object Object) IsDir() bool {
	return object.Type == DIR
}

func (object Object) GetType() string {
	if object.IsDir() {
		return "dir"
	} else {
		return "file"
	}
}

func (user User) GetRootDir() (*Object, error) {
	dir := &Object{}
	result := db.Where("user_id = ? and parent_id is null and type = ?", user.ID, DIR).First(dir)
	return dir, result.Error
}

func (user User) CreateRootDir() (*Object, error) {
	rootFolder := &Object{
		Type:     DIR,
		Name:     "/",
		UserID:   user.ID,
		ParentID: nil,
		FullPath: "/",
	}
	err := rootFolder.Create()
	return rootFolder, err
}

func (object *Object) IsRootDir() bool {
	return object.ParentID == nil && object.Type == DIR && object.Name == "/"
}

func (object *Object) Create() error {
	return db.Create(object).Error
}

func (object *Object) Update() error {
	return db.Save(object).Error
}

func (user User) NewFile(parent Object, name string, size uint64, MIMEType string) (*Object, error) {
	file := &Object{
		Type:     FILE,
		Name:     name,
		UserID:   user.ID,
		ParentID: &parent.ID,
		FullPath: path.Join(parent.FullPath, name),
		Size:     &size,
		MIMEType: &MIMEType,
	}
	err := file.Create()
	return file, err
}

func (user User) NewDir(parent Object, name string) (*Object, error) {
	dir := &Object{
		Type:     DIR,
		Name:     name,
		UserID:   user.ID,
		ParentID: &parent.ID,
		FullPath: path.Join(parent.FullPath, name),
	}
	err := dir.Create()
	return dir, err
}

func (object Object) GetChildObjects() ([]Object, error) {
	var objects []Object
	result := db.Where("user_id = ? and parent_id = ?", object.UserID, object.ID).
		Order("type desc").
		Find(&objects)
	return objects, result.Error
}

func (user User) GetObjectByID(ID uint) (*Object, error) {
	object := &Object{}
	result := db.Where("user_id = ?", user.ID).First(object, ID)
	return object, result.Error
}

func (user User) DeleteObjectByID(ID uint) error {
	err := db.Where("user_id = ?", user.ID).Delete(&Object{}, ID).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	} else {
		return nil
	}
}

func (object *Object) Delete() error {
	if object.IsDir() {
		if err := object.DeleteChildObjects(); err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
	}
	return db.Delete(object).Error
}

func (object Object) DeleteChildObjects() error {
	return db.Where("user_id = ? and full_path like ?", object.UserID, object.FullPath+"/%").Delete(&Object{}).Error
}

func (user User) GetObjectByNameAndParentID(name string, parentID uint, objectTypes ...ObjectType) (*Object, error) {
	if len(objectTypes) == 0 {
		objectTypes = []ObjectType{FILE, DIR}
	}

	file := &Object{}
	result := db.Where("user_id = ? and name = ? and parent_id = ? and type in (?)", user.ID, name, parentID, objectTypes).First(file)
	return file, result.Error
}

func (user User) GetObjectByFullPath(fullPath string, objectTypes ...ObjectType) (*Object, error) {
	if len(objectTypes) == 0 {
		objectTypes = []ObjectType{FILE, DIR}
	}

	dir := &Object{}
	result := db.Where("user_id = ? and full_path = ? and type in (?)", user.ID, fullPath, objectTypes).First(dir)
	return dir, result.Error
}

func (object *Object) Rename(fullPath string) error {
	if object.IsDir() {
		if err := object.RenameChildObjects(fullPath); err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
	}
	object.FullPath = fullPath
	object.Name = path.Base(fullPath)
	return object.Update()
}

func (object *Object) RenameChildObjects(fullPath string) error {
	return db.Model(&Object{}).Where("user_id = ? and full_path like ?", object.UserID, object.FullPath+"/%").
		Update("full_path", gorm.Expr("replace(full_path, ?, ?)", object.FullPath, fullPath)).
		Error
}
