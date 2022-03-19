package rbtutils

import (
	"fmt"
	"go-red-black-tree/global"
	"go-red-black-tree/rbtmodels"
)

// Delete 删除节点
// @deleteNode 待删除节点
func Delete(deleteNode *rbtmodels.RBTNode) {

	// 删除唯一的根
	if deleteNode == global.Root && deleteNode.Left == nil && deleteNode.Right == nil {
		global.Root = nil
		return
	}

	/*
	 * 普通二叉树，删除的时候，用前驱节点or后继节点替换(内容，不替换颜色)，然后删掉替代节点(替代节点可能有单叶子，把这个单叶子提上来)
	 * 找前驱，如没有左子树(删除时候用不到这种情形)，就向父亲不断找，直到向左拐弯的地方。
	 * 找后继，如没有右子树(删除时候用不到这种情形)，就向父亲不断找，直到向右拐弯的地方。
	 * 找后继，如没有右子树(删除时候用不到这种情形)，就向父亲不断找，直到向右拐弯的地方。
	 * [1]找前驱or后继[1.1]找不到前驱/后置，我作为待删除。。[1.2]找到的节点替换我的(内容，不替换颜色)，找到的节点作为待删除。
	 * [2]自己能搞定，相当于234树的三or四节点，删掉低端红色；或删掉高位黑色，唯一儿子变黑上来
	 * [3]自己搞不定，相当于234树的二节点，删掉唯一黑色后，把父亲借下来，234树的兄弟(多节点)顶上去父亲位置
	 * [4]自己搞不定，相当于234树的二节点，删掉唯一黑色后，把父亲借下来，234树的兄弟(二节点)也没得借。需要递归。可能需要递归。走插入节点的双红
	 */

	// [1]找前驱or后继[1.1]找不到前驱/后置，我作为待删除。。[1.2]找到的节点替换我的(内容，不替换颜色)，找到的节点作为待删除。
	avatarNode := deleteNode
	if !(deleteNode.Left == nil || deleteNode.Right == nil) { // deleteNode有双子才需要找前序节点
		avatarNode = Predecessor(deleteNode) // 用前驱节点做替身
		//avatarNode = Successor(deleteNode) // 用后继节点做替身
	}

	if avatarNode == nil {
		fmt.Println("[1.1]找不到前驱/后置，我作为待删除:", deleteNode.Key)
		avatarNode = deleteNode
	} else { // [1.2]找到前驱/后置，找到的节点替换我的(内容，不替换颜色)，找到的节点作为待删除。
		_ = deleteNode.ReplaceInfo(avatarNode)
	}

	// 到这里，待删除节点是是avatarNode
	// [2]avatar自己能搞定，相当于234树的三or四节点，删掉低端红色；或删掉高位黑色，唯一儿子变黑上来
	// [2.1]avatar是末端，avatar是红.删除==》结束（avatar是末端我是黑是最难的，不在这里讨论）
	// [2.2]avatar是次末端，avatar和下级必然不同色，2个下级。avatar复制某个(我决定找左)下级内容,染黑(无论原来啥颜色)，删除avatar左儿子，右儿子变红(无论原来啥颜色)==》结束
	// [2.3]avatar是次末端，avatar和下级必然不同色，只有一个下级。avatar复制下级内容，染黑(无论原来啥颜色)==》结束

	if avatarNode.Left == nil && avatarNode.Right == nil && avatarNode.IsRed {
		// [2.1]avatar是末端，avatar是红.删除==》结束（avatar是末端我是黑是最难的，不在这里讨论）
		if avatarNode.Parent.Left == avatarNode {
			avatarNode.Parent.Left = nil // 上级摘钩就行，下级就会等待GC
			avatarNode.Parent = nil      // 上级摘钩就行，下级就会等待GC
		} else {
			avatarNode.Parent.Right = nil // 上级摘钩就行，下级就会等待GC
			avatarNode.Parent = nil       // 上级摘钩就行，下级就会等待GC
		}
		return
	}

	//if avatarNode.Left != nil && avatarNode.Right != nil {
	//	// [2.2]avatar是次末端，avatar和下级必然不同色，2个下级。（实战不会两个儿子的，写上吧）
	//	// avatar复制某个(我决定找左)下级内容,染黑(无论原来啥颜色)，删除avatar左儿子，右儿子变红(无论原来啥颜色)==》结束
	//	_ = avatarNode.ReplaceInfo(avatarNode.Left)
	//	avatarNode.Left.Parent = nil
	//	avatarNode.Left = nil
	//	avatarNode.IsRed = false
	//	avatarNode.Right.IsRed = true
	//	return
	//}

	// avatar必然只有0or1个儿子，
	// [2.3]avatar是次末端，avatar和下级必然不同色，只有一个下级。avatar复制下级内容，染黑(无论原来啥颜色)==》结束
	if avatarNode.Left != nil { // 唯一儿子在左
		_ = avatarNode.ReplaceInfo(avatarNode.Left)
		avatarNode.Left.Parent = nil
		avatarNode.Left = nil
		avatarNode.IsRed = false
		return
	} else if avatarNode.Right != nil { // 唯一儿子在右
		_ = avatarNode.ReplaceInfo(avatarNode.Right)
		avatarNode.Right.Parent = nil
		avatarNode.Right = nil
		avatarNode.IsRed = false
		return
	}
	// 别瞎return ，给后面[3][4]单黑无下级节点留着口子

	// [3][4]到这里，我是无下级黑节点
	var brotherNode *rbtmodels.RBTNode        // 是时候定义兄弟节点了
	if avatarNode.Parent.Left == avatarNode { // 我在爸爸左手
		brotherNode = avatarNode.Parent.Right
	} else { // 我在爸爸右手
		brotherNode = avatarNode.Parent.Left
	}
	// [3]自己搞不定，相当于234树的二节点，删掉唯一黑色后，把父亲借下来，234树的兄弟(多节点)顶上去父亲位置

	if brotherNode.Left != nil && brotherNode.Right != nil { // 兄弟有2个侄子
		/* [3.1]avatar单黑，兄弟节点有双侄子。父亲为轴向avatar旋转，把父亲借下来，兄弟顶上去父亲位置。
		*       父亲节点用原来父亲节点的颜色，兄弟一级黑色，子侄一级红色。==>结束
		*         (50X)            (50X)  |         (50X)           (50X)
		*         /                /      |         /               /
		*       (30X)            (40X)    |       (30X)           (20X)
		*       /   \            /   \    |       /   \           /   \
		*   A(20B) (40X)      (30B) (45B) |    (20X) A(40B)    (15X) A(30B)
		*          /   \         \        |    /   \                  /
		*       (35X)  (45X)     (35R)    | (15X)  (25X)           (25X)
		*         A是末端，所以兄弟下面不会再有了
		 */
		brotherNode.IsRed = avatarNode.Parent.IsRed // 先染色再旋转，父亲节点用原来父亲节点的颜色
		avatarNode.Parent.IsRed = false             // 先染色再旋转，父亲(下一步占我的位置）黑色
		if avatarNode.Parent.Left == avatarNode {   // 我在爸爸左手
			brotherNode.Left.IsRed = true
			brotherNode.Right.IsRed = false
			_ = LeftRotate(avatarNode.Parent) // 左旋转
			avatarNode.Parent.Left = nil
			avatarNode.Parent = nil
		} else { // 我在爸爸右手
			brotherNode.Left.IsRed = false
			brotherNode.Right.IsRed = true
			_ = RightRotate(avatarNode.Parent) // 右旋转
			avatarNode.Parent.Right = nil
			avatarNode.Parent = nil
		}
		return
	}

	if (brotherNode.Left != nil && brotherNode.Right == nil) || (brotherNode.Left == nil && brotherNode.Right != nil) { // 兄弟有1个侄子（没有侄子情形，属于[4]）.玛德，#golang# 居然没有逻辑运算异或。（#golang是世界上最好的语言# ）
		/* [3.2]avatar单黑，兄弟节点有单侄子。单侄子需要调成一顺。然后父亲为轴向avatar旋转，把父亲借下来，兄弟顶上去父亲位置。
		*       父亲节点用原来父亲节点的颜色，兄弟一级黑色，没有子侄了。==>结束
		*         (50X)            (50X)             (50X)  |         (50X)           (50X)        (50X)
		*         /                /                 /      |         /               /            /
		*       (30X)            (30X)             (35X)    |       (30X)           (30X)        (25X)
		*       /   \            /   \             /   \    |       /   \           /   \        /   \
		*   A(20B) (40X)     A(20B) (35X)       (30B) (40B) |    (20X) A(40B)    (25X) A(40B) (20X)  (30B)
		*          /                     \                  |        \           /
		*       (35X)                    (40X)              |        (25X)    (20X)
		*         A是末端，所以兄弟下面不会再有了
		 */
		if avatarNode.Parent.Left == avatarNode { // [3.2.1]兄弟有1个侄子，我在父亲左手
			if brotherNode.Left != nil { // 侄子在兄弟的左手(不顺)，斩不定色，兄弟为轴右转
				_ = RightRotate(brotherNode)
			}
			// 父亲、兄弟、侄子顺了，先定色，父亲为轴右转。==》结束
			avatarNode.Parent.Right.IsRed = avatarNode.Parent.IsRed // 原兄弟，继承父亲颜色
			avatarNode.Parent.IsRed = false                         // 原父亲，黑色
			avatarNode.Parent.Right.Right.IsRed = false             // 原侄子，黑色
			_ = LeftRotate(avatarNode.Parent)
			avatarNode.Parent.Left = nil
			avatarNode.Parent = nil
			return
		} else { // [3.2.2]兄弟有1个侄子，我在父亲右手
			if brotherNode.Right != nil { // 侄子在兄弟的右手(不顺)，斩不定色，兄弟为轴左转
				_ = LeftRotate(brotherNode)
			}
			// 父亲、兄弟、侄子顺了，先定色，父亲为轴左转。==》结束
			avatarNode.Parent.Left.IsRed = avatarNode.Parent.IsRed // 原兄弟，继承父亲颜色
			avatarNode.Parent.IsRed = false                        // 原父亲，黑色
			avatarNode.Parent.Left.Left.IsRed = false              // 原侄子，黑色
			_ = RightRotate(avatarNode.Parent)
			avatarNode.Parent.Right = nil
			avatarNode.Parent = nil
			return
		}

		//return
	}

	// [4]自己搞不定，相当于234树的二节点，234树的兄弟(二节点)也没得借，删掉唯一黑色后，兄弟变红
	// [4.1]父红==》变黑==》结束。可能需要递归。走插入节点的双红
	// [4.2]父黑==》父亲作为传参，走递归的删除后调整 FixAfterDelete()
	if avatarNode.Left == nil &&
		avatarNode.Right == nil &&
		!avatarNode.IsRed &&
		brotherNode.Left == nil &&
		brotherNode.Right == nil {
		/* [4]自己搞不定，相当于234树的二节点，删掉唯一黑色后，把父亲借下来，234树的兄弟(二节点)也没得借。
		*       父亲节点用原来父亲节点的颜色，兄弟一级黑色，没有子侄了。==>结束
		*
		*               (70X)
		*               /
		*            (50X)
		*         /         \
		*      (30X)       (60X)
		*     /   \         /   \
		* A(20B) (40B)  A(55X) (65X)
		*
		 */
		if avatarNode.Parent.IsRed { // 父亲红
			avatarNode.Parent.IsRed = false           // 父亲黑
			if avatarNode.Parent.Left == avatarNode { // avatar 在左
				avatarNode.Parent.Right.IsRed = true // 兄弟红
				avatarNode.Parent.Left = nil         // avatar 删除
			} else { // avatar 在左
				avatarNode.Parent.Left.IsRed = true // 兄弟红
				avatarNode.Parent.Right = nil       // avatar 删除
			}
			return
		} else { // 父亲黑
			if avatarNode.Parent.Left == avatarNode { // avatar 在左
				avatarNode.Parent.Right.IsRed = true // 兄弟红
				avatarNode = avatarNode.Parent       // 转换一下防失联
				avatarNode.Left = nil                // 实际的avatar 删除
			} else { // avatar 在左
				avatarNode.Parent.Left.IsRed = true // 兄弟红
				avatarNode = avatarNode.Parent      // 转换一下防失联
				avatarNode.Right = nil              // 实际的avatar 删除
			}
			FixAfterDelete(avatarNode) // 开始递归调整(实际上是父亲)
		}
	}

	if avatarNode == nil {
		fmt.Println("删除操作本节点内无法处理，又找不到替换节点，按说不应该出现这个提示的")
		return
	}
	fmt.Println()
	return
}

