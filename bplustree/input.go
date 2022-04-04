// Package bplustree
// @Title B树工具包
// @Description  和插入节点有关的操作
// @Author  https://github.com/BrotherSam66/
// @Update
package bplustree

import (
	"errors"
	"fmt"
	"go-b-tree-bplus-tree/bplustree/bplustreeconst"
	"go-b-tree-bplus-tree/bplustree/bplustreeglobal"
	"go-b-tree-bplus-tree/bplustree/bplustreemodels"
	"math/rand"
)

// Inputs 连续插入节点
// @author https://github.com/BrotherSam66/
func Inputs() {

	for {
		var key int
		fmt.Println("\n请输入KEY，按回车键(空按回车随机,10XX填充1~XX,20XX填充XX~1，-1退出)：")
		_, _ = fmt.Scanln(&key)

		if key == -1 {
			return
		}
		if key == 0 {
			key = rand.Intn(bplustreeconst.MaxKey)
			fmt.Println(key)
		}
		if key > 2000 {
			//if key > 2046 {
			//	fmt.Println("最大2046，否则溢出....")
			//	continue
			//}
			endKey := key - 2000
			for i := endKey; i > 0; i-- {
				Insert(i, "")
			}

			ShowTree(bplustreeglobal.Root)
			continue
		}
		if key > 1000 {
			//if key > 1046 {
			//	fmt.Println("最大1046，否则溢出....")
			//	continue
			//}
			endKey := key - 1000
			for i := 1; i <= endKey; i++ {
				Insert(i, "")
			}

			ShowTree(bplustreeglobal.Root)
			continue
		}
		//if key > 99 || key < 1 {
		//	fmt.Println("必须是0~~99")
		//	continue
		//}
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
		payload = fmt.Sprintf("p%d", key)
	}
	/*
		B+树插入方法（我自己琢磨的）
		[1]树root=nil。==》创建节点=root，==》结束
		[2]完美找到节点(一定是末梢/叶子节点)(节点包含本key)，==》修改家电payload==》结束
		[3]找到应该插入的末梢/叶子节点。==》末梢插入key。==》向[3.1]递归
		[3.1]若被插入关键字的结点(可能是末梢or分支)，key数量超标阶数 M，==》将该结点分裂为两个结点，
			左结点包含⌈Min⌉。将⌈Min⌉的关键字上移至其爸爸结点。
			假设爸爸结点中包含的关键字个数 <= M。==》结束
			关键字个数 > M。==》用爸爸节点向[3.1]递归，直到root

			如果插入的关键字比当前结点中的最大值还大，破坏了B+树中从根结点到当前结点的所有索引值，此时需要及时修正后，再做其他操作。
			例如，在图 1 的 B+树种插入关键字 100，由于其值比 97 还大，插入之后，
			从根结点到该结点经过的所有结点中的所有值都要由 97 改为 100。改完之后再做分裂操作。
	*/

	// [1]树root=nil。==》创建节点=root，==》结束
	if bplustreeglobal.Root == nil { // 原树为空树，新加入的转为根
		// 根 创建成唯一节点
		bplustreeglobal.Root = bplustreemodels.NewBPTreeLeaf(nil, key, payload)
		// Sqt 指向
		bplustreeglobal.Sqt = bplustreeglobal.Root
		return
	}

	// 从root开始查找附加的位置
	tempNode, isTarget, err := Search(key)
	if err != nil {
		fmt.Println("Search 没找到or查找错误，error == ", err)
		return
	}

	// [2]完美找到节点(一定是末梢/叶子节点)(节点包含本key)，==》修改家电payload==》结束
	if isTarget { // 完美找到
		// 寻找替换的位置
		_, realPosition := FindKeyPosition(&tempNode.Key, key)
		if realPosition > -1 { // 准确命中，替换
			tempNode.Payload[realPosition] = payload
		} else {
			fmt.Println("不正常，input FindKeyPosition 准确查到叶子，不应该走到这里", key)
		}
		return
	}

	// [3]找到应该插入的末梢/叶子节点。==》末梢插入key。==》向[3.1]递归
	// 求末级分支插key点位
	insertPosition, realPosition := FindKeyPosition(&tempNode.Key, key)
	if realPosition > -1 {
		fmt.Println("input realPosition > -1不正常，input准备插入末端分支，不应该走到这里", key)
		return
	}
	if insertPosition < 0 { // 在尾巴追加，会回来一个刚好越界切片位置的值
		fmt.Println("insertPosition < 0不正常，input准备插入末端分支，不应该走到这里")
		return
	}

	// 末级分支插入key+payload
	InsertKeyInLeaf(tempNode, key, payload, insertPosition)
	if len(tempNode.Key) != len(tempNode.Payload) {
		fmt.Println("Insert() 不正常，len(tempNode.Key)!=  len(tempNode.Payload) ，不应该走到这里", key)
	}
	if len(tempNode.Key) > bplustreeconst.M { // [3.1]若被key数量超标阶数 M，==》分裂==>用新左节点向[3.1]递归
		newLeftNode, _ := SplitTo2Node(tempNode) // 分裂，返回左半扇
		if newLeftNode == nil {
			fmt.Println("Insert() 不正常，newLeftNode==nil ，不应该走到这里", key)
		}
		_ = InsertOneNode(newLeftNode) // 用左半扇向上插入，进入递归
	}

	return
}

