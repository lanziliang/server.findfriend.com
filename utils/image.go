package utils

//图片处理
import (
	"code.google.com/p/graphics-go/graphics"
	"github.com/revel/revel"
	"image"
	"image/jpeg"
	"os"
)

//缩略图
//参数说明：mode{1:xy固定 2:x固定等比例缩放 3:y固定等比例缩放}
func Thumbnail(img_path, old_name, new_name string, mode, x, y int) error {
	file, err := os.Open(img_path + old_name)
	if err != nil {
		revel.ERROR.Println(err)
		return err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		revel.ERROR.Println(err)
		return err
	}
	bound := img.Bounds()
	dx := bound.Dx()
	dy := bound.Dy()

	var newdx, newdy int
	switch mode {
	case 1:
		newdx = x
		newdy = y
	case 2:
		newdx = x
		newdy = x * dy / dx
	case 3:
		newdx = y * dx / dy
		newdy = y
	}

	// 缩略图的大小
	dst := image.NewRGBA(image.Rect(0, 0, newdx, newdy))
	// 产生缩略图
	err = graphics.Scale(dst, img)
	if err != nil {
		revel.ERROR.Println(err)
		return err
	}

	toimg, err := os.Create(img_path + new_name)
	if err != nil {
		revel.ERROR.Println(err)
		return err
	}
	defer toimg.Close()

	err = jpeg.Encode(toimg, dst, &jpeg.Options{15})
	if err != nil {
		revel.ERROR.Println(err)
		return err
	}

	return nil
}
