const tokenKey = "sid"


export function setToken(token: string | undefined) {
    // 写入本地缓存
    localStorage.setItem(tokenKey, String(token))
}

export function clearAllCaches() {
    localStorage.clear()
}

export function getToken(): string | null {
    return localStorage.getItem(tokenKey)
}