/**
 * @Author: zjj
 * @Date: 2025/4/9
 * @Desc:
**/

package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func DoMerge(trunk, branch string, commit bool) error {
	tempCommitMessageFile := fmt.Sprintf("commit_msg_%d.txt", time.Now().Unix())
	// 获取分支的所有提交记录
	commits, err := getCommits(branch, trunk)
	if err != nil {
		fmt.Println("Error getting commits:", err)
		return err
	}
	if len(commits) == 0 {
		fmt.Println("Error neet merge commits:", err)
		return err
	}
	// 更新主干目录
	err = updateTrunk(trunk)
	if err != nil {
		fmt.Println("Error updating trunk:", err)
		return err
	}

	// 打开临时文件以记录所有提交信息
	file, err := os.Create(tempCommitMessageFile)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return err
	}
	defer file.Close()
	writer := bufio.NewWriter(file)

	// 获取提交信息并写入临时文件
	for _, commit := range commits {
		message, err := getCommitMessage(commit, branch)
		if err != nil {
			fmt.Printf("Error getting commit message for %s: %v\n", commit, err)
			return err
		}
		_, err = writer.WriteString(message + "\n")
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return err
		}
	}
	writer.Flush()

	for _, commit := range commits {
		// 先测试 merge
		err = testMergeCommit(commit, branch, trunk)
		if err != nil {
			fmt.Printf("Error test merging commit %s: %v\n", commit, err)
		}
		// merge 提交到主干
		err = mergeCommit(commit, branch, trunk)
		if err != nil {
			fmt.Printf("Error merging commit %s: %v\n", commit, err)
			return err
		}
	}

	// 提交合并后的更改
	if commit {
		err = commitMergedChanges(tempCommitMessageFile, trunk)
		if err != nil {
			fmt.Println("Error committing merged changes:", err)
			return err
		}
	}

	fmt.Println("Merge and commit completed successfully.")
	return nil
}

// 获取分支的所有提交记录
func getCommits(source, target string) ([]string, error) {
	// C:\\bgg\\u3dv1_server\\ynjg_server_u3d\\branches\\20241226 C:\\bgg\\u3dv1_server\\ynjg_server_u3d\\trunk
	cmd := exec.Command("svn", "mergeinfo", source, target, "--show-revs", "eligible")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(cmd.String())
		fmt.Println(string(output))
		return nil, err
	}

	var commits []string
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "r") {
			parts := strings.Fields(line)
			if len(parts) > 0 {
				commit := strings.TrimPrefix(parts[0], "r")
				commits = append(commits, commit)
			}
		}
	}

	// 对提交记录进行从小到大的排序
	sort.Strings(commits)

	return commits, nil
}

// 将指定提交merge到主干
func testMergeCommit(commit, branch, trunk string) error {
	cmd := exec.Command("svn", "merge", "--dry-run", "-c", commit, branch, trunk)
	fmt.Printf("run cmd: %s\n", cmd.String())
	output, err := cmd.CombinedOutput()
	if err != nil {
		// 检查输出中是否包含冲突信息
		if strings.Contains(string(output), "Conflicts") || strings.Contains(string(output), "Tree conflict") {
			return fmt.Errorf("Test Merge conflict detected for commit %s: %v\n%s", commit, err, string(output))
		}
		return fmt.Errorf("%s: %v", string(output), err)
	}
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()
		// 检查输出中是否包含冲突信息
		if strings.Contains(strings.ToLower(string(output)), strings.ToLower("Conflicts")) || strings.Contains(strings.ToLower(string(output)), strings.ToLower("Tree conflict")) {
			return fmt.Errorf("%s", string(output))
		}
		if strings.HasPrefix(line, "Error") {
			parts := strings.Fields(line)
			if len(parts) > 0 {
				return fmt.Errorf("%s", string(output))
			}
		}
	}
	return nil
}
func mergeCommit(commit, branch, trunk string) error {
	cmd := exec.Command("svn", "merge", "--accept", "theirs-full", "-c", commit, branch, trunk)
	output, err := cmd.CombinedOutput()
	if err != nil {
		// 检查输出中是否包含冲突信息
		if strings.Contains(string(output), "Conflicts") || strings.Contains(string(output), "Tree conflict") {
			return fmt.Errorf("Merge conflict detected for commit %s: %v\n%s", commit, err, string(output))
		}
		return fmt.Errorf("%s: %v", string(output), err)
	}
	return nil
}

// 提交合并后的更改
func commitMergedChanges(messageFile, trunk string) error {
	cmd := exec.Command("svn", "commit", "-F", messageFile, trunk)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %v", string(output), err)
	}
	return nil
}

// 新增函数获取相对URL
func getRelativeURL(path string) (string, error) {
	cmd := exec.Command("svn", "info", "--show-item", "relative-url", path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("获取相对URL失败: %v\n%s", err, string(output))
	}
	return strings.TrimLeft(strings.TrimSpace(string(output)), "^/"), nil
}

// 获取指定提交的提交信息
func getCommitMessage(commit, branch string) (string, error) {
	relativeURL, err := getRelativeURL(branch)
	if err != nil {
		return "", err
	}

	cmd := exec.Command("svn", "log", "--xml", "-r", commit, branch)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	// 定义XML结构体
	type LogEntry struct {
		Revision string `xml:"revision,attr"`
		Author   string `xml:"author"`
		Date     string `xml:"date"`
		Msg      string `xml:"msg"`
	}

	type Log struct {
		Entries []LogEntry `xml:"logentry"`
	}

	var log Log
	if err := xml.Unmarshal(output, &log); err != nil {
		return "", fmt.Errorf("XML解析失败: %v", err)
	}

	if len(log.Entries) == 0 {
		return "", fmt.Errorf("未找到修订版本%s的日志信息", commit)
	}

	entry := log.Entries[0]
	// 构建规范的提交信息格式
	message := fmt.Sprintf("Merged revision(s) %s from %s:\n%s\n........", entry.Revision, relativeURL, strings.TrimSpace(entry.Msg))
	return message, nil
}

// 更新主干目录
func updateTrunk(trunk string) error {
	cmd := exec.Command("svn", "update", trunk)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %v", string(output), err)
	}
	return nil
}
