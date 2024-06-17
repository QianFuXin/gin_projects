package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	router := gin.Default()

	// 设置图片存储的目录
	uploadDir := "./uploaded_images"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		err := os.Mkdir(uploadDir, 0755)
		if err != nil {
			return
		}
	}

	// 返回上传图片的页面
	router.GET("/", func(c *gin.Context) {
		html := `
<!DOCTYPE html>
<html lang="zh">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>图片上传页面</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            display: flex;
            justify-content: center;
            align-items: center;
            height: 100vh;
            margin: 0;
            background-color: #f4f4f4;
        }
        .upload-container {
            background: white;
            padding: 20px;
            border-radius: 8px;
            box-shadow: 0 4px 8px rgba(0,0,0,0.1);
        }
        .upload-btn {
            background-color: #4CAF50;
            color: white;
            padding: 10px 20px;
            border: none;
            border-radius: 4px;
            cursor: pointer;
            transition: background-color 0.3s;
        }
        .upload-btn:hover {
            background-color: #45a049;
        }
        input[type="file"] {
            margin-bottom: 10px;
        }
    </style>
</head>
<body>
    <div class="upload-container">
        <h2>上传图片</h2>
        <form action="/upload" method="post" enctype="multipart/form-data">
            <input type="file" id="file-upload" name="file" accept="image/*">
            <button type="submit" class="upload-btn">上传图片</button>
        </form>
    </div>
</body>
</html>
`
		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, html)
	})

	// 处理图片上传
	router.POST("/upload", func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("Get form err: %s", err.Error()))
			return
		}

		src, err := file.Open()
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("File open err: %s", err.Error()))
			return
		}
		defer func(src multipart.File) {
			err := src.Close()
			if err != nil {

			}
		}(src)

		hash := sha256.New()
		if _, err := io.Copy(hash, src); err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error computing hash: %s", err.Error()))
			return
		}

		filename := hex.EncodeToString(hash.Sum(nil)) + filepath.Ext(file.Filename)
		fileph := filepath.Join(uploadDir, filename)
		if _, err := os.Stat(fileph); os.IsNotExist(err) {
			if err := c.SaveUploadedFile(file, fileph); err != nil {
				c.String(http.StatusInternalServerError, fmt.Sprintf("Upload file err: %s", err.Error()))
				return
			}
		}
		htmlTemplate := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<title>上传成功</title>
<style>
    body {
        font-family: Arial, sans-serif;
        display: flex;
        justify-content: center;
        align-items: center;
        height: 100vh;
        background-color: #f0f0f0;
    }
    #copyButton {
        padding: 10px;
        background-color: #4CAF50;
        color: white;
        border: none;
        border-radius: 5px;
        cursor: pointer;
        margin-right: 10px;
    }
    a {
        padding: 10px;
        background-color: #007BFF;
        color: white;
        text-decoration: none;
        border-radius: 5px;
    }
</style>
<script>
    function copyURL() {
        var copyText = document.getElementById("urlField");
        copyText.select();
        document.execCommand("copy");
        alert("Copied the URL: " + copyText.value);
    }
</script>
</head>
<body>
    <div>
        <p>文件上传成功！</p>
        <input type="text" value="http://%s/files/%s" id="urlField" readonly style="width:300px;">
        <button id="copyButton" onclick="copyURL()">复制URL</button>
        <a href="http://%s/files/%s" target="_blank">访问文件</a>
    </div>
</body>
</html>
`, c.Request.Host, filename, c.Request.Host, filename)
		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, htmlTemplate)

	})

	// 提供图片访问
	router.Static("/files", uploadDir)

	err := router.Run(":5000")
	if err != nil {
		return
	}
}
