// Package btreeutils
// @Title B树工具包
// @Description  和删除节点有关的操作
// @Author  https://github.com/BrotherSam66/
// @Update
package btreeutils

import (
	"errors"
	"fmt"
	"go-b-tree-bplus-tree/btreemodels"
	"go-b-tree-bplus-tree/global"
	"go-b-tree-bplus-tree/globalconst"
	"math/rand"
)

// Deletes 连续删除节点
// @author https://github.com/BrotherSam66/
func Deletes() {

	for {
		var key int
		fmt.Println("请输入KEY，按回车键(空按回车随机,-1退出)：")
		_, _ = fmt.Scanln(&key)

		if key == -1 {
			return
		}
		if key == 0 {
			key = rand.Intn(global.MaxKey)
			fmt.Println(key)
		}

		if key > 99 || key < 1 {
			fmt.Println("必须是0~~99")
			continue
		}
		Delete(key)
		ShowTree(global.Root)
	}
}

// Delete 删除节点
// @key 键值
// @author https://github.com/BrotherSam66/
func Delete(key int) {

	// 从root开始查找附加的位置
	tempNode, isTarget, err := Search(key)
	if err != nil {
		fmt.Println("查找错误，error == ", err)
		return
	}
	if !isTarget {
		fmt.Println("没找到！ ")
		return
	}

	// 非叶子，查找可替换的叶子节点的KEY值，交换。前序or后继均可，优先前序，前序节点数量<=globalconst.Min不容易删除就定死用后继
	// 查到key在tempNode准确位置，deletePoint
	deletePoint := 0
	for deletePoint = 0; deletePoint < tempNode.KeyNum; deletePoint++ {
		if tempNode.Key[deletePoint] == key { // 准确命中，只可能是新创建节点情形
			break
		}
	}
	if deletePoint >= tempNode.KeyNum {
		fmt.Println("发生某种错误，找到KEY又不存在了！ ")
		return
	}
	avatarNode := tempNode
	if tempNode.Child[0] != nil { // 不是叶子
		avatarNode, _ = PredecessorOrSuccessor(tempNode, key, true) // 用前驱节点做替身
		// 考察avatarNode可简易删除？
		if avatarNode.KeyNum <= globalconst.Min { // 不能简易删除
			avatarNode, _ = PredecessorOrSuccessor(tempNode, key, false) // 用后继节点做替身
			// 删除数据 对换 后继（一个神奇的语句）
			tempNode.Key[deletePoint], avatarNode.Key[0] = avatarNode.Key[0], tempNode.Key[deletePoint]
			tempNode.Payload[deletePoint], avatarNode.Payload[0] = avatarNode.Payload[0], tempNode.Payload[deletePoint]
			deletePoint = 0
		} else {
			// 删除数据 对换 前驱（一个神奇的语句）
			tempNode.Key[deletePoint], avatarNode.Key[avatarNode.KeyNum-1] = avatarNode.Key[avatarNode.KeyNum-1], tempNode.Key[deletePoint]
			tempNode.Payload[deletePoint], avatarNode.Payload[avatarNode.KeyNum-1] = avatarNode.Payload[avatarNode.KeyNum-1], tempNode.Payload[deletePoint]
			deletePoint = avatarNode.KeyNum - 1
		}
	}

	// 到这里，KEY在叶子上，就开始删除的递归流程

	_ = DeleteOneKey(avatarNode, key, deletePoint)

	return
}

// DeleteOneKey 删除一个叶子上的KEY，可能要递归
// @avatar 节点
// @key 拟删除键值
// @deletePoint 拟删除键值的位置
// @author https://github.com/BrotherSam66/
func DeleteOneKey(avatar *btreemodels.BTreeNode, key int, deletePoint int) (err error) {
	if avatar.Key[deletePoint] != key {
		err = errors.New("奇怪啊，指定的位置deletePoint键值不吻合啊")
		fmt.Println("奇怪啊，指定的位置deletePoint键值不吻合啊")
		return
	}

	// 删除掉这个key
	_ = MoveKeysLeft(avatar, deletePoint, -1)

	// 检查合法性，可能要递归
	if avatar.KeyNum < globalconst.Min { // avatar节点过短，需要调整，可能递归
		_ = FixAfterDelete(avatar)
	}

	return
}

// MoveKeysLeft 叶子的KEY，排队左移
// @n 节点
// @leftPoint 左面端点
// @rightPoint 右面端点，-1表示最右
// @author https://github.com/BrotherSam66/
func MoveKeysLeft(n *btreemodels.BTreeNode, leftPoint int, rightPoint int) (err error) {
	// 求rightPoint
	isRightPointFull := false // 是满员的
	if rightPoint < 0 {
		rightPoint = n.KeyNum - 1
	}
	if rightPoint >= globalconst.M-2 { // 是满员的，尾部单独处理
		isRightPointFull = true
		rightPoint--
	}
	for i := leftPoint; i < rightPoint; i++ { // 逐个左移
		n.Key[i] = n.Key[i+1]
		n.Payload[i] = n.Payload[i+1]
		// 叶子，Child都是nil
	}
	if isRightPointFull { // 满员，尾部单独处理
		n.Key[rightPoint+1] = 0
		n.Payload[rightPoint+1] = ""
	}
	n.KeyNum--
	return
}

