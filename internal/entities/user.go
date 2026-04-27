package entities

type User struct {
	UUID      string `gorm:"column:uuid;type:uuid;primaryKey;default:gen_random_uuid()"`
	Nickname  string `gorm:"column:nickname;uniqueIndex;not null"`
	Name      string `gorm:"column:name;not null"`
	Surname   string `gorm:"column:surname;not null"`
	AvatarUrl string `gorm:"column:avatar_url"`
}

func (User) TableName() string { return "users" }
