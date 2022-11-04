package httpService

import (
	"dream/manhua"
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func StartFileServer(p int, path string) {
	port = strconv.FormatInt(int64(p), 10)
	rootPath = path
	r := setupRouter()
	r.Run(":" + port)
}

var port string
var rootPath string

//go:embed static/*
var fs embed.FS

type FileInfoModel struct {
	FileName    string `json:"fileName"`
	DownloadUrl string `json:"downloadUrl"`
	FilePath    string `json:"filePath"`
	UpdateTime  string `json:"updateTime"`
	Size        string `json:"size"`
}

func setupRouter() *gin.Engine {
	// 初始化 Gin 框架默认实例，该实例包含了路由、中间件以及配置信息
	r := gin.Default()
	r.GET("static/*filepath", func(ctx *gin.Context) {
		s := http.FileServer(http.FS(fs))
		s.ServeHTTP(ctx.Writer, ctx.Request)
	})
	r.Static("/download", rootPath)
	// 加载网页模板
	t, _ := template.ParseFS(fs, "static/*.html")
	r.SetHTMLTemplate(t)
	r.GET("upload", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", gin.H{})
	})
	r.GET("file", func(ctx *gin.Context) {
		files := fetchFiles()
		ctx.HTML(http.StatusOK, "file.html", gin.H{
			"result": files,
		})
	})
	// 文件上传
	r.POST("/upload", func(ctx *gin.Context) {
		form, err := ctx.MultipartForm()
		if err != nil {
			ctx.JSON(http.StatusExpectationFailed, err)
		}
		f := form.File["file"]
		for _, file := range f {
			filePath := rootPath + "/" + file.Filename
			ctx.SaveUploadedFile(file, filePath)
			ctx.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file.Filename))
			ctx.File(filePath)
		}
	})
	r.GET("file/fetch", func(ctx *gin.Context) {
		files := fetchFiles()
		ctx.JSON(http.StatusOK, files)
	})
	return r
}

func fetchFiles() []FileInfoModel {
	var files []FileInfoModel
	root := rootPath
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		relateFilePath := strings.Replace(path, rootPath, "", 1)
		outPath := "/download" + relateFilePath
		var fileModel FileInfoModel
		fileModel.FilePath = relateFilePath
		fileModel.FileName = info.Name()
		fileModel.DownloadUrl = outPath
		fileModel.UpdateTime = info.ModTime().Format("2006-01-02 15:04:05")
		fileModel.Size = getFileSize(info.Size())
		files = append(files, fileModel)
		return nil
	})
	return files
}

func getFileSize(size int64) string {
	if size < 1024 {
		return strconv.FormatInt(size, 10) + " B"
	}
	kb := float64(size) / 1024.0
	if kb < 1024 {
		return fmt.Sprintf("%.2f", kb) + " KB"
	}
	mb := kb / 1024.0
	if mb < 1024 {
		return fmt.Sprintf("%.2f", mb) + " MB"
	}
	gb := mb / 1024.0
	if gb < 1024 {
		return fmt.Sprintf("%.2f", gb) + " GB"
	}
	tb := gb / 1024.0
	return fmt.Sprintf("%.2f", tb) + " TB"
}

func StartManhuaService(port int) {
	r := gin.Default()
	r.GET("search", func(ctx *gin.Context) {
		page, _ := strconv.Atoi(ctx.Query("page"))
		ctx.JSON(http.StatusOK, manhua.Search(ctx.Query("key"), page).Data)
	})
	r.POST("chapter", func(c *gin.Context) {
		var params map[string]string
		err := c.BindJSON(&params)
		if err != nil {
			// 返回错误信息
			c.String(http.StatusBadRequest, err.Error())
			// 执行退出
			c.Abort()
		}
		// 返回的 code 和 对应的参数星系
		c.JSON(http.StatusOK, manhua.GetChapterList(params["url"]))
	})
	r.POST("detail", func(c *gin.Context) {
		var params map[string]string
		err := c.BindJSON(&params)
		if err != nil {
			// 返回错误信息
			c.String(http.StatusBadRequest, err.Error())
			// 执行退出
			c.Abort()
		}
		// 返回的 code 和 对应的参数星系
		c.JSON(http.StatusOK, manhua.GetChapterDetail(params["url"]))
	})
	r.Run(":" + strconv.FormatInt(int64(port), 10))
}
