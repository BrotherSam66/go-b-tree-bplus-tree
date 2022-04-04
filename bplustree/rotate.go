package bplustree

import (
	"errors"
	"fmt"
	"go-b-tree-bplus-tree/btree/btreeconst"
	"go-b-tree-bplus-tree/btree/btreemodels"
)

// RightRotate 右旋
// @param p Position 旋转的出发节点。P(50)向上把父清的一个Key挤下来给兄弟系欸但
// 右旋，必然是①P的最右key，②挤下来父节点P.Patent右侧的Key，③下来给P右兄弟的0位，其他向后排挤，
// ④P的最右儿子，改为P右兄弟的最左[0位]孩子
// @author https://github.com/BrotherSam66/
/*
 *假设：5阶，最大4个KEY、最小2个KEY，
 *  (20  |     60                        |      80|)     |  (20  |     50                        |      80|)      |
 *  /     \                               \        \     |  /     \                               \        \      |
 *(?1)(30  |  40   |    50   |)           (70|||)  (?3)  |(?1)(30  |  40   ||)           (60    |70    ||)  (?3)  |
 *   /      \       \         \           /      \       |   /      \       \           /        \      \         |
 *(21|29||)(31|39||)(41|49||)(51|59||) (61|69||)(71|79||)|(21|29||)(31|39||)(41|49||) (51|59||) (61|69||)(71|79||)|
 */
func RightRotate(p *btreemodels.BTreeNode) (err error) {
	fmt.Println("RightRotate")
	if p == nil {
		err = errors.New("出错，本节点是空！")
		fmt.Println(err.Error())
		return
	}
	if p.KeyNum < 3 {
		err = errors.New("出错，本节点KeyNum小于3，借不出去！")
		fmt.Println(err.Error())
		return
	}
	parent := p.Parent // 父亲
	if parent == nil {
		err = errors.New("出错，本节点父亲是空！")
		fmt.Println(err.Error())
		return
	}
	// 查找p在父节点位置，定下来谁下来
	upKey := p.Key[p.KeyNum-1] // 准备上去的Key值
	var downPosition int       // 准备下来Key的位置
	for downPosition = 0; downPosition < parent.KeyNum; downPosition++ {
		if upKey < parent.Key[downPosition] { // 小于，说明刚刚越过了.downPosition=upKey位置。不可能循环完了，否则没有右兄弟了
			break
		}
	}
	if downPosition > parent.KeyNum-1 {
		err = errors.New(fmt.Sprintf("RightRotate出错，父亲找到的downPosition:%d不对！", downPosition))
		fmt.Println(err.Error())
		return
	}

	// 查找右兄弟，右兄弟满员否
	rightBrother := parent.Child[downPosition+1]
	if rightBrother.KeyNum >= btreeconst.M-1 {
		err = errors.New("出错，右兄弟是满员的！")
		fmt.Println(err.Error())
		return
	}

	// 开始大搬家
	downKey := parent.Key[downPosition]                                     // 准备下来，并入右兄弟的key来自父节点
	downPayload := parent.Payload[downPosition]                             // 准备下来，并入右兄弟的载荷来自父节点
	moveChild := p.Child[p.KeyNum]                                          // 准备下来，并入右兄弟的child来自左兄弟，最右手
	_ = MoveKeysRight(rightBrother, 0, -1, downKey, downPayload, moveChild) // 右兄弟左侧挤进去了
	parent.Key[downPosition] = upKey                                        // 父位的Key
	parent.Payload[downPosition] = p.Payload[p.KeyNum-1]                    // 父位的载荷。
	// （父位不动孩子）
	//_ = EraseKeys(p, p.KeyNum-1, p.KeyNum-1) // 只是抹掉最后1位
	return
}

// LeftRotate 左旋
// @param p Position 旋转的出发节点。P(50)向上把父清的一个Key挤下来给兄弟系欸但
// 右旋，必然是①P的最右key，②挤下来父节点P.Patent右侧的Key，③下来给P右兄弟的0位，其他向后排挤，
// ④P的最右儿子，改为P右兄弟的最左[0位]孩子
// @author https://github.com/BrotherSam66/
func LeftRotate(p *btreemodels.BTreeNode) (err error) {
	fmt.Println("LeftRotate")

	if p == nil {
		err = errors.New("出错，本节点是空！")
		fmt.Println(err.Error())
		return
	}
	if p.KeyNum < 3 {
		err = errors.New("出错，本节点KeyNum小于3，借不出去！")
		fmt.Println(err.Error())
		return
	}
	parent := p.Parent // 父亲
	if parent == nil {
		err = errors.New("出错，本节点父亲是空！")
		fmt.Println(err.Error())
		return
	}
	// 查找p在父节点位置，定下来谁下来
	upKey := p.Key[0]    // 准备上去的Key值
	var downPosition int // 准备下来Key的位置
	for downPosition = 0; downPosition < parent.KeyNum; downPosition++ {
		if upKey < parent.Key[downPosition] { // 小于，说明刚刚越过了.downPosition=upKey位置-1。
			break
		}
	}
	downPosition-- // downPosition=upKey位置-1。
	if downPosition > parent.KeyNum-1 {
		err = errors.New(fmt.Sprintf("leftRotate出错，父亲找到的downPosition:%d不对！", downPosition))
		fmt.Println(err.Error())
		return
	}

	// 查找左兄弟，左兄弟满员否
	leftBrother := parent.Child[downPosition]
	if leftBrother.KeyNum >= btreeconst.M-1 {
		err = errors.New("出错，左兄弟是满员的！")
		fmt.Println(err.Error())
		return
	}

	// 开始大搬家
	leftBrother.Key[leftBrother.KeyNum] = parent.Key[downPosition]         // 准备下来，并入左兄弟的key来自父节点
	leftBrother.Payload[leftBrother.KeyNum] = parent.Payload[downPosition] // 准备下来，并入左兄弟的载荷来自父节点
	leftBrother.Child[leftBrother.KeyNum+1] = p.Child[0]                   // 准备下来，并入左兄弟的child来自左兄弟，最右手
	leftBrother.KeyNum++
	parent.Key[downPosition] = upKey            // 父位的Key
	parent.Payload[downPosition] = p.Payload[0] // 父位的载荷。
	// （父位不动孩子）
	_ = MoveKeysLeft(p, 0, -1, 0, "", nil) // 只是抹掉0位
	return
}
