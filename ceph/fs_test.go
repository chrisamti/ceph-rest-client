package ceph_test

import (
    "github.com/chrisamti/ceph-rest-client/ceph"
    "net/http"
    "testing"
    "time"
)

func TestClient_ListFS(t *testing.T) {
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

    status, fs, errFS := client.ListFS()

    if errFS != nil {
        t.Error(errFS)
    }

    if status != http.StatusOK {
        t.Errorf("expected http state 200 - got %d", status)
    }

    if len(fs) == 0 {
        t.Error("expected at least one ceph fs - got zero.")
    } else {
        t.Logf("cephfs id: %d", fs[0].ID)
    }

}

func TestClient_GetFS(t *testing.T) {
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

    status, fs, errFS := client.ListFS()

    if errFS != nil {
        t.Error(errFS)
    }

    if status != http.StatusOK {
        t.Errorf("expected http state 200 - got %d", status)
    }

    if len(fs) == 0 {
        t.Fatal("expected at least one ceph fs - got zero.")
    } else {
        t.Logf("cephfs id: %d", fs[0].ID)
    }

    status, _, errFS = client.GetFS(fs[0].ID)

    if errFS != nil {
        t.Error(errFS)
    }

    if status != http.StatusOK {
        t.Errorf("expected http state 200 - got %d", status)
    }

}

func TestClient_GetRootDirectory(t *testing.T) {
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

    status, fs, errFS := client.ListFS()

    if errFS != nil {
        t.Error(errFS)
    }

    if status != http.StatusOK {
        t.Errorf("expected http state 200 - got %d", status)
    }

    if len(fs) == 0 {
        t.Fatal("expected at least one ceph fs - got zero.")
    } else {
        t.Logf("cephfs id: %d", fs[0].ID)
    }

    var rd ceph.Directory
    status, rd, errFS = client.GetRootDirectory(fs[0].ID)

    if errFS != nil {
        t.Error(errFS)
    } else {
        t.Logf("%v", rd)
    }

    if status != http.StatusOK {
        t.Errorf("expected http state 200 - got %d", status)
    }

}

func TestClient_ListDir(t *testing.T) {
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

    status, fs, errFS := client.ListFS()

    if errFS != nil {
        t.Error(errFS)
    }

    if status != http.StatusOK {
        t.Errorf("expected http state 200 - got %d", status)
    }

    if len(fs) == 0 {
        t.Fatal("expected at least one ceph fs - got zero.")
    } else {
        t.Logf("cephfs id: %d", fs[0].ID)
    }

    var rd ceph.Directory
    status, rd, errFS = client.GetRootDirectory(fs[0].ID)

    if errFS != nil {
        t.Error(errFS)
    } else {
        t.Logf("%v", rd)
    }

    if status != http.StatusOK {
        t.Errorf("expected http state 200 - got %d", status)
    }

    var dir interface{}
    status, dir, errFS = client.ListDir(fs[0].ID, rd.Path, 20)
    if errFS != nil {
        t.Error(errFS)
    } else {
        t.Logf("%v", dir)
    }

    if status != http.StatusOK {
        t.Errorf("expected http state 200 - got %d", status)
    }
}

func TestClient_CreateDir(t *testing.T) {
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

    status, fs, errFS := client.ListFS()

    if errFS != nil {
        t.Error(errFS)
    }

    if status != http.StatusOK {
        t.Errorf("expected http state 200 - got %d", status)
    }

    if len(fs) == 0 {
        t.Fatal("expected at least one ceph fs - got zero.")
    } else {
        t.Logf("cephfs id: %d", fs[0].ID)
    }

    var rd ceph.Directory
    status, rd, errFS = client.GetRootDirectory(fs[0].ID)

    if errFS != nil {
        t.Error(errFS)
    } else {
        t.Logf("%v", rd)
    }

    if status != http.StatusOK {
        t.Errorf("expected http state 200 - got %d", status)
    }

    status, err = client.CreateDir(fs[0].ID, "/fs_test")
    if err != nil {
        t.Fatal(err)
    }

    if status != http.StatusCreated &&
        status != http.StatusAccepted &&
        status != http.StatusOK {
        t.Errorf("expected http state 201 or 202 - got %d", status)
    }

    status, err = client.CreateDir(fs[0].ID, "/fs_test/sub1/sub2")
    if err != nil {
        t.Fatal(err)
    }

    if status != http.StatusCreated &&
        status != http.StatusAccepted &&
        status != http.StatusOK {
        t.Errorf("expected http state 201 or 202 - got %d", status)
    }

    var dir []ceph.Directory
    status, dir, errFS = client.ListDir(fs[0].ID, rd.Path, 20)
    if errFS != nil {
        t.Error(errFS)
    } else {
        for _, d := range dir {
            t.Logf("dir: %s", d.Path)
        }
        t.Logf("%v", dir)
    }

    if status != http.StatusOK {
        t.Errorf("expected http state 200 - got %d", status)
    }
}

