// Package btreeutils
// @Title B树工具包
// @Description  显示节点
// @Author  https://github.com/BrotherSam66/
// @Update
package btreeutils

import (
	"fmt"
	"go-b-tree-bplus-tree/btreemodels"
)

// ShowTree 逐层显示这个树，
// @b 展示树的开始位置（可以不从root开始展示）
// @Author  https://github.com/BrotherSam66/
func ShowTree(b *btreemodels.BTreeNode) {
	if b == nil { // 先判断，后定义变量，免得浪费内存
		fmt.Println("这个树/分支是空的")
		return
	}
	//fmt.Printf("%#v\n", b)

	// ShowTemp 数据。数据可能是nil。最多10层，每层最多1000 数据
	var showTemp [10][1000]*btreemodels.BTreeNode
	//totalLevel := 0             // 总层数
	//nowLevel := 0               // 当前层数
	//nnn := Name
	nowColumn := 0 // 当前列
	fmt.Printf("\n展示树：(父|KEY数)左腿/KEY-PAYLOAD\\右腿|KEY-PAYLOAD\\右腿|...")
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

			// 在本行填写上一行所有节点的所有儿子
			showTemp[i][nowColumn] = showTemp[i-1][j].Child[0] // 最左边的左腿
			nowColumn++
			for k := 0; k < showTemp[i-1][j].KeyNum; k++ { // 其他的所有右腿
				showTemp[i][nowColumn] = showTemp[i-1][j].Child[k+1]
				nowColumn++
			}
		}
	}

}

// ShowOneNode 展示单个节点
// @Author  https://github.com/BrotherSam66/
func ShowOneNode(n *btreemodels.BTreeNode) {
	if n == nil {
		fmt.Printf("()nil")
		return
	}
	fmt.Printf(" , ") // 左右分割

	// show父节点
	if n.Parent == nil {
		fmt.Printf("(nil|%d)", n.KeyNum)
	} else {
		fmt.Printf("(%d|%d)", n.Parent.Key[0], n.KeyNum)
	}

	// 先show最左的左腿
	if n.Child[0] == nil {
		fmt.Printf("nil/")
	} else {
		fmt.Printf("%d/", n.Child[0].Key[0])
	}
	// 循环show本节点所有Key+右腿
	for i := 0; i < n.KeyNum; i++ {
		// key
		fmt.Printf("%d", n.Key[i])
		// 载荷
		if n.Payload[i] != "" {
			fmt.Printf("-%s", n.Payload[i])
		}
		// 右腿
		if n.Child[i+1] == nil {
			fmt.Printf("\\nil|")
		} else {
			fmt.Printf("\\%d|", n.Child[i+1].Key[0])
		}
	}
}
