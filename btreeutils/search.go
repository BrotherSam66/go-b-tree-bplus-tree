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
// @Author  https://github.com/BrotherSam66/
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

// PredecessorOrSuccessor 找前驱or后继Key。比我稍小的最大Key，比我大的最小key
// @key 键值
// @avatar 找到的替身节点指针
// @Author  https://github.com/BrotherSam66/
func PredecessorOrSuccessor(tempNode *btreemodels.BTreeNode, key int, isPredecessor bool) (avatar *btreemodels.BTreeNode, err error) {
	if tempNode.Child[0] == nil {
		err = errors.New("已经是叶子了，不可以找前驱or后继")
		fmt.Println("已经是叶子了，不可以找前驱or后继")
		return
	}

	// 精准找到拟删除KEY的位置，deletePoint
	deletePoint := 0
	for deletePoint = 0; deletePoint < tempNode.KeyNum; deletePoint++ {
		if tempNode.Key[deletePoint] == key { // 准确命中，只可能是新创建节点情形
			break
		}
	}
	if deletePoint >= tempNode.KeyNum {
		fmt.Println("发生某种错误，找到KEY又不存在了！ ")
		return
	}

	if isPredecessor { // 前驱
		avatar = tempNode.Child[deletePoint] // 命中点左边的腿
	} else { // 后继
		avatar = tempNode.Child[deletePoint+1] // 命中点右边的腿
	}

	for { // 递归循环
		if avatar.Child[0] == nil { // 到叶子，就算找到了，前驱Key=》尾巴，后继Key=头
			return
		}
		// 下移一层
		if isPredecessor { // 前驱
			avatar = avatar.Child[avatar.KeyNum] // 不断向最右一条腿找
		} else { // 后继
			avatar = avatar.Child[0] // 不断向最左一条腿找
		}
	}
}
