package types

// CreateContainerRequest is request form of create container
type CreateContainerRequest struct {
	Image   string   `json:"image"`
	PWD     string   `json:"pwd"`
	ENV     []string `json:"env"`
	MNT     []string `json:"mnt"`
	Target  []string `json:"target"`
	Network []string `json:"network"`
}

// CreateContainerResponse is the response for creating container
type CreateContainerResponse struct {
	ID  string `json:"id"`
	OK  bool   `json:"ok"`
	Msg string `json:"msg"`
}

// ConnectContainerRequest contains msg to be written to the container labeled by id
type ConnectContainerRequest struct {
	ID  string `json:"id"`
	Msg string `json:"msg"`
}

// ConnectContainerResponse is stores data from container
type ConnectContainerResponse struct {
	OK  bool   `json:"ok"`
	Msg string `json:"msg"`
}

// ResizeContainerRequest is request for resize container
type ResizeContainerRequest struct {
	ID     string `json:"id"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

// ResizeContainerResponse is response for resize container
type ResizeContainerResponse struct {
	OK  bool   `json:"ok"`
	Msg string `json:"msg"`
}
