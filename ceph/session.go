package ceph

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
)

// see https://docs.ceph.com/en/pacific/mgr/ceph_api/index.html

// ErrUserNameEmpty is returned if param username in Login() is empty.
var ErrUserNameEmpty = errors.New("param user can not be empty")

// ErrPasswordEmpty is returned if param password in Login() is empty.
var ErrPasswordEmpty = errors.New("param user can not be empty")

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Token struct {
	Token string `json:"token"`
}

type Auth struct {
	Token       string `json:"token"`
	Username    string `json:"username"`
	Permissions struct {
		CephFS            []string `json:"cephfs"`
		ConfigOpt         []string `json:"config-opt"`
		DashboardSettings []string `json:"dashboard-settings"`
		Grafana           []string `json:"grafana"`
		Hosts             []string `json:"hosts"`
		Iscsi             []string `json:"iscsi"`
		Log               []string `json:"log"`
		Manager           []string `json:"manager"`
		Monitor           []string `json:"monitor"`
		NfsGanesha        []string `json:"nfs-ganesha"`
		Osd               []string `json:"osd"`
		Pool              []string `json:"pool"`
		Prometheus        []string `json:"prometheus"`
		RbdImage          []string `json:"rbd-image"`
		RbdMirroring      []string `json:"rbd-mirroring"`
		Rgw               []string `json:"rgw"`
		User              []string `json:"user"`
	} `json:"permissions"`
	PwdExpirationDate interface{} `json:"pwdExpirationDate"`
	Sso               bool        `json:"sso"`
	PwdUpdateRequired bool        `json:"pwdUpdateRequired"`
}

type Server struct {
	Address            string
	Port               uint
	Protocol           string
	APIPath            string
	InsecureSkipVerify bool
}

func (server *Server) getURL(subPath string) string {
	return fmt.Sprintf("%s://%s:%d/%s/%s",
		server.Protocol,
		server.Address,
		server.Port,
		server.APIPath,
		subPath)
}

type Session struct {
	Client *resty.Client
	Server Server
	Auth   Auth
}

var (
	defaultHeaders = map[string]string{
		"Accept": "application/vnd.ceph.api.v1.0+json",
	}

	defaultHeaderJson = map[string]string{
		"Accept":       "application/vnd.ceph.api.v1.0+json",
		"Content-type": "application/json",
	}
)

func NewSession(server Server) (session *Session, err error) {

	session = &Session{
		Client: resty.New(),
		Server: server,
	}

	// session.Client.SetCookieJar(nil)

	if server.InsecureSkipVerify {
		session.Client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}

	return session, nil
}

// Login log in to ceph rest api (https://docs.ceph.com/en/latest/mgr/ceph_api/#post--api-auth)
func (s *Session) Login(username, password string) (status int, err error) {
	var resp *resty.Response

	authBody := Credentials{
		Username: username,
		Password: password,
	}

	if username == "" {
		return 0, ErrUserNameEmpty
	}

	if password == "" {
		return 0, ErrPasswordEmpty
	}

	resp, err = s.Client.R().
		SetHeaders(defaultHeaderJson).
		SetBody(authBody).
		SetResult(&s.Auth).
		Post(s.Server.getURL("auth"))

	if err != nil {
		return 0, err
	}

	if !resp.IsSuccess() {
		return resp.StatusCode(), fmt.Errorf("could not login: %v", resp.Error())
	}

	return resp.StatusCode(), err
}

// Logout from ceph rest api.
func (s *Session) Logout() (err error) {
	var resp *resty.Response

	resp, err = s.Client.R().SetHeaders(defaultHeaders).Post(s.Server.getURL("auth/logout"))

	if err != nil {
		return err
	}

	if !resp.IsSuccess() {
		return fmt.Errorf(resp.String())
	}

	s.Auth.Token = ""

	return err
}