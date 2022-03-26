package ceph

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/go-resty/resty/v2"
)

var (
	// ErrImageSpecIsEmpty is returned if param imageSpec is empty.
	ErrImageSpecIsEmpty = errors.New("param imageSpec can not be empty")

	// ErrPoolNameIsEmpty is returned if param poolName is empty.
	ErrPoolNameIsEmpty = errors.New("param poolName can not be empty")

	// ErrImageNameIsEmpty is returned if param imageName is empty.
	ErrImageNameIsEmpty = errors.New("param imageName can not be empty")

	// ErrSnapshotNameIsEmpty is returned if param imageName is empty.
	ErrSnapshotNameIsEmpty = errors.New("param snapShotName can not be empty")

	// ErrCreateImageAlreadyExists is return if image to be created already exists.
	ErrCreateImageAlreadyExists = errors.New("RBD image already exists (error creating image)")

	// ErrNameSpaceNameIsEmpty is returned if param nameSpace is empty.
	ErrNameSpaceNameIsEmpty = errors.New("param nameSpace can not be empty")

	// ErrEditImageAlreadyExists is return if image to be renamed already exists.
	ErrEditImageAlreadyExists = errors.New("RBD image already exists (error renaming image)")

	// ErrNameSpaceAlreadyExists is returned if a namespace already exist for a rbd pool.
	ErrNameSpaceAlreadyExists = errors.New("namespace already exists")
)

const (
	RBDImageAlreadyExists  = "17"
	NameSpaceAlreadyExists = "namespace_already_exists"
)

// RBDConfiguration implements struct for some rbd configuration values.
type RBDConfiguration struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Source int    `json:"source"`
}

// RBD implements struct returned from GET /api/block/image/{image_spec}
// --> https://docs.ceph.com/en/latest/mgr/ceph_api/#get--api-block-image-image_spec.
type RBD struct {
	Size            uint64             `json:"size"`
	ObjSize         uint64             `json:"obj_size"`
	NumObjs         uint               `json:"num_objs"`
	Order           uint               `json:"order"`
	BlockNamePrefix string             `json:"block_name_prefix"`
	Name            string             `json:"name"`
	UniqueID        string             `json:"unique_id"`
	ID              string             `json:"id"`
	ImageFormat     int                `json:"image_format"`
	PoolName        string             `json:"pool_name"`
	Namespace       *string            `json:"namespace"`
	Features        int                `json:"features"`
	FeaturesName    []string           `json:"features_name"`
	Timestamp       time.Time          `json:"timestamp"`
	StripeCount     *uint              `json:"stripe_count"`
	StripeUnit      *uint64            `json:"stripe_unit"`
	DataPool        interface{}        `json:"data_pool"`
	Parent          interface{}        `json:"parent"`
	Snapshots       []interface{}      `json:"snapshots"`
	TotalDiskUsage  uint64             `json:"total_disk_usage"`
	DiskUsage       uint64             `json:"disk_usage"`
	Configuration   []RBDConfiguration `json:"configuration"`
}

type RBDQosConfig struct {
	RbdQosBpsLimit       uint `json:"rbd_qos_bps_limit"`
	RbdQosIopsLimit      uint `json:"rbd_qos_iops_limit"`
	RbdQosReadBpsLimit   uint `json:"rbd_qos_read_bps_limit"`
	RbdQosReadIopsLimit  uint `json:"rbd_qos_read_iops_limit"`
	RbdQosWriteBpsLimit  uint `json:"rbd_qos_write_bps_limit"`
	RbdQosWriteIopsLimit uint `json:"rbd_qos_write_iops_limit"`
	RbdQosBpsBurst       uint `json:"rbd_qos_bps_burst"`
	RbdQosIopsBurst      uint `json:"rbd_qos_iops_burst"`
	RbdQosReadBpsBurst   uint `json:"rbd_qos_read_bps_burst"`
	RbdQosReadIopsBurst  uint `json:"rbd_qos_read_iops_burst"`
	RbdQosWriteBpsBurst  uint `json:"rbd_qos_write_bps_burst"`
	RbdQosWriteIopsBurst uint `json:"rbd_qos_write_iops_burst"`
}

// RBDList implements struct received from GET /api/block/image.
// --> https://docs.ceph.com/en/latest/mgr/ceph_api/#get--api-block-image.
type RBDList []struct {
	Status   int    `json:"status"`
	Value    []RBD  `json:"value"`
	PoolName string `json:"pool_name"`
}

