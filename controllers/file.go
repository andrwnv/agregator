package controllers

import (
	"fmt"
	"github.com/andrwnv/event-aggregator/misc"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type FileController struct {
	DownloadPath    string
	httpContentType map[string]string
}

func handleSaveError(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
		"error": "Unable to save the file(s)",
	})
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
	dirName := ctx.Query("uuid")
	targetPath := filepath.Join(c.DownloadPath, dirName, fileName)

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

func (c *FileController) UploadAvatarMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, parseErr := misc.ExtractJwtPayload(ctx)
		if parseErr {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Cant extract info from claims",
			})
			return
		}

		file, err := ctx.FormFile("file")
		if err != nil {
			ctx.Next()
			return
		}

		ext := filepath.Ext(file.Filename)
		newFileName := fmt.Sprintf("%s%s", uuid.New(), ext)

		pathForSave := path.Join(c.DownloadPath, user.ID)
		if os.MkdirAll(pathForSave, os.ModePerm) != nil {
			handleSaveError(ctx)
			return
		}

		if err := ctx.SaveUploadedFile(file, path.Join(pathForSave, newFileName)); err != nil {
			handleSaveError(ctx)
			return
		}

		ctx.Set("file-name", newFileName)
		ctx.Next()
	}
}

func (c *FileController) UploadImagesMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, parseErr := misc.ExtractJwtPayload(ctx)
		if parseErr {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Cant extract info from claims",
			})
			return
		}

		form, err := ctx.MultipartForm()
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Cant extract images",
			})
			return
		}

		files := form.File["files"]
		if len(files) == 0 {
			ctx.Next()
			return
		}

		pathForSave := path.Join(c.DownloadPath, user.ID)
		if os.MkdirAll(pathForSave, os.ModePerm) != nil {
			handleSaveError(ctx)
			return
		}

		for _, file := range files {
			ext := filepath.Ext(file.Filename)
			newFileName := fmt.Sprintf("%s%s", uuid.New(), ext)

			if err := ctx.SaveUploadedFile(file, path.Join(pathForSave, newFileName)); err != nil {
				handleSaveError(ctx)
				return
			}
		}

		ctx.Next()
	}
}
