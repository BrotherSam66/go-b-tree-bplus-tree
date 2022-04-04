// Package bplustree
// @Title B树工具包
// @Description  查找节点
// @Author  https://github.com/BrotherSam66/
// @Update
package bplustree

import (
	"errors"
	"fmt"
	"go-b-tree-bplus-tree/bplustree/bplustreeglobal"
	"go-b-tree-bplus-tree/bplustree/bplustreemodels"
)

// Search 查找节点
// @key 键值
// @tempNode 找到的节点指针（可能是适合插入的位置）
// @isTarget 找到的是命中的节点
// @Author  https://github.com/BrotherSam66/
func Search(key int) (tempNode *bplustreemodels.BPTreeNode, isTarget bool, err error) {
	if bplustreeglobal.Root == nil {
		fmt.Println("这个树/分支是空的")
		return nil, false, errors.New("这个树/分支是空的！")
	}

	tempNode = bplustreeglobal.Root // 临时的指针
	var i int                       // 循环外定义，是因为循环后要用到这个变量
	for {                           // 递归循环
		// 循环本层关键字key[]
		for i = 0; i < len(tempNode.Key); i++ { // 必须看.Key。因为可能是叶子，需要判断是否完美找到
			if key == tempNode.Key[i] { // 准确命中，可能不是否叶子，
				if len(tempNode.Child) == 0 || tempNode.Child[0] == nil { // 是末端叶子
					return tempNode, true, nil
				}
				break
			} else if key < tempNode.Key[i] { // 小于，说明刚刚越过了，向这个tempNode.key的左腿递归
				//tempNode=tempNode.Child[i] // 后面有这句
				break
			}
			// 到这里：可能①会向后找；可能②KeyNum循环结束，下级得到的i是最右key的右边，向。
		}
		// tempNode 我是叶子，就算是不完美找到了，返回
		if len(tempNode.Child) == 0 || tempNode.Child[0] == nil {
			return tempNode, false, nil
		}
		// 下移一层
		if i >= len(tempNode.Child) { // for溢出，说明新 key 大于所有原 key，向最右节点查找
			i = len(tempNode.Child) - 1
		}
		tempNode = tempNode.Child[i] // 如果是①break过来的，这是正确Key的左腿；如果②是循环结束过来的，这是最后一个KEY的右腿。刚刚好
	}
	return
}

// FindKeyPosition 找插入或者替换key的位置
// @intSlice 被查找切片
// @key 键值
// @insertPosition 拟插入位置，realPosition>-1就无效
// @realPosition 准确命中的位置。==-1 表示没找到
// @Author  https://github.com/BrotherSam66/
func FindKeyPosition(intSlice *[]int, key int) (insertPosition, realPosition int) {
	var i int
	for i = 0; i < len(*intSlice); i++ {
		if (*intSlice)[i] == key { // 完美找到
			return -1, i
		}
		if (*intSlice)[i] > key { // 插入的位置（刚过把的位置）
			return i, -1
		}
	}
	return i, -1 // 越过最右，key附加在尾巴
}
