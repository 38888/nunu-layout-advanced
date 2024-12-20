package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        int64  `gorm:"primarykey"`
	UserId    int64  `gorm:"unique;not null"`
	Nickname  string `gorm:"not null"`
	Password  string `gorm:"not null"`
	Email     string `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (u *User) TableName() string {
	return "users"
}
