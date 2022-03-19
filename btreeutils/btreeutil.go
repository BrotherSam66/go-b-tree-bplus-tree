package rbtutils

import (
	"errors"
	"fmt"
	"go-red-black-tree/global"
	"go-red-black-tree/rbtmodels"
	"math"
	"math/rand"
)

// RBTInputs 红黑树连续插入节点
func RBTInputs() {
	//RBTCreat()
	for {
		var key int
		fmt.Println("请输入KEY，按回车键(空按回车随机,10XX填充1~XX，-1退出)：")
		_, _ = fmt.Scanln(&key)

		if key == -1 {
			return
		}
		if key == 0 {
			key = rand.Intn(global.MaxKey)
			fmt.Println(key)
		}
		if key > 1000 {
			if key > 1046 {
				fmt.Println("最大1046，否则溢出....")
				continue
			}
			endKey := key - 1000
			for i := 1; i <= endKey; i++ {
				Insert(i, "")
			}

			ShowTreeColor(global.Root)
			continue
		}
		if key > 99 || key < 1 {
			fmt.Println("必须是0~~99")
			continue
		}
		Insert(key, "")
		ShowTreeColor(global.Root)
	}
}

// RBTDeletes 红黑树连续删除节点
func RBTDeletes() {

	for {
		var key int
		fmt.Println("请输入KEY，按回车键(-1退出)：")
		_, _ = fmt.Scanln(&key)

		if key == -1 {
			return
		}
		if key == 0 {
			key = rand.Intn(global.MaxKey)
			fmt.Println(key)
		}
		if key > 99 || key < 1 {
			fmt.Println("必须是0~~99")
			continue
		}
		node, err := Find(key)
		if err != nil {
			fmt.Println("查找错误，error == ", err)
			continue
		}
		Delete(node)
		ShowTreeColor(global.Root)
	}
}

// ShowTreeColor 彩色逐层显示这个树，回头废止。感觉需要递归、队列
func ShowTreeColor(r *rbtmodels.RBTNode) {
	var data [10][1000]*rbtmodels.RBTNode     // 数据。数据可能是nil。最多10层，每层最多1000数据
	var dataTemp [10][1000]*rbtmodels.RBTNode // 临时数据，回头某列全nil就不显示。数据可能是nil。最多10层，每层最多1000数据
	totalLevel := 1                           // 总层数
	//nowLevel := 0               // 当前层数
	//nnn := Name
	//nowColumn := 0 // 当前列
	if r == nil {
		fmt.Println("这个树/分支是空的")
	}
	data[0][0] = r // 来的最高位指针

	// 循环，把每个节点指针放入对应层的队列，每个上级节点占死2个下级，没有这个子节点就空着
	for i := 1; i < len(data); i++ { // 循环每一层
		countNotNil := 0                                      // 本层非nil个数，==0 表示上一层是最后一层
		for j := 0; j < int(math.Pow(2, float64(i-1))); j++ { // 上层应有的元素数量，遍历，本层翻倍
			if data[i-1][j] != nil {
				if data[i-1][j].Left != nil {
					countNotNil++
					data[i][j*2] = data[i-1][j].Left // 上层左儿子，放入
				}
				if data[i-1][j].Right != nil {
					countNotNil++
					data[i][j*2+1] = data[i-1][j].Right // 上层右儿子，放入
				}
			}
		}
		if countNotNil == 0 { // 本层无元素，中断，退出
			break
		}
		totalLevel++ // 总层数
	}

	// 二次循环，数据导入dataTemp，一组占3列，父亲骑在2个儿子中间
	for i := 0; i < totalLevel; i++ { // 循环每一层
		blankMiddleLen := int(math.Pow(2, float64(totalLevel-i))) - 1 // 中间空
		blankLeftLen := (blankMiddleLen - 1) / 2

		for j := 0; j < int(math.Pow(2, float64(i))*1.5); j++ { // 上层应有的元素数量，遍历，本层翻倍
			dataTemp[i][blankLeftLen+j*(blankMiddleLen+1)+1] = data[i][j]
			//fmt.Printf("i== %d ,j== %d ,totalLevel== %d ,blankMiddleLen== %d ,blankLeftLen== %d \n", i, j, totalLevel, blankMiddleLen, blankLeftLen)
		}
	}

	// 三次循环，把每层数据展示出来，
	for i := 0; i < totalLevel; i++ { // 循环每一层
		for j := 0; j < int(math.Pow(2, float64(totalLevel))); j++ { // 应有的元素数量

			// 查本列是否全nil
			isAllNil := true // 本列全nil
			for k := 0; k < totalLevel; k++ {
				if dataTemp[k][j] != nil {
					isAllNil = false
				}
			}
			if !isAllNil { // 只有本列不是全nil，才show一下
				ShowOneNodeColorNew(dataTemp[i][j])
			}
			//ShowOneNodeColorNew(dataTemp[i][j], totalLevel, i, j)
		}

		fmt.Println()
	}

	//// 三次循环，把每层数据展示出来，
	//for i := 1; i < totalLevel+1; i++ { // 循环每一层
	//	//nowColumn = 0                                         // 当前列
	//	for j := 0; j < int(math.Pow(2, float64(i-1))); j++ { // 上层应有的元素数量，遍历，本层翻倍
	//		ShowOneNodeColor(data[i-1][j], totalLevel, i, j)
	//	}
	//	fmt.Println()
	//}
}

