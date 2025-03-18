- 组网
```shell
curl -s https://install.zerotier.com | sudo bash
zerotier-cli info
zerotier-cli join 创建好的网络ID
```
- zerotier - moon 自己的服务器组网
```shell
# 去到目录下
# 服务端
cd /var/lib/zerotier-one
# 生成 Moon
zerotier-idtool initmoon identity.public >> moon.json
# 编辑 moon.json. 找到 roots -> stableEndpoints, 指定为当前云服务器的外网ip与监听端口号 例如
"stableEndpoints": ["127.0.0.1/9993"]
# 生成 .moon 文件
zerotier-idtool genmoon moon.json
# 得到一个输出结果
wrote 000000xxxxxxxxx.moon (signed world with timestamp 1742300694141)
# 创建目录 并且挪进去
mkdir moons.d
mv 000000xxxxxxxxx.moon moons.d/
# 重启服务
systemctl restart zerotier-one.service
# 查看是否配制成功
zerotier-cli listmoons
# id 就是 moonId
"id": "000000xxxxxxxxx"
# 查看自己的 ztaddr. 冒号左边的就是网络
"identity": "d8xxxxxx:xxxxxx"

# 客户端 
# 加入网络 join 
zerotier-cli join [moonId]

# 查看网络
zerotier-cli listpeers

# 找到和 ztaddr 对应的节点ID ztaddr
# 加入网络
zerotier-cli orbit [节点ID] [节点ID]  # 节点ID ztaddr d8xxxxxx

# 离开网络
zerotier-cli deorbit [节点ID] #离开某个Moon节点

# 查看网络 zerotier-cli listpeers 会显示加入成功
200 listpeers d8xxxxxx - -1 1.14.2 MOON

# 查看加入的所有网络
zerotier-cli listnetworks
```
