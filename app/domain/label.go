package domain

import (
	"errors"
	"time"
)

const (
	// LabelStatusActive is downloadable status for its label
	LabelStatusActive = LabelStatus(iota)
	// LabelStatusInactive is not downloadable status for its label
	LabelStatusInactive = LabelStatus(iota)
)

type labelSentencesByLang map[LangID]*LangSentence

// SentenceVerified records last verified date and operator
type SentenceVerified struct {
	VerifiedAt   time.Time
	VerifiedUser UserID
}

// LangSentence manages sentence and verified info
type LangSentence struct {
	Lang            LangID
	Sentence        string
	LastUpdatedAt   time.Time
	LastUpdatedUser UserID
	LastVerified    *SentenceVerified
}

// Label manages sentences by lang
type Label struct {
	ID        LabelID
	Tenant    TenantID
	Key       string
	Note      string
	Tags      []TagID
	Status    LabelStatus
	Sentences labelSentencesByLang
	CreatedAt time.Time
}

// LabelRepository is interface for Label repository
type LabelRepository interface {
	DispenseID() LabelID
	Find(id string) *Label
	Save(l *Label)
}

// LabelFactory provides builder for Label
type LabelFactory struct {
	lRepo LabelRepository
}

// NewLabelFactory returns LabelFactory object
func NewLabelFactory(r LabelRepository) *LabelFactory {
	return &LabelFactory{lRepo: r}
}

// Build returns initialized Label object
func (f *LabelFactory) Build(t TenantID, key string) *Label {
	return &Label{
		ID:        f.lRepo.DispenseID(),
		Tenant:    t,
		Key:       key,
		Tags:      []TagID{},
		Status:    LabelStatusActive,
		Sentences: labelSentencesByLang{},
		CreatedAt: time.Now(),
	}
}

func (m labelSentencesByLang) Fill(ln LangID, s string, u UserID) labelSentencesByLang {
	ns := labelSentencesByLang{}
	for k, v := range m {
		if k != ln {
			ns[k] = v
		}
	}
	ns[ln] = &LangSentence{
		Lang:            ln,
		Sentence:        s,
		LastUpdatedAt:   time.Now(),
		LastUpdatedUser: u,
	}
	return ns
}

func (m labelSentencesByLang) Verify(ln LangID, u UserID) (labelSentencesByLang, bool) {
	ns := labelSentencesByLang{}
	verified := false
	for k, v := range m {
		if k == ln {
			ns[k] = v.Verify(u)
			verified = true
		} else {
			ns[k] = v
		}
	}

	return ns, verified
}

// FillLangSentence stores sentence by lang. IsVerified flag is forced to set false
func (l *Label) FillLangSentence(ln LangID, s string, u UserID) *Label {
	return &Label{
		ID:        l.ID,
		Tenant:    l.Tenant,
		Key:       l.Key,
		Note:      l.Note,
		Tags:      l.Tags,
		Sentences: l.Sentences.Fill(ln, s, u),
		Status:    l.Status,
		CreatedAt: l.CreatedAt,
	}
}

// GetSentence returns langSentence object
func (l *Label) GetSentence(ln LangID) (*LangSentence, error) {
	if s, ok := l.Sentences[ln]; ok {
		return s, nil
	}
	return nil, errors.New("not found")
}

// Verify returns itself but filled verified data
func (s *LangSentence) Verify(u UserID) *LangSentence {
	return &LangSentence{
		Lang:            s.Lang,
		Sentence:        s.Sentence,
		LastUpdatedAt:   s.LastUpdatedAt,
		LastUpdatedUser: s.LastUpdatedUser,
		LastVerified: &SentenceVerified{
			VerifiedAt:   time.Now(),
			VerifiedUser: u,
		},
	}
}

// VerifyByLang record verified date and operator for specified lang
func (l *Label) VerifyByLang(ln LangID, u UserID) (*Label, bool) {
	ns, verified := l.Sentences.Verify(ln, u)
	lbl := &Label{
		ID:        l.ID,
		Tenant:    l.Tenant,
		Key:       l.Key,
		Note:      l.Note,
		Tags:      l.Tags,
		Sentences: ns,
		Status:    l.Status,
		CreatedAt: l.CreatedAt,
	}

	return lbl, verified
}

// Activate enables to download label
func (l *Label) Activate() *Label {
	return &Label{
		ID:        l.ID,
		Tenant:    l.Tenant,
		Key:       l.Key,
		Note:      l.Note,
		Tags:      l.Tags,
		Sentences: l.Sentences,
		Status:    LabelStatusActive,
		CreatedAt: l.CreatedAt,
	}
}

// Deactivate disables to download label
func (l *Label) Deactivate() *Label {
	return &Label{
		ID:        l.ID,
		Tenant:    l.Tenant,
		Key:       l.Key,
		Note:      l.Note,
		Tags:      l.Tags,
		Sentences: l.Sentences,
		Status:    LabelStatusInactive,
		CreatedAt: l.CreatedAt,
	}

}
