// 过滤一下 option 里具体的某个 key
export function replaceOps(optList?: Array<lb.ListOption_Option>, newOpt?: lb.ListOption_Option): (Array<lb.ListOption_Option>) {
    let newOptList: Array<lb.ListOption_Option>
    // 过滤一下那个枚举
    newOptList = optList!.filter((item) => {
        return item.key !== newOpt!.key
    })
    newOptList.push(newOpt!)
    console.log("newOptList", newOptList)
    return newOptList
}