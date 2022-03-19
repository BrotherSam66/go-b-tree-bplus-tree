package rbtutils

import (
	"fmt"
	"go-red-black-tree/bstmodels"
	"go-red-black-tree/global"
	"go-red-black-tree/rbtmodels"
)

// PreOrder 前序遍历：中左右 就是先访问根节点，再访问左节点，最后访问右节点，
func PreOrder(node *bstmodels.Hero) {
	if node != nil {
		//fmt.Printf("No:%d;Label:%s;Left:%v;Right:%v\n", node.No, node.Label, node.Left, node.Right)
		fmt.Println(node.No, node.Name, node.Left, node.Right)
		PreOrder(node.Left)
		PreOrder(node.Right)
	}
	return
}

// InfixOrder 中序遍历：左中右 所谓的中序遍历就是先访问左节点，再访问根节点，最后访问右节点，
func InfixOrder(node *bstmodels.Hero) {
	if node != nil {
		InfixOrder(node.Left)
		//fmt.Printf("No:%d;Label:%s;Left:%v;Right:%v\n", node.No, node.Label, node.Left, node.Right)
		fmt.Println(node.No, node.Name, node.Left, node.Right)
		InfixOrder(node.Right)
	}
	return
}

// PostOrder 后序遍历：左右中 所谓的后序遍历就是先访问左节点，再访问右节点，最后访问根节点。
func PostOrder(node *bstmodels.Hero) {
	if node != nil {
		PostOrder(node.Left)
		PostOrder(node.Right)
		//fmt.Printf("No:%d;Label:%s;Left:%v;Right:%v\n", node.No, node.Label, node.Left, node.Right)
		fmt.Println(node.No, node.Name, node.Left, node.Right)
	}
	return
}

// LevelOrder 层序遍历：按层，左右
// 弄一个指针切片，仿队列，①显示left，②left进队列，③显示right，④right进队列；取队列下一个指针；
func LevelOrder(node *bstmodels.Hero) {
	if node == nil {
		fmt.Println("这是个空树！")
		return
	}
	// 定义一些准全局便利性+函数
	nodeQueue := make([]bstmodels.Hero, 0, 100) // 切片仿队列
	queueHead := 0                              // 队列的头
	queueTail := 0                              // 队列的尾巴
	travel := func() {
		curNode := nodeQueue[queueHead] // 当前、刚取出来的节点
		queueHead++                     // 队列头修正
		fmt.Println(curNode)
		if curNode.Left != nil {
			nodeQueue = append(nodeQueue, *curNode.Left) // 压入队列
			queueTail++                                  // 队列尾巴修正
		}
		if curNode.Right != nil {
			nodeQueue = append(nodeQueue, *curNode.Right) // 压入队列
			queueTail++                                   // 队列尾巴修正
		}
	}

	// 开始程序
	nodeQueue = append(nodeQueue, *node) // 压入队列
	queueTail++                          // 队列尾巴修正
	for queueTail-queueHead > 0 {
		travel()
	}
}

// RBTCreat 红黑树创建
func RBTCreat() {
	// 定义
	global.Root = rbtmodels.NewRBTNode(false, 1, "1宋江", nil, nil, nil)
	DemoPush(global.Root, true, rbtmodels.NewRBTNode(false, 2, "2卢俊义", nil, nil, nil))
	DemoPush(global.Root, false, rbtmodels.NewRBTNode(true, 3, "3吴用", nil, nil, nil))
	DemoPush(global.Root.Left, true, rbtmodels.NewRBTNode(false, 4, "4公孙胜", nil, nil, nil))
	DemoPush(global.Root.Left, false, rbtmodels.NewRBTNode(true, 5, "5关胜", nil, nil, nil))
	DemoPush(global.Root.Right, true, rbtmodels.NewRBTNode(false, 6, "6林冲", nil, nil, nil))
	DemoPush(global.Root.Right, false, rbtmodels.NewRBTNode(false, 7, "7秦明", nil, nil, nil))
	DemoPush(global.Root.Left.Left, true, rbtmodels.NewRBTNode(false, 8, "8呼延灼", nil, nil, nil))
	DemoPush(global.Root.Left.Left, false, rbtmodels.NewRBTNode(false, 9, "9华融", nil, nil, nil))
	DemoPush(global.Root.Left.Right, true, rbtmodels.NewRBTNode(false, 10, "10柴进", nil, nil, nil))
	DemoPush(global.Root.Left.Right, false, rbtmodels.NewRBTNode(false, 11, "11李应", nil, nil, nil))
	DemoPush(global.Root.Right.Left, true, rbtmodels.NewRBTNode(false, 12, "12朱仝", nil, nil, nil))
	DemoPush(global.Root.Right.Left, false, rbtmodels.NewRBTNode(false, 13, "13鲁智深", nil, nil, nil))
	DemoPush(global.Root.Right.Right, true, rbtmodels.NewRBTNode(false, 14, "14武松", nil, nil, nil))
	DemoPush(global.Root.Right.Right, false, rbtmodels.NewRBTNode(false, 15, "15董平", nil, nil, nil))

}

// DemoPush 简易附加在尾部，回头废止
func DemoPush(r *rbtmodels.RBTNode, isLeft bool, san *rbtmodels.RBTNode) {
	if isLeft { // 附加在左边
		r.Left = san
	} else { // 附加在右边
		r.Right = san
	}
	san.Parent = r
}
