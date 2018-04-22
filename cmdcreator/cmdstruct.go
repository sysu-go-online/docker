package cmdcreator

//********************************************
// Author : huziang
//   cmd结构
//********************************************

// Command 解析客户端发送的信息
type Command struct {
	Command     string   `json:"command"`
	PWD         string   `json:"pwd"`
	ENV         []string `json:"env"`
	UserName    string   `json:"user"`
	ProjectName string   `json:"project"`
	Entrypoint  []string `json:"entrypoint"`
}
