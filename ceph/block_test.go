package ceph_test

import (
	"fmt"
	"github.com/chrisamti/ceph-rest-client/ceph"
	"net/http"
	"testing"
	"time"
)

var (
	username = "test-user"
	password = "XJEGy5yWrYxu758"
	server   = "192.168.21.31"
)

func getServer() ceph.Server {
	return ceph.Server{
		Address:            server,
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
		PoolName:      "test-pool-1",
		Namespace:     nil,
		Name:          "rest-client-one-img-1",
		Size:          1073741824,
		ObjSize:       0,
		StripeUnit:    nil,
		StripeCount:   nil,
		DataPool:      nil,
		Configuration: struct{}{},
	}

	status, err = client.CreateBlockImage(rbd, 0)

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

	status, err = client.DeleteBlockImage("test-pool-1", nil, "rest-client-one-img-1", 0)

	if err != nil {
		t.Error(err)
	}

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

	// create some block devices
	for i := 0; i < 10; i++ {
		_, err := func(i int) (int, error) {
			rbd := ceph.RBDCreate{
				Features:      nil,
				PoolName:      "test-pool-1",
				Namespace:     nil,
				Name:          fmt.Sprintf("rest-client-test-%d", i),
				Size:          1073741824,
				ObjSize:       0,
				StripeUnit:    nil,
				StripeCount:   nil,
				DataPool:      nil,
				Configuration: struct{}{},
			}
			return client.CreateBlockImage(rbd, 0)
		}(i)
		if err != nil {
			t.Error(err)
		}
	}

	status, block, errBlock := client.ListBlockImage("test-pool-1")

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

	time.Sleep(10 * time.Second)

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

	// clean --> delete created image
	for j := 0; j < 10; j++ {
		_, err := func(i int) (int, error) {
			var imageName = fmt.Sprintf("rest-client-test-%d", i)
			return client.DeleteBlockImage("test-pool-1", nil, imageName, 0)
		}(j)

		if err != nil {
			t.Error(err)
		}
	}

}

func TestClient_UpdateBlockImage(t *testing.T) {
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

	rbd := ceph.RBDCreate{
		Features:      nil,
		PoolName:      "test-pool-1",
		Namespace:     nil,
		Name:          "rest-client-one-img-1",
		Size:          1073741824,
		ObjSize:       0,
		StripeUnit:    nil,
		StripeCount:   nil,
		DataPool:      nil,
		Configuration: struct{}{},
	}

	statusCreate, errCreate := client.CreateBlockImage(rbd, 0)

	if errCreate != nil {
		t.Error(errCreate)
	}

	if statusCreate != http.StatusCreated {
		t.Errorf("expected http state 201 - got %d", statusCreate)
	}

	// modify

	var rbdUpdate = ceph.RBDUpdate{
		Features:      nil,
		Name:          "rest-client-one-img-1",
		Size:          int64(rbd.Size * 2),
		Configuration: struct{}{},
	}

	statusModify, errModify := client.UpdateBlockImage("test-pool-1", nil, "rest-client-one-img-1", rbdUpdate, 0)

	if errModify != nil {
		t.Error(err)
	}

	if statusModify != http.StatusOK {
		t.Errorf("expected http state 200 - got %d", statusCreate)
	}
}

//func TestClient_GetBlockImage(t *testing.T) {
//	client, err := ceph.New(getServer())
//
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	statusLogin, errLogin := client.Session.Login(username, password)
//	if errLogin != nil {
//		t.Error(errLogin)
//	}
//
//	if statusLogin != http.StatusCreated {
//		t.Fatalf("could not login - expected http state 201 - got %d", statusLogin)
//	}
//
//	status, rbd, errRbd := client.GetBlockImage("k14-pool01/vm-101071-disk-0")
//
//	if errRbd != nil {
//		t.Error(errRbd)
//	}
//
//	if status != http.StatusOK {
//		t.Errorf("expected http state 200 - got %d", status)
//	}
//
//	t.Logf("ID: %s", rbd.ID)
//
//}
