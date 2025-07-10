package entities

type ModelHasPermissions struct {
	PermissionID uint       `gorm:"primaryKey"`
	ModelType    string     `gorm:"size:255;primaryKey"`
	ModelID      uint       `gorm:"primaryKey"`
	Permission   Permission `gorm:"foreignKey:PermissionID"` // tambah foreignKey
}

func (ModelHasPermissions) TableName() string {
	return "model_has_permissions"
}