// FixAfterDelete  删除后，定义“双黑节点”，找到“红”卸载一个黑出去，可能递归
// @avatar black+black双黑属性节点
func FixAfterDelete(avatar *rbtmodels.RBTNode) {
	/* avatar 实际是双黑节点，就是前面
	* [3]avatar 为黑色，这时候是递归，avatar和兄弟都可能有多级子节点
	* [3.1]avatar 黑兄红侄
	* [3.2]avatar 黑兄黑侄红父
	* [3.3]avatar 黑兄黑侄黑父
	* [3.4]avatar 红兄（黑侄黑父）
	 */
	//fmt.Printf("进入FixAfterDelete avatar==%d", avatar.Key)
	ShowTreeColor(global.Root)
	//err := errors.New("出错，本节点是空！")
	// [3.0] 双黑节点是根==》结束???不对吧，还是不平衡啊
	if avatar == global.Root {
		return
	}

	isAvatarOnRight := avatar.Parent.Right == avatar // avatar 在右边
	parant := avatar.Parent                          // 父亲
	brother := parant.Right                          // 定义胸怀兄弟
	if isAvatarOnRight {
		brother = parant.Left
	}
	brothersLeftSon := brother.Left
	var brothersLeftSonIsRed bool
	if brothersLeftSon == nil {
		brothersLeftSonIsRed = true
	} else {
		brothersLeftSonIsRed = brothersLeftSon.IsRed
	}
	brothersRightSon := brother.Right
	var brothersRightSonIsRed bool
	if brothersRightSon == nil {
		brothersRightSonIsRed = true
	} else {
		brothersRightSonIsRed = brothersRightSon.IsRed
	}

	// [3.1]avatar黑兄红侄
	/* [3.1.1]avatar在右，黑兄红侄.侄子左红右随意。==》原兄升父位继承父的颜色/原父下来变黑(填补双黑值)/原红侄黑==》父轴，右旋，==》结束
	*  [3.1.1.J]镜像，avatar在左
	*          (50Y)       |        (30Y)           ||        (50Y)           |          (70Y)       |
	*        /       \     |      /      \          ||      /      \          |        /       \     |
	*     (30B)     A(60B) |   (20B)     (50B)      ||  A(40B)     (70B)      |     (50B)      (80B) |
	*     /    \     / \   |   / \      /   \       ||   / \      /   \       |     /    \     / \   |
	*  (20R)  (40X)(?5)(?6)|(?1)(?2) (40X)  (60B)   ||(?1)(?2) (60X)  (80R)   |  (40R)  (60X)(?5)(?6)|
	*   / \    / \         |         / \      / \   ||         / \      / \   |   / \    / \         |
	*(?1)(?2)(?3)(?4)      |       (?3)(?4) (?5)(?6)||       (?3)(?4) (?5)(?6)|(?1)(?2)(?3)(?4)      |
	 */
	if (isAvatarOnRight && !brother.IsRed && brothersLeftSonIsRed) || // 在右
		(!isAvatarOnRight && !brother.IsRed && brothersRightSonIsRed) { // 镜像
		brother.IsRed = parant.IsRed // 原兄升父位继承父的颜色
		parant.IsRed = false         // 原父下来变黑(填补双黑值)
		if isAvatarOnRight {         // 在右
			brothersLeftSon.IsRed = false // 原红侄黑
			_ = RightRotate(parant)
		} else { // 镜像
			brothersRightSon.IsRed = false // 原红侄黑
			_ = LeftRotate(parant)
		}
		return // 结束
	}

	/* [3.1.2]avatar在右，黑兄红侄.侄子左黑右红。==>原父下来变黑(填补双黑值)/原红侄升父位继承父的颜色==》兄轴左旋==》父轴，右旋，==》结束
	*  [3.1.2.J]镜像avatar在左，黑兄红侄.侄子右黑左红。==>原父下来变黑(填补双黑值)/原红侄升父位继承父的颜色==》兄轴右旋==》父轴，左旋，==》结束
	*          (50Y)       |          (50Y)      |         (40Y)        ||       (50Y)           |       (50Y)         |         (60Y)        |
	*        /       \     |        /      \     |       /      \       ||     /      \          |     /     \         |       /      \       |
	*     (30B)     A(60B) |     (40R)    A(60B) |    (30B)     (50B)   || A(40B)     (70B)      | A(40B)   (60R)      |    (5B)     (70B) c  |
	*     /    \     / \   |     /    \    / \   |    / \      /   \    ||  / \       /    \     |  /  \    / \        |    / \      /   \    |
	*  (20B)  (40R)(?5)(?6)|  (30B)  (?4)(?5)(?6)| (20B)(?3) (?4) (60B) ||(?1)(?2) (60R)  (80B)  |(?1)(?2)(?3)(70B)    | (40B)(?3) (?4) (80B) |
	*   / \    / \         |   / \               |  / \           /  \  ||          / \    / \   |             / \     |  / \           /  \  |
	*(?1)(?2)(?3)(?4)      |(20B)(?3)            |(?1)(?2)      (?5)(?6)||       (?3)(?4)(?5)(?6)|          (?4)(80B)  |(?1)(?2)      (?5)(?6)|
	*                      |  / \                |                      ||                       |               / \   |                      |
	*                      |(?1)(?2)             |                      ||                       |             (?5)(?6)|                      |
	 */
	if (isAvatarOnRight && !brother.IsRed && !brothersLeftSonIsRed && brothersRightSonIsRed) || // 在右
		(!isAvatarOnRight && !brother.IsRed && brothersLeftSonIsRed && !brothersRightSonIsRed) { // 镜像
		parant.IsRed = false // 原父下来变黑(填补双黑值)
		if isAvatarOnRight { // 在右
			brothersRightSon.IsRed = parant.IsRed // 原红侄升父位继承父的颜色
			_ = LeftRotate(brother)               // 兄轴左旋
			_ = RightRotate(parant)               // 父轴，右旋
		} else { // 镜像
			brothersLeftSon.IsRed = parant.IsRed // 原红侄升父位继承父的颜色
			_ = RightRotate(brother)             // 镜像
			_ = LeftRotate(parant)               // 镜像
		}
		return // 结束
	}

	/* [3.2]avatar在左右都行，黑兄黑侄红父。==》父黑兄红==》结束
	* 对应2-3-4树删除操作中兄弟节点为2节点，父节点至少是个3节点，父节点key下移与兄弟节点合并。
	* 黑兄红侄.侄子左红右随意。==》父轴，右旋，父下来变黑(填补双黑值)/升父位的兄继承父的颜色/左侄黑==》结束
	*          (50R)       |          (50B)       |
	*        /       \     |        /       \     |
	*     (30B)     A(60B) |     (30R)     A(60B) |
	*     /    \     / \   |     /    \     / \   |
	*  (20B)  (40B)(?5)(?6)|  (20B)  (40B)(?5)(?6)|
	*   / \    / \         |   / \    / \         |
	*(?1)(?2)(?3)(?4)      |(?1)(?2)(?3)(?4)      |
	 */
	if parant.IsRed && !brother.IsRed && !brothersLeftSonIsRed && !brothersRightSonIsRed {
		parant.IsRed = false // 父黑
		brother.IsRed = true // 兄红
		return               // 结束
	}

	/* [3.3]avatar在左右都行，黑兄黑侄黑父。==》兄红==》父亲做为avatar==》递归
	*          (50B)       |         A(50B)       ||        (30B)           |       A(30B)           |
	*        /       \     |        /       \     ||      /      \          |      /      \          |
	*     (30B)     A(60B) |     (30R)      (60B) ||  A(20B)     (50B)      |   (20B)     (50R)      |
	*     /    \     / \   |     /    \     / \   ||   / \      /   \       |   / \      /   \       |
	*  (20B)  (40B)(?5)(?6)|  (20B)  (40B)(?5)(?6)||(?1)(?2) (40B)  (60B)   |(?1)(?2) (40B)  (60B)   |
	*   / \    / \         |   / \    / \         ||         / \      / \   |         / \      / \   |
	*(?1)(?2)(?3)(?4)      |(?1)(?2)(?3)(?4)      ||       (?3)(?4) (?5)(?6)|       (?3)(?4) (?5)(?6)|
	 */
	if !parant.IsRed && !brother.IsRed && !brothersLeftSonIsRed && !brothersRightSonIsRed {
		brother.IsRed = true   // 兄红
		FixAfterDelete(parant) // 父亲作为双黑，递归
		return                 // 结束
	}

	/* [3.4]avatar在右，红兄黑侄黑父。==>原兄黑，原父红==》父轴右旋==》avatar的新兄弟就是黑的==》用avatar递归（红父黑兄[3.2]）
	*  [3.4J]镜像avatar在左，红兄黑侄黑父。==>原兄黑，原父红==》父轴左旋==》avatar的新兄弟就是黑的==》用avatar递归（红父黑兄[3.2]）
	*  [注意]就地讨论avatar侄子是不是存在，否则递归下去默认侄子是全的
	*          (50B)       |        (30B)           ||        (50B)           |          (70B)       |
	*        /       \     |      /      \          ||      /      \          |        /       \     |
	*     (30R)     A(60B) |   (20B)     (50R)      ||  A(40B)     (70R)      |     (50R)      (80B) |
	*     /    \     / \   |   / \      /   \       ||   / \      /   \       |     /    \     / \   |
	*  (20B)  (40B)(?5)(?6)|(?1)(?2) (40B) A(60B)   ||(?1)(?2) (60B)  (80B)   | A(40B)  (60B)(?5)(?6)|
	*   / \    / \         |         / \      / \   ||         / \      / \   |   / \    / \         |
	*(?1)(?2)(?3)(?4)      |       (?3)(?4) (?5)(?6)||       (?3)(?4) (?5)(?6)|(?1)(?2)(?3)(?4)      |
	 */
	if !parant.IsRed && brother.IsRed && !brothersLeftSonIsRed && !brothersRightSonIsRed {
		brother.IsRed = false // 原父下来变黑(填补双黑值)
		parant.IsRed = true   // 原父下来变黑(填补双黑值)
		if isAvatarOnRight {  // 在右
			_ = RightRotate(parant) // 父轴，右旋
			// avatar依然在右
			parant = avatar.Parent // 父亲
			brother = parant.Left  // 定义兄弟
		} else { // 镜像
			_ = LeftRotate(parant) // 镜像
			// avatar依然在左
			parant = avatar.Parent // 父亲
			brother = parant.Right // 定义兄弟
		}
		// 这时，红父黑兄，侄子有没有都两说
		if brother.Left != nil && brother.Right != nil { // 有双侄子
			FixAfterDelete(avatar) // 有双侄子 avatar作为双黑，递归
			return
		}
		// 没侄子，单侄子，都不可能，avatar双黑，配不上
		fmt.Println("[3.4]没侄子，单侄子，都不可能，avatar双黑，配不上")

		return // 结束
	}
}