// RBDCreate implements struct send to ceph for rbd image creation on POST /api/block/image.
// --> https://docs.ceph.com/en/latest/mgr/ceph_api/#post--api-block-image
type RBDCreate struct {
	Features      []string      `json:"features"`
	PoolName      string        `json:"pool_name"`
	Namespace     *string       `json:"namespace"`
	Name          string        `json:"name"`
	Size          uint64        `json:"size"`
	ObjSize       uint64        `json:"obj_size"`
	StripeUnit    *uint64       `json:"stripe_unit"`
	StripeCount   *uint         `json:"stripe_count"`
	DataPool      *string       `json:"data_pool"`
	Configuration *RBDQosConfig `json:"configuration"`
}

// CreateImageSpec creates a valid ceph rbd image name spec.
func CreateImageSpec(poolName string, nameSpace *string, imgName string) (string, error) {
	if poolName == "" {
		return "", ErrPoolNameIsEmpty
	}

	if imgName == "" {
		return "", ErrImageNameIsEmpty
	}

	return PathJoin(poolName, nameSpace, imgName), nil
}

// RBDCopy implements struct needed on rbd image copy operations.
// See https://docs.ceph.com/en/latest/mgr/ceph_api/#post--api-block-image-image_spec-copy.
type RBDCopy struct {
	Configuration *RBDQosConfig `json:"configuration"`
	DataPool      *string       `json:"data_pool"`
	DestImageName string        `json:"dest_image_name"`
	DestNameSpace *string       `json:"dest_namespace"`
	DestPoolName  string        `json:"dest_pool_name"`
	Features      []string      `json:"features"`
	ObjSize       uint64        `json:"obj_size"`
	SnapShotName  string        `json:"snapshot_name,omitempty"`
	StripeCount   *uint         `json:"stripe_count"`
	StripeUnit    *uint64       `json:"stripe_unit"`
}

// RBDError implements error struct returned.
type RBDError struct {
	Detail    string `json:"detail"`
	Code      string `json:"code"`
	Component string `json:"component"`
	Status    int    `json:"status"`
	Task      struct {
		Name     string `json:"name"`
		Metadata struct {
			PoolName  string      `json:"pool_name"`
			Namespace interface{} `json:"namespace"`
			ImageName string      `json:"image_name"`
		} `json:"metadata"`
	} `json:"task"`
}

// RBDUpdate implements struct send to ceph for rbd image updates on PUT /api/block/image/{image_spec}.
// --> https://docs.ceph.com/en/latest/mgr/ceph_api/#put--api-block-image-image_spec
type RBDUpdate struct {
	Features      []string `json:"features"`
	Name          string   `json:"name"`
	Size          int64    `json:"size"`
	Configuration struct{} `json:"configuration"`
}

// ListBlockImage gets a list of RBD block images (https://docs.ceph.com/en/latest/mgr/ceph_api/#get--api-block-image)
func (c *Client) ListBlockImage(poolName string) (status int, rbdList RBDList, err error) {
	var resp *resty.Response

	client := *c.Session.Client

	if poolName != "" {
		resp, err = client.R().
			SetHeaders(defaultHeaderJson).
			SetQueryParam("pool_name", poolName).
			SetResult(&rbdList).
			Get(c.Session.Server.getURL("block/image"))
	} else {
		resp, err = client.R().
			SetHeaders(defaultHeaderJson).
			SetResult(&rbdList).
			Get(c.Session.Server.getURL("block/image"))

	}

	if !resp.IsSuccess() {
		return resp.StatusCode(), nil, fmt.Errorf("%v", resp.RawResponse)
	}

	return resp.StatusCode(), rbdList, err

}

// GetBlockImage gets an RBD block image (https://docs.ceph.com/en/latest/mgr/ceph_api/#get--api-block-image-image_spec)
func (c *Client) GetBlockImage(imageSpec string) (status int, rbd RBD, err error) {
	var resp *resty.Response

	if imageSpec == "" {
		return 0, rbd, ErrImageSpecIsEmpty
	}

	client := *c.Session.Client

	resp, err = client.
		R().
		SetHeaders(defaultHeaderJson).
		SetResult(&rbd).
		Get(c.Session.Server.getURL(fmt.Sprintf("block/image/%s", url.QueryEscape(imageSpec))))

	if !resp.IsSuccess() {
		return resp.StatusCode(), rbd, fmt.Errorf("could not get image %v: %v", imageSpec, resp.Error())
	}

	return resp.StatusCode(), rbd, err
}

