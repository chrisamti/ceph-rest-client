package ceph_test

import (
	"github.com/chrisamti/ceph-rest-client/ceph"
	"net/http"
	"testing"
)

var (
	username = "test-user"
	password = "XJEGy5yWrYxu758"
)

func getServer() ceph.Server {
	return ceph.Server{
		Address:            "10.10.2.23",
		Port:               8443,
		Protocol:           "https",
		APIPath:            "api",
		InsecureSkipVerify: true,
	}
}

func TestClient_CreateBlockImage(t *testing.T) {
	client, err := ceph.New(getServer())

	if err != nil {
		t.Fatal(err)
	}

	status, errLogin := client.Session.Login(username, password)
	if errLogin != nil {
		t.Error(errLogin)
	}

	if status != http.StatusCreated {
		t.Fatalf("could not login - expected http state 201 - got %d", status)
	}

	rbd := ceph.RBDCreate{
		Features:      nil,
		PoolName:      "k14-pool01",
		Namespace:     nil,
		Name:          "rest-client-test-1",
		Size:          1073741824,
		ObjSize:       0,
		StripeUnit:    nil,
		StripeCount:   nil,
		DataPool:      nil,
		Configuration: struct{}{},
	}

	status, err = client.CreateBlockImage(rbd)

	if err != nil {
		t.Error(err)
	}

	if status != http.StatusCreated {
		t.Errorf("expected http state 201 - got %d", status)
	}

}

func TestClient_DeleteBlockImage(t *testing.T) {
	client, err := ceph.New(getServer())

	if err != nil {
		t.Fatal(err)
	}

	status, errLogin := client.Session.Login(username, password)
	if errLogin != nil {
		t.Error(errLogin)
	}

	if status != http.StatusCreated {
		t.Fatalf("could not login - expected http state 201 - got %d", status)
	}

	status, err = client.DeleteBlockImage("k14-pool01/rest-client-test-1")

	if status != http.StatusNoContent {
		t.Errorf("expected http state 204 - got %d", status)
	}

}

func TestClient_ListBlockImage(t *testing.T) {

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

	status, block, errBlock := client.ListBlockImage("k14-csi-0")

	if errBlock != nil {
		t.Error(errBlock)
	}

	if status != http.StatusOK {
		t.Errorf("expected http state 200 - got %d", status)
	}

	if len(block) == 0 {
		t.Error("expected more than 0 block images")
	}

	for _, v := range block {
		t.Logf("pool: %s", v.PoolName)
		t.Logf("status: %d", v.Status)
		for _, b := range v.Value {
			t.Logf("\tID: %s\tname: %s", b.ID, b.Name)
		}
	}

	status, block, errBlock = client.ListBlockImage("")

	if errBlock != nil {
		t.Error(errBlock)
	}

	if status != http.StatusOK {
		t.Errorf("expected http state 200 - got %d", status)
	}

	if len(block) == 0 {
		t.Error("expected more than 0 block images")
	}

	for _, v := range block {
		t.Logf("pool: %s", v.PoolName)
		t.Logf("status: %d", v.Status)
		for _, b := range v.Value {
			t.Logf("\tID: %s\tname: %s", b.ID, b.Name)
		}
	}
}

func TestClient_GetBlockImage(t *testing.T) {
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

	status, rbd, errRbd := client.GetBlockImage("k14-pool01/vm-101071-disk-0")

	if errRbd != nil {
		t.Error(errRbd)
	}

	if status != http.StatusOK {
		t.Errorf("expected http state 200 - got %d", status)
	}

	t.Logf("ID: %s", rbd.ID)

}
