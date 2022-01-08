package ceph

import (
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/go-resty/resty/v2"
)

var ErrImageSpecIsEmpty = errors.New("param imageSpec can not be empty")

// RBD implements struct returned from https://docs.ceph.com/en/latest/mgr/ceph_api/#get--api-block-image-image_spec.
type RBD struct {
	Size            int64         `json:"size"`
	ObjSize         int           `json:"obj_size"`
	NumObjs         int           `json:"num_objs"`
	Order           int           `json:"order"`
	BlockNamePrefix string        `json:"block_name_prefix"`
	Name            string        `json:"name"`
	UniqueID        string        `json:"unique_id"`
	ID              string        `json:"id"`
	ImageFormat     int           `json:"image_format"`
	PoolName        string        `json:"pool_name"`
	Namespace       interface{}   `json:"namespace"`
	Features        int           `json:"features"`
	FeaturesName    []string      `json:"features_name"`
	Timestamp       time.Time     `json:"timestamp"`
	StripeCount     int           `json:"stripe_count"`
	StripeUnit      int           `json:"stripe_unit"`
	DataPool        interface{}   `json:"data_pool"`
	Parent          interface{}   `json:"parent"`
	Snapshots       []interface{} `json:"snapshots"`
	TotalDiskUsage  int           `json:"total_disk_usage"`
	DiskUsage       int           `json:"disk_usage"`
	Configuration   []struct {
		Name   string `json:"name"`
		Value  string `json:"value"`
		Source int    `json:"source"`
	} `json:"configuration"`
}

// RBDList implements struct returned from https://docs.ceph.com/en/latest/mgr/ceph_api/#get--api-block-image.
type RBDList []struct {
	Status   int    `json:"status"`
	Value    []RBD  `json:"value"`
	PoolName string `json:"pool_name"`
}

type RBDCreate struct {
	Features      []string    `json:"features"`
	PoolName      string      `json:"pool_name"`
	Namespace     interface{} `json:"namespace"`
	Name          string      `json:"name"`
	Size          int         `json:"size"`
	ObjSize       int         `json:"obj_size"`
	StripeUnit    interface{} `json:"stripe_unit"`
	StripeCount   interface{} `json:"stripe_count"`
	DataPool      interface{} `json:"data_pool"`
	Configuration struct{}    `json:"configuration"`
}

// ListBlockImage gets a list of RBD block images (https://docs.ceph.com/en/latest/mgr/ceph_api/#get--api-block-image)
func (c *Client) ListBlockImage(poolName string) (status int, rbdList RBDList, err error) {
	var resp *resty.Response

	if poolName != "" {
		resp, err = c.Session.Client.R().
			SetHeaders(defaultHeaderJson).
			SetQueryParam("pool_name", poolName).
			SetResult(&rbdList).
			Get(c.Session.Server.getURL("block/image"))
	} else {
		resp, err = c.Session.Client.R().
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

	resp, err = c.Session.Client.R().
		SetHeaders(defaultHeaderJson).
		SetResult(&rbd).
		Get(c.Session.Server.getURL(fmt.Sprintf("block/image/%s", url.QueryEscape(imageSpec))))

	if !resp.IsSuccess() {
		return resp.StatusCode(), rbd, fmt.Errorf("could not get image %v: %v", imageSpec, resp.Error())
	}

	return resp.StatusCode(), rbd, err
}

// CreateBlockImage creates an RBD image (https://docs.ceph.com/en/latest/mgr/ceph_api/#post--api-block-image)
func (c *Client) CreateBlockImage(rbdCreate RBDCreate) (status int, err error) {
	var resp *resty.Response

	resp, err = c.Session.Client.R().
		SetHeaders(defaultHeaderJson).
		SetBody(rbdCreate).
		Post(c.Session.Server.getURL("block/image"))

	if !resp.IsSuccess() {
		return resp.StatusCode(), fmt.Errorf("could not create image: %v on pool %v: %v ", rbdCreate.Name, rbdCreate.PoolName, resp.Error())
	}

	return resp.StatusCode(), nil
}

// DeleteBlockImage deletes an RBD image defined with imageSpec
// (https://docs.ceph.com/en/latest/mgr/ceph_api/#delete--api-block-image-image_spec)
func (c *Client) DeleteBlockImage(imageSpec string) (status int, err error) {
	var resp *resty.Response

	if imageSpec == "" {
		return 0, ErrImageSpecIsEmpty
	}

	resp, err = c.Session.Client.R().
		SetHeaders(defaultHeaderJson).
		Delete(c.Session.Server.getURL(fmt.Sprintf("block/image/%s", url.QueryEscape(imageSpec))))

	if !resp.IsSuccess() {
		return resp.StatusCode(), fmt.Errorf("%v", resp.RawResponse)
	}

	return resp.StatusCode(), err
}