// CreateBlockImage creates an RBD image (https://docs.ceph.com/en/latest/mgr/ceph_api/#post--api-block-image)
func (c *Client) CreateBlockImage(rbdCreate RBDCreate, counter uint) (status int, err error) {

	if counter > c.MaxIterations {
		return 0, ErrMaxIterationsExceeded
	}

	counter++

	var (
		resp      *resty.Response
		exception Exception
	)

	// create copy of client
	client := *c.Session.Client

	resp, err = client.
		SetRetryCount(10).
		SetRetryWaitTime(10 * time.Second).
		AddRetryCondition(c.retryConditionCheckForAccepted).
		R().
		SetHeaders(defaultHeaderJson).
		SetBody(rbdCreate).
		Post(c.Session.Server.getURL("block/image"))

	if err != nil {
		return 0, err
	}

	if !resp.IsSuccess() {
		if resp.StatusCode() == http.StatusBadRequest {
			err = client.JSONUnmarshal(resp.Body(), &exception)
			if err == nil {
				c.Logger.Debugf("err %s (%s)", exception.Code, exception.Detail)
				if exception.Code == RBDImageAlreadyExists {
					return resp.StatusCode(), ErrEditImageAlreadyExists
				}
			}
		}

		return resp.StatusCode(), fmt.Errorf("could not create image: %v on pool %v: %v ", rbdCreate.Name, rbdCreate.PoolName, resp.Error())
	}

	status = resp.StatusCode()

	// check task state
	switch status {
	case http.StatusCreated, http.StatusAccepted:
		lookForTask := Task{
			Name: "rbd/create",
			MetaData: MetaData{
				PoolName:  rbdCreate.PoolName,
				Namespace: rbdCreate.Namespace,
				ImageName: rbdCreate.Name,
				ImageSpec: "", // ImageSpec is always empty for rbd/create
			},
		}

		lookForTask, err = c.WaitForTaskIsDone(lookForTask)
		if err != nil {
			return 0, err
		}

		if !lookForTask.Success {
			// try to create again
			c.Logger.Debugf("call CreateBlockImage again with counter %d", counter)
			return c.CreateBlockImage(rbdCreate, counter)
		} else {
			status = http.StatusCreated
		}
	}

	return status, nil
}

// CopyBlockImage create a copy of existing rbd.
// See https://docs.ceph.com/en/latest/mgr/ceph_api/#post--api-block-image-image_spec-copy.
func (c *Client) CopyBlockImage(poolName string, nameSpace *string, imageName string, dst RBDCopy, counter uint) (status int, err error) {
	if counter > c.MaxIterations {
		return 0, ErrMaxIterationsExceeded
	}

	counter++

	var resp *resty.Response
	var imageSpec string

	imageSpec, err = CreateImageSpec(poolName, nameSpace, imageName)

	if err != nil {
		return 0, err
	}

	client := *c.Session.Client

	resp, err = client.
		SetRetryCount(10).
		SetRetryWaitTime(10 * time.Second).
		AddRetryCondition(c.retryConditionCheckForAccepted).
		R().
		SetHeaders(defaultHeaderJson).
		SetBody(dst).
		Post(c.Session.Server.getURL(fmt.Sprintf("block/image/%s/copy", url.QueryEscape(imageSpec))))

	if !resp.IsSuccess() {
		return resp.StatusCode(), fmt.Errorf("%v", resp.RawResponse)
	}

	status = resp.StatusCode()

	switch status {
	case http.StatusCreated, http.StatusAccepted:
		lookForTask := Task{
			Name: "rbd/copy",
			MetaData: MetaData{
				ImageSpec: imageSpec, // only ImageSpec is needed on copy
			},
		}

		lookForTask, err = c.WaitForTaskIsDone(lookForTask)

		if err != nil {
			return 0, err
		}

		if !lookForTask.Success {
			// try delete again...
			c.Logger.Debugf("calling CopyBlockImage with counter %d", counter)
			return c.CopyBlockImage(poolName, nameSpace, imageName, dst, counter)
		} else {
			status = http.StatusNoContent
			err = nil
		}
	}

	return status, err
}

// DeleteBlockImage deletes an RBD image defined with imageSpec
// (https://docs.ceph.com/en/latest/mgr/ceph_api/#delete--api-block-image-image_spec)
func (c *Client) DeleteBlockImage(poolName string, nameSpace *string, imageName string, counter uint) (status int, err error) {

	if counter > c.MaxIterations {
		return 0, ErrMaxIterationsExceeded
	}

	counter++

	var (
		resp      *resty.Response
		imageSpec string
	)

	imageSpec, err = CreateImageSpec(poolName, nameSpace, imageName)
	if err != nil {
		return 0, err
	}

	client := *c.Session.Client

	resp, err = client.
		SetRetryCount(10).
		SetRetryWaitTime(10 * time.Second).
		AddRetryCondition(c.retryConditionCheckForAccepted).
		R().
		SetHeaders(defaultHeaderJson).
		Delete(c.Session.Server.getURL(fmt.Sprintf("block/image/%s", url.QueryEscape(imageSpec))))

	if !resp.IsSuccess() {
		return resp.StatusCode(), fmt.Errorf("%v", resp.RawResponse)
	}

	status = resp.StatusCode()

	switch status {
	case http.StatusAccepted, http.StatusNoContent, http.StatusBadRequest:
		lookForTask := Task{
			Name: "rbd/delete",
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
			return c.DeleteBlockImage(poolName, nameSpace, imageName, counter)
		} else {
			status = http.StatusNoContent
			err = nil
		}
	}

	return status, err
}

