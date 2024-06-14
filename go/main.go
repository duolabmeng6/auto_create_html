package main

import (
	"fmt"
	"github.com/CloudyKit/jet/v6"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"os"
	"path/filepath"
)

var views = jet.NewSet(
	jet.NewOSFileSystemLoader("./views"),
	jet.InDevelopmentMode(),
)

func main() {
	e := echo.New()

	// 中间件
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// 定义动态路由
	e.GET("*", func(c echo.Context) error {
		requestedPath := c.Param("*")
		templatePath := requestedPath
		if requestedPath == "" {
			templatePath = "index.html"
		}

		// 检查模板文件是否存在
		if _, err := os.Stat(filepath.Join("./views", templatePath)); err == nil {
			t, err := views.GetTemplate(templatePath)
			if err != nil {
				return c.String(http.StatusInternalServerError, fmt.Sprintf("Error loading template: %v", err))
			}

			// 打开输出文件
			outputFile, err := os.Create(filepath.Join("./dist", templatePath))
			if err != nil {
				return c.String(http.StatusInternalServerError, fmt.Sprintf("Error creating output file: %v", err))
			}
			defer outputFile.Close()

			// 执行模板并写入文件
			err = t.Execute(outputFile, nil, nil)
			if err != nil {
				return c.String(http.StatusInternalServerError, fmt.Sprintf("Error executing template: %v", err))
			}

			// 将渲染结果输出到浏览器
			return c.File(filepath.Join("./dist", templatePath))
		}

		// 如果模板文件不存在，交给静态文件处理
		return c.File(filepath.Join("./dist", templatePath))
	})
	// 静态文件目录
	e.Static("/static", "./dist")

	// 启动服务器
	e.Logger.Fatal(e.Start(":3000"))
}
