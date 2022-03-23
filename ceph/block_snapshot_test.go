package ceph_test

import (
    "errors"
    "fmt"
    "github.com/chrisamti/ceph-rest-client/ceph"
    "net/http"
    "testing"
    "time"
)

func TestClient_CreateBlockSnapShot(t *testing.T) {
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

    var name = fmt.Sprintf("test-for-snapshot-%s", time.Now().Format(time.RFC3339))

    rbd := ceph.RBDCreate{
        Features:      nil,
        PoolName:      "test-pool-1",
        Namespace:     nil,
        Name:          name,
        Size:          1073741824 * 10,
        ObjSize:       0,
        StripeUnit:    nil,
        StripeCount:   nil,
        DataPool:      nil,
        Configuration: nil,
    }

    statusCreate, errCreate := client.CreateBlockImage(rbd, 0)

    if errCreate != nil {
        t.Error(errCreate)
    }

    if statusCreate != http.StatusCreated {
        t.Errorf("expected http state 201 - got %d", statusCreate)
    }

    // create snapshot without name --> should auto generate snapshot name
    status, errSnapshot := client.CreateBlockSnapShot(rbd.PoolName, rbd.Namespace, rbd.Name, "", 0)

    if !errors.Is(errSnapshot, ceph.ErrSnapshotNameIsEmpty) {
        t.Errorf("exptected error %v - got error %v", ceph.ErrSnapshotNameIsEmpty, errSnapshot)
    }

    // create snapshot without name --> should auto generate snapshot name
    snapshotName := fmt.Sprintf("%s-snap-1", name)
    status, errSnapshot = client.CreateBlockSnapShot(rbd.PoolName, rbd.Namespace, rbd.Name, snapshotName, 0)

    if errSnapshot != nil {
        t.Error(errSnapshot)
    }

    if status != http.StatusCreated {
        t.Errorf("expected http state 201 - got %d", status)
    }

}
