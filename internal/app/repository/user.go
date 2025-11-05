package repository

import (
	"fmt"
	"lr4/internal/app/ds"
)

func (r *Repository) CreateUser(user *ds.User) error {
	var existingUser ds.User
	if err := r.DB.Where("login = ?", user.Login).First(&existingUser).Error; err == nil {
		return fmt.Errorf("user already exists")
	}

	return r.DB.Create(user).Error
}

func (r *Repository) GetUserByLogin(login string) (*ds.User, error) {
	var user ds.User
	if err := r.DB.Where("login = ?", login).First(&user).Error; err != nil {
		return nil, fmt.Errorf("user not found")
	}
	return &user, nil
}

func (r *Repository) GetUserByID(userID int) (*ds.User, error) {
	var user ds.User
	if err := r.DB.Where("user_id = ?", userID).First(&user).Error; err != nil {
		return nil, fmt.Errorf("user not found")
	}
	return &user, nil
}

func (r *Repository) UpdateUser(userID int, login, password string) (*ds.User, error) {
	var user ds.User
	if err := r.DB.Where("user_id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}

	// Check if login is being changed and if it's already taken
	if login != "" && login != user.Login {
		var count int64
		r.DB.Model(&ds.User{}).Where("login = ? AND user_id != ?", login, userID).Count(&count)
		if count > 0 {
			return nil, fmt.Errorf("login already taken")
		}
		user.Login = login
	}

	if password != "" {
		user.Password = password
	}

	if err := r.DB.Save(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) RegisterUser(input ds.ChangeUserDTO) error {
	user := &ds.User{
		Login:    input.Login,
		Password: input.Password,
	}

	return r.DB.Create(user).Error
}

func (r *Repository) LoginUser(login, password string) (*ds.User, error) {
	var user ds.User
	err := r.DB.Where("login = ? AND password = ?", login, password).First(&user).Error
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}
	return &user, nil
}
