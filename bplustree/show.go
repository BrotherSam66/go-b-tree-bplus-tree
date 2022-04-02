// Package bplustree
// @Title B树工具包
// @Description  显示节点
// @Author  https://github.com/BrotherSam66/
// @Update
package bplustree

import (
	"fmt"
	"go-b-tree-bplus-tree/bplustree/bplustreeglobal"
	"go-b-tree-bplus-tree/bplustree/bplustreemodels"
)

// ShowTree 逐层显示这个树，
// @b 展示树的开始位置（可以不从root开始展示）
// @Author  https://github.com/BrotherSam66/
func ShowTree(b *bplustreemodels.BPTreeNode) {
	if b == nil { // 先判断，后定义变量，免得浪费内存
		fmt.Println("这个树/分支是空的")
		return
	}
	//fmt.Printf("%#v\n", b)

	// ShowTemp 数据。数据可能是nil。最多10层，每层最多1000 数据
	var showTemp [10][1000]*bplustreemodels.BPTreeNode
	var showTempLeaf [10000]*bplustreemodels.BPTreeLeaf
	//totalLevel := 0             // 总层数
	//nowLevel := 0               // 当前层数
	//nnn := Name
	nowColumn := 0     // 当前列
	nowColumnLeaf := 0 // 叶子队列当前列
	fmt.Printf("\n展示树：(父|KEY数)KEY\\腿|KEY\\腿|...")
	//ShowOneNode(global.Root)
	//return
	showTemp[0][0] = b // 来的最高位指针

	for i := 1; i < len(showTemp); i++ { // 循环每一层
		fmt.Println("")              // 先来一个换行
		if showTemp[i-1][0] == nil { // 上一层全nil，结束显示
			return
		}
		nowColumn = 0 // 当前列
		for j := 0; j < len(showTemp[0]); j++ {
			if showTemp[i-1][j] == nil { // 本行遍历结束
				break
			}

			ShowOneNode(showTemp[i-1][j]) // 显示遍历到的上一行的这个节点。显示没换行

			if len(showTemp[i-1][j].Child) > 0 { // 下级是分支节点
				// 在本行填写上一行所有节点的所有儿子
				for k := 0; k < len(showTemp[i-1][j].Key); k++ { // 其他的所有右腿
					showTemp[i][nowColumn] = showTemp[i-1][j].Child[k]
					nowColumn++
				}
			} else { // 下级是叶子节点
				// 在叶子队列填写上一行所有节点的所有叶子儿子
				for k := 0; k < len(showTemp[i-1][j].Key); k++ { // 其他的所有右腿
					showTempLeaf[nowColumnLeaf] = showTemp[i-1][j].Leaf[k]
					nowColumnLeaf++
				}
			}
		}
	}

	// 最后，展示叶子
	for j := 0; j < len(showTemp[0]); j++ {
		if showTempLeaf[j] == nil { // 本行遍历结束
			break
		}
		ShowOneLeaf(showTempLeaf[j]) // 显示遍历到的上一行的这个节点。显示没换行
	}

	// 最最后，展示Sqt头
	fmt.Printf("\n展示Sqt头 : ")
	ShowOneLeaf(bplustreeglobal.Sqt)
}

// ShowOneNode 展示单个节点
// @Author  https://github.com/BrotherSam66/
func ShowOneNode(n *bplustreemodels.BPTreeNode) {
	if n == nil {
		fmt.Printf("()nil")
		return
	}
	fmt.Printf(" , ") // 左右分割

	// show父节点
	if n.Parent == nil {
		fmt.Printf("(nil|%d)", len(n.Key))
	} else {
		fmt.Printf("(%d|%d)", n.Parent.Key[0], len(n.Key))
	}

	// 循环show本节点所有Key+右腿
	for i := 0; i < len(n.Key); i++ {
		// key
		fmt.Printf("%d", n.Key[i])
		// 腿 or 叶子
		if len(n.Leaf) == 0 { // 下级是分支
			if n.Child[i] == nil {
				fmt.Printf("\\nil|")
			} else {
				fmt.Printf("\\%d|", n.Child[i].Key[0])
			}
		} else { // 下级是叶子
			if n.Leaf[i] == nil {
				fmt.Printf("\\nil|")
			} else {
				fmt.Printf("\\%d|", n.Leaf[i].Key)
			}
		}

	}
}

// ShowOneLeaf 展示单个节点
// @Author  https://github.com/BrotherSam66/
func ShowOneLeaf(n *bplustreemodels.BPTreeLeaf) {
	if n == nil {
		fmt.Printf("()nil")
		return
	}
	fmt.Printf(" , ") // 左右分割

	// show父节点
	if n.Parent == nil {
		fmt.Printf("(nil)")
	} else {
		fmt.Printf("(%d)", n.Parent.Key[0])
	}

	// key
	fmt.Printf("%d", n.Key)
	// 载荷
	if n.Payload != "" {
		fmt.Printf("-%s", n.Payload)
	}
	// 右兄弟
	if n.RightBrother == nil {
		fmt.Printf("->nil")
	} else {
		fmt.Printf("->%d", n.RightBrother.Key)
	}

}
