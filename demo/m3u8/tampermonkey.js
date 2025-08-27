// ==UserScript==
// @name         Bilibili Float Ball Title Panel with M3U8 Parser & Collector (Pair Mode + Toast)
// @namespace    http://tampermonkey.net/
// @version      2.0
// @description  B站悬浮球：支持拖拽、展开面板、动态解析m3u8并生成 采集/播放 按钮成对排列，带toast提示~
// @author       You
// @match        *://*.bilibili.com/*
// @grant        none
// ==/UserScript==

(function () {
    'use strict';

    // ================== Toast 系统 ==================
    const toastContainer = document.createElement("div");
    toastContainer.style.position = "fixed";
    toastContainer.style.bottom = "20px";
    toastContainer.style.right = "20px";
    toastContainer.style.display = "flex";
    toastContainer.style.flexDirection = "column";
    toastContainer.style.gap = "10px";
    toastContainer.style.zIndex = "100000";
    document.body.appendChild(toastContainer);

    function showToast(message, success = true) {
        const toast = document.createElement("div");
        toast.innerText = message;
        toast.style.padding = "10px 16px";
        toast.style.borderRadius = "6px";
        toast.style.color = "#fff";
        toast.style.fontSize = "14px";
        toast.style.background = success ? "rgba(0,180,90,0.9)" : "rgba(220,0,60,0.9)";
        toast.style.boxShadow = "0 2px 6px rgba(0,0,0,0.3)";
        toast.style.opacity = "0";
        toast.style.transform = "translateY(20px)";
        toast.style.transition = "all 0.3s ease";

        toastContainer.appendChild(toast);

        // 动画进入
        requestAnimationFrame(() => {
            toast.style.opacity = "1";
            toast.style.transform = "translateY(0)";
        });

        // 自动消失
        setTimeout(() => {
            toast.style.opacity = "0";
            toast.style.transform = "translateY(20px)";
            setTimeout(() => toast.remove(), 300);
        }, 3000);
    }

    // ================= 悬浮球 =================
    const ball = document.createElement("div");
    ball.style.position = "fixed";
    ball.style.top = "100px";
    ball.style.right = "20px";
    ball.style.width = "50px";
    ball.style.height = "50px";
    ball.style.borderRadius = "50%";
    ball.style.background = "linear-gradient(135deg, #ff6699, #ff3366)";
    ball.style.boxShadow = "0 4px 12px rgba(0,0,0,0.3)";
    ball.style.cursor = "grab";
    ball.style.zIndex = "99999";
    ball.style.display = "flex";
    ball.style.alignItems = "center";
    ball.style.justifyContent = "center";
    ball.style.color = "#fff";
    ball.style.fontSize = "20px";
    ball.innerText = "🎬";

    // ================= 悬浮面板 =================
    const panel = document.createElement("div");
    panel.style.position = "fixed";
    panel.style.minWidth = "280px";
    panel.style.background = "rgba(0,0,0,0.85)";
    panel.style.color = "#fff";
    panel.style.padding = "12px";
    panel.style.borderRadius = "8px";
    panel.style.boxShadow = "0 4px 12px rgba(0,0,0,0.3)";
    panel.style.zIndex = "99999";
    panel.style.display = "none"; // 默认隐藏
    panel.style.fontSize = "14px";

    // 按钮容器 (竖排)
    const btnContainer = document.createElement("div");
    btnContainer.style.display = "flex";
    btnContainer.style.flexDirection = "column"; // 竖排按钮
    btnContainer.style.gap = "8px";
    panel.appendChild(btnContainer);

    function createBtn(text, handler, bg = "#ff6699", hoverBg = "#ff3366") {
        const btn = document.createElement("button");
        btn.innerText = text;
        btn.style.padding = "8px";
        btn.style.border = "none";
        btn.style.borderRadius = "5px";
        btn.style.cursor = "pointer";
        btn.style.background = bg;
        btn.style.color = "#fff";
        btn.style.fontSize = "13px";
        btn.style.textAlign = "center";
        btn.addEventListener("mouseenter", () => (btn.style.background = hoverBg));
        btn.addEventListener("mouseleave", () => (btn.style.background = bg));
        btn.onclick = handler;
        return btn;
    }

    document.body.appendChild(ball);
    document.body.appendChild(panel);

    // ================= 点击展开/收起（带自适应位置） =================
    let isPanelVisible = false;
    let isDragging = false;

    ball.addEventListener("click", () => {
        if (isDragging) return;
        isPanelVisible = !isPanelVisible;

        if (isPanelVisible) {
            const rect = ball.getBoundingClientRect();
            const screenWidth = window.innerWidth;

            panel.style.top = rect.top + "px";

            if (rect.right + 300 > screenWidth) {
                panel.style.left = rect.left - 300 + "px";
                panel.style.right = "auto";
            } else {
                panel.style.left = rect.right + 10 + "px";
                panel.style.right = "auto";
            }

            panel.style.display = "block";
        } else {
            panel.style.display = "none";
        }
    });

    // ================= 拖拽功能 =================
    ball.addEventListener("mousedown", (e) => {
        e.preventDefault();
        ball.style.cursor = "grabbing";

        let shiftX = e.clientX - ball.getBoundingClientRect().left;
        let shiftY = e.clientY - ball.getBoundingClientRect().top;
        isDragging = false;

        function moveAt(pageX, pageY) {
            ball.style.left = pageX - shiftX + "px";
            ball.style.top = pageY - shiftY + "px";
            ball.style.right = "auto";
        }

        function onMouseMove(e) {
            isDragging = true;
            moveAt(e.pageX, e.pageY);
        }

        document.addEventListener("mousemove", onMouseMove);

        document.addEventListener("mouseup", () => {
            document.removeEventListener("mousemove", onMouseMove);
            ball.style.cursor = "grab";
            setTimeout(() => {
                isDragging = false;
            }, 50);
        }, {once: true});
    });

    ball.ondragstart = () => false;

    // ================= m3u8 解析逻辑 =================
    (function parseM3u8AndInjectButtons() {
        console.log("开始解析 m3u8 文件...")
        const filename = document.querySelector('.order-first .mt-4 h1')?.textContent.trim() || "unknown";
        const videoDoc = document.querySelector('.order-first');
        if (!videoDoc) return;

        const prefix = 'https://surrit.com/';
        const suffix = '/playlist.m3u8';

        const nodeValue = document.evaluate(
            '/html/body/script[5]/text()',
            document,
            null,
            XPathResult.FIRST_ORDERED_NODE_TYPE,
            null
        ).singleNodeValue?.textContent;

        if (!nodeValue) return;

        const index = nodeValue.indexOf('seek');
        if (index !== -1 && index - 32 >= 0) {
            const first32Chars = nodeValue.substring(index - 38, index - 2);
            const url = prefix + first32Chars + suffix;
            console.log("解析到的 uuid 文件地址:", first32Chars)
            fetch(url)
                .then(resp => resp.text())
                .then(text => {
                    const lines = text.split('\n');
                    lines.forEach(line => {
                        if (line.trim() && !line.startsWith('#')) {
                            const fileInfo = {
                                filename: filename,
                                url: prefix + first32Chars + '/' + line.trim(),
                                uuid: first32Chars,
                                god_num:  "",
                            };

                            // 按钮1：采集
                            const collectBtn = createBtn(`采集 ${fileInfo.filename}`, () => {
                                const payload = {
                                    player_url: fileInfo.url,
                                    name: fileInfo.filename,
                                    god_num: "",
                                    uuid: first32Chars,
                                };
                                fetch("https://oldbai.top/m3u8/video/add", {
                                    method: "POST",
                                    headers: {"Content-Type": "application/json"},
                                    body: JSON.stringify(payload)
                                }).then(res => res.text())
                                    .then(() => {
                                        showToast(`📡 已采集: ${fileInfo.filename}`, true);
                                    })
                                    .catch(err => {
                                        console.error("采集失败:", err);
                                        showToast("❌ 采集失败", false);
                                    });
                            }, "#228be6", "#1864ab");

                            // 按钮2：播放
                            const playBtn = createBtn(`播放 ${fileInfo.filename}`, () => {
                                const targetUrl = `https://oldbai.top/onlinem3u8?m3u8Url=${encodeURIComponent(fileInfo.url)}&proxyUrl=https://oldbai.top/m3u8`;
                                window.open(targetUrl, "_blank");
                            }, "#ff6699", "#ff3366");

                            // 交替添加
                            btnContainer.appendChild(collectBtn);
                            btnContainer.appendChild(playBtn);
                        }
                    });
                });
        }
    })();

})();
