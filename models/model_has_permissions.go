package models

type ModelHasPermissions struct {
	PermissionID uint   `gorm:"primaryKey"`
	ModelType    string `gorm:"size:255;primaryKey"`
	ModelID      uint   `gorm:"primaryKey"` // Ini akan menjadi foreign key ke user
	Permission   Permission
}
