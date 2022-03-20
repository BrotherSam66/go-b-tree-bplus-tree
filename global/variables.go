// Package global
// @Title B树全局变量
// @Description  全局变量
// @Author  https://github.com/BrotherSam66/
// @Update
package global

import (
	"go-b-tree-bplus-tree/btreemodels"
)

// Root 根
var Root *btreemodels.BTreeNode

// KeyLen 彩色显示树，每个KEY字节长度 todo 根据输入的数字最大值，动态调整这个
var KeyLen int = 2

// MaxKey 最大key值
var MaxKey int = 100
