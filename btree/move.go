package btree

import (
	"errors"
	"fmt"
	"go-b-tree-bplus-tree/btree/btreeconst"
	"go-b-tree-bplus-tree/btree/btreemodels"
)

// MoveKeysLeft 叶子的KEY，排队左移
// @n 节点
// @leftPosition 左面端点
// @rightPosition 右面端点，-1表示最右
// @tailKey 尾部准备补进来键值
// @tailPayLoad 尾部准备补进来载荷
// @tailChild 尾部准备补进来孩子
// @author https://github.com/BrotherSam66/
func MoveKeysLeft(n *btreemodels.BTreeNode, leftPosition int, rightPosition int, tailKey int, tailPayLoad string, tailChild *btreemodels.BTreeNode) (err error) {
	if rightPosition == -1 {
		rightPosition = n.KeyNum - 1
	}
	endPosition := rightPosition         // 循环结束位
	if rightPosition == btreeconst.M-2 { // 说明满员，且从尾巴移动，需要少循环1位并补尾巴
		endPosition--
	}
	for i := leftPosition; i <= endPosition; i++ { // 逐个左移
		n.Key[i] = n.Key[i+1]
		n.Payload[i] = n.Payload[i+1]
		n.Child[i] = n.Child[i+1]
	}
	n.Child[endPosition] = n.Child[endPosition+1] // 右腿，补一下

	if rightPosition == btreeconst.M-2 { // 说明满员，且从尾巴移动，需要少循环1位并补尾巴
		n.Key[rightPosition] = tailKey
		n.Payload[rightPosition] = tailPayLoad
		n.Child[rightPosition+1] = tailChild // 只补右腿，左腿前面处理好了
	}
	if tailChild == nil { // 没有补充的尾巴，key数量就少1
		n.KeyNum--
	}
	return
}

// MoveKeysLeftWithoutLeftChild 叶子的KEY，排队左移,最左腿不动
// @n 节点
// @leftPosition 左面端点
// @rightPosition 右面端点，-1表示最右
// @tailKey 尾部准备补进来键值
// @tailPayLoad 尾部准备补进来载荷
// @tailChild 尾部准备补进来孩子
// @author https://github.com/BrotherSam66/
func MoveKeysLeftWithoutLeftChild(n *btreemodels.BTreeNode, leftPosition int, rightPosition int, tailKey int, tailPayLoad string, tailChild *btreemodels.BTreeNode) (err error) {
	if rightPosition == -1 {
		rightPosition = n.KeyNum - 1
	}
	endPosition := rightPosition         // 循环结束位
	if rightPosition == btreeconst.M-2 { // 说明满员，且从尾巴移动，需要少循环1位并补尾巴
		endPosition--
	}
	for i := leftPosition; i <= endPosition; i++ { // 逐个左移
		n.Key[i] = n.Key[i+1]
		n.Payload[i] = n.Payload[i+1]
		n.Child[i+1] = n.Child[i+2]
	}

	if rightPosition == btreeconst.M-2 { // 说明满员，且从尾巴移动，需要少循环1位并补尾巴
		n.Key[rightPosition] = tailKey
		n.Payload[rightPosition] = tailPayLoad
		n.Child[rightPosition+1] = tailChild // 只补右腿，左腿前面处理好了
	}
	if tailChild == nil { // 没有补充的尾巴，key数量就少1
		n.KeyNum--
	}
	return
}

