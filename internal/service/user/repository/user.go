package repository

import "github.com/just4n4cc/tp-sem2-db/internal/models"

type User struct {
	Id       int    `db:"id"`
	Nickname string `db:"nickname"`
	Fullname string `db:"fullname"`
	About    string `db:"about"`
	Email    string `db:"email"`
}

func jsonToDbModel(u *models.User) *User {
	return &User{
		Id:       0,
		Nickname: u.Nickname,
		Fullname: u.Fullname,
		About:    u.About,
		Email:    u.Email,
	}
}

func dbToJsonModel(u *User) *models.User {
	return &models.User{
		Nickname: u.Nickname,
		Fullname: u.Fullname,
		About:    u.About,
		Email:    u.Email,
	}
}
