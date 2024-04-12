package client

import (
	"blog-server/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"path"
	"time"
)

type FileApi struct {
	userService service.UserService
	roleService service.RoleService
	postService service.PostService
}

// 上传图像
func (a FileApi) Upload(ctx *gin.Context) {
	file, header, err := ctx.Request.FormFile("file")

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "格式错误",
		})
		return
	}

	filename := header.Filename
	ext := path.Ext(filename)
	// 用上传时间作为文件名
	name := "image_" + time.Now().Format("20060102150405")
	newFilename := name + ext
	out, err := os.Create("static/images/" + newFilename)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "创建错误",
		})
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "复制错误",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{"filePath": "/images/" + newFilename},
		"msg":  "上传成功",
	})
}

// 上传富文本编辑器中的图像
func (a FileApi) RichEditorUpload(ctx *gin.Context) {
	formData, _ := ctx.MultipartForm()
	files := formData.File["wangeditor-uploaded-image"]
	var url []string

	for _, file := range files {
		ext := path.Ext(file.Filename)
		name := "image_" + time.Now().Format("20060102150405")
		newFilename := name + ext
		dst := path.Join("./static/images", newFilename)
		fileurl := "/images/" + newFilename
		url = append(url, fileurl)
		if err := ctx.SaveUploadedFile(file, dst); err != nil {
			fmt.Println(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"errno":   1,
				"message": "上传失败",
			})
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"errno": 0,
		"data": gin.H{
			"url": url[0],
		},
	})
}
