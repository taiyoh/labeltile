package infra

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"text/template"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type ExportedOAuth2Google = OAuth2Google

func (o *ExportedOAuth2Google) FillConfig(conf *oauth2.Config) {
	o.oauth2Conf = conf
}

func (o *ExportedOAuth2Google) FillUserInfoURL(url string) {
	o.userInfoURL = url
}

func TestNewOAuth2Google(t *testing.T) {
	o := NewOAuth2Google("foo", "bar", "baz")
	u := o.AuthCodeURL("hoge")
	if !strings.HasPrefix(u, google.Endpoint.AuthURL) {
		t.Error("wrong Auth URL")
	}
	parsedURL, _ := url.Parse(u)
	q := parsedURL.Query()
	type testStruct struct {
		Key      string
		Expected string
	}
	for _, s := range []testStruct{
		testStruct{Key: "client_id", Expected: "foo"},
		testStruct{Key: "redirect_uri", Expected: "baz"},
		testStruct{Key: "scope", Expected: "profile email"},
		testStruct{Key: "state", Expected: "hoge"},
	} {
		if v := q.Get(s.Key); v != s.Expected {
			t.Errorf("wrong %s -> %s", s.Key, v)
		}
	}
}

func TestOAuth2GoogleGetTokenInfo(t *testing.T) {
	tStatusCode := 200
	tAccessToken := "hogefugatoken"
	tExpiresIn := 3600
	tTmpl := template.Must(template.New("response").Parse("access_token={{.accessToken}}&expires_in={{.expiresIn}}&scope=profile&scope=email&token_type=bearer"))
	tsToken := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(tStatusCode)
		w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
		b := bytes.NewBuffer([]byte{})
		tTmpl.Execute(b, map[string]interface{}{
			"accessToken": tAccessToken,
			"expiresIn":   tExpiresIn,
		})
		w.Write(b.Bytes())
	}))
	defer tsToken.Close()

	iJSON := `{"email":"foo@example.com","id":"1111111111111"}`
	iTmpl := template.Must(template.New("response").Parse("{{.json}}"))
	tsInfo := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		b := bytes.NewBuffer([]byte{})
		iTmpl.Execute(b, map[string]string{
			"json": iJSON,
		})
		w.Write(b.Bytes())
	}))
	defer tsInfo.Close()

	conf := &oauth2.Config{
		ClientID:     "CLIENT_ID",
		ClientSecret: "CLIENT_SECRET",
		RedirectURL:  "REDIRECT_URL",
		Scopes: []string{
			"profile",
			"email",
		},
		Endpoint: oauth2.Endpoint{
			AuthURL:  tsToken.URL + "/auth",
			TokenURL: tsToken.URL + "/token",
		},
	}

	o := &ExportedOAuth2Google{}
	o.FillConfig(conf)
	o.FillUserInfoURL(tsInfo.URL + "/userinfo/me")

	t.Run("status != 200", func(t *testing.T) {
		tStatusCode = 500
		if _, err := o.GetTokenInfo("code"); err == nil {
			t.Error("fetch token failed")
		}
	})

	t.Run("token expired", func(t *testing.T) {
		tStatusCode = 200
		tExpiresIn = -100
		if _, err := o.GetTokenInfo("code"); err == nil {
			t.Error("invalid token")
		}
	})

	t.Run("broken json", func(t *testing.T) {
		tExpiresIn = 3600
		iJSON = `{""}`
		if _, err := o.GetTokenInfo("code"); err == nil {
			t.Error("json should be broken")
		}
	})
	t.Run("fetched valid token json", func(t *testing.T) {
		tExpiresIn = 3600
		iJSON = `{"email":"foo@example.com","id":"111111"}`
		now := time.Now()
		tk, err := o.GetTokenInfo("code")
		if err != nil {
			t.Error("error found:", err)
		}
		if tk.Email() != "foo@example.com" {
			t.Error("wrong email fetched")
		}
		if tk.UserID() != "111111" {
			t.Error("wrong user id fetched")
		}
		if tk.ExpiresAt().Unix()-now.Unix() != 3600 {
			t.Error("wrong expire_in fetched")
		}
	})

}
