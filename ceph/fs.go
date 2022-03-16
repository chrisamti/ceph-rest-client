package ceph

import (
    "fmt"
    "github.com/go-resty/resty/v2"
)

// MdsMap implements ceph fs Mds Mapping struct.
type MdsMap struct {
    Epoch                     int                    `json:"epoch"`
    Flags                     int                    `json:"flags"`
    EverAllowedFeatures       int                    `json:"ever_allowed_features"`
    ExplicitlyAllowedFeatures int                    `json:"explicitly_allowed_features"`
    TableServer               int                    `json:"tableserver"`
    Root                      int                    `json:"root"`
    SessionTimeout            int                    `json:"session_timeout"`
    SessionAutoClose          int                    `json:"session_autoclose"`
    RequiredClientFeatures    map[string]interface{} `json:"required_client_features"`
    MaxFileSize               int64                  `json:"max_file_size"`
    LastFailure               int                    `json:"last_failure"`
    LastFailureOsdEpoch       int                    `json:"last_failure_osd_epoch"`
    Compat                    struct {
        Compat   map[string]interface{} `json:"compat"`
        RoCompat map[string]interface{} `json:"ro_compat"`
        InCompat map[string]interface{} `json:"incompat"`
    } `json:"compat"`
    MaxMds             int                    `json:"max_mds"`
    In                 []int                  `json:"in"`
    Up                 map[string]interface{} `json:"up"`
    Failed             []interface{}          `json:"failed"`
    Damaged            []interface{}          `json:"damaged"`
    Stopped            []interface{}          `json:"stopped"`
    Info               map[string]interface{} `json:"info"`
    DataPools          []int                  `json:"data_pools"`
    MetadataPool       int                    `json:"metadata_pool"`
    Enabled            bool                   `json:"enabled"`
    FsName             string                 `json:"fs_name"`
    Balancer           string                 `json:"balancer"`
    StandbyCountWanted int                    `json:"standby_count_wanted"`
    Created            string                 `json:"created"`
    Modified           string                 `json:"modified"`
}

type FS struct {
    MdsMap MdsMap `json:"mdsmap"`
    ID     int    `json:"id"`
}

// Quota implements a ceph fs quota.
type Quota struct {
    MaxBytes int    `json:"max_bytes"`
    MaxFiles int    `json:"max_files"`
    Path     string `json:"path"`
}

// Directory implements a ceph fs directory.
type Directory struct {
    Name      string        `json:"name"`
    Path      string        `json:"path"`
    Parent    string        `json:"parent"`
    Snapshots []interface{} `json:"snapshots"`
    Quotas    Quota         `json:"quotas"`
}

// ListFS gets all possible ceph fs available (https://docs.ceph.com/en/latest/mgr/ceph_api/#get--api-cephfs).
func (c *Client) ListFS() (status int, list []FS, err error) {
    var resp *resty.Response

    client := *c.Session.Client

    resp, err = client.R().
        SetHeaders(defaultHeaders).
        SetResult(&list).
        Get(c.Session.Server.getURL("cephfs"))

    if !resp.IsSuccess() {
        return resp.StatusCode(), nil, fmt.Errorf("%v", resp.RawResponse)
    }

    return resp.StatusCode(), list, err
}

// GetFS gets a specific ceph fs by id (https://docs.ceph.com/en/latest/mgr/ceph_api/#get--api-cephfs).
func (c *Client) GetFS(id int) (status int, fs interface{}, err error) {

    var resp *resty.Response

    client := *c.Session.Client

    resp, err = client.R().
        SetHeaders(defaultHeaders).
        Get(c.Session.Server.getURL(fmt.Sprintf("cephfs/%d", id)))

    if !resp.IsSuccess() {
        return resp.StatusCode(), nil, fmt.Errorf("%v", resp.RawResponse)
    }

    return resp.StatusCode(), nil, err
}