// InsertKeyInLeaf 末级分支插入key+payload，满了也插，
// @n 被插入节点
// @key 拟插入的键值
// @payload 拟插入的载荷
// @insertPosition 拟插入的位置
// @author https://github.com/BrotherSam66/
func InsertKeyInLeaf(n *bplustreemodels.BPTreeNode, key int, payload string, insertPosition int) (err error) {
	n.Key = append(n.Key, 0)          // 只是扩容
	n.Payload = append(n.Payload, "") // 只是扩容
	for i := len(n.Key) - 1; i > insertPosition; i-- {
		n.Key[i] = n.Key[i-1]
		n.Payload[i] = n.Payload[i-1]
	}
	n.Key[insertPosition] = key         // 插入Key
	n.Payload[insertPosition] = payload // 插入payload

	// 如果插入的key是最大的，飘在最右，需要向上修正每级爸爸的最右key
	if key == n.Key[len(n.Key)-1] {
		if n.Parent != nil { // 到root了，修正结束
			tempNode := n.Parent
			for tempNode.Key[len(tempNode.Key)-1] < key {
				tempNode.Key[len(tempNode.Key)-1] = key
				if tempNode.Parent == nil { // 到root了，修正结束
					break
				} else { // 向爸爸递归
					tempNode = tempNode.Parent
				}
			}
		}
	}

	return

}

// InsertOneNode 左半扇最大Key插入父亲
// @insertNode 拟插入的节点
// @author https://github.com/BrotherSam66/
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
	insertPosition, realPosition := FindKeyPosition(&parent.Key, insertNode.Key[len(insertNode.Key)-1])
	if realPosition > -1 {
		err = errors.New("realPosition > -1不正常，input准备插入末端分支，不应该走到这里")
		fmt.Println(err.Error())
		return
	}

	// 插入的是分支（这个分支，只是带着左半扇分支node，需要把最大KEY上移就好，）

	// 被插入
	parent.Key = append(parent.Key, 0)       // 只是扩容
	parent.Child = append(parent.Child, nil) // 只是扩容
	// parent 不可能有 payload载荷的
	for i := len(parent.Key) - 1; i > insertPosition; i-- { // 需要插入的位置，向右挤
		parent.Key[i] = parent.Key[i-1]
		parent.Child[i] = parent.Child[i-1]
	}
	parent.Key[insertPosition] = insertNode.Key[len(insertNode.Key)-1] // key=下级的最大KEY
	parent.Child[insertPosition] = insertNode                          // 上级下指
	// 下级 insertNode 上指，在分裂里就做好了

	// 查看是否需要递归 // 需要分裂？
	if len(parent.Key) > bplustreeconst.M {
		newLeftNode, _ := SplitTo2Node(parent) // 分裂，返回左半扇
		_ = InsertOneNode(newLeftNode)         // 用左半扇向上插入，进入递归
	}
	return

}

