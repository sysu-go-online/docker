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