// GetRootDirectory gets the ceph fs root directory (https://docs.ceph.com/en/latest/mgr/ceph_api/#get--api-cephfs-fs_id-get_root_directory).
func (c *Client) GetRootDirectory(id int) (status int, rootDir Directory, err error) {
    var resp *resty.Response

    client := *c.Session.Client

    resp, err = client.R().
        SetHeaders(defaultHeaders).
        SetResult(&rootDir).
        Get(c.Session.Server.getURL(fmt.Sprintf("cephfs/%d/get_root_directory", id)))

    if !resp.IsSuccess() {
        return resp.StatusCode(), rootDir, fmt.Errorf("%v", resp.RawResponse)
    }

    return resp.StatusCode(), rootDir, err
}

// ListDir gets a list if ceph fs directories (https://docs.ceph.com/en/latest/mgr/ceph_api/#get--api-cephfs-fs_id-ls_dir).
func (c *Client) ListDir(id int, path string, depth uint) (status int, dir []Directory, err error) {
    var resp *resty.Response

    client := *c.Session.Client

    resp, err = client.R().
        SetHeaders(defaultHeaders).
        SetResult(&dir).
        SetQueryParam("path", path).
        SetQueryParam("depth", fmt.Sprintf("%d", depth)).
        Get(c.Session.Server.getURL(fmt.Sprintf("cephfs/%d/ls_dir", id)))

    if !resp.IsSuccess() {
        return resp.StatusCode(), dir, fmt.Errorf("%v", resp.RawResponse)
    }

    return resp.StatusCode(), dir, err
}

// CreateDir creates a ceph fs directory (https://docs.ceph.com/en/latest/mgr/ceph_api/#post--api-cephfs-fs_id-tree)
func (c *Client) CreateDir(id int, path string) (status int, err error) {
    var resp *resty.Response

    client := *c.Session.Client

    body := struct {
        Path string `json:"path"`
    }{Path: path}

    resp, err = client.R().
        SetHeaders(defaultHeaders).
        SetBody(body).
        Post(c.Session.Server.getURL(fmt.Sprintf("cephfs/%d/tree", id)))

    if !resp.IsSuccess() {
        return resp.StatusCode(), fmt.Errorf("%v", resp.RawResponse)
    }

    return resp.StatusCode(), err
}

// DeleteDir remove a directory from ceph fs (https://docs.ceph.com/en/latest/mgr/ceph_api/#delete--api-cephfs-fs_id-tree).
func (c *Client) DeleteDir(id int, path string) (status int, err error) {
    var resp *resty.Response

    client := *c.Session.Client

    resp, err = client.R().
        SetHeaders(defaultHeaders).
        SetQueryParam("path", path).
        Delete(c.Session.Server.getURL(fmt.Sprintf("cephfs/%d/tree", id)))

    if !resp.IsSuccess() {
        return resp.StatusCode(), fmt.Errorf("%v", resp.RawResponse)
    }

    return resp.StatusCode(), err
}

// GetQuota gets ceph fs quota for given path (https://docs.ceph.com/en/latest/mgr/ceph_api/#get--api-cephfs-fs_id-quota).
func (c *Client) GetQuota(id int, path string) (status int, quotas Quota, err error) {
    var resp *resty.Response

    client := *c.Session.Client

    resp, err = client.R().
        SetHeaders(defaultHeaders).
        SetResult(&quotas).
        SetQueryParam("path", path).
        Get(c.Session.Server.getURL(fmt.Sprintf("cephfs/%d/quota", id)))

    if !resp.IsSuccess() {
        return resp.StatusCode(), quotas, fmt.Errorf("%v", resp.RawResponse)
    }

    return resp.StatusCode(), quotas, err
}

// SetQuota sets ceph fs quota defined by path (https://docs.ceph.com/en/latest/mgr/ceph_api/#put--api-cephfs-fs_id-quota).
func (c *Client) SetQuota(id int, quota Quota) (status int, err error) {
    var resp *resty.Response

    client := *c.Session.Client

    resp, err = client.R().
        SetHeaders(defaultHeaders).
        SetBody(quota).
        Put(c.Session.Server.getURL(fmt.Sprintf("cephfs/%d/quota", id)))

    if !resp.IsSuccess() {
        return resp.StatusCode(), fmt.Errorf("%v", resp.RawResponse)
    }

    return resp.StatusCode(), err
}
