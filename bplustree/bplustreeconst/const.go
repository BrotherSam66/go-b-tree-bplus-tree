package bplustreeconst

// M M阶B+树
var M int = 5

// Min 每个节点至少有的成员/关键字个数（M一半向下舍入）
var Min int = (M + 1) / 2

// KeyLen 彩色显示树，每个KEY字节长度 todo 根据输入的数字最大值，动态调整这个
var KeyLen int = 2

// MaxKey 最大key值
var MaxKey int = 100
