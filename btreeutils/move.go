package btreeutils

import (
	"errors"
	"fmt"
	"go-b-tree-bplus-tree/btreemodels"
	"go-b-tree-bplus-tree/globalconst"
)

// MoveKeysLeft 叶子的KEY，排队左移
// @n 节点
// @leftPoint 左面端点
// @rightPoint 右面端点，-1表示最右
// @tailKey 尾部准备补进来键值
// @tailPayLoad 尾部准备补进来载荷
// @tailChild 尾部准备补进来孩子
// @author https://github.com/BrotherSam66/
func MoveKeysLeft(n *btreemodels.BTreeNode, leftPoint int, rightPoint int, tailKey int, tailPayLoad string, tailChild *btreemodels.BTreeNode) (err error) {
	if rightPoint == -1 {
		rightPoint = n.KeyNum - 1
	}
	endPoint := rightPoint             // 循环结束位
	if rightPoint == globalconst.M-2 { // 说明满员，且从尾巴移动，需要少循环1位并补尾巴
		endPoint--
	}
	for i := leftPoint; i <= endPoint; i++ { // 逐个左移
		n.Key[i] = n.Key[i+1]
		n.Payload[i] = n.Payload[i+1]
		n.Child[i] = n.Child[i+1]
	}
	n.Child[endPoint] = n.Child[endPoint+1] // 右腿，补一下

	if rightPoint == globalconst.M-2 { // 说明满员，且从尾巴移动，需要少循环1位并补尾巴
		n.Key[rightPoint] = tailKey
		n.Payload[rightPoint] = tailPayLoad
		n.Child[rightPoint+1] = tailChild // 只补右腿，左腿前面处理好了
	}
	if tailChild == nil { // 没有补充的尾巴，key数量就少1
		n.KeyNum--
	}
	return
}

// MoveKeysLeftWithoutLeftChild 叶子的KEY，排队左移,最左腿不动
// @n 节点
// @leftPoint 左面端点
// @rightPoint 右面端点，-1表示最右
// @tailKey 尾部准备补进来键值
// @tailPayLoad 尾部准备补进来载荷
// @tailChild 尾部准备补进来孩子
// @author https://github.com/BrotherSam66/
func MoveKeysLeftWithoutLeftChild(n *btreemodels.BTreeNode, leftPoint int, rightPoint int, tailKey int, tailPayLoad string, tailChild *btreemodels.BTreeNode) (err error) {
	if rightPoint == -1 {
		rightPoint = n.KeyNum - 1
	}
	endPoint := rightPoint             // 循环结束位
	if rightPoint == globalconst.M-2 { // 说明满员，且从尾巴移动，需要少循环1位并补尾巴
		endPoint--
	}
	for i := leftPoint; i <= endPoint; i++ { // 逐个左移
		n.Key[i] = n.Key[i+1]
		n.Payload[i] = n.Payload[i+1]
		n.Child[i+1] = n.Child[i+2]
	}

	if rightPoint == globalconst.M-2 { // 说明满员，且从尾巴移动，需要少循环1位并补尾巴
		n.Key[rightPoint] = tailKey
		n.Payload[rightPoint] = tailPayLoad
		n.Child[rightPoint+1] = tailChild // 只补右腿，左腿前面处理好了
	}
	if tailChild == nil { // 没有补充的尾巴，key数量就少1
		n.KeyNum--
	}
	return
}

