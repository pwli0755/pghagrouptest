// Package main ...
package main

import (
	"context"
	"gorm.io/driver/postgres"
	"gorm.io/plugin/dbresolver"

	"gorm.io/gorm"
)

var db *gorm.DB

func InitGorm() {
	var err error
	dsnPrimary := "postgres://haha_user:secret@127.0.0.1:5433,127.0.0.1:5432/haha?sslmode=disable&target_session_attrs=primary"
	dsnPreferStandby := "postgres://haha_user:secret@127.0.0.1:5433,127.0.0.1:5432/haha?sslmode=disable&target_session_attrs=prefer-standby"
	//dsn := "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
	db, err = gorm.Open(postgres.Open(dsnPrimary), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.Use(dbresolver.Register(dbresolver.Config{
		Sources:  []gorm.Dialector{postgres.Open(dsnPrimary)},
		Replicas: []gorm.Dialector{postgres.Open(dsnPreferStandby)},
		Policy:   dbresolver.RandomPolicy{},
	}))
	db = db.Debug()
}

type User struct {
	Name string
	Age  int
}

func CreateUserGorm(ctx context.Context, db *gorm.DB) error {
	return db.Model(&User{}).WithContext(ctx).Create(&User{
		Name: "gorm",
		Age:  18,
	}).Error
}

func QueryUserGorm(ctx context.Context, db *gorm.DB) (user User, err error) {
	err = db.Model(&User{}).WithContext(ctx).Where(User{Name: "gorm"}).First(&user).Error
	return
}