// SplitTo2Node >M，满员了，左右分裂，左边是新创建节点.注意 Sqt 数据头
// @n 拟分裂的节点（兼做分裂后的右半部分）
// @retNode 向左分裂出来的节点
// @author https://github.com/BrotherSam66/
func SplitTo2Node(n *bplustreemodels.BPTreeNode) (retNode *bplustreemodels.BPTreeNode, err error) {

	// 左边是新创建节点
	newLeftNode := &bplustreemodels.BPTreeNode{}
	if len(n.Child) > 0 && n.Child[0] != nil { // 被分裂的是分支，如果是末梢/叶子，就没有child
		newLeftNode = bplustreemodels.NewBPTreeNode(bplustreeconst.Min) // 捎带完成 Child+Key 扩容，为copy准备
		copy(newLeftNode.Key, n.Key[:bplustreeconst.Min])               // key 深拷贝
		copy(newLeftNode.Child, n.Child[:bplustreeconst.Min])           // 左组指向下级分支的Child  深拷贝
		for i := 0; i < bplustreeconst.Min; i++ {
			n.Child[i].Parent = newLeftNode // 左组下级分支的父级上联
		}
	} else { // 叶子的。叶子才有payload
		newLeftNode.Key = make([]int, bplustreeconst.Min)         // Key 扩容，为copy准备
		newLeftNode.Payload = make([]string, bplustreeconst.Min)  // Payload 扩容，为copy准备
		copy(newLeftNode.Key, n.Key[:bplustreeconst.Min])         // key 深拷贝
		copy(newLeftNode.Payload, n.Payload[:bplustreeconst.Min]) // payload 深拷贝
	}

	// 指向爸爸的
	newLeftNode.Parent = n.Parent

	// 旧节点裁掉左边Min
	n.Key = n.Key[bplustreeconst.Min:] // key

	if len(n.Child) > 0 && n.Child[0] != nil { // 被分裂的如果是末梢/叶子，就没有child
		n.Child = n.Child[bplustreeconst.Min:] // 右组指向下级分支的
	} else { // 叶子哦
		n.Payload = n.Payload[bplustreeconst.Min:] // key
	}

	retNode = newLeftNode // 返回的节点

	if len(n.Child) == 0 || n.Child[0] == nil { // 叶子分裂，才需要重组 兄弟指向
		// 横向 兄弟 指向：新左--》新右
		newLeftNode.RightBrother = n

		// 处理 Sqt 数据头（先做，是有道理的）
		if bplustreeglobal.Sqt == n { // Sqt 数据头 原先指向被分裂节点，说明n是整棵树最左，需要重新指向
			bplustreeglobal.Sqt = newLeftNode
			return
		}

		// 横向 兄弟 指向：左左--》新左
		// 找 左左，step1：n 向上递归，①在爸爸中的排位非最左，找到；（不会找不到的，因为Sqt 数据头处理过了，新左不会是全树最左）
		// 到这里，n 肯定有爸爸
		tempNode := n
		for tempNode.Parent != nil { // 确保到root 停止
			if tempNode.Key[0] < newLeftNode.Key[bplustreeconst.Min-1] {
				break
			}
			tempNode = tempNode.Parent
		}
		// 找 左左，step2：找newLeftNode在某级祖宗 tempNode.Key 的位置，要的是这个位置的 左一位
		var i int
		for i = len(tempNode.Key) - 1; i >= 0; i-- {
			if tempNode.Key[i] < newLeftNode.Key[bplustreeconst.Min-1] {
				break // 这个i就是左一位
			}
		}
		if i < 0 {
			err = errors.New("找 左左，step2：找newLeftNode在某级祖宗 tempNode.Key 的位置,左越界，不应该走到这里")
			fmt.Println(err.Error())
			//retNode = newLeftNode
			return
		}
		// 找 左左，step3：排位左兄弟，一溜向右下，找到末梢/叶子，就是
		tempNode = tempNode.Child[i]
		for len(tempNode.Child) > 0 && tempNode.Child[0] != nil { // 还不到末端/叶子，就向右下
			tempNode = tempNode.Child[len(tempNode.Child)-1]
		}

		// 终于找到的左左，把新建的作为右兄弟
		tempNode.RightBrother = newLeftNode
	}

	return
}
