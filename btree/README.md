# go-b-tree-bplus-tree
####B Tree and B+ Tree implementation for go
####B树和B+树。欢迎引用，欢迎star。
####引用方法：把btree，或bplustree子目录完整地copy到您的项目里，然后进行您希望的修改。
```
I Insert插入数据
S Show完整的树
F Find查找数据
D Delete删除数据
E Exit退出
请输入指令，按回车键：

展示树：(父|KEY数)左腿/KEY-PAYLOAD\右腿|KEY-PAYLOAD\右腿|...
 , (nil|1)3/9-9\12|
 , (9|2)1/3-3\4|6-6\7| , (9|3)10/12-12\13|15-15\16|18-18\19|
 , (3|2)nil/1-1\nil|2-2\nil| , (3|2)nil/4-4\nil|5-5\nil| , (3|2)nil/7-7\nil|8-8\nil| , (12|2)nil/10-10\nil|11-11\nil| , (12|2)nil/13-13\nil|14-14\nil| , (12|2)nil/16-16\nil|17-17\nil| , (12|2)nil/19-19\nil|20-20\nil|
```