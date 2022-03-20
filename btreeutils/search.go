// Package btreeutils
// @Title B树工具包
// @Description  查找节点
// @Author  https://github.com/BrotherSam66/
// @Update
package btreeutils

import (
	"errors"
	"fmt"
	"go-b-tree-bplus-tree/btreemodels"
	"go-b-tree-bplus-tree/global"
)

// Search 查找节点
// @key 键值
// @tempNode 找到的节点指针（可能是适合插入的位置）
// @isTarget 找到的是命中的节点
func Search(key int) (tempNode *btreemodels.BTreeNode, isTarget bool, err error) {
	if global.Root == nil {
		fmt.Println("这个树/分支是空的")
		return nil, false, errors.New("这个树/分支是空的！")
	}

	tempNode = global.Root // 临时的指针
	var i int              // 循环外定义，是因为循环后要用到这个变量
	for {                  // 递归循环
		// 循环本层关键字key[]
		for i = 0; i < tempNode.KeyNum; i++ {
			if key == tempNode.Key[i] { // 准确命中，完美找到，不管是否叶子，返回
				return tempNode, true, nil
			} else if key < tempNode.Key[i] { // 小于，说明刚刚越过了，向这个tempNode.key的左腿递归
				//tempNode=tempNode.Child[i] // 后面有这句
				break
			}
			// 到这里：可能①会向后找；可能②KeyNum循环结束，下级得到的i是最右key的右边，向。
		}
		if tempNode.Child[0] == nil { // tempNode是叶子，就算是找到了，返回
			return
		}
		// 下移一层
		tempNode = tempNode.Child[i] // 如果是①break过来的，这是正确Key的左腿；如果②是循环结束过来的，这是最后一个KEY的右腿。刚刚好
	}
	return

}
