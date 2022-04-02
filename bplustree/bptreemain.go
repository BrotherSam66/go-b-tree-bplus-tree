// Package bplustree
// @Title B树工具包
// @Description  总纲
// @Author  https://github.com/BrotherSam66/
// @Update
package bplustree

import (
	"fmt"
	"go-b-tree-bplus-tree/bplustree/bplustreeglobal"
	"os"
	"strings"
)

// BPlusTreeDemo B+树 演示
// @author https://github.com/BrotherSam66/
func BPlusTreeDemo() {
	for {
		var command string
		fmt.Println("I Insert插入数据")
		fmt.Println("S Show完整的树")
		fmt.Println("F Find查找数据")
		fmt.Println("D Delete删除数据")
		//fmt.Println("Q PreOrder 前序遍历")
		//fmt.Println("Z InfixOrder 中序遍历")
		//fmt.Println("H PostOrder 后序遍历")
		fmt.Println("E Exit退出")
		fmt.Println("请输入指令，按回车键：")
		_, _ = fmt.Scanln(&command)
		command = strings.ToUpper(command)

		switch command {
		case "I":
			Inputs()
		case "S":
			ShowTree(bplustreeglobal.Root)
		case "F":
			var key int
			fmt.Println("请输入KEY，按回车键(0退出)：")
			_, _ = fmt.Scanln(&key)
			node, isTarget, err := Search(key)
			if err != nil || !isTarget {
				fmt.Println("没找到or查找错误，error == ", err)
			} else {
				ShowOneNode(node)
				fmt.Println()
			}
		case "D":
			//Deletes()
		case "E":
			os.Exit(0)
		default:
			fmt.Println("输入错误")

		}

	}
}
