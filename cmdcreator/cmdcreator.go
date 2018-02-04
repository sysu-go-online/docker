package cmdcreator

// Command 解析客户端发送的信息
type Command struct {
	Command    string   `json:"command"`
	Entrypoint []string `json:"entrypoint"`
	PWD        string   `json:"pwd"`
	ENV        []string `json:"env"`
	UserName   string   `json:"user"`
}
