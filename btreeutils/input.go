package btreeutils

import (
	"errors"
	"fmt"
	"go-b-tree-bplus-tree/btreemodels"
	"go-b-tree-bplus-tree/global"

	"math/rand"
)

// Inputs 连续插入节点
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

			ShowTreeColor(global.Root)
			continue
		}
		if key > 99 || key < 1 {
			fmt.Println("必须是0~~99")
			continue
		}
		Insert(key, "")
		ShowTreeColor(global.Root)
	}
}

// Insert 加入节点
// @key 插入的键值
// @payload 插入的载荷值
func Insert(key int, payload string) {
	if global.Root == nil { // 原树为空树，新加入的转为根
		global.Root = btreemodels.NewBTreeNode(nil, 1, key, payload)
		return
	}

	// 从root开始查找附加的位置
	tempParent := global.Root // 临时的父亲，移动的指针
	// 搜索新加元素应该的位置

	var isToLeft bool // 新加节点在tempParent的左儿子吗？
	for {
		if tempParent.Key > key { // 新来数值小，向左搜索
			if tempParent.Left == nil { // 左为空，左就是new位置，跳出循环
				isToLeft = true
				break
			}
			tempParent = tempParent.Left
		} else if tempParent.Key < key { // 新来数值大，向右搜索
			if tempParent.Right == nil { // 右为空，右就是new位置，跳出循环
				isToLeft = false
				break
			}
			tempParent = tempParent.Right
		} else { // 相等，就更新标签，完成任务退出
			tempParent.Label = label
			return
		}
	}

	// 找到位置了，开始拼装。global.NewUpNode是拟增加的节点（也可能是下级旋转上升上来的随机色节点）
	global.NewUpNode = rbtmodels.NewRBTNode(true, key, label, tempParent, nil, nil)
	if isToLeft { // 拼装在左儿子
		tempParent.Left = global.NewUpNode
	} else { // 拼装在右儿子
		tempParent.Right = global.NewUpNode
	}
	ShowTreeColor(global.Root)
	FixAfterInsert() // 拼装后，要调整，包括旋转+变色，可能递归

	return
}

