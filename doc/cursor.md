- 魔法使用问答
- ctrl + shift + p => 输入 open user settings (json) 写入以下内容
```
{
    "git.autofetch": true,
    "workbench.editor.enablePreview": false,
    "http.proxy": "http://127.0.0.1:7890",
    "http.proxyStrictSSL": false,
    "http.proxySupport": "override",
    "http.noProxy": [],
    "cursor.general.disableHttp2": true
}
```
