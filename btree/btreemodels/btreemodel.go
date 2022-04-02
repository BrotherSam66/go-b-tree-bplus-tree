// Package btreemodels
// @Title B树模型
// @Description  模型和构造函数
// @Author  https://github.com/BrotherSam66/
// @Update
package btreemodels

import (
	"go-b-tree-bplus-tree/btree/btreeconst"
)

// BTreeNode B树结构模型
// @Author  https://github.com/BrotherSam66/
type BTreeNode struct {
	Parent  *BTreeNode               //指向父节点的指针
	KeyNum  int                      //关键字个数
	Key     [btreeconst.M - 1]int    //关键字向量
	Payload [btreeconst.M - 1]string // 本Key的负载信息
	Child   [btreeconst.M]*BTreeNode //子树指针向量
}

// NewBTreeNode 构造函数
// @parent 指向父亲的指针
// @keyNum 本节点的KEY数量
// @key 本节点第一个KEY的值
// @payload 本节点的载荷（用来承载节点的信息）
// @Author https://github.com/BrotherSam66/
func NewBTreeNode(parent *BTreeNode, keyNum int, key int, payload string) *BTreeNode {
	return &BTreeNode{
		Parent:  parent,
		KeyNum:  keyNum,                     // 新创建得通常都会是1
		Key:     [btreeconst.M - 1]int{key}, // 暂时放在[0]位。todo 这里数组长度和初始化元素个数不符合可能有问题
		Payload: [btreeconst.M - 1]string{payload},
		//Payload: [btreeconst.M - 1]string{fmt.Sprintf("%d", key)},
	}
}

/*
一、B树的定义
	1.定义任意非叶子结点最多只有M个儿子，且M>2；
	2.根结点的儿子数为[2, M]；
	3.除根结点以外的非叶子结点的儿子数为[M/2, M]；
	4.每个结点存放至少M/2-1（取上整）和至多M-1个关键字；（至少2个关键字）
	5.非叶子结点的关键字个数=指向儿子的指针个数-1；
	6.非叶子结点的关键字：K[1], K[2], …, K[M-1]，且K[i] < K[i+1]；
	7.非叶子结点的指针：P[1], P[2], …, P[M]，其中P[1]指向关键字小于K[1]的子树，P[M]指向关键字大于K[M-1]的子树，其它P[i]指向关键字属于(K[i-1], K[i])的子树；
	8.所有叶子结点位于同一层；

二、B树的插入操作
插入操作是指插入一个关键字（也可以是一个复杂的自定义结构体）。如果B树中已存在需要插入的关键字，则不再插入。若B树不存在这个关键字,则一定是在叶子结点中进行插入操作。
1、根据要插入的关键字的值，找到叶子结点并插入。
2、判断当前结点关键字的个数是否小于等于m-1，若满足则结束，否则进行第3步。
3、以结点中间的关键字为中心分裂成左右两部分，然后将这个中间的关键字插入到父结点中，这个关键字的左子树指向分裂后的左半部分，这个关键字的右子支指向分裂后的右半部分，然后将当前结点指向父结点，继续进行第3步。

三、B树的删除操作
删除操作是指删除该B树中的某个节点中的指定关键字，如果B树中不存对应的关键字，则删除失败。
1、如果当前需要删除的关键字位于非叶子结点上，则用后继最小关键字覆盖要删除的关键字，然后在后继关键字所在的子支中删除该后继关键字。此时后继关键字一定位于叶子结点上，这个过程和二叉搜索树删除结点的方式类似。删除这个记录后执行第2步
2、该结点关键字个数大于等于Math.ceil(m/2)-1，结束删除操作，否则执行第3步。
3、如果兄弟结点关键字个数大于Math.ceil(m/2)-1，则父结点中的关键字下移到该结点，兄弟结点中的一个关键字上移，删除操作结束。
否则，将父结点中的关键字下移与当前结点及它的兄弟结点中的关键字合并，形成一个新的结点。原父结点中的key的两个孩子指针就变成了一个孩子指针，指向这个新结点。然后当前结点的指针指向父结点，重复上第2步。
有些结点它可能即有左兄弟，又有右兄弟，那么我们任意选择一个兄弟结点进行操作即可。

向上分裂
合并操作
B树的的展示
*/
