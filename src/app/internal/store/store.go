package store

// TODO: update anything that includes User
type User struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	FirstName string `                  json:"first_name"`
	LastName  string `                  json:"last_name"`
	Email     string `                  json:"email"`
	Password  string `                  json:"-"`
}

type Macro struct {
	ID      uint   `gorm:"primaryKey"        json:"id"`
	MacroID string `                         json:"macro_id"`
	Name    string `                         json:"name"`
	Content string `                         json:"content"`
	UserID  uint   `                         json:"user_id"`
	User    User   `gorm:"foreignKey:UserID" json:"user"`
}

type Session struct {
	ID        uint   `gorm:"primaryKey"        json:"id"`
	SessionID string `                         json:"session_id"`
	UserID    uint   `                         json:"user_id"`
	User      User   `gorm:"foreignKey:UserID" json:"user"`
}

// TODO: implement additional methods
type UserStore interface {
	CreateUser(email, password string) error
	GetUser(email string) (*User, error)
	// additional methods
	// NOTE: may need context?
	UpdateUserFirstName(userID, fname string) error
	UpdateUserLastName(userID, lname string) error
	UpdateUserEmail(userID, email string) error
	UpdateUserPassword(userID, passwordhash string) error
	DeleteUser(userID string) error
}

// TODO: implement MacroStore methods
// NOTE: may need context?
type MacroStore interface {
	CreateMacro(name, content string) error
	GetAllMacrosFromUser(userID string) ([]*Macro, error)
	GetMacroFromName(name string) (*Macro, error)
	UpdateMacroName(macroID, userID string) error
	UpdateMacroContent(macroID, userID string) error
	DeleteMacro(macroID, userID string) error
}

// TODO: may need to add delete session and refresh session
type SessionStore interface {
	CreateSession(session *Session) (*Session, error)
	GetUserFromSession(sessionID, userID string) (*User, error)
}
