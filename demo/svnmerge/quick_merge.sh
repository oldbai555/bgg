#!/bin/bash
# -*- coding: utf-8 -*-

echo "==========SVN合并开始=========="
rm -rf commit_msg_*.txt

echo "========== 1.开始合并服务端代码 =========="
svn revert C:\\bgg\\u3dv1_server\\ynjg_server_u3d\\branches\\20241226 -R
svn revert C:\\bgg\\u3dv1_server\\ynjg_server_u3d\\trunk -R
echo "========== 处理变动完成 =========="

./svnmerge.exe C:\\bgg\\u3dv1_server\\ynjg_server_u3d\\branches\\20241226 C:\\bgg\\u3dv1_server\\ynjg_server_u3d\\trunk

echo "========== 2.合并协议文件 =========="
svn revert C:\\jjyz\\ynjg_u3d\\proto -R
svn revert C:\\bgg\\proto -R
echo "========== 处理变动完成 =========="

./svnmerge.exe C:\\jjyz\\ynjg_u3d\\proto C:\\bgg\\proto

echo "==========SVN合并结束=========="
findList=$(find ./ -type f -name "commit_msg_*.txt" -printf '%T+ %p\n' | sort -r | cut -d' ' -f2-)
# 遍历findList
for file in $findList; do
    iconv.exe -f iso-8859-1 -t UTF-8//TRANSLIT ${file} -o ${file}
    echo "==========处理文件：$file 结束=========="
done

read -p "任意键关闭......" input
