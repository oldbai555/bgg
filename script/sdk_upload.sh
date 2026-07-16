#!/usr/bin/env bash

FILE_PATH="${1}"
API_KEY=""
API_SECRET=""
BASE_URL=""

echo "======== SDK 文件上传测试 ========"
echo "BASE_URL : $BASE_URL"
echo "API_KEY  : $API_KEY"
echo "文件路径 : $FILE_PATH"
echo "================================="

set -x
curl -sS -w "\nHTTP_STATUS:%{http_code}\n" \
  -X POST "${BASE_URL}/sdk/file/upload" \
  -H "X-API-Key: ${API_KEY}" \
  -H "X-API-Secret: ${API_SECRET}" \
  -H "Accept: application/json" \
  -F "file=@${FILE_PATH}"
set +x

echo
echo "请求已完成（请根据返回的 JSON 与 HTTP_STATUS 判断是否成功）。"

