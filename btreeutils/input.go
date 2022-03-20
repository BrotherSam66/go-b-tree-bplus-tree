// Package btreeutils
// @Title B树工具包
// @Description  和插入节点有关的操作
// @Author  https://github.com/BrotherSam66/
// @Update
package btreeutils

import (
	"fmt"
	"go-b-tree-bplus-tree/btreemodels"
	"go-b-tree-bplus-tree/global"
	"go-b-tree-bplus-tree/globalconst"

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

			ShowTree(global.Root)
			continue
		}
		if key > 99 || key < 1 {
			fmt.Println("必须是0~~99")
			continue
		}
		Insert(key, "")
		ShowTree(global.Root)
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
	if global.Root == nil { // 原树为空树，新加入的转为根
		global.Root = btreemodels.NewBTreeNode(nil, 1, key, payload)
		return
	}

	// 从root开始查找附加的位置
	tempNode, isTarget, err := Search(key)
	if err != nil {
		fmt.Println("没找到or查找错误，error == ", err)
		return
	}
	if isTarget { // 拟插入的key存在，替换payload就好
		// 寻找替换的位置
		for j := 0; j < tempNode.KeyNum; j++ {
			if tempNode.Key[j] == key { // 准确命中，只可能是新创建节点情形
				tempNode.Payload[j] = payload // 拟插入的key存在，替换payload就好
				fmt.Println("拟插入的key存在，替换payload就好")
			}
		}
		return
	}

	// 到这里，找到的必然是叶子节点。拟插入的key放在这里，可能递归
	_ = InsertOneNode(tempNode, btreemodels.NewBTreeNode(nil, 1, key, payload))

	return
}

// InsertOneNode 插入一个节点，可能要递归
// @n 被插入的节点
// @insertNode 拟插入的节点，要么①新节点，儿子均nil；要么②下层满员把中间节点挤上来（上来前key放在Key[1]，把下层分裂，作为我的2个儿子，放在Child[0,1]）
// @author https://github.com/BrotherSam66/
func InsertOneNode(n *btreemodels.BTreeNode, insertNode *btreemodels.BTreeNode) (err error) {
	// 寻找插入的位置，拟插入放在这个点前面
	keyPoint := 0
	for keyPoint = 0; keyPoint < n.KeyNum; keyPoint++ {
		if insertNode.Key[0] == n.Key[keyPoint] { // 准确命中，只可能是新创建节点
			n.Payload = insertNode.Payload
			return
		} else if insertNode.Key[0] < n.Key[keyPoint] { // 说明已经找过头了,结束本节点循环，key插在i前面
			break
		}
		// 到这里：可能①会向后找；可能②KeyNum循环结束，得到的i是最右key的右边，拟插入key本组最大。
	}

	// 到这里：i表示了拟插入key的位置。insertNode可能是不带孩子的新创建节点，也可能是下层挤上来的带2个孩子的节点(不会凭空上来，有一条腿是要替换原来的父节点的，我们指定用左腿)
	// 强行插入，无论是否满员，溢出的在Tail里
	keyTail, payloadTail, childTail, _ := InsertOneKey(n, insertNode, keyPoint)
	// 分析本节点是否需要裂变
	if n.KeyNum < globalconst.M-1 { // 被插入节点不满员，不用递归
		n.KeyNum++
		return
	}

	// 到这里，被插入节点满员，就需要分裂，需要递归了
	// 开始对本节点分裂，分裂成3个，升起的是M/2位置的
	upNode, isUpRoot, _ := SplitTo3Node(n, keyTail, payloadTail, childTail)

	// 这里只是把中间节点升起来，拟插入下一级，带着两条腿，进入下一层递归。（如果本节点是root，升起来的就是新root就结束）
	if isUpRoot { // 说明升起来的是单root
		n.Parent = upNode    // 左儿子重新认爹
		global.Root = upNode // 重新指定根节点n.Parent = {*go-b-tree-bplus-tree/btreemodels.BTreeNode | 0xc000120500}
		return
	} else { // 不是root升起来的。递归...
		tempNode := n.Parent              // 原来被插入的节点的爹作为新的被插入的节点，拿来递归的
		upNode.Child[0].Parent = tempNode // 上升节点的两个儿子指向上升节点拟插入的节点
		upNode.Child[1].Parent = tempNode // 上升节点的两个儿子指向上升节点拟插入的节点
		//n.Parent = upNode                   // 原来被插入的节点当up节点的左儿子
		_ = InsertOneNode(tempNode, upNode) // 递归
		return
	}
	// 不可能到这里
}

