package main

import (
	"bytes"
	"crypto"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/stretchr/goweb/context"
	_ "github.com/ziutek/mymysql/godrv"
	"menteslibres.net/gosexy/checksum"
	"net/http"
	"strings"
)

const (
	GITHUB_API         = "https://api.github.com"
	USER_AGENT         = "robxu9, Reservoir Build Server | Auth"
	ErrAuthFailed      = fmt.Errorf("Failed to authenticate with given tickets.")
	AuthUserNameHeader = "X-Reservoir-User"
)

var client *http.Client

func init() {
	client = &http.Client{}
}

func Auth_GetUser(c *context.Context) string {
	return c.HttpRequest().Header.Get(AuthUserNameHeader)
}

type AuthHandler struct {
	SuccessHandler *http.Handler
}

func (a *AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, fmt.Errorf("No authorzation header given."), http.StatusUnauthorized)
		return
	}

	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) != 2 {
		http.Error(w, fmt.Errorf("Malformed authorization header."), http.StatusBadRequest)
		return
	}

	authString, err := base64.StdEncoding.DecodeString(authHeaderParts[1])

	if err != nil {
		http.Error(w, fmt.Errorf("Malformed token provided."), http.StatusBadRequest)
		return
	}

	authed := false
	var reason error = nil

	switch authHeaderParts[0] {
	case "Basic":
		basicAuth := strings.Split(authString, ":")
		user, err := Auth_BasicCheck(basicAuth[0], basicAuth[1])
		if err != nil {
			reason = err
		} else {
			authed = true
			r.Header.Del(AuthUserNameHeader)
			r.Header.Add(AuthUserNameHeader, user)
		}
	case "Hash":
		basicAuth := strings.Split(authString, ":")
		user, err := Auth_HashedCheck(basicAuth[0], basicAuth[1])
		if err != nil {
			reason = err
		} else {
			authed = true
			r.Header.Del(AuthUserNameHeader)
			r.Header.Add(AuthUserNameHeader, user)
		}
	case "Token":
		user, err := Auth_TokenCheck(authString)
		if err != nil {
			reason = err
		} else {
			authed = true
			r.Header.Del(AuthUserNameHeader)
			r.Header.Add(AuthUserNameHeader, user)
		}
	default:
		reason = fmt.Errorf("Couldn't auth with procedure.")
	}

	if !authed {
		http.Error(w, fmt.Errorf("%s auth failed: %s", authHeaderParts[0], reason.Error()), http.StatusUnauthorized)
		return
	}

	a.SuccessHandler.ServeHTTP(w, r)
}

func Auth_BasicCheck(username, password string) (string, error) {
	hashed := checksum.String(password, crypto.SHA512)
	return Auth_HashedCheck(username, hashed)
}

func Auth_HashedCheck(username, hash string) (string, error) {
	db, err := Model_DB()
	if err != nil {
		return "", err
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT username FROM authentication WHERE username='?' and password='?' LIMIT 1")
	if err != nil {
		return "", err
	}

	row := stmt.QueryRow(username, hash)

	var dbusername string
	err = row.Scan(&dbusername)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", ErrAuthFailed
		}
		return "", err
	}

	return dbusername, nil
}

func Auth_TokenCheck(token string) (string, error) {
	db, err := Model_DB()
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT username FROM authentication WHERE token='?' LIMIT 1")
	if err != nil {
		return "", err
	}

	row := stmt.QueryRow(token)

	var dbusername string
	err = row.Scan(&dbusername)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", ErrAuthFailed
		}
		return "", err
	}

	return dbusername, nil
}

type GithubAuth struct {
	Id         uint64
	Url        string
	Scopes     []string
	Token      string
	App        map[string]string
	Note       string
	Note_URL   string
	Updated_At string
	Created_At string
	User       GithubAuthUser
}

type GithubAuthUser struct {
	Login       string
	Id          uint64
	Avatar_URL  string
	Gravatar_ID string
	URL         string
}

func Auth_GithubCheck(token string) (*GithubAuthUser, bool, error) {
	url := "/applications/:client_id/tokens/:access_token"

	config := make(map[string]map[string]string)
	err := Config_GetConfig("auth", &config)
	if err != nil {
		return nil, false, err
	}

	client_id := config[Config_Environment]["client_id"]
	client_secret := config[Config_Environment]["client_secret"]

	url = strings.Replace(url, ":client_id", client_id, -1)
	url = strings.Replace(url, ":access_token", token, -1)

	req := http.NewRequest("GET", GITHUB_API+url, nil)
	req.SetBasicGithubAuth(client_id, client_secret)
	req.Header.Set("User-Agent", USER_AGENT)

	resp, err := client.Do(req)
	if err != nil {
		return nil, false, err
	}

	if resp.StatusCode == 404 {
		return nil, false, nil
	}

	var b bytes.Buffer

	_, err = b.ReadFrom(resp.Body)
	if err != nil {
		return nil, true, err
	}

	var GithubAuthResp GithubAuth

	err = json.Unmarshal(b.Bytes(), &authResp)
	if err != nil {
		return nil, true, err
	}

	return &authResp.User, true, nil

}
