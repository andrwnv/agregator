package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"strings"
)

type FileController struct {
	DownloadPath    string
	httpContentType map[string]string
}

func NewFileController(rootPath string) *FileController {
	return &FileController{
		DownloadPath: rootPath,
		httpContentType: map[string]string{
			".jpg": "image/jpeg",
			".png": "image/png",
			".gif": "image/gif",
		},
	}
}

func (c *FileController) GetImage(ctx *gin.Context) {
	fileName := ctx.Param("filename")
	targetPath := filepath.Join(c.DownloadPath, fileName)

	if !strings.HasPrefix(filepath.Clean(targetPath), c.DownloadPath) {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "Look like you attacking me",
		})
		return
	}

	fileExt := filepath.Ext(fileName)
	if fileExt == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Cant find file extension",
		})
		return
	}

	ctx.Header("Content-Disposition", "inline")
	ctx.Header("Content-Type", c.httpContentType[fileExt])
	ctx.File(targetPath)

	ctx.Status(http.StatusNotModified)
}