func TestClient_DeleteDir(t *testing.T) {
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

    status, fs, errFS := client.ListFS()

    if errFS != nil {
        t.Error(errFS)
    }

    if status != http.StatusOK {
        t.Errorf("expected http state 200 - got %d", status)
    }

    if len(fs) == 0 {
        t.Fatal("expected at least one ceph fs - got zero.")
    } else {
        t.Logf("cephfs id: %d", fs[0].ID)
    }

    var rd ceph.Directory
    status, rd, errFS = client.GetRootDirectory(fs[0].ID)

    if errFS != nil {
        t.Error(errFS)
    } else {
        t.Logf("%v", rd)
    }

    if status != http.StatusOK {
        t.Errorf("expected http state 200 - got %d", status)
    }

    status, err = client.CreateDir(fs[0].ID, "/fs_test")
    if err != nil {
        t.Fatal(err)
    }

    if status != http.StatusCreated &&
        status != http.StatusAccepted &&
        status != http.StatusOK {
        t.Errorf("expected http state 201 or 202 - got %d", status)
    }

    status, err = client.CreateDir(fs[0].ID, "/fs_test/sub1/sub2")
    if err != nil {
        t.Fatal(err)
    }

    if status != http.StatusCreated &&
        status != http.StatusAccepted &&
        status != http.StatusOK {
        t.Errorf("expected http state 200, 201 or 202 - got %d", status)
    }

    status, err = client.DeleteDir(fs[0].ID, "/fs_test/sub1/sub2")

    if err != nil {
        t.Fatal(err)
    }

    if status != http.StatusOK &&
        status != http.StatusNoContent &&
        status != http.StatusAccepted {
        t.Errorf("expected http state 202 or 204 - got %d", status)
    }

    status, err = client.DeleteDir(fs[0].ID, "/fs_test/sub1")

    if err != nil {
        t.Fatal(err)
    }

    if status != http.StatusOK &&
        status != http.StatusNoContent &&
        status != http.StatusAccepted {
        t.Errorf("expected http state 200, 202 or 204 - got %d", status)
    }

    status, err = client.DeleteDir(fs[0].ID, "/fs_test")

    if err != nil {
        t.Fatal(err)
    }

    if status != http.StatusOK &&
        status != http.StatusNoContent &&
        status != http.StatusAccepted {
        t.Errorf("expected http state 200, 202 or 204 - got %d", status)
    }

    var dir []ceph.Directory
    status, dir, errFS = client.ListDir(fs[0].ID, rd.Path, 20)
    if errFS != nil {
        t.Error(errFS)
    } else {
        for _, d := range dir {
            t.Logf("dir: %s", d.Path)
        }
        t.Logf("%v", dir)
    }

    if status != http.StatusOK {
        t.Errorf("expected http state 200 - got %d", status)
    }
}

func TestClient_CreateDeleteSnapShot(t *testing.T) {
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

    status, fs, errFS := client.ListFS()

    if errFS != nil {
        t.Error(errFS)
    }

    if status != http.StatusOK {
        t.Errorf("expected http state 200 - got %d", status)
    }

    if len(fs) == 0 {
        t.Fatal("expected at least one ceph fs - got zero.")
    } else {
        t.Logf("cephfs id: %d", fs[0].ID)
    }

    var rd ceph.Directory
    status, rd, errFS = client.GetRootDirectory(fs[0].ID)

    if errFS != nil {
        t.Error(errFS)
    } else {
        t.Logf("%v", rd)
    }

    if status != http.StatusOK {
        t.Errorf("expected http state 200 - got %d", status)
    }

    status, err = client.CreateDir(fs[0].ID, "/fs_test")
    if err != nil {
        t.Fatal(err)
    }

    if status != http.StatusCreated &&
        status != http.StatusAccepted &&
        status != http.StatusOK {
        t.Errorf("expected http state 201 or 202 - got %d", status)
    }

    status, err = client.CreateDir(fs[0].ID, "/fs_test/sub1/sub2")
    if err != nil {
        t.Fatal(err)
    }

    if status != http.StatusCreated &&
        status != http.StatusAccepted &&
        status != http.StatusOK {
        t.Errorf("expected http state 201 or 202 - got %d", status)
    }

    var dir []ceph.Directory
    status, dir, errFS = client.ListDir(fs[0].ID, rd.Path, 20)
    if errFS != nil {
        t.Error(errFS)
    } else {
        for _, d := range dir {
            t.Logf("dir: %s", d.Path)
        }
        t.Logf("%v", dir)
    }

    if status != http.StatusOK {
        t.Errorf("expected http state 200 - got %d", status)
    }

    // create a snapshot
    snap := ceph.SnapShot{
        Name: "test_snap_1",
        Path: "/fs_test/sub1/sub2",
    }

    t.Logf("creating snapshot %v", snap)
    status, err = client.CreateSnapShot(fs[0].ID, snap)

    if err != nil {
        t.Error(err)
    }

    t.Logf("status: %d", status)
    if status != http.StatusCreated &&
        status != http.StatusAccepted &&
        status != http.StatusOK {
        t.Errorf("expected http state 200, 201 or 202 - got %d", status)
    }

    time.Sleep(time.Second)

    // delete snapshot
    t.Logf("deleting snapshot %v", snap)
    status, err = client.DeleteSnapShot(fs[0].ID, snap)
    if err != nil {
        t.Error(err)
    }

    t.Logf("status: %d", status)
    if status != http.StatusCreated &&
        status != http.StatusAccepted &&
        status != http.StatusOK {
        t.Errorf("expected http state 200, 201 or 202 - got %d", status)
    }

}