// MoveBlockImageToTrash moves ceph rbd image to ceph trash.
// --> https://docs.ceph.com/en/latest/mgr/ceph_api/#post--api-block-image-image_spec-move_trash
// Attention: the documentation claims the response code on moving an image to the trash returns 201. But at least on
// ceph version 16.2.7 (f9aa029788115b5df5eeee328f584156565ee5b7) pacific (stable) 200 is returned.
func (c *Client) MoveBlockImageToTrash(poolName string, nameSpace *string, imageName string, delay time.Duration, counter uint) (status int, err error) {
	if counter > c.MaxIterations {
		return 0, ErrMaxIterationsExceeded
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

	delayPost := struct {
		Delay float64 `json:"delay"`
	}{Delay: delay.Seconds()}

	resp, err = client.
		SetRetryCount(10).
		SetRetryWaitTime(10 * time.Second).
		AddRetryCondition(c.retryConditionCheckForAccepted).
		R().
		SetHeaders(defaultHeaderJson).
		SetBody(delayPost).
		Post(c.Session.Server.getURL(fmt.Sprintf("block/image/%s/move_trash", url.QueryEscape(imageSpec))))

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
	case http.StatusOK, http.StatusCreated, http.StatusBadRequest:
		lookForTask := Task{
			Name: "rbd/trash/move",
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
			return c.MoveBlockImageToTrash(poolName, nameSpace, imageName, delay, counter)
		} else {
			status = http.StatusOK
			err = nil
		}
	}

	return status, err
}

// UpdateBlockImage updates ceph rbd image (name, size etc al).
// --> https://docs.ceph.com/en/latest/mgr/ceph_api/#put--api-block-image-image_spec
func (c *Client) UpdateBlockImage(poolName string, nameSpace *string, imageName string, rbdUpdate RBDUpdate, counter uint) (status int, err error) {
	if counter > c.MaxIterations {
		return 0, ErrMaxIterationsExceeded
	}

	counter++

	var (
		resp      *resty.Response
		imageSpec string
		exception Exception
	)

	// check rbdUpdate
	// TODO: check all other rbdUpdate struct values for validity.
	if rbdUpdate.Name == "" {
		return 0, ErrImageNameIsEmpty
	}

	imageSpec, err = CreateImageSpec(poolName, nameSpace, imageName)

	if err != nil {
		return 0, err
	}

	client := *c.Session.Client

	resp, err = client.
		SetRetryCount(10).
		SetRetryWaitTime(10 * time.Second).
		AddRetryCondition(c.retryConditionCheckForAccepted).
		R().
		SetHeaders(defaultHeaderJson).
		SetBody(rbdUpdate).
		Put(c.Session.Server.getURL(fmt.Sprintf("block/image/%s", url.QueryEscape(imageSpec))))

	if !resp.IsSuccess() {
		if resp.StatusCode() == http.StatusBadRequest {
			err = client.JSONUnmarshal(resp.Body(), &exception)
			if err == nil {
				c.Logger.Debugf("err %s (%s)", exception.Code, exception.Detail)
				if exception.Code == "17" {
					return resp.StatusCode(), ErrCreateImageAlreadyExists
				}
			}
		}

		return resp.StatusCode(), fmt.Errorf("%v", resp.RawResponse)
	}

	status = resp.StatusCode()

	switch status {
	case http.StatusAccepted, http.StatusOK, http.StatusBadRequest:
		lookForTask := Task{
			Name: "rbd/edit",
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
			return c.UpdateBlockImage(poolName, nameSpace, imageName, rbdUpdate, counter)
		} else {
			status = http.StatusOK
			err = nil
		}
	}

	return status, err

}

func (c *Client) retryConditionCheckForAccepted(r *resty.Response, _ error) bool {
	switch r.StatusCode() {
	case http.StatusOK,
		http.StatusNoContent,
		http.StatusNotFound,
		http.StatusCreated,
		http.StatusAccepted,
		http.StatusBadRequest:
		// no retry needed
		c.Logger.Debugf("http status: %d --> no retry for %s", r.StatusCode(), r.Request.URL)
		return false
	}

	c.Logger.Debugf("http status: %d --> retry for %s needed", r.StatusCode(), r.Request.URL)
	return true // retry on all other status codes.
}
