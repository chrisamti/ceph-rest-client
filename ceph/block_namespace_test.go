package ceph_test

import (
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/chrisamti/ceph-rest-client/ceph"
)

func TestClient_CreateBlockNameSpaceInPool(t *testing.T) {
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

	// create namespace
	var nameSpace = fmt.Sprintf("test-namespace-%s", time.Now().Format(time.RFC3339))
	var pool = "test-pool-1"

	status, errNameSpace := client.CreateBlockNameSpaceInPool(pool, nameSpace)

	if errNameSpace != nil {
		t.Error(errNameSpace)
	}

	if status != http.StatusCreated {
		t.Errorf("expected status 201 - got %d", status)
	}

	// create same namespace again
	status, errNameSpace = client.CreateBlockNameSpaceInPool(pool, nameSpace)

	if !errors.Is(errNameSpace, ceph.ErrNameSpaceAlreadyExists) {
		t.Errorf("expected err %v - got %v", ceph.ErrSnapshotNameIsEmpty, errNameSpace)
	}

	if status != http.StatusBadRequest {
		t.Errorf("expected status 400 - got %d", status)
	}

	// cleanup
	status, errNameSpace = client.DeleteBlockNameSpaceInPool(pool, nameSpace)

	if errNameSpace != nil {
		t.Error(errNameSpace)
	}

	if status != http.StatusNoContent {
		t.Errorf("expected status 204 - got %d", status)
	}

}

func TestClient_GetBlockNameSpaceListInPool(t *testing.T) {
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

}
