// Package bplustree
// @Title B树工具包
// @Description  和插入节点有关的操作
// @Author  https://github.com/BrotherSam66/
// @Update
package bplustree

import (
	"errors"
	"fmt"
	"go-b-tree-bplus-tree/bplustree/bplustreeglobal"
	"go-b-tree-bplus-tree/bplustree/bplustreemodels"
	"math/rand"
)

// Inputs 连续插入节点
// @author https://github.com/BrotherSam66/
func Inputs() {

	for {
		var key int
		fmt.Println("请输入KEY，按回车键(空按回车随机,10XX填充1~XX，-1退出)：")
		_, _ = fmt.Scanln(&key)

		if key == -1 {
			return
		}
		if key == 0 {
			key = rand.Intn(bplustreeglobal.MaxKey)
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

			ShowTree(bplustreeglobal.Root)
			continue
		}
		if key > 99 || key < 1 {
			fmt.Println("必须是0~~99")
			continue
		}
		Insert(key, "")
		ShowTree(bplustreeglobal.Root)
	}
}

// Insert 加入节点
// @key 插入的键值
// @payload 插入的载荷值
// @author https://github.com/BrotherSam66/
func Insert(key int, payload string) {
	if payload == "" {
		payload = fmt.Sprintf("%d", key)
	}
	/*
		B+树插入方法（我自己琢磨的）
		[1]树root=nil。==》创建叶子，==》创建分支=root，==》上下链接好。
		[2]完美找到叶子，==》修改叶子payload
		[3]找到应该插入的末级分支。==》创建叶子==》末级分支插入key。
		[3.1]若被插入关键字所在的结点，其含有关键字数目等于阶数 M，则需要将该结点分裂为两个结点，
			左结点包含⌈Min⌉。将⌈Min⌉的关键字上移至其双亲结点。
			假设其双亲结点中包含的关键字个数小于 M，则插入操作完成。。
			==》否则用父亲节点递归[3.1]，直到root
		如果插入的关键字比当前结点中的最大值还大，破坏了B+树中从根结点到当前结点的所有索引值，此时需要及时修正后，再做其他操作。
			例如，在图 1 的 B+树种插入关键字 100，由于其值比 97 还大，插入之后，
			从根结点到该结点经过的所有结点中的所有值都要由 97 改为 100。改完之后再做分裂操作。
	*/

	// [1]树root=nil。==》创建叶子，==》创建分支=root，==》上下链接好。
	if bplustreeglobal.Root == nil { // 原树为空树，新加入的转为根
		// 创建叶子
		newLeaf := bplustreemodels.NewBPTreeLeaf(bplustreeglobal.Root, nil, key, payload)
		// Sqt 指向
		bplustreeglobal.Sqt = newLeaf
		// 根 创建成唯一分支
		bplustreeglobal.Root = bplustreemodels.NewBPTreeNode(nil, key, newLeaf)
		return
	}

	// 从root开始查找附加的位置
	tempNode, isTarget, err := Search(key)
	if err != nil {
		fmt.Println("没找到or查找错误，error == ", err)
		return
	}

	// [2]完美找到叶子，==》修改叶子payload
	if isTarget { // 拟插入的key存在，替换payload就好
		// 寻找替换的位置
		_, realPoint := FindKeyPoint(&tempNode.Key, key)
		if realPoint > -1 { // 准确命中，只可能是新创建节点情形
			tempNode.Leaf[realPoint].Payload = payload // 拟插入的key存在，替换下级叶子payload就好
		} else {
			fmt.Println("不正常，input准确查到叶子，不应该走到这里")
		}
		return
	}

	// [3]找到应该插入的末级分支。==》创建叶子==》末级分支插入key。
	// 创建叶子
	newLeaf := bplustreemodels.NewBPTreeLeaf(tempNode, nil, key, payload)
	// 处理 Sqt 数据头
	if key < bplustreeglobal.Sqt.Key { // 新入Key 全树最小。换Sqt
		bplustreeglobal.Sqt = newLeaf
	}
	// 求末级分支插key点位
	insertPoint, realPoint := FindKeyPoint(&tempNode.Key, key)
	if realPoint > -1 {
		fmt.Println("realPoint > -1不正常，input准备插入末端分支，不应该走到这里")
		return
	}
	if insertPoint < 0 { // todo 表示在尾巴追加
		fmt.Println("insertPoint < 0不正常，input准备插入末端分支，不应该走到这里")
		return
	}
	// 末级分支插入key+下级叶子
	InsertOneLeaf(tempNode, newLeaf, insertPoint)

	return
}

