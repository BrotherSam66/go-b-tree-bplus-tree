package global

import "go-red-black-tree/rbtmodels"

var Name string
var Root *rbtmodels.RBTNode      // RBT根
var NewUpNode *rbtmodels.RBTNode // RBT刚升上去的节点
var KeyLen int = 2               // 彩色显示树，每个KEY字节长度 todo 根据输入的数字最大值，动态调整这个
var MaxKey int = 100             // 最大key值
