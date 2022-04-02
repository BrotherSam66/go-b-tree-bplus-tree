// Package btreeconst
// @Title B树全局常量
// @Description  全局常量
// @Author  https://github.com/BrotherSam66/
// @Update
package btreeconst

// M M阶B树
const M int = 5

// Min 每个节点至少有的成员/关键字个数（M一半向上舍入-1，6:最大6个孩子/最大5个成员/最小2个成员；5:最大5个孩子/最大4个成员/最小2个成员。）
const Min int = (M+1)/2 - 1
