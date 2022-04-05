package bplustree

import (
	"errors"
	"fmt"
	"go-b-tree-bplus-tree/bplustree/bplustreeconst"
	"go-b-tree-bplus-tree/bplustree/bplustreemodels"
)

// BorrowFromLeft 向左兄弟借
// @param p Position 旋转的出发节点。P(50)向上把父清的一个Key挤下来给兄弟系欸但
// 右旋，必然是①P的最右key，②挤下来父节点P.Patent右侧的Key，③下来给P右兄弟的0位，其他向后排挤，
// ④P的最右儿子，改为P右兄弟的最左[0位]孩子
// @author https://github.com/BrotherSam66/
/*
 *假设：5阶，最大5个KEY、最小3个KEY，60-->移到右边
 *  (20  |     60           |      80|)      |  (20  |     50           |      80|)      |
 *  /           \                    \       |  /           \                    \       |
 *(?1)(30 | 40 | 50   |60)        (70|80)    |(?1)(30 | 40 | 50 )   (60   |   70|80)     |
 *    /      \    \      \         /    \    |    /      \    \       \       /    \     |
 *(21|30)(31|40)(41|50)(51|60) (61|70)(71|80)|(21|30)(31|40)(41|50) (51|60)(61|70)(71|80)|
 */
func BorrowFromLeft(brother, avatar *bplustreemodels.BPTreeNode) (err error) {
	//fmt.Println("BorrowFromLeft")
	if brother == nil {
		err = errors.New("出错，本节点是空！")
		fmt.Println(err.Error())
		return
	}
	if len(brother.Key) < bplustreeconst.Min {
		err = errors.New("出错，兄弟节点KeyNum小于Min，借不出去！")
		fmt.Println(err.Error())
		return
	}
	parent := brother.Parent // 父亲
	if parent == nil {
		err = errors.New("出错，本节点父亲是空！")
		fmt.Println(err.Error())
		return
	}

	// 计算 avatarPosition int avatar 在父亲的排位
	_, avatarPosition := FindKeyPosition(&parent.Key, avatar.Key[len(avatar.Key)-1])
	if avatarPosition < 0 {
		fmt.Println("不正常，BorrowFromLeft FindKeyPosition 准确查到叶子，不应该走到这里", avatar.Key[len(avatar.Key)-1])
	}

	// 开始大搬家
	if len(avatar.Child) > 0 && avatar.Child[0] != nil { // 被分裂的是分支，如果是末梢/叶子，就没有child
		// avatar 接收 key+payload+child
		avatar.Key = append([]int{brother.Key[len(brother.Key)-1]}, avatar.Key...)                                 // key
		avatar.Child = append([]*bplustreemodels.BPTreeNode{brother.Child[len(brother.Child)-1]}, avatar.Child...) // Child
		// 如果是非叶子，下级的上指
		brother.Child[len(brother.Key)-1].Parent = avatar
		// brother 剪除一个 key+payload+child
		brother.Key = brother.Key[:len(brother.Key)-1]
		brother.Child = brother.Child[:len(brother.Child)-1]
	} else { // 叶子的。叶子才有payload
		// avatar 接收 key+payload+child
		avatar.Key = append([]int{brother.Key[len(brother.Key)-1]}, avatar.Key...)                    // key
		avatar.Payload = append([]string{brother.Payload[len(brother.Payload)-1]}, avatar.Payload...) // Child
		// brother 剪除一个 key+payload+child
		brother.Key = brother.Key[:len(brother.Key)-1]
		brother.Payload = brother.Payload[:len(brother.Payload)-1]
	}
	// parent的key换一下
	parent.Key[avatarPosition-1] = brother.Key[len(brother.Key)-1] // 左边 上指 key值(不可能是父亲的最右，不用递归)

	return
}

// BorrowFromRight 向右兄弟借
// @param p Position 旋转的出发节点。P(50)向上把父清的一个Key挤下来给兄弟系欸但
// 右旋，必然是①P的最右key，②挤下来父节点P.Patent右侧的Key，③下来给P右兄弟的0位，其他向后排挤，
// ④P的最右儿子，改为P右兄弟的最左[0位]孩子
// @author https://github.com/BrotherSam66/
func BorrowFromRight(brother, avatar *bplustreemodels.BPTreeNode) (err error) {
	//fmt.Println("BorrowFromLeft")
	if brother == nil {
		err = errors.New("出错，本节点是空！")
		fmt.Println(err.Error())
		return
	}
	if len(brother.Key) < bplustreeconst.Min {
		err = errors.New("出错，兄弟节点KeyNum小于Min，借不出去！")
		fmt.Println(err.Error())
		return
	}
	parent := brother.Parent // 父亲
	if parent == nil {
		err = errors.New("出错，本节点父亲是空！")
		fmt.Println(err.Error())
		return
	}

	// 计算 avatarPosition int avatar 在父亲的排位
	_, avatarPosition := FindKeyPosition(&parent.Key, avatar.Key[len(avatar.Key)-1])
	if avatarPosition < 0 {
		fmt.Println("不正常，BorrowFromRight FindKeyPosition 准确查到叶子，不应该走到这里", avatar.Key[len(avatar.Key)-1])
	}

	// 开始大搬家
	if len(avatar.Child) > 0 && avatar.Child[0] != nil { // 被分裂的是分支，如果是末梢/叶子，就没有child
		// avatar 接收 key+payload+child
		avatar.Key = append(avatar.Key, brother.Key[0])       // key
		avatar.Child = append(avatar.Child, brother.Child[0]) // Child
		// 如果是非叶子，下级的上指
		brother.Child[0].Parent = avatar
		// brother 剪除一个 key+payload+child
		brother.Key = brother.Key[1:]
		brother.Child = brother.Child[1:]
	} else { // 叶子的。叶子才有payload
		// avatar 接收 key+payload+child
		avatar.Key = append(avatar.Key, brother.Key[0])             // key
		avatar.Payload = append(avatar.Payload, brother.Payload[0]) // Child
		// brother 剪除一个 key+payload+child
		brother.Key = brother.Key[1:]
		brother.Payload = brother.Payload[1:]
	}
	// parent的key换一下
	parent.Key[avatarPosition] = avatar.Key[len(avatar.Key)-1] // 左边就是我 上指 key值(不可能是父亲的最右，不用递归)
	return
}
