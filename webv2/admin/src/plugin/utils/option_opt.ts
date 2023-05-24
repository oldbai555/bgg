import type {Options_Opt} from "../api/model/lb";

// 过滤一下 option 里具体的某个 key
export function replaceOps(optList: Options_Opt[], newOpt: Options_Opt): (Options_Opt[]) {
    let newOptList: Options_Opt[] = optList
    // 过滤一下那个枚举
    newOptList = optList.filter((item) => {
        return item.key !== newOpt.key
    })
    newOptList.push(newOpt)
    console.log("newOptList", newOptList)
    return newOptList
}