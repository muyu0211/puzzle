package image

import (
	"Puzzle/api/image"
	"Puzzle/internal/service"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"os"
)

// Image 控制器对象
type Image struct{}

// NewImage 创建控制器对象
func NewImage() *Image {
	return &Image{}
}

// Index 项目首页, 路由: GET :'/puzzle'
func (i *Image) Index(ctx context.Context, req *image.IndexReq) (res *image.IndexRes, err error) {
	var (
		r = g.RequestFromCtx(ctx)
	)
	// 跳转到主页面
	_ = r.Response.WriteTpl("puzzle/include/main.html")
	return
}

// UploadImage 上传图片, 路由: POST :'/upload'
func (i *Image) UploadImage(ctx context.Context, req *image.UploadImageReq) (res *image.UploadImageRes, err error) {
	res = &image.UploadImageRes{}
	var (
		r        = g.RequestFromCtx(ctx)
		img      = r.GetUploadFile("image")
		cols     = req.Cols
		rows     = req.Rows
		fileName = ""
	)

	// 对行列进行判断
	if cols > 8 || rows > 8 {
		res.Code = 400
		res.Message = "列数或行数不能超过8."
		return
	} else if cols <= 0 || rows <= 0 {
		res.Code = 401
		res.Message = "列数或行数不能小于等于0."
		return
	}
	if img == nil {
		res.Code = 402
		res.Message = "没有接收到传递的图像."
		return
	} else {
		fileName = "resource/public/upload/" + img.Filename
		_, saveErr := img.Save("resource/public/upload/")
		if saveErr != nil {
			res.Code = 403
			res.Message = "图片保存失败."
			return
		}
	}

	// 调用service层的业务逻辑
	res.Code = 200
	res.Message = "上传成功"
	res.RandomImages = service.Image().RandomSplitImage(ctx, req)

	// 将图片数据保存至Session
	_ = r.Session.Set("image", fileName)
	_ = r.Session.Set("cols", cols)
	_ = r.Session.Set("rows", rows)
	_ = r.Session.Set("randomImages", res.RandomImages)
	return
}

// StartPuzzle 跳转至拼图页面, 路由: POST :'/start'
func (i *Image) StartPuzzle(ctx context.Context, req *image.StartReq) (res *image.StartRes, err error) {
	var (
		r = g.RequestFromCtx(ctx)
	)
	fmt.Println("start:", req.RandomImages)
	// 跳转到拼图页面
	_ = r.Response.WriteTpl("puzzle/include/puzzle.html")
	return
}

// GetImage 获取原始图片, 路由: GET :'/getImage'
func (i *Image) GetImage(ctx context.Context, req *image.GetImageReq) (res *image.GetImageRes, err error) {
	resMap, err := service.Image().GetOriginImage(ctx)
	var (
		r            = g.RequestFromCtx(ctx)
		fileName     = resMap["fileName"].(string)
		cols         = resMap["cols"]
		rows         = resMap["rows"]
		randomImages = resMap["randomImages"]
	)
	// 读取图片数据
	imgData, readErr := os.ReadFile(fileName)
	if readErr != nil {
		res.Code = 404
		res.Message = "图片未找到"
		return
	}

	// 将二进制数据转换为 Base64 编码的字符串
	imgBase64 := base64.StdEncoding.EncodeToString(imgData)
	response := map[string]interface{}{
		"image":        fmt.Sprintf("data:image/jpeg;base64,%s", imgBase64),
		"cols":         cols,
		"rows":         rows,
		"randomImages": randomImages,
	}

	// 设置响应头为图片格式（如 PNG）
	r.Response.Header().Set("Content-Type", "application/json")
	// 返回图片数据
	r.Response.WriteJson(response)
	return
}

// GetShuffles 获取拼图后的图片, 路由: GET :'/getShuffles'
func (i *Image) GetShuffles(ctx context.Context, req *image.GetShufflesReq) (res *image.GetShufflesRes, err error) {
	resMap, err := service.Image().GetOriginImage(ctx)
	var (
		r            = g.RequestFromCtx(ctx)
		fileName     = resMap["fileName"].(string)
		cols, _      = resMap["cols"].(json.Number).Int64()
		rows, _      = resMap["rows"].(json.Number).Int64()
		randomImages = resMap["randomImages"].([]interface{})
		imageBytes   []byte
	)

	// 获取打乱后的图片，以字节的形式返回
	imageBytes, err = service.Image().GetShuffledImage(ctx, fileName, int(rows), int(cols), randomImages)

	// 将二进制数据转换为 Base64 编码的字符串
	imgBase64 := base64.StdEncoding.EncodeToString(imageBytes)
	response := map[string]interface{}{
		"image": fmt.Sprintf("data:image/jpeg;base64,%s", imgBase64),
	}
	// 设置响应头为图片格式（如 PNG）
	r.Response.Header().Set("Content-Type", "application/json")
	// 返回图片数据
	r.Response.WriteJson(response)
	return
}

// CheckPuzzle 检查拼图是否拼好, 路由: POST :'/checkPuzzle'
func (i *Image) CheckPuzzle(ctx context.Context, req *image.CheckPuzzleReq) (res *image.CheckPuzzleRes, err error) {
	r := g.RequestFromCtx(ctx)
	res = &image.CheckPuzzleRes{}
	for i := 1; i < len(req.PuzzleImages); i++ {
		if req.PuzzleImages[i-1] > req.PuzzleImages[i] {
			res.Code = 400
			res.Message = "拼图失败"
			return
		}
	}
	res.Code = 200
	res.Message = "拼图成功"

	// 拼图成功后将存储在upload的图片删除
	data, _ := r.Session.Data()
	fmt.Println(data["image"].(string))
	removeErr := os.Remove(data["image"].(string))
	if removeErr != nil {
		res.Code = 200
		res.Message = "拼图成功"
		err = fmt.Errorf("删除图片时出错: %w", removeErr)
	}
	fmt.Println("图片删除成功:", data["image"])
	return
}
