# ceph-rest-client

A rest client for ceph rest api. 

See https://docs.ceph.com/en/latest/mgr/ceph_api/#specification.

Not all endpoints are implemented. 
Still under development. 

Implemented ceph rest endpoints: 

## AUTH
- https://docs.ceph.com/en/latest/mgr/ceph_api/#post--api-auth
- https://docs.ceph.com/en/latest/mgr/ceph_api/#post--api-auth-logout

## CEPHFS
- https://docs.ceph.com/en/latest/mgr/ceph_api/#get--api-cephfs
- https://docs.ceph.com/en/latest/mgr/ceph_api/#get--api-cephfs-fs_id
- https://docs.ceph.com/en/latest/mgr/ceph_api/#get--api-cephfs-fs_id-get_root_directory
- https://docs.ceph.com/en/latest/mgr/ceph_api/#get--api-cephfs-fs_id-ls_dir
- https://docs.ceph.com/en/latest/mgr/ceph_api/#get--api-cephfs-fs_id-quota
- https://docs.ceph.com/en/latest/mgr/ceph_api/#put--api-cephfs-fs_id-quota
- https://docs.ceph.com/en/latest/mgr/ceph_api/#delete--api-cephfs-fs_id-snapshot
- https://docs.ceph.com/en/latest/mgr/ceph_api/#post--api-cephfs-fs_id-snapshot
- https://docs.ceph.com/en/latest/mgr/ceph_api/#delete--api-cephfs-fs_id-tree

## RBD
- https://docs.ceph.com/en/latest/mgr/ceph_api/#get--api-block-image
- https://docs.ceph.com/en/latest/mgr/ceph_api/#post--api-block-image
- https://docs.ceph.com/en/latest/mgr/ceph_api/#delete--api-block-image-image_spec
- https://docs.ceph.com/en/latest/mgr/ceph_api/#put--api-block-image-image_spec
- https://docs.ceph.com/en/latest/mgr/ceph_api/#post--api-block-image-image_spec-copy
- https://docs.ceph.com/en/latest/mgr/ceph_api/#post--api-block-image-image_spec-move_trash