// FixAfterInsert  拼装后，要调整，包括旋转+变色，可能递归
// global.NewUpNode是拟增加的节点（也可能是下级旋转上升上来的随机色节点）
func FixAfterInsert() {
	err := errors.New("出错，本节点是空！")

	// [1]新加节点or上升上来的节点是root，改黑==》结束
	if global.NewUpNode == global.Root {
		global.Root.IsRed = false
		return
	}

	// [2]（二三四树原来有1个节点），新加一个红，上黑下红，不变
	if global.NewUpNode.Parent.IsRed == false { // 新加节点or上升上来的 的父亲黑色就不用旋转 ==》结束
		return
	}

	// [3.1] 父红，叔红(不能空，空算黑)， ==》爷红，父叔黑，爷爷变为当前节点 ==》递归
	/*    gB            gR
	 *   /  \          /  \
	 * flR  urR  ==> flB  urB
	 *   \             \
	 *   srR           srR
	 */
	if global.NewUpNode.Parent.Parent.Left != nil && global.NewUpNode.Parent.Parent.Right != nil { // 确保有叔叔
		if global.NewUpNode.Parent.Parent.Left.IsRed && global.NewUpNode.Parent.Parent.Right.IsRed {
			global.NewUpNode.Parent.Parent.Left.IsRed = false  // 父叔黑
			global.NewUpNode.Parent.Parent.Right.IsRed = false // 父叔黑
			global.NewUpNode.Parent.Parent.IsRed = true        // 爷红
			global.NewUpNode = global.NewUpNode.Parent.Parent  // 爷爷变为当前节点 ==》递归
			FixAfterInsert()                                   // 递归
			return
		}
	}

	// 到这里，叔叔必然黑(或空)
	if global.NewUpNode.Parent.Parent.Left == global.NewUpNode.Parent { // [4.1] 父在爷左手
		// [4.1] 父在爷左手
		// [4.1.1] 父flR红，叔黑(空也算黑)，我在右， ==》以父flR为P左旋，原父flR做当前系欸但 ==》递归
		// [4.1.2] 父srR红，叔黑(空也算黑)，我在左(其实[4.1.1]递归过来就是这个)， ==》父黑爷红，以爷爷gB为P右旋，
		//       ==》新爷爷变红，父叔边喝，爷爷作为新节点递归
		/*   gB                   gB              srR                 srR
		 *   /  \     flR左旋     /  \   gB右旋    /   \  父黑爷孙红    /   \
		 *  flR urB    ==>    srR   urB  ==>  flR     gB  ==>      flB   gB
		 *   \                /                        \                  \
		 *   srR            flR                        urB                urB
		 */
		if global.NewUpNode.Parent.Right == global.NewUpNode { // [2.1.1]我在爸爸右手，flR左旋
			err = LeftRotate(global.NewUpNode.Parent)
			if err != nil {
				fmt.Println(err)
			}
			global.NewUpNode = global.NewUpNode.Left // 模拟新加的基准点，向左下移一下，递归
			FixAfterInsert()                         // 递归
			return
		}
		// 到这里一定是[4.1.2] 父srR红，叔黑(空也算黑)，我在左(其实[4.1.1]递归过来就是这个)， ==》父黑爷红，以爷爷gB为P右旋，
		//		==》原爷爷的右手做当前节点，黑就结束，红就递归原
		err = RightRotate(global.NewUpNode.Parent.Parent) // 以爷爷gB为P右旋。
		if err != nil {
			fmt.Println(err)
		}
		global.NewUpNode.IsRed = false             // =》我flR变黑
		global.NewUpNode = global.NewUpNode.Parent // 我的父亲(新的爷爷是红)，作为新节点递归
		FixAfterInsert()                           // 递归
		return
	} else { // [4.2] 父在爷右手
		// [4.2] 父在爷左手
		// [4.2.1] 父frR红，叔黑(空也算黑)，我在左， ==》以父frR为P右旋，原父frR做当前系欸但 ==》递归
		// [4.2.2] 父srR红，叔黑(空也算黑)，我在左(其实[4.1.1]递归过来就是这个)， ==》父黑爷红，以爷爷gB为P右旋，
		//       ==》新爷爷变红，父叔边喝，爷爷作为新节点递归

		/* [2.2.1]右三，爷右左，黑红红=》父亲支点右旋=》爷右右，黑红红
		 * [2.2.2]右三，爷右右，黑红红=》爷爷支点左旋=》上黑两下红。
		 *
		 *    gB             gB             slR                  slR
		 *   / \   frR右旋   /  \    gB左旋  /   \   父黑爷孙红    /    \
		 * ulB  frR  ==>   ulB slR   ==>  gB     frR   ==>     gB     frB
		 *     /                \        /
		 *    slR               frR    ulB
		 */
		if global.NewUpNode.Parent.Left == global.NewUpNode { // [4.2.1] 父frR红，叔黑(空也算黑)，我在左，
			err = RightRotate(global.NewUpNode.Parent)
			if err != nil {
				fmt.Println(err)
			}
			global.NewUpNode = global.NewUpNode.Right // 模拟新加的基准点，向右下移一下
			FixAfterInsert()                          // 递归
			return
		}
		// 到这里一定是[4.2.2] 父slR红，叔黑(空也算黑)，我在右(其实[4.2.1]递归过来就是这个)， ==》父黑爷红，以爷爷gB为P左旋，
		//		==》原爷爷的右手做当前节点，黑就结束，红就递归原
		err = LeftRotate(global.NewUpNode.Parent.Parent) // 以爷爷gB为P左旋，
		if err != nil {
			fmt.Println(err)
		}

		global.NewUpNode.IsRed = false             // =》我frR变黑
		global.NewUpNode = global.NewUpNode.Parent // 我的父亲(新的爷爷是红)，作为新节点递归
		FixAfterInsert()                           // 递归
		return
	}
}
