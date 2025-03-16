package model

type Platform struct {
	Id   uint16 `gorm:"uniqueIndex:idx_platform"`
	Name string `gorm:"primaryKey; uniqueIndex:idx_platform"`
}

func (Platform) TableName() string {
	return "platform"
}
