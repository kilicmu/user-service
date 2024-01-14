package database

import "time"

type UserInfo struct {
	Id            int64     `gorm:"column:id;primary_key;autoIncrement;<-:create" json:"id"`
	UId           string    `gorm:"column:uid;primary_key;not null;uniqueIndex;<-:create" json:"uid"`
	Name          string    `gorm:"column:name" json:"name"`
	Password      string    `gorm:"column:password" json:"password"`
	Email         string    `gorm:"column:email;not null;unique" json:"email"`
	Phone         string    `gorm:"column:phone" json:"phone"`
	Picture       string    `gorm:"column:picture" json:"picture"`
	EmailVerified bool      `gorm:"column:email_verified" json:"email_verified"`
	Channel       string    `gorm:"column:channel" json:"channel"`
	Identify      string    `gorm:"column:identify" json:"identify"`
	CreateAt      time.Time `gorm:"column:create_at;autoCreateTime" json:"create_at"`
	UpdateAt      time.Time `gorm:"column:update_at;autoUpdateTime" json:"update_at"`
}
