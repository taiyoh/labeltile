package labeltile

import "testing"

func TestLabel(t *testing.T) {
	l := NewLabel(TenantID("tenant"), "foo", CategoryID("bar"))
	if l.Key != "foo" {
		t.Error("label.key should be 'foo'")
	}
	if string(l.Category) != "bar" {
		t.Error("label.category should be 'bar'")
	}
	if !l.Active {
		t.Error("label.active should be true")
	}
	if len(l.Sentences) > 0 {
		t.Error("label.sentences should be nothing")
	}
	l.Deactivate()
	if l.Active {
		t.Error("deactivate makes Active flag false")
	}
	l.Activate()
	if !l.Active {
		t.Error("activate makes Active flag true")
	}

	if ok := l.VerifyByLang(LangID("ja"), UserID("1")); ok {
		t.Error("no sentences exists")
	}

	l.FillLangSentence(LangID("ja"), "hoge", UserID("1"))
	if len(l.Sentences) != 1 {
		t.Error("label.sentences should be nothing")
	}

	if _, err := l.GetSentence(LangID("en")); err == nil {
		t.Error("lang:en not registered")
	}

	s, err := l.GetSentence(LangID("ja"))
	if err != nil {
		t.Error("lang:ja should exists")
	}

	if s.IsVerified {
		t.Error("sentence is not verified")
	}
	if !s.LastVerifiedAt.IsZero() {
		t.Error("LastVerifiedAt is not recorded")
	}
	if string(s.LastVerifiedUser) != "" {
		t.Error("LastVerifiedUser is not recorded")
	}

	l.VerifyByLang(LangID("ja"), UserID("1"))
	if !s.IsVerified {
		t.Error("sentence is verified")
	}
	if s.LastVerifiedAt.IsZero() {
		t.Error("LastVerifiedAt is recorded")
	}
	if string(s.LastVerifiedUser) == "" {
		t.Error("LastVerifiedUser is recorded")
	}

	l.FillLangSentence(LangID("ja"), "fuga", UserID("2"))
	var s2 *langSentence
	s2, _ = l.GetSentence(LangID("ja"))
	if s2.Sentence != "fuga" {
		t.Error("sentence should be 'fuga'")
	}
	if s2.IsVerified {
		t.Error("sentence is not verified")
	}
	if s2.LastVerifiedAt.Unix() != s.LastVerifiedAt.Unix() {
		t.Error("LastVerifiedAt is carried old record")
	}
	if s2.LastVerifiedUser != s.LastVerifiedUser {
		t.Error("LastVerifiedUser is carried old record")
	}

}