// InsertOneLeaf 插入一个Key，满了也插，溢出在Tail里
func InsertOneLeaf(n *bplustreemodels.BPTreeNode, insertLeaf *bplustreemodels.BPTreeLeaf, insertPoint int) (err error) {
	if n.Key == nil {
		n.Key = []int{}
	}
	n.Key = append(n.Key, 0)     // 只是扩容
	n.Leaf = append(n.Leaf, nil) // 只是扩容
	for i := len(n.Key) - 1; i > insertPoint; i-- {
		n.Key[i] = n.Key[i-1]
		n.Leaf[i] = n.Leaf[i-1]
	}
	n.Key[insertPoint] = insertLeaf.Key // 插入Key
	n.Leaf[insertPoint] = insertLeaf    // 上级下指叶子
	insertLeaf.Parent = n               // 叶子上指

	// 需要分裂？
	if len(n.Key) > bplustreeglobal.M {
		insertNode, _ := SplitTo2Node(n) // 分裂，返回左半扇
		_ = InsertOneNode(insertNode)    // 左半扇最大Key插入父亲
	}
	return

}

// InsertOneNode 左半扇最大Key插入父亲
func InsertOneNode(insertNode *bplustreemodels.BPTreeNode) (err error) {
	if insertNode.Parent == nil { // 说明是root在分裂
		// 新建分支，Key是我的老大+原root的老大，下级是我和原root
		newRoot := &bplustreemodels.BPTreeNode{}
		oldRoot := bplustreeglobal.Root
		newRoot.Key = []int{insertNode.Key[len(insertNode.Key)-1], oldRoot.Key[len(oldRoot.Key)-1]}
		newRoot.Child = []*bplustreemodels.BPTreeNode{insertNode, oldRoot}
		// root指向新分支
		bplustreeglobal.Root = newRoot
		// 把我和原root 的父级指向新分支
		insertNode.Parent = bplustreeglobal.Root
		oldRoot.Parent = bplustreeglobal.Root
		return
	}

	parent := insertNode.Parent //  父亲
	// 找插入位置；
	insertPoint, realPoint := FindKeyPoint(&parent.Key, insertNode.Key[len(insertNode.Key)-1])
	if realPoint > -1 {
		err = errors.New("realPoint > -1不正常，input准备插入末端分支，不应该走到这里")
		fmt.Println(err.Error())
		return
	}
	// 插入的是分支（这个分支，只是带着左半扇分支node，需要把最大KEY上移就好，）

	parent.Key = append(parent.Key, 0)       // 只是扩容
	parent.Child = append(parent.Child, nil) // 只是扩容
	for i := len(parent.Key) - 1; i > insertPoint; i-- {
		parent.Key[i] = parent.Key[i-1]
		parent.Child[i] = parent.Child[i-1]
	}
	parent.Key[insertPoint] = insertNode.Key[len(insertNode.Key)-1] // key=下级的最大KEY
	parent.Child[insertPoint] = insertNode                          // 上级下指
	// 下级 insertNode 上指，在分裂里就做好了

	// 查看是否需要递归 // 需要分裂？
	if len(parent.Key) > bplustreeglobal.M {
		newLeftNode, _ := SplitTo2Node(parent) // 分裂，返回左半扇
		_ = InsertOneNode(newLeftNode)         // 用左半扇向上插入，进入递归
	}
	return

}

// SplitTo2Node >M，满员了，左右分裂，
func SplitTo2Node(n *bplustreemodels.BPTreeNode) (retNode *bplustreemodels.BPTreeNode, err error) {
	// 左边包含Min位置的元素，然后用这个元素向上插入
	if n.Parent == nil { // 说明是root在分裂
		// todo ？？？？
	}

	// 左边是新创建节点
	newLeftNode := &bplustreemodels.BPTreeNode{}
	newLeftNode.Parent = n.Parent // 父亲的指向
	newLeftNode.Key = n.Key[:bplustreeglobal.Min]

	// 右边搬家过来Min
	if len(n.Leaf) > 0 { // 是末级级分支，带的是叶子
		newLeftNode.Leaf = n.Leaf[:bplustreeglobal.Min] // 一组指向叶子的
		for i := 0; i < bplustreeglobal.Min+1; i++ {
			n.Leaf[i].Parent = newLeftNode // 一组叶子的父级上联
		}
	} else { // 非末级级分支，带的是分支
		newLeftNode.Child = n.Child[:bplustreeglobal.Min] // 一组指向下级分支的
		for i := 0; i < bplustreeglobal.Min+1; i++ {
			n.Child[i].Parent = newLeftNode // 一组下级分支的父级上联
		}
	}

	// 旧节点裁掉左边Min
	n.Key = n.Key[bplustreeglobal.Min:]
	if len(n.Leaf) > 0 { // 是末级级分支，带的是叶子
		n.Leaf = n.Leaf[bplustreeglobal.Min:] // 一组指向叶子的
	} else { // 非末级级分支，带的是分支
		n.Child = n.Child[bplustreeglobal.Min:] // 一组指向下级分支的
	}

	retNode = newLeftNode
	return
}
