// Package bplustreemodels
// @Title B树模型
// @Description  模型和构造函数
// @Author  https://github.com/BrotherSam66/
// @Update
package bplustreemodels

/*
   一个k阶的B+树具有如下几个特征：
   有k个子树的中间节点包含有k个元素（B树中是k-1个元素），每个元素不保存数据，只用来索引，所有数据都保存在叶子节点。
   所有的叶子结点中包含了全部元素的信息，及指向含这些元素记录的指针，且叶子结点本身依关键字的大小自小而大顺序链接。
   所有的中间节点元素都同时存在于子节点，在子节点元素中是最大（或最小）元素。
   叶子节点有链接右侧叶子节点的单向链接。
*/

// BPTreeNode B+树结构模型
// @Author  https://github.com/BrotherSam66/
type BPTreeNode struct {
	Parent *BPTreeNode   //指向父节点的指针
	Key    []int         //关键字向量
	Child  []*BPTreeNode // 子树指针向量
	Leaf   []*BPTreeLeaf // 指向叶子节点（如果有的话）
}

// BPTreeLeaf B+树结构模型
// @Author  https://github.com/BrotherSam66/
type BPTreeLeaf struct {
	Parent       *BPTreeNode //指向父节点的指针
	Key          int
	Payload      string      // 叶子的载荷信息
	RightBrother *BPTreeLeaf // 叶子指向右兄弟指针向量
}

// NewBPTreeNode 构造函数
// @parent 指向父亲的指针
// @key 本节点第一个KEY的值
// @payload 本节点的载荷（用来承载节点的信息）
// @Author https://github.com/BrotherSam66/
func NewBPTreeNode(parent *BPTreeNode, key int, leaf *BPTreeLeaf) *BPTreeNode {
	return &BPTreeNode{
		Parent: parent,
		Key:    []int{key},
		Leaf:   []*BPTreeLeaf{leaf},
	}
}

// NewBlankBPTreeNode 构造函数
// @parent 指向父亲的指针
// @key 本节点第一个KEY的值
// @payload 本节点的载荷（用来承载节点的信息）
// @Author https://github.com/BrotherSam66/
func NewBlankBPTreeNode() *BPTreeNode {
	return &BPTreeNode{
		Parent: parent,
		Key:    []int{key},
		Leaf:   []*BPTreeLeaf{leaf},
	}
}

// NewBPTreeLeaf 构造函数
// @parent 指向父亲的指针
// @key 本节点第一个KEY的值
// @payload 本节点的载荷（用来承载节点的信息）
// @Author https://github.com/BrotherSam66/
func NewBPTreeLeaf(parent *BPTreeNode, rightBrother *BPTreeLeaf, key int, payload string) *BPTreeLeaf {
	return &BPTreeLeaf{
		Parent:       parent,
		Key:          key,
		Payload:      payload,
		RightBrother: rightBrother,
	}
}