// ShowOneNodeColorNew 彩色展示单个节点
func ShowOneNodeColorNew(n *rbtmodels.RBTNode) {

	// blank=blankLeft+n*(位长global.KeyLen+1)
	blankNil := "" // 空节点，也占位置
	for k := 0; k < global.KeyLen; k++ {
		blankNil = blankNil + " "
	}

	if n == nil {
		fmt.Printf("%s", blankNil)
	} else {
		//其中0x1B是标记，[开始定义颜色，1代表高亮，40代表黑色背景，32代表绿色前景，0代表恢复默认颜色。
		red := 31
		black := 0
		if n.IsRed {
			fmt.Printf("%c[1;0;%dm%02d%c[0m", 0x1B, red, n.Key, 0x1B)
		} else {
			fmt.Printf("%c[1;0;%dm%02d%c[0m", 0x1B, black, n.Key, 0x1B)
		}
	}

}

// ShowOneNodeColor 彩色展示单个节点
func ShowOneNodeColor(n *rbtmodels.RBTNode, totalLevel, i, j int) {
	// blank=blankLeft+n*(位长global.KeyLen+1)
	blankNil := "" // 空节点，也占位置
	for k := 0; k < global.KeyLen+1; k++ {
		blankNil = blankNil + " "
	}
	blankLeftHead := ""                                                             //  最左边                                                  // 总体左边空
	blankMiddleLen := int(math.Pow(2, float64(totalLevel-i)))*(global.KeyLen+1) - 3 // 中间空
	blankLeftLen := blankMiddleLen / 2
	blankLeft := blankLeftHead // 最左空的
	for k := 0; k < blankLeftLen; k++ {
		blankLeft = blankLeft + " "
	}
	blankMiddle := "" // 中间空的
	for k := 0; k < blankMiddleLen; k++ {
		blankMiddle = blankMiddle + " "
	}

	if j == 0 { // 本列第一个
		fmt.Printf("%s", blankLeft)
	} else {
		fmt.Printf("%s", blankMiddle)
	}

	if n == nil {
		fmt.Printf("%s", blankNil)
	} else {
		//其中0x1B是标记，[开始定义颜色，1代表高亮，40代表黑色背景，32代表绿色前景，0代表恢复默认颜色。
		red := 31
		black := 0
		if n.IsRed {
			fmt.Printf("%c[1;0;%dm%02d%c[0m", 0x1B, red, n.Key, 0x1B)
		} else {
			fmt.Printf("%c[1;0;%dm%02d%c[0m", 0x1B, black, n.Key, 0x1B)
		}
	}

}

// ShowTree 逐层显示这个树，回头废止。感觉需要递归、队列
func ShowTree(r *rbtmodels.RBTNode) {
	var data [10][1000]*rbtmodels.RBTNode // 数据。数据可能是nil。最多10层，每层最多1000数据
	//totalLevel := 0             // 总层数
	//nowLevel := 0               // 当前层数
	//nnn := Name
	nowColumn := 0 // 当前列
	fmt.Printf("\n展示树：[左子]本(父)[右子]")
	if r == nil {
		fmt.Println("这个树/分支是空的")
	}
	data[0][0] = r // 来的最高位指针

	for i := 1; i < len(data); i++ { // 循环每一层
		fmt.Println("") // 先来一个换行
		nowColumn = 0   // 当前列
		for j := 0; j < len(data[0]); j++ {
			if data[i-1][j] == nil { // 本行遍历结束
				break
			}

			ShowOneNode(data[i-1][j]) // 显示遍历到的上一行的这个节点。显示没换行

			if data[i-1][j].Left != nil { // 如果有，在下一行填入左节点
				data[i][nowColumn] = data[i-1][j].Left
				nowColumn++
			}
			if data[i-1][j].Right != nil { // 如果有，在下一行填入右节点
				data[i][nowColumn] = data[i-1][j].Right
				nowColumn++
			}
		}
	}

}

// ShowOneNode 展示单个节点
func ShowOneNode(n *rbtmodels.RBTNode) {
	if n == nil {
		fmt.Println("[]nil[]")
		return
	}
	fmt.Printf(" , ")  // 左右分割
	if n.Left == nil { // 左儿子KEY
		fmt.Printf("[ ]")
	} else {
		fmt.Printf("[%d]", n.Left.Key)
	}

	fmt.Printf("%d", n.Key) // 本节点KEY
	if n.IsRed == true {    // 本节点是红色
		fmt.Printf("R")
	} else { // 黑色
		fmt.Printf("B")
	}

	if n.Parent == nil { // 父节点KEY
		fmt.Printf("( )")
	} else {
		fmt.Printf("(%d)", n.Parent.Key)
	}

	if n.Right == nil { // 右节点KEY
		fmt.Printf("[ ]")
	} else {
		fmt.Printf("[%d]", n.Right.Key)
	}
}

// Find 查找节点
// @key 键值
// @n 找到的节点指针
func Find(key int) (ret *rbtmodels.RBTNode, err error) {
	if global.Root == nil {
		fmt.Println("这个树/分支是空的")
		return nil, errors.New("没找到！")
	}
	tempNode := global.Root
	for { // 递归循环
		if tempNode.Key == key { // 等于，找到了
			return tempNode, nil
		} else if tempNode.Key > key { // 大于，向左
			if tempNode.Left == nil { // 到nil了，返没找到
				return nil, errors.New("没找到！")
			}
			tempNode = tempNode.Left // 向左下递归
		} else { // 小于，向右
			if tempNode.Right == nil { // 到nil了，返没找到
				return nil, errors.New("没找到！")
			}
			tempNode = tempNode.Right // 向右下递归
		}
	}
}
