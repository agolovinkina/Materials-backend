package ds

type User struct {
	UserID      uint   `gorm:"primaryKey;column:user_id"`
	Login       string `gorm:"type:varchar(150);unique;not null"`
	Password    string `gorm:"type:varchar(128);not null"`
	IsModerator bool   `gorm:"type:boolean;not null;default:false"`
}
