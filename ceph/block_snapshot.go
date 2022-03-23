package ceph

import (
    "fmt"
    "github.com/go-resty/resty/v2"
    "net/http"
    "net/url"
    "time"
)

// CreateBlockSnapShot creates a snapshot on an RBD image.
// see --> https://docs.ceph.com/en/latest/mgr/ceph_api/#post--api-block-image-image_spec-snap
func (c *Client) CreateBlockSnapShot(poolName string, nameSpace *string, imageName, snapShotName string, counter uint) (status int, err error) {

    if counter > c.MaxIterations {
        return 0, ErrMaxIterationsExceeded
    }

    if snapShotName == "" {
        return 0, ErrSnapshotNameIsEmpty
    }

    counter++

    var (
        resp      *resty.Response
        imageSpec string
        exception Exception
    )

    imageSpec, err = CreateImageSpec(poolName, nameSpace, imageName)

    if err != nil {
        return 0, err
    }

    client := *c.Session.Client

    jsonBody := struct {
        SnapshotName string `json:"snapshot_name"`
    }{SnapshotName: snapShotName}

    resp, err = client.
        SetRetryCount(10).
        SetRetryWaitTime(10 * time.Second).
        AddRetryCondition(c.retryConditionCheckForAccepted).
        R().
        SetHeaders(defaultHeaderJson).
        SetBody(jsonBody).
        Post(c.Session.Server.getURL(fmt.Sprintf("block/image/%s/snap", url.QueryEscape(imageSpec))))

    if !resp.IsSuccess() {
        if resp.StatusCode() == http.StatusBadRequest {
            err = client.JSONUnmarshal(resp.Body(), &exception)
            if err == nil {
                c.Logger.Debugf("err %s (%s)", exception.Code, exception.Detail)
                if exception.Code == RBDImageAlreadyExists {
                    return resp.StatusCode(), ErrCreateImageAlreadyExists
                }
            }
        }

        return resp.StatusCode(), fmt.Errorf("%v", resp.RawResponse)
    }

    status = resp.StatusCode()

    switch status {
    case http.StatusAccepted, http.StatusCreated, http.StatusBadRequest:
        lookForTask := Task{
            Name: "rbd/snap/create",
            MetaData: MetaData{
                ImageSpec: imageSpec, // only ImageSpec is needed on delete
            },
        }

        lookForTask, err = c.WaitForTaskIsDone(lookForTask)

        if err != nil {
            return 0, err
        }

        if !lookForTask.Success {
            // try delete again...
            c.Logger.Debugf("calling DeleteBlockImage with counter %d", counter)
            return c.CreateBlockSnapShot(poolName, nameSpace, imageName, snapShotName, counter)
        } else {
            status = http.StatusCreated
            err = nil
        }
    }

    return status, err
}
