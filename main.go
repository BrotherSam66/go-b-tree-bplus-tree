// Package main
// @Title B树入口
// @Description  主程序
// @Author  https://github.com/BrotherSam66/
// @Update
package main

import (
	"go-b-tree-bplus-tree/bplustree"
	"go-b-tree-bplus-tree/btree"

	"math/rand"
	"time"
)

func main() {
	//a := []int{0, 1, 2, 3, 4, 5, 6, 7}
	//a = append(a[:3], a[3+1:]...)
	//fmt.Println("a == ", a) // [0 1 2 4 5 6 7]
	rand.Seed(time.Now().Unix())
	bplustree.BPlusTreeDemo()
	btree.BTreeDemo()
}
