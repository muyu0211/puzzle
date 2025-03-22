package service

import (
	"Puzzle/api/image"
	"context"
)

// IImage				1.定义接口
type IImage interface {
	GetOriginImage(ctx context.Context) (map[string]interface{}, error)
	RandomSplitImage(ctx context.Context, req *image.UploadImageReq) []int
	GetShuffledImage(ctx context.Context, fileName string, row, col int, randomImages []interface{}) ([]byte, error)
}

// localImage			2.定义接口变量
var localImage IImage

// RegisterImage 		3.定义一个接口实现的注册方法
func RegisterImage(i IImage) {
	localImage = i
}

// Image			 	4.定义一个获取接口实例的函数
func Image() IImage {
	if localImage == nil {
		panic("IImage接口未实现")
	}
	return localImage
}
