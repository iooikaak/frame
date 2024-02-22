#项目介绍
该项目实现了基于流量的负载均衡，返回指定算法的node，均衡算法分别有:平均，随机，加权

#快速开始
##实例化需要进行负载均衡的节点列表
var (
mockNodeList = NodeList{
&Node{
Id:         1,
IP:         "127.0.0.1",
Port:       8080,
Weight:     10,
},
&Node{
Id:         2,
IP:         "127.0.0.1",
Port:       8081,
Weight:     10,
},
&Node{
Id:         3,
IP:         "127.0.0.1",
Port:       8082,
Weight:     10,
},
}
)

##选择指定的算法，获取到指定节点
var defaultBanlce = "random" //average 平均算法，random 随机算法， weight 权重算法
b := NewBalance(mockNodeList, defaultBanlce)
defer node.AddRequestNum(-1)
n := b.Get(context.Background())
