package users

import (
	"github.com/raedahgroup/fileman/files"
	"github.com/raedahgroup/fileman/rules"
)

// ViewMode describes a view mode.
type ViewMode string

const (
	ListViewMode   ViewMode = "list"
	MosaicViewMode ViewMode = "mosaic"
)

// User describes a user.
type User struct {
	ID           uint          `storm:"id,increment" json:"id"`
	Username     string        `storm:"unique" json:"username"`
	Password     string        `json:"password"`
	Scope        string        `json:"scope"`
	Locale       string        `json:"locale"`
	LockPassword bool          `json:"lockPassword"`
	ViewMode     ViewMode      `json:"viewMode"`
	Perm         Permissions   `json:"perm"`
	Sorting      files.Sorting `json:"sorting"`
	Rules        []rules.Rule  `json:"rules"`
}

// GetRules implements rules.Provider.
func (u *User) GetRules() []rules.Rule {
	return u.Rules
}

var checkableFields = []string{
	"Username",
	"Password",
	"Scope",
	"ViewMode",
	"Sorting",
	"Rules",
}
