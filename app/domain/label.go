package domain

import (
	"errors"
	"time"
)

// Label manages sentence by lang
type Label struct {
	ID        LabelID
	Tenant    TenantID
	Key       string
	Note      string
	Category  CategoryID
	Active    bool
	Sentences map[LangID]*langSentence
	CreatedAt time.Time
}

type langSentence struct {
	Lang             LangID
	Sentence         string
	IsVerified       bool
	LastUpdatedAt    time.Time
	LastUpdatedUser  UserID
	LastVerifiedAt   time.Time
	LastVerifiedUser UserID
}

// NewLabel returns initialized Label object
func NewLabel(t TenantID, key string, catID CategoryID) *Label {
	return &Label{
		Tenant:    t,
		Key:       key,
		Category:  catID,
		Active:    true,
		Sentences: map[LangID]*langSentence{},
		CreatedAt: time.Now(),
	}
}

// FillLangSentence stores sentence by lang. IsVerified flag is forced to set false
func (l *Label) FillLangSentence(ln LangID, s string, u UserID) {
	if ls, exists := l.Sentences[ln]; exists {
		l.Sentences[ln] = &langSentence{
			Lang:             ln,
			Sentence:         s,
			LastUpdatedAt:    time.Now(),
			LastUpdatedUser:  u,
			LastVerifiedAt:   ls.LastVerifiedAt,
			LastVerifiedUser: ls.LastVerifiedUser,
		}
	} else {
		l.Sentences[ln] = &langSentence{
			Lang:            ln,
			Sentence:        s,
			LastUpdatedAt:   time.Now(),
			LastUpdatedUser: u,
		}
	}
}

// GetSentence returns langSentence object
func (l *Label) GetSentence(ln LangID) (*langSentence, error) {
	if s, ok := l.Sentences[ln]; ok {
		return s, nil
	}
	return nil, errors.New("not found")
}

// VerifyByLang exchange IsVerified flag true, and record verified date and user for specified lang
func (l *Label) VerifyByLang(ln LangID, u UserID) bool {
	if ls, exists := l.Sentences[ln]; exists {
		ls.IsVerified = true
		ls.LastVerifiedAt = time.Now()
		ls.LastVerifiedUser = u
		return true
	}
	return false
}

// Activate enables to download label
func (l *Label) Activate() {
	l.Active = true
}

// Deactivate disables to download label
func (l *Label) Deactivate() {
	l.Active = false
}