// MoveKeysRight 叶子的KEY，排队右移
// @n 节点
// @leftPoint 左面端点
// @rightPoint 右面端点，-1表示最右
// @headKey 头部准备补进来键值
// @headPayLoad 头部准备补进来载荷
// @headChild 头部准备补进来孩子
// @author https://github.com/BrotherSam66/
func MoveKeysRight(n *btreemodels.BTreeNode, leftPoint int, rightPoint int, headKey int, headPayLoad string, headChild *btreemodels.BTreeNode) (err error) {
	// 只解决不满的，在最左侧加一个的
	if leftPoint != 0 {
		err = errors.New("出错，leftPoint必须是0！")
		fmt.Println(err.Error())
		return
	}
	if rightPoint != -1 {
		err = errors.New("出错，rightPoint必须是-1！")
		fmt.Println(err.Error())
		return
	}
	if n.KeyNum >= globalconst.M-1 {
		err = errors.New("出错，本节点满的，加不进来！")
		fmt.Println(err.Error())
		return
	}

	n.Child[n.KeyNum+2] = n.Child[n.KeyNum+1] // 最右腿，补一下
	for i := n.KeyNum - 1; i >= 0; i-- {      // 逐个左移
		n.Key[i+1] = n.Key[i]
		n.Payload[i+1] = n.Payload[i]
		n.Child[i+1] = n.Child[i]
	}

	n.Key[0] = headKey
	n.Payload[0] = headPayLoad
	n.Child[0] = headChild

	n.KeyNum++
	return
}

// Merge3Nodes 三个节点合并（合并到leftSan）
// @leftSan 准备接收合并的节点
// @parent 父节点，只下来一个Key
// @rightSan 准备被合并的节点
// @avatarPoint 父节点下来Key的位置
// @author https://github.com/BrotherSam66/
func Merge3Nodes(leftSan *btreemodels.BTreeNode, parent *btreemodels.BTreeNode, rightSan *btreemodels.BTreeNode, avatarPoint int) (err error) {
	if leftSan.KeyNum+1+rightSan.KeyNum > globalconst.M-1 {
		err = errors.New("三个节点Key叠加起来溢出！")
		fmt.Println(err.Error())
		return
	}
	if parent.KeyNum <= 1 { // 父亲剩一个key了
		if parent.Parent == nil { // 父亲剩一个key && 是root，减少1个层级
			//global.Root = leftSan
			leftSan.Parent = nil
		} else { // 不可以借走父亲的key
			err = errors.New("父亲剩一个key，又不是root，不能借！")
			fmt.Println(err.Error())
			return
		}
	}
	// parent 那个key的数据，复制到leftSan(不用带腿)，leftSan.KeyNum++
	leftSan.Key[leftSan.KeyNum] = parent.Key[avatarPoint]
	leftSan.Payload[leftSan.KeyNum] = parent.Payload[avatarPoint]
	leftSan.KeyNum++

	// parent 那个key的数据，删除，保留左腿，后面向左排挤1位，
	_ = MoveKeysLeftWithoutLeftChild(parent, avatarPoint, -1, 0, "", nil)

	// rightSan 所有key的数据，复制到leftSan(多一条最左腿补充parent没带下来的)，leftSan.KeyNUM+=rightSan.KeyNUM，
	_ = Merge2Nodes(leftSan, rightSan)

	return
}

// Merge2Nodes 三个节点合并（合并到leftSan，结合点的腿用rightSan的）
// @leftSan 准备接收合并的节点
// @rightSan 准备被合并的节点
// @author https://github.com/BrotherSam66/
func Merge2Nodes(leftSan *btreemodels.BTreeNode, rightSan *btreemodels.BTreeNode) (err error) {
	if leftSan.KeyNum+rightSan.KeyNum > globalconst.M-1 {
		err = errors.New("2个节点Key叠加起来溢出！")
		fmt.Println(err.Error())
		return
	}

	// 先处理多出来的左腿
	leftSan.Child[leftSan.KeyNum] = rightSan.Child[0]
	// 循环处理剩下的3要素
	for i := 0; i < rightSan.KeyNum; i++ {
		leftSan.Key[leftSan.KeyNum+i] = rightSan.Key[i]
		leftSan.Payload[leftSan.KeyNum+i] = rightSan.Payload[i]
		leftSan.Child[leftSan.KeyNum+i+1] = rightSan.Child[i+1]
		if rightSan.Child[i+1] != nil { // 下级孙子的上指向也要调整
			rightSan.Child[i+1].Parent = leftSan
		}
	}

	leftSan.KeyNum = leftSan.KeyNum + rightSan.KeyNum // 不用解释吧
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
