package dbstore

import (
	"gorm.io/gorm"
	"goth/internal/hash"
	"goth/internal/store"
	"time"
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

// TODO: add tests to mock.go for additional methods

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
	err := s.db.Model(&store.User{}).Where("email = ?", email).First(&user).Error

	if err != nil {
		return nil, err
	}
	return &user, err
}

func (s *UserStore) UpdateUserFirstName(userID uint, fname string) error {

	err := s.db.Model(&store.User{}).Where("id = ? AND active = ?", userID, true).
		Update("first_name", fname).
		Error

	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) UpdateUserLastName(userID uint, lname string) error {

	err := s.db.Model(&store.User{}).Where("id = ? AND active = ?", userID, true).
		Update("last_name", lname).
		Error

	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) UpdateUserEmail(userID uint, email string) error {

	err := s.db.Model(&store.User{}).Where("id = ? AND active = ?", userID, true).
		Update("email", email).
		Error

	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) UpdateUserPassword(userID uint, password string) error {

	hashedPassword, err := s.passwordhash.GenerateFromPassword(password)

	if err != nil {
		return err
	}

	return s.db.Model(&store.User{}).
		Where("id = ?", userID).
		Update("password", hashedPassword).
		Error
}

func (s *UserStore) DeleteUser(userID uint) error {

	err := s.db.Delete(&store.User{}, userID).Error

	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) SetIsActive(userID uint) error {

	err := s.db.Model(&store.User{}).Where("id = ?", userID).Update("active", true).Error

	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) SetInactive(userID uint) error {

	err := s.db.Model(&store.User{}).Where("id = ?", userID).Update("active", false).Error

	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) IncrementLoginAttempts(user *store.User) error {
	if user.LoginAttempts >= 5 {
		err := s.SetLockOut(user)
		if err != nil {
			return err
		}
	} else {
		attempts := user.LoginAttempts + 1
		err := s.db.Model(&store.User{}).Where("id = ?", user.ID).Update("login_attempts", attempts).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *UserStore) ResetLoginAttempts(user *store.User) error {
	err := s.db.Model(&store.User{}).Where("id = ?", user.ID).Update("login_attempts", 0).Error
	if err != nil {
		return err
	}

	err = s.db.Model(&store.User{}).Where("id = ?", user.ID).Update("locked_out", false).Error
	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) SetLockOut(user *store.User) error {
	err := s.db.Model(&store.User{}).Where("id = ?", user.ID).Update("locked_out", true).Error
	if err != nil {
		return err
	}

	duration, err := time.ParseDuration("5m")
	if err != nil {
		return err
	}

	err = s.db.Model(&store.User{}).
		Where("id = ?", user.ID).
		Update("lockout_duration", time.Now().Add(duration)).
		Error
	if err != nil {
		return err
	}

	return nil
}
