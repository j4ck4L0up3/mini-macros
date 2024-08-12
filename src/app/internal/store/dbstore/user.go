package dbstore

import (
	"goth/internal/hash"
	"goth/internal/store"

	"gorm.io/gorm"
)

type UserStore struct {
	db           *gorm.DB
	passwordhash hash.PasswordHash
}

type NewUserStoreParams struct {
	DB           *gorm.DB
	PasswordHash hash.PasswordHash
}

func NewUserStore(params NewUserStoreParams) *UserStore {
	return &UserStore{
		db:           params.DB,
		passwordhash: params.PasswordHash,
	}
}

// TODO: add additional interface methods (see store.go)

func (s *UserStore) CreateUser(fname, lname, email, password string) error {

	hashedPassword, err := s.passwordhash.GenerateFromPassword(password)
	if err != nil {
		return err
	}

	return s.db.Create(&store.User{
		FirstName: fname,
		LastName:  lname,
		Email:     email,
		Password:  hashedPassword,
	}).Error
}

func (s *UserStore) GetUser(email string) (*store.User, error) {

	var user store.User
	err := s.db.Where("email = ?", email).First(&user).Error

	if err != nil {
		return nil, err
	}
	return &user, err
}

func (s *UserStore) UpdateUserFirstName(userID uint, fname string) error {

	var user store.User
	err := s.db.Where(&store.User{ID: userID, Active: true}).
		First(&user).
		Update("first_name", fname).
		Error

	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) UpdateUserLastName(userID uint, lname string) error {

	var user store.User
	err := s.db.Where(&store.User{ID: userID, Active: true}).
		First(&user).
		Update("last_name", lname).
		Error

	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) UpdateUserEmail(userID uint, email string) error {

	var user store.User
	err := s.db.Where(&store.User{ID: userID, Active: true}).
		First(&user).
		Update("email", email).
		Error

	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) UpdateUserPassword(userID uint, password string) error {

	var user store.User
	hashedPassword, err := s.passwordhash.GenerateFromPassword(password)

	if err != nil {
		return err
	}

	return s.db.Where("id = ?", userID).First(&user).Update("password", hashedPassword).Error
}

func (s *UserStore) DeleteUser(userID uint) error {

	err := s.db.Delete(&store.User{}, userID).Error

	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) SetIsActive(userID uint) error {

	var user store.User
	err := s.db.Where("id = ?", userID).First(&user).Update("active", true).Error

	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) SetInactive(userID uint) error {

	var user store.User
	err := s.db.Where("id = ?", userID).First(&user).Update("active", false).Error

	if err != nil {
		return err
	}

	return nil
}
