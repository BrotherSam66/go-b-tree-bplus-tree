package bplustree

import (
	"errors"
	"fmt"
	"go-b-tree-bplus-tree/bplustree/bplustreeconst"
	"go-b-tree-bplus-tree/bplustree/bplustreeglobal"
	"go-b-tree-bplus-tree/bplustree/bplustreemodels"
)

// MoveKeysLeft 排队左移(删除一个KEY)
// @n 节点
// @leftPosition 左面端点
// @author https://github.com/BrotherSam66/
func MoveKeysLeft(n *bplustreemodels.BPTreeNode, leftPosition int) (err error) {

	if len(n.Child) > 0 && n.Child[0] != nil { // 被分裂的是分支，如果是末梢/叶子，就没有child
		n.Key = append(n.Key[:leftPosition], n.Key[leftPosition+1:]...)       // key 不必深拷贝
		n.Child = append(n.Child[:leftPosition], n.Child[leftPosition+1:]...) // 左组指向下级分支的Child  不必深拷贝
		// 下级指上我的，不用动
	} else { // 叶子的。叶子才有payload
		n.Key = append(n.Key[:leftPosition], n.Key[leftPosition+1:]...)             // key 不必深拷贝
		n.Payload = append(n.Payload[:leftPosition], n.Payload[leftPosition+1:]...) // Payload 不必深拷贝
	}

	// 删除的是节点的最有，需要向上递归调整最右KEY元素（这个不在这里做，在这里return后做）

	return
}

// Merge2Nodes 节点合并（合并到leftSon）
// @leftSon 准备接收合并的节点
// @rightSon 准备被合并的节点
// @author https://github.com/BrotherSam66/
/*
 *假设：5阶，最大5个KEY、最小3个KEY，70向左合并
 *  (20   |   50      |      80|)     |  (20               |      80|)     |
 *  /           \             \       |  /                         \       |
 *(?1)(30 | 40 | 50)       (70|80)    |(?1)(30 | 40 | 50    |    70|80)    |
 *    /      \    \         /    \    |    /      \    \         /    \    |
 *(21|30)(31|40)(41|50) (61|70)(71|80)|(21|30)(31|40)(41|50) (61|70)(71|80)|
 */
func Merge2Nodes(leftSon, rightSon *bplustreemodels.BPTreeNode, leftPosition int) (err error) {
	if len(leftSon.Key)+len(rightSon.Key) > bplustreeconst.M {
		err = errors.New("2个节点Key叠加起来溢出！")
		fmt.Println(err.Error())
		return
	}

	// 开始大搬家（合并到leftSon）
	if len(rightSon.Child) > 0 && rightSon.Child[0] != nil { // 是分支，如果是末梢/叶子，就没有child
		// 向 leftSon 合并 key+payload+child
		leftSon.Key = append(leftSon.Key, rightSon.Key...)       // key
		leftSon.Child = append(leftSon.Child, rightSon.Child...) // Child
		// rightSon 的孩子没上指 leftSon （叶子没有）
		for i := 0; i < len(rightSon.Child); i++ {
			rightSon.Child[i].Parent = leftSon
		}
	} else { // 叶子的。叶子才有payload
		// 向 leftSon 合并 key+payload+child
		leftSon.Key = append(leftSon.Key, rightSon.Key...)             // key
		leftSon.Payload = append(leftSon.Payload, rightSon.Payload...) // Child

		// 调整 Sqt 数据 横向指针
		leftSon.RightBrother = rightSon.RightBrother

	}

	// 父亲 剪掉 左son上指的那个 key+ 右son上指的child
	leftSon.Parent.Key = append(leftSon.Parent.Key[:leftPosition], leftSon.Parent.Key[leftPosition+1:]...)
	leftSon.Parent.Child = append(leftSon.Parent.Child[:leftPosition+1], leftSon.Parent.Child[leftPosition+2:]...)

	// 如果 父亲剩一个腿儿 + 父亲是root ，剪掉父亲，层数降级
	if leftSon.Parent.Parent == nil && len(leftSon.Parent.Key) < 2 {
		leftSon.Parent = nil
		bplustreeglobal.Root = leftSon
	}
	return
}

// ModifyParentKeyByOldestSon 修改/删除/增加了达儿子，向上递归，修改爸爸的key
// @n 节点
// @key 修订后的key
// @oldKey 被修订的key（这个一定有，必须有）
// @author https://github.com/BrotherSam66/
func ModifyParentKeyByOldestSon(n *bplustreemodels.BPTreeNode, key, oldKey int) (err error) {
	if n.Parent == nil {
		fmt.Println("不正常，ModifyParentKeyByOldestSon n.Parent == nil", oldKey)
		return
	}
	tempNode := n
	for tempNode.Parent != nil { // 不到root了，修正持续
		tempNode = tempNode.Parent
		insertPosition, modifyPosition := FindKeyPosition(&tempNode.Key, oldKey)
		if insertPosition > -1 {
			fmt.Println("不正常，ModifyParentKeyByOldestSon FindKeyPosition 准确查到叶子，不应该走到这里", oldKey)
		}
		// 替换
		tempNode.Key[modifyPosition] = key
		// 是否在尾巴，不在尾巴就结束
		if modifyPosition < len(tempNode.Key)-1 {
			return
		}
	}
	return
}
