package glb

import (
	"encoding/base64"
	"fmt"
	"strings"

	"gitlab.com/gitlab-org/api/client-go"
)

type GitSubmoduleInfo struct {
	ProjectId     string `json:"id"`
	ProjectName   string `json:"name"`
	SubmodulePath string `json:"path"`
}

// 递归获取项目的所有子模块ID，支持指定分支（为空时默认 main）
func (glb *GitLabClient) GetSubmoduleRecursive(projectId, branch string) ([]GitSubmoduleInfo, error) {
	var result []GitSubmoduleInfo
	visited := make(map[string]bool)
	var recursive func(string) error

	recursive = func(pid string) error {
		if visited[pid] {
			return nil
		}
		visited[pid] = true

		ref := branch
		if ref == "" {
			ref = "main"
		}
		file, _, err := glb.Client.RepositoryFiles.GetFile(pid, ".gitmodules", &gitlab.GetFileOptions{Ref: &ref}, nil)
		if err != nil {
			// 没有.gitmodules文件则跳过
			return nil
		}
		// base64解码内容
		contentBytes, err := base64.StdEncoding.DecodeString(file.Content)
		if err != nil {
			return err
		}
		content := string(contentBytes)

		// 解析.gitmodules内容
		modules := parseGitmodules(content)
		for _, m := range modules {
			// 获取子项目ID
			project, _, err := glb.Client.Projects.GetProject(m.ProjectPath, nil)
			if err != nil {
				continue
			}
			info := GitSubmoduleInfo{
				ProjectId:     fmt.Sprintf("%v", project.ID),
				ProjectName:   project.Name,
				SubmodulePath: m.Path,
			}
			result = append(result, info)
			// 递归查找子模块
			if err := recursive(fmt.Sprintf("%v", project.ID)); err != nil {
				return err
			}
		}
		return nil
	}

	err := recursive(projectId)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// parseGitmodules 解析.gitmodules内容，返回子模块路径和项目路径
func parseGitmodules(content string) []struct{ Path, ProjectPath string } {
	var modules []struct{ Path, ProjectPath string }
	var path, url string
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "path = ") {
			path = strings.TrimPrefix(line, "path = ")
		} else if strings.HasPrefix(line, "url = ") {
			url = strings.TrimPrefix(line, "url = ")
			// 只处理gitlab项目url
			if strings.Contains(url, "/") {
				parts := strings.Split(url, "/")
				projectPath := parts[len(parts)-2] + "/" + strings.TrimSuffix(parts[len(parts)-1], ".git")
				modules = append(modules, struct{ Path, ProjectPath string }{path, projectPath})
			}
		}
	}
	return modules
}
