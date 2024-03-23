sed -i 's/\r//' .env
export $(xargs <.env)

#使用说明，用来提示输入参数
usage() {
  echo "Usage: sh baixctl.sh [OPTION]"
  echo "gc | genclient 根据 proto 生成客户端代码"
  echo "gs | genserver 根据 proto 生成服务端代码"
  echo "gsc | gensc 根据 proto 生成客户端 服务端代码"
  echo "gg | gengingateway 根据 proto 生成网关代码"
  echo "gt | gents 根据 proto 生成 ts 代码"
  echo "ap | addproto 新增 proto 文件"
  echo "ar | addrpc 新增 rpc 方法"
  echo "ac | addCurdRpc 新增 curd rpc 方法以及 message. 示例 sh bin.sh addCurdRpc lbbill.proto Bill,BillCategory"
  echo "asc | addCurdSysRpc 新增系统 curd rpc 方法以及 message. 示例 sh bin.sh addCurdRpc lbbill.proto Bill,BillCategory"
  exit 1
}

function genclient() {
  ./baixctl genclient -p $1
}
function genserver() {
  ./baixctl genserver -p $1
}
function gsc() {
  genclient $1
  genserver $1
}
function gengingateway() {
  ./baixctl gengingateway -p $1
}
function gents() {
  # baixctl gen_ts_vue -p lbddz -o /Users/zhangjianjun/work/lb/github.com/oldbai555/bgg/webv2/admin/
  ./baixctl gen_ts_vue -p $1 -o $2
}
function addrpc() {
  ./baixctl genAddRpc -p $1 -r $2
}
function addrpc() {
  ./baixctl genAddRpc -p $1 -r $2
}
function addCurdRpc() {
  ./baixctl genAddCurdRpc -p $1 -m $2
}
function addCurdSysRpc() {
  ./baixctl genAddCurdRpc -p $1 -m $2 -s true
}

case "$1" in
"gc" | "genclient")
  genclient "$2"
  ;;
"gs" | "genserver")
  genserver "$2"
  ;;
"gsc" | "gensc")
  gsc "$2"
  ;;
"gg" | "gengingateway")
  gengingateway "$2"
  ;;
"gt" | "gents")
  gents "$2" "$3"
  ;;
"ar" | "addrpc")
  addrpc "$2" "$3"
  ;;
"ap" | "addproto")
  addproto "$2"
  ;;
"ac" | "addCurdRpc")
  addCurdRpc "$2" "$3"
  ;;
"asc" | "addCurdSysRpc")
  addCurdSysRpc "$2" "$3"
  ;;
*)
  usage
  ;;

esac
