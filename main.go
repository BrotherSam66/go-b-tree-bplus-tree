package main

import (
	"fmt"
	"go-red-black-tree/global"
	"go-red-black-tree/rbtutils"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {
	rand.Seed(time.Now().Unix())
	//global.Name = "我的名字"
	//rbtutils.RBTDemo()

	// 运算哈希表DEMO
	//bstutils.BinaryTreeDemo()
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
			rbtutils.RBTInputs()
		case "S":
			rbtutils.ShowTreeColor(global.Root)
		case "F":
			var key int
			fmt.Println("请输入KEY，按回车键(0退出)：")
			_, _ = fmt.Scanln(&key)
			node, err := rbtutils.Find(key)
			if err != nil {
				fmt.Println("查找错误，error == ", err)
			} else {
				rbtutils.ShowOneNode(node)
				fmt.Println()
			}
		case "D":
			rbtutils.RBTDeletes()
		case "E":
			os.Exit(0)
		default:
			fmt.Println("输入错误")

		}

	}
}
