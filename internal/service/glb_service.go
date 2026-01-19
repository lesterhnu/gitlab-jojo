package service

import (
	"ops/pkg/glb"
)

type GitlabService struct {
	*Service
	GitlabClient *glb.GitLabClient
}

func NewGitlabService(
	service *Service,
	gitlabClient *glb.GitLabClient,
) *GitlabService {
	return &GitlabService{
		Service:      service,
		GitlabClient: gitlabClient,
	}
}
