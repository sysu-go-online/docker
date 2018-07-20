package server

//********************************************
// Author : huziang
//   新建客户端，在此文件添加路由转发
//********************************************

import (
	"github.com/codegangsta/martini"
	"github.com/unrolled/render"
)

// NewServer 新建客户端
func NewServer() *martini.ClassicMartini {
	formatter := render.New(render.Options{
		IndentJSON: true,
	})
	m := martini.Classic()

	initRoutes(m, formatter)

	return m
}

// 初始化路由
func initRoutes(m *martini.ClassicMartini, formatter *render.Render) {
	m.Get("/", HandleConnection(formatter))
	// m.Get("/api/auth", HandleAuth(formatter))
	// m.Get("/test", TestFunciton(formatter))
}