// FixAfterDelete 删除后调整
// @avatar 节点
// @author https://github.com/BrotherSam66/
func FixAfterDelete(avatar *btreemodels.BTreeNode) (err error) {
	// 如果该节点递归、上升到了root，结束
	if avatar == global.Root {
		return
	}
	// 2）该结点key个数大于等于Math.ceil(m/2)-1，结束删除操作，否则执行第3步。
	if avatar.KeyNum >= globalconst.Min {
		return
	}
	// 3）如果兄弟结点key个数大于Math.ceil(m/2)-1，则父结点中的key下移到该结点，兄弟结点中的一个key上移，删除操作结束。
	// 找出avatar的左右兄弟
	leftBrother := avatar  // 临时定义
	rightBrother := avatar // 临时定义
	parent := avatar.Parent
	// 找到 avatar 在父亲的排位
	i := 0
	for i = 0; i < parent.KeyNum; i++ {
		if avatar.Key[0] < parent.Key[i] { // 小于，说明刚刚越过了，(用avatar任何Key都行)
			break
		}
	}

	// 找到兄弟后直接借KEY
	isSuccess := false
	if i == 0 { // 在最左
		rightBrother = parent.Child[1]
		isSuccess, _ = TryBorrowBrotherKey(avatar, rightBrother, false)
		if isSuccess {
			return
		}
	} else if i >= parent.KeyNum { // 在最右
		leftBrother = parent.Child[parent.KeyNum-2]
		isSuccess, _ = TryBorrowBrotherKey(avatar, leftBrother, true)
		if isSuccess {
			return
		}
	} else { // 居中，有左右2个兄弟
		rightBrother = parent.Child[i+1]
		isSuccess, _ = TryBorrowBrotherKey(avatar, rightBrother, false)
		if isSuccess {
			return
		}
		leftBrother = parent.Child[i-1]
		isSuccess, _ = TryBorrowBrotherKey(avatar, leftBrother, true)
		if isSuccess {
			return
		}
	}

	// 到这里，就是兄弟借不来。将父结点中的key下移与当前结点及它的兄弟结点中的key合并，形成一个新的结点。
	// 原父结点中的key的两个孩子指针就变成了一个孩子指针，指向这个新结点。然后当前结点的指针指向父结点，重复上第2步。

	return
}

// TryBorrowBrotherKey 尝试向兄弟借KEY
// @avatar 本节点
// @brother 兄弟节点
// @isLeftBrother 左兄弟or右兄弟
// @isSuccess 借节点成功了吗？
// @author https://github.com/BrotherSam66/
func TryBorrowBrotherKey(avatar *btreemodels.BTreeNode, brother *btreemodels.BTreeNode, isLeftBrother bool) (isSuccess bool, err error) {
	if brother.KeyNum <= globalconst.Min { // 兄弟太短，没得借
		return
	}
	// 3）如果兄弟结点key个数大于Math.ceil(m/2)-1，则父结点中的key下移到该结点，兄弟结点中的一个key上移，删除操作结束。
	/*
	 *假设：5阶，最大4个KEY、最小2个KEY，
	 *  (20|60             |              80|nil)|  (20|50             |              80|nil)|
	 *  /   \              \                \    |  /   \              \                \    |
	 *(?1)(30|40|50|nil) (70|nil|nil|nil)  (?3)  |(?1)(30|40|nil|nil) (60|70|nil|nil)  (?3)  |
	 *
	 *(70)右边刚删掉(75)，(60)下来并入(70)，(50)上去，填补(60)
	 */
	isSuccess = true
	return
}

/*
 *    B树的删除操作 https://www.cnblogs.com/nullzx/p/8729425.html
 *    删除操作是指，根据key删除记录，如果B树中的记录中不存对应key的记录，则删除失败。
 * 1）如果当前需要删除的key位于非叶子结点上，则用后继key（这里的后继key均指后继记录的意思）覆盖要删除的key，
 * 然后在后继key所在的子支中删除该后继key。此时后继key一定位于叶子结点上，这个过程和二叉搜索树删除结点的方式类似。删除这个记录后执行第2步
 * 2）该结点key个数大于等于Math.ceil(m/2)-1，结束删除操作，否则执行第3步。
 * 3）如果兄弟结点key个数大于Math.ceil(m/2)-1，则父结点中的key下移到该结点，兄弟结点中的一个key上移，删除操作结束。
 *    否则，将父结点中的key下移与当前结点及它的兄弟结点中的key合并，形成一个新的结点。原父结点中的key的两个孩子指针就变成了一个孩子指针，
 * 指向这个新结点。然后当前结点的指针指向父结点，重复上第2步。
 *    有些结点它可能即有左兄弟，又有右兄弟，那么我们任意选择一个兄弟结点进行操作即可。
 */

// 上下交换与i个key，明确上下分别交换第几位key

// 左右合并函数

// 左右分裂函数 明确 从哪个位置分裂

// 元素向左移动

// 元素向右移动
//_=MoveKeys(n,)
// 部分元素 迁徙（key位序一般会不一样）
