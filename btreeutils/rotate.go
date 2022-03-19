package rbtutils

import (
	"errors"
	"go-red-black-tree/global"
	"go-red-black-tree/rbtmodels"
)

// LeftRotate 左旋+变色
// @param p point 旋转的支点
/*
 *   parent                 parent
 *     |                      |
 *     p                      pr
 *    / \        ===>     	 /  \
 *  pl   pr                 p    rr
 *      /  \               / \
 *     rl   rr            pl  rl
 */
func LeftRotate(p *rbtmodels.RBTNode) (err error) {
	if p == nil {
		return errors.New("出错，本节点是空！")
	}
	if p.Right == nil {
		return errors.New("出错，本节点右儿子是空！")
	}

	parent := p.Parent // 父亲
	pSelf := p         // 本身
	pr := p.Right      // 右儿子（升级）
	rl := pr.Left      // 右儿子的左孙子（断枝重连）

	//// 下来的P ==》红 ；上去的pr ==》黑
	//p.IsRed = true
	//pr.IsRed = false
	//fmt.Printf("以%d为轴左转|", p.Key)
	// 下方需要判断 p 是否root
	if parent != nil { // p不是root，，还需要分析，p在父亲的左还是右
		if parent.Left == pSelf { // p 在父亲左手
			parent.Left = pr // 1.1 父亲左指向：
		} else { // p 在父亲右手
			parent.Right = pr // 1.2 父亲右指向：
		}
		pr.Parent = parent // 2.1 (升级的)pr上指向：
	} else { // p是root
		global.Root = pr         // 1.1 + 1.2 父亲的指向。其实是root指向pr
		global.Root.Parent = nil // 2.1 (升级的)pr上指向：
	}

	// 下方p是否root都要执行
	pr.Left = pSelf // 2.2 (升级的)pr左指向：
	// 2.3 (升级的)pr右指向：不动
	pSelf.Parent = pr // 3.1 pSelf上指向：
	// 3.2 pSelf左指向：不动
	pSelf.Right = rl // 3.3 pSelf右指向：rl
	if rl != nil {
		rl.Parent = pSelf // 4.1 rl上指向：
	}
	// 4.2 rl左指向：不动
	// 4.3 rl右指向：不动

	return err
}

// RightRotate 右旋+变色
// @param p point 旋转的支点
/*
 *     parent                 parent
 *       |                      |
 *       p                     pl
 *      / \        ===>    	  /  \
 *    pl   pr                LL   P
 *   /  \                        / \
 *  ll   lr                     lr  pr
 */
func RightRotate(p *rbtmodels.RBTNode) (err error) {
	if p == nil {
		return errors.New("出错，本节点是空！")
	}
	if p.Left == nil {
		return errors.New("出错，本节点左儿子是空！")
	}

	parent := p.Parent // 父亲
	pSelf := p         // 本身
	pl := p.Left       // 左儿子（升级）
	lr := pl.Right     // 左儿子的右孙子（断枝重连）
	//fmt.Printf("以%d为轴右转|", p.Key)
	//// 下来的P ==》红 ；上去的pl ==》黑
	//p.IsRed = true
	//pl.IsRed = false

	// 下方需要判断 p 是否root
	if parent != nil { // p不是root，，还需要分析，p在父亲的左还是右
		if parent.Left == pSelf { // p 在父亲左手
			parent.Left = pl // 1.1 父亲左指向：
		} else { // p 在父亲右手
			parent.Right = pl // 1.2 父亲右指向：
		}
		pl.Parent = parent // 2.1 (升级的)pl上指向：
	} else { // p是root
		global.Root = pl         // 1.1 + 1.2 父亲的指向。其实是root指向pr
		global.Root.Parent = nil // 2.1 (升级的)pr上指向：
	}

	// 下方p是否root都要执行
	pl.Right = pSelf // 2.2 (升级的)pl右指向：
	// 2.3 (升级的)pl左指向：不动
	pSelf.Parent = pl // 3.1 pSelf上指向：
	// 3.2 pSelf右指向：不动
	pSelf.Left = lr // 3.3 pSelf左指向：lr
	if lr != nil {
		lr.Parent = pSelf // 4.1 lr上指向：
	}
	// 4.2 lr左指向：不动
	// 4.3 lr右指向：不动

	return err
}
