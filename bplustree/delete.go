// Package bplustree
// @Title B树工具包
// @Description  和删除节点有关的操作
// @Author  https://github.com/BrotherSam66/
// @Update
package bplustree

//
//// Deletes 连续删除节点
//// @author https://github.com/BrotherSam66/
//func Deletes() {
//
//	for {
//		var key int
//		fmt.Println("请输入KEY，按回车键(空按回车随机,-1退出)：")
//		_, _ = fmt.Scanln(&key)
//
//		if key == -1 {
//			return
//		}
//		if key == 0 {
//			key = rand.Intn(btreeglobal.MaxKey)
//			fmt.Println(key)
//		}
//
//		if key > 99 || key < 1 {
//			fmt.Println("必须是0~~99")
//			continue
//		}
//		Delete(key)
//		ShowTree(bplustreeglobal.Root)
//	}
//}
//
//// Delete 删除节点
//// @key 键值
//// @author https://github.com/BrotherSam66/
//func Delete(key int) {
//
//	// 从root开始查找附加的位置
//	tempNode, isTarget, err := Search(key)
//	if err != nil {
//		fmt.Println("查找错误，error == ", err)
//		return
//	}
//	if !isTarget {
//		fmt.Println("没找到！ ")
//		return
//	}
//
//	// 非叶子，查找可替换的叶子节点的KEY值，交换。前序or后继均可，优先前序，前序节点数量<=btreeconst.Min不容易删除就定死用后继
//	// 查到key在tempNode准确位置，deletePosition
//	deletePosition := 0
//	for deletePosition = 0; deletePosition < tempNode.KeyNum; deletePosition++ {
//		if tempNode.Key[deletePosition] == key { // 准确命中，只可能是新创建节点情形
//			break
//		}
//	}
//	if deletePosition >= tempNode.KeyNum {
//		fmt.Println("发生某种错误，找到KEY又不存在了！ ")
//		return
//	}
//	avatarNode := tempNode
//	if tempNode.Child[0] != nil { // 不是叶子
//		avatarNode, _ = PredecessorOrSuccessor(tempNode, key, true) // 用前驱节点做替身
//		// 考察avatarNode可简易删除？
//		if avatarNode.KeyNum <= btreeconst.Min { // 不能简易删除
//			avatarNode, _ = PredecessorOrSuccessor(tempNode, key, false) // 用后继节点做替身
//			// 删除数据 对换 后继（一个神奇的语句）
//			tempNode.Key[deletePosition], avatarNode.Key[0] = avatarNode.Key[0], tempNode.Key[deletePosition]
//			tempNode.Payload[deletePosition], avatarNode.Payload[0] = avatarNode.Payload[0], tempNode.Payload[deletePosition]
//			deletePosition = 0
//		} else {
//			// 删除数据 对换 前驱（一个神奇的语句）
//			tempNode.Key[deletePosition], avatarNode.Key[avatarNode.KeyNum-1] = avatarNode.Key[avatarNode.KeyNum-1], tempNode.Key[deletePosition]
//			tempNode.Payload[deletePosition], avatarNode.Payload[avatarNode.KeyNum-1] = avatarNode.Payload[avatarNode.KeyNum-1], tempNode.Payload[deletePosition]
//			deletePosition = avatarNode.KeyNum - 1
//		}
//	}
//
//	// 到这里，KEY在叶子上，就开始删除的递归流程
//
//	_ = DeleteOneKey(avatarNode, key, deletePosition)
//
//	return
//}
//
//// DeleteOneKey 删除一个叶子上的KEY，可能要递归
//// @avatar 节点
//// @key 拟删除键值
//// @deletePosition 拟删除键值的位置
//// @author https://github.com/BrotherSam66/
//func DeleteOneKey(avatar *btreemodels.BTreeNode, key int, deletePosition int) (err error) {
//	if avatar.Key[deletePosition] != key {
//		err = errors.New("奇怪啊，指定的位置deletePosition键值不吻合啊")
//		fmt.Println("奇怪啊，指定的位置deletePosition键值不吻合啊")
//		return
//	}
//
//	// 删除掉这个key
//	_ = MoveKeysLeft(avatar, deletePosition, -1, 0, "", nil)
//
//	// 检查合法性，可能要递归
//	if avatar.KeyNum < btreeconst.Min && avatar.Parent != nil { // avatar节点过短 && 不是root，需要调整，可能递归
//		_ = FixAfterDelete(avatar)
//	}
//
//	if avatar.KeyNum == 0 && avatar.Parent == nil { // avatar节点过短 && 不是root，需要调整，可能递归
//		btreeglobal.Root = nil
//	}
//
//	return
//}
//
//// EraseKeys 抹除部分KEY，必须是右侧的 todo 主要是分裂的时候用
//// @n 节点
//// @leftPosition 左面端点
//// @rightPosition 右面端点，-1表示最右
//// @author https://github.com/BrotherSam66/
//func EraseKeys(n *btreemodels.BTreeNode, leftPosition int, rightPosition int) (err error) {
//	if n == nil {
//		err = errors.New("出错，n是nil！")
//		fmt.Println(err.Error())
//		return
//	}
//	if leftPosition <= 0 {
//		err = errors.New("出错，leftPosition必须是>0！")
//		fmt.Println(err.Error())
//		return
//	}
//	if rightPosition < leftPosition {
//		err = errors.New("出错，rightPosition < leftPosition")
//		fmt.Println(err.Error())
//		return
//	}
//
//	for i := leftPosition; i <= rightPosition; i++ {
//		n.Key[i] = 0
//		n.Payload[i] = ""
//		n.Child[i+1] = nil
//	}
//
//	n.KeyNum = n.KeyNum - 1 - rightPosition + leftPosition
//	return
//}
//
//// FixAfterDelete 删除后调整
//// @avatar 递归的节点
//// @author https://github.com/BrotherSam66/
//func FixAfterDelete(avatar *btreemodels.BTreeNode) (err error) {
//	// 如果该节点递归、上升到了root，结束
//	if avatar.Parent == nil {
//		btreeglobal.Root = avatar
//
//		return
//	}
//	// 2）该结点key个数大于等于Math.ceil(m/2)-1，结束删除操作，否则执行第3步。
//	if avatar.KeyNum >= btreeconst.Min || avatar.Parent == nil {
//		return
//	}
//	// 3）如果兄弟结点key个数大于Math.ceil(m/2)-1，则父结点中的key下移到该结点，兄弟结点中的一个key上移，删除操作结束。
//	// 找出avatar的左右兄弟
//	leftBrother := avatar  // 临时定义
//	rightBrother := avatar // 临时定义
//	parent := avatar.Parent
//	// 找到 avatar 在父亲的排位
//	avatarPosition := 0
//	for avatarPosition = 0; avatarPosition < parent.KeyNum; avatarPosition++ {
//		if avatar.Key[0] < parent.Key[avatarPosition] { // 小于，说明刚刚越过了，(用avatar任何Key都行)
//			break
//		}
//	}
//
//	// 找到兄弟后直接借KEY
//	isSuccess := false
//	if avatarPosition == 0 { // 在最左
//		rightBrother = parent.Child[1]
//		isSuccess, _ = TryBorrowBrotherKey(rightBrother, false)
//		if isSuccess {
//			return
//		}
//	} else if avatarPosition >= parent.KeyNum { // 在最右
//		leftBrother = parent.Child[avatarPosition-1]
//		isSuccess, _ = TryBorrowBrotherKey(leftBrother, true)
//		if isSuccess {
//			return
//		}
//	} else { // 居中，有左右2个兄弟
//		rightBrother = parent.Child[avatarPosition+1]
//		isSuccess, _ = TryBorrowBrotherKey(rightBrother, false)
//		if isSuccess {
//			return
//		}
//		leftBrother = parent.Child[avatarPosition-1]
//		isSuccess, _ = TryBorrowBrotherKey(leftBrother, true)
//		if isSuccess {
//			return
//		}
//	}
//
//	// 到这里，就是兄弟借不来。将父结点中的key下移与当前结点及它的兄弟结点中的key合并，形成一个新的结点。
//	// 原父结点中的key的两个孩子指针就变成了一个孩子指针，指向这个新结点。然后当前结点的指针指向父结点，重复上第2步。
//	/*
//	 *假设：5阶，最大4个KEY、最小2个KEY，
//	 *  (20|60             |              80|nil)|  (20|50             |              80|nil)|
//	 *  /   \              \                \    |  /   \              \                \    |
//	 *(?1)(30|40|nil|nil) (70|nil|nil|nil)  (?3)  |(?1)(30|40|nil|nil) (60|70|nil|nil)  (?3)  |
//	 *
//	 *向父亲借(60)形成(30|40|60|70)，父亲指向20的右腿，(20)去递归
//	 */
//	if avatarPosition == 0 { // 在最左,只能用右兄弟
//		rightBrother = parent.Child[1]
//		_ = Merge3Nodes(avatar, parent, rightBrother, avatarPosition) // 三个节点合并
//	} else { // 优先用左兄弟
//		leftBrother = parent.Child[avatarPosition-1]
//		_ = Merge3Nodes(leftBrother, parent, avatar, avatarPosition-1) // 三个节点合并
//	}
//	if parent.KeyNum == 0 && parent.Parent == nil {
//		btreeglobal.Root = avatar
//		return
//	}
//	_ = FixAfterDelete(parent) // 递归了
//	return
//}
//
//// TryBorrowBrotherKey 尝试向兄弟借KEY，只是判断能不能
//// @avatar 本节点
//// @brother 兄弟节点
//// @isLeftBrother 左兄弟or右兄弟
//// @isSuccess 借节点成功了吗？
//// @author https://github.com/BrotherSam66/
//func TryBorrowBrotherKey(brother *btreemodels.BTreeNode, isLeftBrother bool) (isSuccess bool, err error) {
//	if brother.KeyNum <= btreeconst.Min { // 兄弟太短，没得借
//		return // 不算error，isSuccess=false就可
//	}
//	// 3）如果兄弟结点key个数大于Math.ceil(m/2)-1，则父结点中的key下移到该结点，兄弟结点中的一个key上移，删除操作结束。
//	/*
//	 *假设：5阶，最大4个KEY、最小2个KEY，
//	 *  (20|60             |              80|nil)|  (20|50             |              80|nil)|
//	 *  /   \              \                \    |  /   \              \                \    |
//	 *(?1)(30|40|50|nil) (70|nil|nil|nil)  (?3)  |(?1)(30|40|nil|nil) (60|70|nil|nil)  (?3)  |
//	 *
//	 *(70)右边刚删掉(75)，(60)下来并入(70)，(50)上去，填补(60)
//	 */
//	if isLeftBrother { // 是企图向左兄弟借
//		_ = RightRotate(brother) // 右转，就算完整借完
//	} else { // 是企图向右兄弟借
//		_ = LeftRotate(brother) // 左转，就算完整借完
//	}
//	isSuccess = true
//	return
//}
//
///*
// *    B树的删除操作 https://www.cnblogs.com/nullzx/p/8729425.html
// *    删除操作是指，根据key删除记录，如果B树中的记录中不存对应key的记录，则删除失败。
// * 1）如果当前需要删除的key位于非叶子结点上，则用后继key（这里的后继key均指后继记录的意思）覆盖要删除的key，
// * 然后在后继key所在的子支中删除该后继key。此时后继key一定位于叶子结点上，这个过程和二叉搜索树删除结点的方式类似。删除这个记录后执行第2步
// * 2）该结点key个数大于等于Math.ceil(m/2)-1，结束删除操作，否则执行第3步。
// * 3）如果兄弟结点key个数大于Math.ceil(m/2)-1，则父结点中的key下移到该结点，兄弟结点中的一个key上移，删除操作结束。
// *    否则，将父结点中的key下移与当前结点及它的兄弟结点中的key合并，形成一个新的结点。原父结点中的key的两个孩子指针就变成了一个孩子指针，
// * 指向这个新结点。然后当前结点的指针指向父结点，重复上第2步。
// *    有些结点它可能即有左兄弟，又有右兄弟，那么我们任意选择一个兄弟结点进行操作即可。
// */
