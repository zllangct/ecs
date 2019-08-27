package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func main()  {
	server:=gin.Default()
	gin.SetMode(gin.DebugMode)

	server.Use(gin.Logger())
	server.Use(gin.Recovery())

	server.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK,"hello,this is collection server.")
	})
	server.POST("/upload_data",DataHandlerBuffer)

	if err:=server.Run(":7708");err!=nil{
		println(err.Error())
	}
}

func DataHandlerBuffer(c *gin.Context)  {
	println(c.Request.Host)
	fileName :=c.PostForm("filename")
	image,_:=c.MultipartForm()
	img:= image.File["image"]
	savePath:="D:/test_resv/"
	if  !Exists(savePath){
		err:=os.MkdirAll(savePath,os.ModePerm)
		if err!=nil {
			println(err)
		}
	}
	err:=c.SaveUploadedFile(img[0],savePath+img[0].Filename)
	if err!=nil {
		println(err)
	}
	println(fileName)
}

func DataHandler(c *gin.Context) {
	println(c.Request.Host)
	fileName :=c.PostForm("filename")
	image,_:=c.FormFile("image")
	
	savePath:="D:/test_resv/"
	if  !Exists(savePath){
		err:=os.MkdirAll(savePath,os.ModePerm)
		if err!=nil {
			println(err)
		}
	}
	err :=c.SaveUploadedFile(image,savePath+image.Filename)
	if err!=nil {
		println(err)
	}
	println(fileName)
}

// 判断所给路径文件/文件夹是否存在
func Exists(path string) bool {
	_, err := os.Stat(path)    //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// 判断所给路径是否为文件
func IsFile(path string) bool {
	return !IsDir(path)
}