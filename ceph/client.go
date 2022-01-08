package ceph

type Client struct {
	Session *Session
}

func New(server Server) (client *Client, err error) {

	client = &Client{}
	client.Session, err = NewSession(server)

	if err != nil {
		return nil, err
	}

	return client, nil
}
