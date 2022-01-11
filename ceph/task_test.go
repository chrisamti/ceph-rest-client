package ceph_test

import (
	"github.com/chrisamti/ceph-rest-client/ceph"
	"net/http"
	"testing"
)

func TestClient_GetTask(t *testing.T) {

	client, err := ceph.New(getServer())

	if err != nil {
		t.Fatal(err)
	}

	statusLogin, errLogin := client.Session.Login(username, password)
	if errLogin != nil {
		t.Error(errLogin)
	}

	if statusLogin != http.StatusCreated {
		t.Fatalf("could not login - expected http state 201 - got %d", statusLogin)
	}

	status, tasks, errTask := client.GetTask()

	if errTask != nil {
		t.Error(errTask)
	}

	t.Logf("got status: %d", status)

	t.Log(tasks)
}