// MoveKeysRight 叶子的KEY，排队右移
// @n 节点
// @leftPosition 左面端点
// @rightPosition 右面端点，-1表示最右
// @headKey 头部准备补进来键值
// @headPayLoad 头部准备补进来载荷
// @headChild 头部准备补进来孩子
// @author https://github.com/BrotherSam66/
func MoveKeysRight(n *btreemodels.BTreeNode, leftPosition int, rightPosition int, headKey int, headPayLoad string, headChild *btreemodels.BTreeNode) (err error) {
	// 只解决不满的，在最左侧加一个的
	if leftPosition != 0 {
		err = errors.New("出错，leftPosition必须是0！")
		fmt.Println(err.Error())
		return
	}
	if rightPosition != -1 {
		err = errors.New("出错，rightPosition必须是-1！")
		fmt.Println(err.Error())
		return
	}
	if n.KeyNum >= btreeconst.M-1 {
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

// Merge3Nodes 三个节点合并（合并到leftSon）
// @leftSon 准备接收合并的节点
// @parent 父节点，只下来一个Key
// @rightSon 准备被合并的节点
// @avatarPosition 父节点下来Key的位置
// @author https://github.com/BrotherSam66/
func Merge3Nodes(leftSon *btreemodels.BTreeNode, parent *btreemodels.BTreeNode, rightSon *btreemodels.BTreeNode, avatarPosition int) (err error) {
	if leftSon.KeyNum+1+rightSon.KeyNum > btreeconst.M-1 {
		err = errors.New("三个节点Key叠加起来溢出！")
		fmt.Println(err.Error())
		return
	}
	if parent.KeyNum <= 1 { // 父亲剩一个key了
		if parent.Parent == nil { // 父亲剩一个key && 是root，减少1个层级
			//global.Root = leftSon
			leftSon.Parent = nil
		} else { // 不可以借走父亲的key
			err = errors.New("父亲剩一个key，又不是root，不能借！")
			fmt.Println(err.Error())
			return
		}
	}
	// parent 那个key的数据，复制到leftSon(不用带腿)，leftSon.KeyNum++
	leftSon.Key[leftSon.KeyNum] = parent.Key[avatarPosition]
	leftSon.Payload[leftSon.KeyNum] = parent.Payload[avatarPosition]
	leftSon.KeyNum++

	// parent 那个key的数据，删除，保留左腿，后面向左排挤1位，
	_ = MoveKeysLeftWithoutLeftChild(parent, avatarPosition, -1, 0, "", nil)

	// rightSon 所有key的数据，复制到leftSon(多一条最左腿补充parent没带下来的)，leftSon.KeyNUM+=rightSon.KeyNUM，
	_ = Merge2Nodes(leftSon, rightSon)

	return
}

// Merge2Nodes 三个节点合并（合并到leftSon，结合点的腿用rightSon的）
// @leftSon 准备接收合并的节点
// @rightSon 准备被合并的节点
// @author https://github.com/BrotherSam66/
func Merge2Nodes(leftSon *btreemodels.BTreeNode, rightSon *btreemodels.BTreeNode) (err error) {
	if leftSon.KeyNum+rightSon.KeyNum > btreeconst.M-1 {
		err = errors.New("2个节点Key叠加起来溢出！")
		fmt.Println(err.Error())
		return
	}

	// 先处理多出来的左腿
	leftSon.Child[leftSon.KeyNum] = rightSon.Child[0]
	// 循环处理剩下的3要素
	for i := 0; i < rightSon.KeyNum; i++ {
		leftSon.Key[leftSon.KeyNum+i] = rightSon.Key[i]
		leftSon.Payload[leftSon.KeyNum+i] = rightSon.Payload[i]
		leftSon.Child[leftSon.KeyNum+i+1] = rightSon.Child[i+1]
		if rightSon.Child[i+1] != nil { // 下级孙子的上指向也要调整
			rightSon.Child[i+1].Parent = leftSon
		}
	}

	leftSon.KeyNum = leftSon.KeyNum + rightSon.KeyNum // 不用解释吧
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
	rightSon := btreemodels.NewBTreeNode(nil, 1, n.Key[btreeconst.M/2+1], n.Payload[btreeconst.M/2+1]) // 缺很多参数没加
	upNode = btreemodels.NewBTreeNode(nil, 1, n.Key[btreeconst.M/2], n.Payload[btreeconst.M/2])        // 缺很多参数没加

	upNode.Child[0] = n                           // 上升的左腿
	upNode.Child[1] = rightSon                    // 上升的右腿
	rightSon.Parent = upNode                      // 右儿子的爹
	rightSon.Child[0] = n.Child[btreeconst.M/2+1] // 预先补上右儿子的左腿，每个Key的右腿后面的循环里补
	if rightSon.Child[0] != nil {
		rightSon.Child[0].Parent = rightSon // 右儿子的每一个孙子. Parent 都要重新指向
	}

	// 补右儿子。例如 btreeconst.M=10最大9个key；n.KeyNum=9目前多插一个共10个key；升起5号，左儿子0~4，右儿子6~9，首KEY6号已加，下面重新循环6~9，9要单独做
	for j := btreeconst.M/2 + 1; j < btreeconst.M-1; j++ { // 例如①M=10，循环6~8；M=9，循环5~7；M=5，循环3~3
		rightSon.Key[j-btreeconst.M/2-1] = n.Key[j]         // 第一个是0，最后一个是(M+1)/2-2
		rightSon.Payload[j-btreeconst.M/2-1] = n.Payload[j] // 第一个是0，最后一个是(M+1)/2-2
		rightSon.Child[j-btreeconst.M/2] = n.Child[j+1]     // 第一个是1，补每个Key右边的腿
		if rightSon.Child[j-btreeconst.M/2] != nil {
			rightSon.Child[j-btreeconst.M/2].Parent = rightSon // 右儿子的每一个孙子. Parent 都要重新指向
		}
	}
	rightSon.Key[(btreeconst.M+1)/2-2] = keyTail         // 补充尾巴，例如9号
	rightSon.Payload[(btreeconst.M+1)/2-2] = payloadTail // 补充尾巴，例如9号
	rightSon.Child[(btreeconst.M+1)/2-1] = childTail     // 补充尾巴，例如9号的右腿
	if rightSon.Child[(btreeconst.M+1)/2-1] != nil {
		rightSon.Child[(btreeconst.M+1)/2-1].Parent = rightSon // 右儿子的每一个孙子. Parent 都要重新指向
	}
	rightSon.KeyNum = (btreeconst.M+1)/2 - 1 // 右儿子key数，自己算，保证是这个

	// n 其实是左儿子 leftSon。擦除掉(已经升上去+分到右儿子)的数据
	for j := btreeconst.M / 2; j < btreeconst.M-1; j++ { // 例如①M=10，循环5~8；M=9，循环4~7；M=5，循环3~3
		n.Key[j] = 0       // 第一个是0，最后一个是(M+1)/2-2
		n.Payload[j] = ""  // 第一个是0，最后一个是(M+1)/2-2
		n.Child[j+1] = nil // 第一个是1，补每个Key右边的孙子
	}
	n.KeyNum = btreeconst.M / 2 // 左儿子key数量。例如①M=10，循环5~9，左边保留0~4，长度5；M=9，循环4~8，左边保留0~3，长度4

	// 这里只是把中间节点升起来，拟插入下一级，带着两条腿，进入下一层递归。（如果本节点是root，升起来的就是新root就结束）
	if n.Parent == nil { // 说明是root升起来的
		isUpRoot = true
	}
	return
}
