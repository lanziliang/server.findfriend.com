package utils

import (
	//"fmt"
	"github.com/revel/revel"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path"
	"strings"
)

var (
	File_Suffix string = ".jpg,.jpeg,.png"
)

//上传文件
func Upload(savePath string, file multipart.File, header *multipart.FileHeader) (string, error) {
	//判断是否是系统的分隔符
	separator := "/"
	if os.IsPathSeparator('\\') {
		separator = "\\"
	} else {
		separator = "/"
	}

	fileName := header.Filename
	defer file.Close()

	//读取文件数据
	bytes, err := ioutil.ReadAll(file)
	if err != nil || len(bytes) == 0 {
		revel.ERROR.Println("Failed to read image:", err)
		return fileName, err
	}

	//文件类型检测
	if !strings.Contains(File_Suffix, path.Ext(fileName)) {
		revel.ERROR.Println("不支持上传该类型的文件!")
		return fileName, err
	}

	//字符串替换 /替换为系统分隔符
	savePath = strings.Replace(savePath, "/", separator, -1)

	//创建目录
	err = os.MkdirAll(savePath, os.ModePerm)
	if err != nil {
		revel.WARN.Println(err)
		return fileName, err
	}

	//保存文件
	err = ioutil.WriteFile(savePath+fileName, bytes, os.ModePerm)
	if err != nil {
		revel.WARN.Println(err)
		return fileName, err
	}

	return fileName, nil
}

//是否文件
func IsFile(filePath string) bool {
	f, e := os.Stat(filePath)
	if e != nil {
		return false
	}
	if f.IsDir() {
		return false
	}
	return true
}
