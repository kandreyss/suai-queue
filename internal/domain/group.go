package domain

type Group struct {
	Number  string `gorm:"primaryKey;type:varchar(7);not null;uniqueIndex"`
	Counter int    `gorm:"default:0"`
}
