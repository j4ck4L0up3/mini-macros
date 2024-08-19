package store

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	FirstName  string `gorm:"size:64"  json:"first_name"`
	LastName   string `gorm:"size:64"  json:"last_name"`
	Email      string `gorm:"size:319" json:"email"`
	Password   string `gorm:"size:256" json:"-"`
	Active     bool   `                json:"active"`
	MacroCount uint   `gorm:"size:16"  json:"macro_count"`
}

type Macro struct {
	gorm.Model
	MacroID string `gorm:"not null;unique"   json:"macro_id"`
	Name    string `gorm:"size:64"           json:"name"`
	Content string `gorm:"not null"          json:"content"`
	UserID  uint   `gorm:"not null"          json:"user_id"`
	User    User   `gorm:"foreignKey:UserID" json:"user"`
}

type Session struct {
	gorm.Model
	SessionID string `gorm:"unique"            json:"session_id"`
	UserID    uint   `gorm:"not null"          json:"user_id"`
	User      User   `gorm:"foreignKey:UserID" json:"user"`
}

type PasswordResetToken struct {
	gorm.Model
	UserID      uint      `gorm:"primaryKey;not null" json:"user_id"`
	Token       string    `gorm:"unique;not null"     json:"token"`
	TokenExpiry time.Time `gorm:"not null"            json:"token_expiry"`
	User        User      `gorm:"foreignKey:UserID"   json:"user"`
}

type UserStore interface {
	CreateUser(fname, lname, email, password string) error
	GetUser(email string) (*User, error)
	UpdateUserFirstName(userID uint, fname string) error
	UpdateUserLastName(userID uint, lname string) error
	UpdateUserEmail(userID uint, email string) error
	UpdateUserPassword(userID uint, passwordhash string) error
	DeleteUser(userID uint) error
	SetIsActive(userID uint) error
	SetInactive(userID uint) error
}

// TODO: create macro.go & implement MacroStore methods
type MacroStore interface {
	CreateMacro(name, content string) error
	GetAllMacrosFromUser(userID string) ([]*Macro, error)
	GetMacroFromName(name string) (*Macro, error)
	UpdateMacroName(macroID, userID string) error
	UpdateMacroContent(macroID, userID string) error
	DeleteMacro(macroID, userID string) error
}

// TODO: may need to add delete session or refresh session
type SessionStore interface {
	CreateSession(session *Session) (*Session, error)
	GetUserFromSession(sessionID, userID string) (*User, error)
}

// TODO: set PasswordResetTokenStore methods
type PasswordResetTokenStore interface {
}
