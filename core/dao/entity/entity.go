package entity

// 用户表
type User struct {
	ID    int    `gorm:"column:id;primaryKey;autoIncrement"`
	Name  string `gorm:"column:name;type:varchar(255);unique;not null"`
	Votes int    `gorm:"column:votes;type:int;default:0"`
}

// 票据表
type Ticket struct {
	ID        int    `gorm:"column:id;primaryKey;autoIncrement"`
	Ticket    string `gorm:"column:ticket;type:varchar(255);unique;not null"`
	CreatedAt int64  `gorm:"column:created_at;autoCreateTime"`
	ExpiresAt int64  `gorm:"column:expires_at;type:int;not null"`
	MaxUsage  int    `gorm:"column:max_usage;type:int;not null"`
	UsedCount int    `gorm:"column:used_count;type:int;default:0"`
}

func (u *User) TableName() string {
	return "User"
}

func (t *Ticket) TableName() string {
	return "Ticket"
}
