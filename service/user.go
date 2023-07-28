package service

import "douyin/db"

func UserExists(username string) (bool, error) {
	return true, nil
}

func UserCreate(username string, password string) (*db.User, error) {
	return &db.User{}, nil
}
