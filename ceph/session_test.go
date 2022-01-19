package ceph_test

import (
	"github.com/chrisamti/ceph-rest-client/ceph"
	"net/http"
	"testing"
)

func TestNewSession(t *testing.T) {
	server := ceph.Server{
		Address:            "192.168.21.30",
		Port:               8443,
		Protocol:           "https",
		APIPath:            "api",
		InsecureSkipVerify: true,
	}
	s, err := ceph.NewSession(server)
	if err != nil {
		t.Error(err)
	}

	status, errLogin := s.Login(username, password)
	if errLogin != nil {
		t.Error(errLogin)
	} else {
		t.Logf("auth token: %s", s.Auth.Token)
	}

	if status != http.StatusCreated {
		t.Errorf("expected http state 201 - got %d", status)
	}

}

func TestSession_Logout(t *testing.T) {
	server := ceph.Server{
		Address:            "10.10.2.23",
		Port:               8443,
		Protocol:           "https",
		APIPath:            "api",
		InsecureSkipVerify: true,
	}
	s, err := ceph.NewSession(server)
	if err != nil {
		t.Error(err)
	}

	status, errLogin := s.Login(username, password)
	if errLogin != nil {
		t.Error(errLogin)
	} else {
		t.Logf("auth token: %s", s.Auth.Token)
	}

	errLogout := s.Logout()

	if errLogout != nil {
		t.Error(errLogout)
	}

	if status != http.StatusCreated {
		t.Errorf("expected http state 201 - got %d", status)
	}

}
