package image

import (
	"Puzzle/api/image"
	"Puzzle/internal/service"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/gogf/gf/v2/frame/g"
	stdimage "image"
	"image/draw"
	"image/jpeg"
	"math"
	"math/rand"
	"os"
)

type sImage struct{}

func init() {
	service.RegisterImage(&sImage{})
}

// GetOriginImage 从Session中读取图片，行列等信息
func (*sImage) GetOriginImage(ctx context.Context) (map[string]interface{}, error) {
	var (
		r = g.RequestFromCtx(ctx)
	)
	data, err := r.Session.Data()
	if err != nil {
		return nil, errors.New("session 数据不完整")
	}

	// 获取session中的数据
	fileName, ok1 := data["image"].(string)
	cols, ok2 := data["cols"]
	rows, ok3 := data["rows"]
	randomImages, ok4 := data["randomImages"]
	if !ok1 || !ok2 || !ok3 || !ok4 {
		return nil, errors.New("session 数据不完整")
	}

	resMap := map[string]interface{}{
		"fileName":     fileName,
		"cols":         cols,
		"rows":         rows,
		"randomImages": randomImages,
	}
	return resMap, nil
}

// RandomSplitImage 根据传递的Rows和Cols进行随机切图
func (*sImage) RandomSplitImage(ctx context.Context, req *image.UploadImageReq) []int {
	var (
		cols = req.Cols
		rows = req.Rows
	)
	randomSplit := make([]int, cols*rows)
	for i := range randomSplit {
		randomSplit[i] = i
	}
	// 打乱切片
	rand.Shuffle(len(randomSplit), func(i, j int) {
		randomSplit[i], randomSplit[j] = randomSplit[j], randomSplit[i]
	})
	return randomSplit
}

// GetShuffledImage 获取打乱的图片
func (*sImage) GetShuffledImage(ctx context.Context, fileName string, row, col int, randomImages []interface{}) ([]byte, error) {
	// 读取图片数据
	imgData, readErr := os.ReadFile(fileName)
	if readErr != nil {
		return nil, errors.New("图片未找到")
	}
	img, _, decodeErr := stdimage.Decode(bytes.NewReader(imgData))
	if decodeErr != nil {
		return nil, errors.New("图片解码失败")
	}
	// 获取图片的尺寸
	var (
		imgBounds = img.Bounds()
		width     = imgBounds.Dx()
		height    = imgBounds.Dy()
	)

	// 计算单个块的宽度和高度
	blockWidth := int(math.Ceil(float64(width) / float64(col)))
	blockHeight := int(math.Ceil(float64(height) / float64(row)))

	// 创建一个新的空白图片（RGBA 格式）
	newImg := stdimage.NewRGBA(stdimage.Rect(0, 0, width, height))

	// 将图片分块并根据随机顺序重新组合
	for idx, blockIdxJson := range randomImages {
		blockIdx, _ := blockIdxJson.(json.Number).Int64()
		// 计算原始块的坐标
		srcX := (int(blockIdx) % col) * blockWidth
		srcY := (int(blockIdx) / col) * blockHeight
		srcRect := stdimage.Rect(srcX, srcY, srcX+blockWidth, srcY+blockHeight)

		// 计算目标块的坐标
		dstX := (idx % col) * blockWidth
		dstY := (idx / col) * blockHeight
		dstRect := stdimage.Rect(dstX, dstY, dstX+blockWidth, dstY+blockHeight)

		// 将块从原始图片拷贝到新图片的目标位置
		draw.Draw(newImg, dstRect, img, srcRect.Min, draw.Src)
	}

	// 创建字节缓冲区
	var buf bytes.Buffer

	// 将新图片编码为 JPEG 格式，并写入缓冲区
	encodeErr := jpeg.Encode(&buf, newImg, nil)
	if encodeErr != nil {
		return nil, errors.New("图片编码失败")
	}

	// 返回字节流
	return buf.Bytes(), nil
}
