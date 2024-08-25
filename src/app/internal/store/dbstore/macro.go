package dbstore

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"goth/internal/store"
)

type MacroStore struct {
	db *gorm.DB
}

type MacroStoreParams struct {
	DB *gorm.DB
}

func NewMacroStore(params MacroStoreParams) *MacroStore {
	return &MacroStore{
		db: params.DB,
	}
}

func (s *MacroStore) CreateMacro(macro *store.Macro) error {

	macro.MacroCookieID = uuid.New().String()

	err := s.db.Create(macro).Error

	if err != nil {
		return err
	}

	return nil
}

func (s *MacroStore) GetAllMacrosFromUser(userID uint) ([]*store.Macro, error) {

	var macros []*store.Macro
	err := s.db.Model(&store.User{}).Where("user_id = ?", userID).Find(&macros).Error

	if err != nil {
		return nil, err
	}

	return macros, nil
}

func (s *MacroStore) GetMacrosFromQuery(query string) ([]*store.Macro, error) {

	var macros []*store.Macro
	err := s.db.Raw("SELECT name FROM macros WHERE to_tsvector(name, content) @@ to_tsquery(?);", query).
		Scan(&macros).
		Error

	if err != nil {
		return nil, err
	}

	return macros, nil
}

func (s *MacroStore) UpdateMacroName(name string, macroID uint, userID uint) error {

	var macro store.Macro
	err := s.db.Model(&store.Macro{}).
		Where("id = ? AND user_id = ?", macroID, userID).
		First(&macro).
		Update("name", name).
		Error

	if err != nil {
		return err
	}

	return nil
}

func (s *MacroStore) UpdateMacroContent(content string, macroID uint, userID uint) error {

	var macro store.Macro
	err := s.db.Model(&store.Macro{}).
		Where("id = ? AND user_id = ?", macroID, userID).
		First(&macro).
		Update("content", content).
		Error

	if err != nil {
		return err
	}

	return nil
}

func (s *MacroStore) DeleteMacro(macroID, userID uint) error {

	err := s.db.Where("id = ? AND user_id = ?", macroID, userID).Delete(&store.Macro{}).Error

	if err != nil {
		return err
	}

	return nil
}

func (s *MacroStore) IncrementClickCount(macroID, userID uint) error {

	err := s.db.Exec(
		"UPDATE macros SET click_count = click_count + 1 WHERE id = ? AND user_id = ?",
		macroID,
		userID,
	).Error

	if err != nil {
		return err
	}

	return nil
}