// Predecessor 找前驱节点。比我稍小的最大节点
func Predecessor(n *rbtmodels.RBTNode) (ret *rbtmodels.RBTNode) {
	if n == nil {
		fmt.Println("纳尼？让我给一个nil找前驱节点？")
		return nil
	}
	if n.Left != nil { // 有左儿子，找左儿子下面的最大
		ret = n.Left
		for {
			if ret.Right == nil { // 右下nil就算找到了
				return ret
			}
			ret = ret.Right // 沿着右手一直向下
		}
	} else { // 没有左儿子，向上找，每个和我比较，到root，返nil
		if n.Parent == nil { // 没左儿子，自己又是root
			fmt.Println("这个真没有，没左儿子，自己又是root")
			return nil
		}
		ret = n
		for {
			if ret.Parent == nil { // 游标移到本身是root，还没有找到，就nil了
				return nil
			} else {
				if ret.Parent.Key < n.Key { // 某个父辈小于我了，这个就是
					return ret.Parent
				}
			}
			ret = ret.Parent // 向上一直找
		}
	}

}

// Successor 找后继节点。比我稍大的最小节点
func Successor(n *rbtmodels.RBTNode) (ret *rbtmodels.RBTNode) {
	if n == nil {
		fmt.Println("纳尼？让我给一个nil找后继节点？")
		return nil
	}
	if n.Right != nil { // 有右儿子，找右儿子下面的最大
		ret = n.Right
		for {
			if ret.Left == nil { // 左下nil就算找到了
				return ret
			}
			ret = ret.Left // 沿着左手一直向下
		}
	} else { // 没有右儿子，向上找，每个和我比较，到root，返nil
		if n.Parent == nil { // 没右儿子，自己又是root
			fmt.Println("这个真没有，没右儿子，自己又是root")
			return nil
		}
		ret = n
		for {
			if ret.Parent == nil { // 游标移到本身是root，还没有找到，就nil了
				return nil
			} else {
				if ret.Parent.Key > n.Key { // 某个父辈小于我了，这个就是
					return ret.Parent
				}
			}
			ret = ret.Parent // 向上一直找
		}
	}

}