// InsertOneKey 插入一个Key，满了也插，溢出在Tail里
// @n 被插入节点
// @insertNode 拟插入节点
// @insertPoint 拟插入位置，新入的占用这个位置
// @keyTail 准备承载Key数组最后一个元素
// @ChildTail  准备承载Child数组最后一个元素
// @payloadTail 准备承载payload数组最后一个元素
// @author https://github.com/BrotherSam66/
/*
 *假设：5阶，最大4个KEY、最小2个KEY，孩子数=KEY数+1，(65)是从60|70中间原来指向节点分裂升上来的
 *  (20|30  |              80)   |  (20|30  |              80)    |  (20|30  |      60|       80)   |
 *  /   \    \                \  |  /   \    \                \   |  /   \    \        \         \  |
 *(?1)(?2) (40|50|60   |70)  (?3)|(?1)(?2) (40|50|60|65 | 70) (?3)|(?1)(?2) (40|50)     (65|70) (?3)|
 *         /   \  \        \     |         /   \  \  \   \   \    |         /   \  \    /   \  \    |
 *       (?4)(?5)(?6) (65) (?7)  |       (?4)(?5)(?6)(?8)(?9)(?7) |       (?4)(?5)(?6) (?8)(?9)(?7) |
 *                    / \        |                                |                                 |
 *                  (?8)(?9)     |                                |                                 |
 *(?8)是(65)原归属节点左半部分，原来就和60|70指针勾连，
 *(?9)是(65)原归属节点右半部分，是新分裂出来的。
 */
func InsertOneKey(n *btreemodels.BTreeNode, insertNode *btreemodels.BTreeNode, insertPoint int) (keyTail int, payloadTail string, childTail *btreemodels.BTreeNode, err error) {
	keyTail = n.Key[globalconst.M-2]         // 数组最后一个元素
	payloadTail = n.Payload[globalconst.M-2] // 数组最后一个元素
	childTail = n.Child[globalconst.M-1]     // 数组最后一个元素

	// 把往后挤走的KEY处理完
	// 例如 globalconst.M-1=9最大9个key；n.KeyNum=6目前6个key；①keyPoint=3表示拟插入要在3这个位置，②keyPoint=6表示拟插入最大
	//for j := n.KeyNum; j > insertPoint; j-- { // 例如①KeyNum=4，insertPoint=0，j=3~1；②KeyNum=4，insertPoint=1，j=3~2
	for j := globalconst.M - 2; j > insertPoint; j-- { // 咬死从数组最后一个元素倒其，可能浪费些算力
		n.Key[j] = n.Key[j-1]
		n.Payload[j] = n.Payload[j-1]
		n.Child[j+1] = n.Child[j] // 搬移的是每个Key的右腿
	}

	// 把拟插入节点放进来
	// 升上来的节点(不会凭空上来，有一条腿是要替换原来的父节点的，我们指定共享左腿，但是插入0位指定共享右腿)。下面有一句是废话
	if insertPoint > globalconst.M-2 { // 插入的是在溢出的尾巴，实际上插入的才是溢出
		keyTail = insertNode.Key[0]         // 溢出的key
		payloadTail = insertNode.Payload[0] // 数组最后一个元素
		childTail = insertNode.Child[1]     // 溢出的右腿。（左腿insertNode上来前已经确保和n的最右腿取值一样了）
	} else { // 插入的不是在尾巴，真实插入
		n.Key[insertPoint] = insertNode.Key[0]
		n.Payload[insertPoint] = insertNode.Payload[0]
		n.Child[insertPoint+1] = insertNode.Child[1] // 右腿
		// todo 还要考虑插入的左右腿
		if insertPoint == 0 {
			n.Child[0] = insertNode.Child[0] // 左腿
		}
	}
	return
}

