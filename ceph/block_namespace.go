package ceph

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/go-resty/resty/v2"
)

type NameSpace struct {
	NameSpace string `json:"namespace"`
	NumImages uint   `json:"num_images,omitempty"`
}

// GetBlockNameSpaceListInPool gets a list of CEPH RBD namespaces inside given pool.
// see --> https://docs.ceph.com/en/latest/mgr/ceph_api/#get--api-block-pool-pool_name-namespace
func (c *Client) GetBlockNameSpaceListInPool(poolName string) (status int, ns []NameSpace, err error) {

	if poolName == "" {
		return 0, ns, ErrPoolNameIsEmpty
	}

	var resp *resty.Response

	client := *c.Session.Client

	resp, err = client.
		SetRetryCount(10).
		SetRetryWaitTime(10 * time.Second).
		AddRetryCondition(c.retryConditionCheckForAccepted).
		R().
		SetHeaders(defaultHeaderJson).
		SetResult(ns).
		Get(c.Session.Server.getURL(fmt.Sprintf("block/pool/%s/namespace/", url.QueryEscape(poolName))))

	if err != nil {
		return 0, ns, err
	}

	if !resp.IsSuccess() {
		return resp.StatusCode(), ns, fmt.Errorf("%v", resp.RawResponse)
	}

	return resp.StatusCode(), ns, err
}

// CreateBlockNameSpaceInPool creates a new namespace for given pool.
// --> https://docs.ceph.com/en/latest/mgr/ceph_api/#post--api-block-pool-pool_name-namespace
func (c *Client) CreateBlockNameSpaceInPool(poolName, nameSpace string) (status int, err error) {

	if poolName == "" {
		return 0, ErrPoolNameIsEmpty
	}

	if nameSpace == "" {
		return 0, ErrNameSpaceNameIsEmpty
	}

	var (
		resp      *resty.Response
		exception Exception
		ns        = NameSpace{
			NameSpace: nameSpace,
		}
	)

	client := *c.Session.Client

	resp, err = client.
		SetRetryCount(10).
		SetRetryWaitTime(10 * time.Second).
		AddRetryCondition(c.retryConditionCheckForAccepted).
		R().
		SetHeaders(defaultHeaderJson).
		SetBody(ns).
		Post(c.Session.Server.getURL(fmt.Sprintf("block/pool/%s/namespace/", url.QueryEscape(poolName))))

	if err != nil {
		return 0, err
	}

	if !resp.IsSuccess() {
		if resp.StatusCode() == http.StatusBadRequest {
			err = client.JSONUnmarshal(resp.Body(), &exception)
			if err == nil {
				c.Logger.Debugf("err %s (%s)", exception.Code, exception.Detail)
				if exception.Code == NameSpaceAlreadyExists {
					return resp.StatusCode(), ErrNameSpaceAlreadyExists
				}
				// return more generic error
				return resp.StatusCode(), fmt.Errorf("could not create namespace: %v on pool %v: %v ", nameSpace, poolName, exception.Detail)
			}
		}

		return resp.StatusCode(), fmt.Errorf("could not create namespace: %v on pool %v: %v ", nameSpace, poolName, resp.Error())
	}

	status = resp.StatusCode()

	return status, err
}

// DeleteBlockNameSpaceInPool creates a new namespace for given pool.
// --> https://docs.ceph.com/en/pacific/mgr/ceph_api/index.html#delete--api-block-pool-pool_name-namespace-namespace
func (c *Client) DeleteBlockNameSpaceInPool(poolName, nameSpace string) (status int, err error) {

	if poolName == "" {
		return 0, ErrPoolNameIsEmpty
	}

	if nameSpace == "" {
		return 0, ErrNameSpaceNameIsEmpty
	}

	var (
		resp      *resty.Response
		exception Exception
	)

	client := *c.Session.Client

	resp, err = client.
		SetRetryCount(10).
		SetRetryWaitTime(10 * time.Second).
		AddRetryCondition(c.retryConditionCheckForAccepted).
		R().
		SetHeaders(defaultHeaderJson).
		Delete(c.Session.Server.getURL(fmt.Sprintf("block/pool/%s/namespace/%s", url.QueryEscape(poolName), url.QueryEscape(nameSpace))))

	if err != nil {
		return 0, err
	}

	if !resp.IsSuccess() {
		if resp.StatusCode() == http.StatusBadRequest {
			err = client.JSONUnmarshal(resp.Body(), &exception)
			if err == nil {
				c.Logger.Debugf("err %s (%s)", exception.Code, exception.Detail)
				if exception.Code == NameSpaceAlreadyExists {
					return resp.StatusCode(), ErrNameSpaceAlreadyExists
				}
				// return more generic error
				return resp.StatusCode(), fmt.Errorf("could not delete namespace: %v on pool %v: %v ", nameSpace, poolName, exception.Detail)
			}
		}

		return resp.StatusCode(), fmt.Errorf("could not delete namespace: %v on pool %v: %v ", nameSpace, poolName, resp.Error())
	}

	status = resp.StatusCode()

	return status, err
}
