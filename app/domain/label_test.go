package domain_test

import (
	"testing"

	"github.com/taiyoh/labeltile/app/domain"
)

func TestLabel(t *testing.T) {
	l := domain.NewLabel(domain.LabelID("1"), domain.TenantID("tenant"), "foo", domain.CategoryID("bar"))
	if l.Key != "foo" {
		t.Error("label.key should be 'foo'")
	}
	if l.Category.String() != "bar" {
		t.Error("label.category should be 'bar'")
	}
	if int(l.Status) != int(domain.LabelStatusActive) {
		t.Error("label.active should be true")
	}
	if len(l.Sentences) > 0 {
		t.Error("label.sentences should be nothing")
	}
	l = l.Deactivate()
	if int(l.Status) != int(domain.LabelStatusInactive) {
		t.Error("deactivate makes Active flag false")
	}
	l = l.Activate()
	if int(l.Status) != int(domain.LabelStatusActive) {
		t.Error("activate makes Active flag true")
	}

	var ok bool
	if l, ok = l.VerifyByLang(domain.LangID("ja"), domain.UserID("1")); ok {
		t.Error("no sentences exists")
	}

	l = l.FillLangSentence(domain.LangID("ja"), "hoge", domain.UserID("1"))
	if len(l.Sentences) != 1 {
		t.Error("label.sentences should be nothing")
	}

	if _, err := l.GetSentence(domain.LangID("en")); err == nil {
		t.Error("lang:en not registered")
	}

	s, err := l.GetSentence(domain.LangID("ja"))
	if err != nil {
		t.Error("lang:ja should exists")
	}

	if s.LastVerified != nil {
		t.Error("sentence is not verified")
	}

	l, ok = l.VerifyByLang(domain.LangID("ja"), domain.UserID("1"))

	if !ok {
		t.Error("verify failed")
	}

	s, _ = l.GetSentence(domain.LangID("ja"))
	if s.LastVerified == nil {
		t.Error("sentence is verified")
	}
	if s.LastVerified.VerifiedAt.IsZero() {
		t.Error("LastVerifiedAt is recorded")
	}
	if s.LastVerified.VerifiedUser.String() == "" {
		t.Error("LastVerifiedUser is recorded")
	}

	l = l.FillLangSentence(domain.LangID("ja"), "fuga", domain.UserID("2"))
	s2, _ := l.GetSentence(domain.LangID("ja"))
	if s2.Sentence != "fuga" {
		t.Error("sentence should be 'fuga'")
	}
	if s2.LastVerified != nil {
		t.Error("sentence is not verified")
	}
}