// SplitTo3Node 左右3分裂，
// @n 被分裂节点，同时也是分裂后的左儿子
// @keyTail Key数组最后一个元素
// @payloadTail payload数组最后一个元素
// @ChildTail  Child数组最后一个元素
// @upNode 准备上升的节点(中间节点)
// @isUpRoot
// @author https://github.com/BrotherSam66/
func SplitTo3Node(n *btreemodels.BTreeNode, keyTail int, payloadTail string, childTail *btreemodels.BTreeNode) (upNode *btreemodels.BTreeNode, isUpRoot bool, err error) {
	rightSan := btreemodels.NewBTreeNode(nil, 1, n.Key[globalconst.M/2+1], n.Payload[globalconst.M/2+1]) // 缺很多参数没加
	upNode = btreemodels.NewBTreeNode(nil, 1, n.Key[globalconst.M/2], n.Payload[globalconst.M/2])        // 缺很多参数没加

	upNode.Child[0] = n                            // 上升的左腿
	upNode.Child[1] = rightSan                     // 上升的右腿
	rightSan.Parent = upNode                       // 右儿子的爹
	rightSan.Child[0] = n.Child[globalconst.M/2+1] // 预先补上右儿子的左腿，每个Key的右腿后面的循环里补
	if rightSan.Child[0] != nil {
		rightSan.Child[0].Parent = rightSan // 右儿子的每一个孙子. Parent 都要重新指向
	}

	// 补右儿子。例如 globalconst.M=10最大9个key；n.KeyNum=9目前多插一个共10个key；升起5号，左儿子0~4，右儿子6~9，首KEY6号已加，下面重新循环6~9，9要单独做
	for j := globalconst.M/2 + 1; j < globalconst.M-1; j++ { // 例如①M=10，循环6~8；M=9，循环5~7；M=5，循环3~3
		rightSan.Key[j-globalconst.M/2-1] = n.Key[j]         // 第一个是0，最后一个是(M+1)/2-2
		rightSan.Payload[j-globalconst.M/2-1] = n.Payload[j] // 第一个是0，最后一个是(M+1)/2-2
		rightSan.Child[j-globalconst.M/2] = n.Child[j+1]     // 第一个是1，补每个Key右边的腿
		if rightSan.Child[j-globalconst.M/2] != nil {
			rightSan.Child[j-globalconst.M/2].Parent = rightSan // 右儿子的每一个孙子. Parent 都要重新指向
		}
	}
	rightSan.Key[(globalconst.M+1)/2-2] = keyTail         // 补充尾巴，例如9号
	rightSan.Payload[(globalconst.M+1)/2-2] = payloadTail // 补充尾巴，例如9号
	rightSan.Child[(globalconst.M+1)/2-1] = childTail     // 补充尾巴，例如9号的右腿
	if rightSan.Child[(globalconst.M+1)/2-1] != nil {
		rightSan.Child[(globalconst.M+1)/2-1].Parent = rightSan // 右儿子的每一个孙子. Parent 都要重新指向
	}
	rightSan.KeyNum = (globalconst.M+1)/2 - 1 // 右儿子key数，自己算，保证是这个

	// n 其实是左儿子 leftSan。擦除掉(已经升上去+分到右儿子)的数据
	for j := globalconst.M / 2; j < globalconst.M-1; j++ { // 例如①M=10，循环5~8；M=9，循环4~7；M=5，循环3~3
		n.Key[j] = 0       // 第一个是0，最后一个是(M+1)/2-2
		n.Payload[j] = ""  // 第一个是0，最后一个是(M+1)/2-2
		n.Child[j+1] = nil // 第一个是1，补每个Key右边的孙子
	}
	n.KeyNum = globalconst.M / 2 // 左儿子key数量。例如①M=10，循环5~9，左边保留0~4，长度5；M=9，循环4~8，左边保留0~3，长度4

	// 这里只是把中间节点升起来，拟插入下一级，带着两条腿，进入下一层递归。（如果本节点是root，升起来的就是新root就结束）
	if n.Parent == nil { // 说明是root升起来的
		isUpRoot = true
	}
	return
}
