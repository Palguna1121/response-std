package entities

type ModelHasRoles struct {
	RoleID    uint   `gorm:"primaryKey"`
	ModelType string `gorm:"size:255;primaryKey"`
	ModelID   uint   `gorm:"primaryKey"`
	Roles     Roles  `gorm:"foreignKey:RoleID"` // Foreign key to Roles
}

func (ModelHasRoles) TableName() string {
	return "model_has_roles"
}
