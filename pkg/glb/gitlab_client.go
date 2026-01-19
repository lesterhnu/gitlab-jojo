package glb

import (
	"github.com/spf13/viper"
	"gitlab.com/gitlab-org/api/client-go"
	"go.uber.org/zap"
)

type GitLabClient struct {
	Client *gitlab.Client
	logger *zap.Logger
}

func NewGitLabClient(conf *viper.Viper, logger *zap.Logger) (*GitLabClient, error) {
	gitlabURL := conf.GetString("gitlab.base_url")
	token := conf.GetString("gitlab.token")
	client, err := gitlab.NewClient(token, gitlab.WithBaseURL(gitlabURL))
	if err != nil {
		logger.Warn("Failed to create GitLab client", zap.Error(err))
		return nil, err
	}
	return &GitLabClient{Client: client, logger: logger}, nil
}
