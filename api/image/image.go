package image

import "github.com/gogf/gf/v2/frame/g"

// IndexReq IndexRes 跳转至项目首页
type IndexReq struct {
	g.Meta `path:"/puzzle" method:"get" summary:"Home page"`
}

type IndexRes struct {
	g.Meta `mime:"text/html" type:"string" example:"<html/>"`
}

// UploadImageReq UploadImageRes 上传拼图图片及划分行列数
type UploadImageReq struct {
	g.Meta `path:"/upload" method:"post" tags:"Image" summary:"上传图片"`

	Rows int `p:"rows" v:"required#缺失rows值" description:"行数"`
	Cols int `p:"cols" v:"required#缺失cols值" description:"列数"`
}

type UploadImageRes struct {
	g.Meta `mime:"text/html" type:"string" example:"<html/>"`

	Code         int    `json:"code" description:"状态码"`
	Message      string `json:"message" description:"提示信息"`
	RandomImages []int  `json:"randomImages" description:"随机打乱的每个图像块的编号"`
}

// StartReq StartRes 跳转至开始拼图界面
type StartReq struct {
	g.Meta `path:"/start" method:"post" summary:"拼图页面"`

	RandomImages []int `p:"randomImages" v:"required#缺失randomImages值" description:"随机打乱的每个图像块的编号"`
}

type StartRes struct {
	g.Meta `mime:"text/html" type:"string" example:"<html/>"`
}

// GetImageReq GetImageRes 从服务器获取原始图片
type GetImageReq struct {
	g.Meta `path:"/getImage" method:"get" summary:"获取Session中的图片"`
}

type GetImageRes struct {
	Code    int    `json:"code" description:"状态码"`
	Message string `json:"message" description:"提示信息"`
}

// GetShufflesReq GetShufflesRes 从服务器获取打乱后的图像块
type GetShufflesReq struct {
	g.Meta `path:"/getShuffles" method:"get" summary:"获取打乱后的图像块"`
}
type GetShufflesRes struct {
	Code    int    `json:"code" description:"状态码"`
	Message string `json:"message" description:"提示信息"`
}

// CheckPuzzleReq CheckPuzzleRes 检查拼图是否成功
type CheckPuzzleReq struct {
	g.Meta `path:"/checkPuzzle" method:"post" summary:"检查拼图是否成功"`

	PuzzleImages []int `p:"puzzleImages" v:"required#缺失puzzleImages值" description:"拼图后的图像块编号"`
}

type CheckPuzzleRes struct {
	Code    int    `json:"code" description:"状态码"`
	Message string `json:"message" description:"提示信息"`
}
