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
	rand.Seed(time.Now().Unix())
	bplustree.BPlusTreeDemo()
	btree.BTreeDemo()
}
