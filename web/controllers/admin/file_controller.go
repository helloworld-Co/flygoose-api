package admin

import (
	"flygoose/configs"
	"flygoose/pkg/tools"
	"flygoose/web/controllers/comm"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"path/filepath"
	"strings"
)

type FileController struct {
	comm.BaseComponent
}

func NewFileController() *FileController {
	return &FileController{}
}

func (c *FileController) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle("POST", "/uploadImg", "UploadImg", iris.LimitRequestBodySize(50*1024*1024)) //上传图片
}

func (c *FileController) UploadImg() {
	f, fh, err := c.Ctx.FormFile("uploadfile")
	if f != nil {
		defer f.Close()
	}

	if err != nil {
		c.RespFailedMessage("读取uploadfile失败：" + err.Error())
		return
	}

	//后缀名
	var suffix string
	suffixArr := fh.Header["Content-Type"]
	if len(suffixArr) > 0 {
		var arr = strings.Split(suffixArr[0], "/")
		if len(arr) > 1 && arr[0] == "image" {
			suffix = arr[1]
		}
	}
	if suffix == "" {
		suffix = ".png"
	}
	//生成文件名
	var absStaticImgDir = filepath.Join(configs.Cfg.ExecuteDir, configs.Cfg.StaticImgDir)
	var newFileName = tools.GenNumberCode(12) + "." + suffix
	var newAbsFilePath = filepath.Join(absStaticImgDir, newFileName)

	_, err = c.Ctx.SaveFormFile(fh, newAbsFilePath)
	if err != nil {
		c.RespFailedMessage("保存文件失败:" + err.Error())
		return
	}

	c.RespSuccess(iris.Map{
		"filename": filepath.Join(configs.Cfg.StaticImgDir, newFileName),
	}, "上传文件成功")
}
