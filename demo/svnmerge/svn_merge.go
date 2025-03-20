package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("请提供要Merge的文件路径")
		return
	}
	branch := os.Args[1]
	branch = filepath.ToSlash(branch)
	trunk := os.Args[2]
	trunk = filepath.ToSlash(trunk)

	tempCommitMessageFile := fmt.Sprintf("commit_message_%d.txt", time.Now().Unix())
	// 获取分支的所有提交记录
	commits, err := getCommits(branch, trunk)
	if err != nil {
		fmt.Println("Error getting commits:", err)
		return
	}
	if len(commits) == 0 {
		fmt.Println("Error neet merge commits:", err)
		return
	}
	// 打开临时文件以记录所有提交信息
	file, err := os.Create(tempCommitMessageFile)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	// 更新主干目录
	err = updateTrunk(trunk)
	if err != nil {
		fmt.Println("Error updating trunk:", err)
		return
	}
	// 获取提交信息并写入临时文件
	for _, commit := range commits {
		message, err := getCommitMessage(commit, branch)
		if err != nil {
			fmt.Printf("Error getting commit message for %s: %v\n", commit, err)
			return
		}
		_, err = writer.WriteString(message + "\n")
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
	}
	writer.Flush()

	for _, commit := range commits {
		// 先测试 merge
		err = testMergeCommit(commit, branch, trunk)
		if err != nil {
			fmt.Printf("Error merging commit %s: %v\n", commit, err)
			return
		}
	}

	for _, commit := range commits {
		// Merge 每个提交到主干
		err = mergeCommit(commit, branch, trunk)
		if err != nil {
			fmt.Printf("Error merging commit %s: %v\n", commit, err)
			return
		}
	}

	// 提交合并后的更改
	err = commitMergedChanges(tempCommitMessageFile, trunk)
	if err != nil {
		fmt.Println("Error committing merged changes:", err)
		return
	}

	fmt.Println("Merge and commit completed successfully.")
}

// 获取分支的所有提交记录
func getCommits(source, target string) ([]string, error) {
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
		if strings.HasPrefix(line, "C    ") {
			parts := strings.Fields(line)
			if len(parts) > 0 {
				return fmt.Errorf("%s", string(output))
			}
		}
	}
	return nil
}
func mergeCommit(commit, branch, trunk string) error {
	cmd := exec.Command("svn", "merge", "-c", commit, branch, trunk)
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

// 获取指定提交的提交信息
func getCommitMessage(commit, branch string) (string, error) {
	cmd := exec.Command("svn", "log", "-r", commit, branch)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	// 确保行尾符一致
	message := strings.ReplaceAll(string(output), "\r\n", "\n")
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
