package glb

import gitlab "gitlab.com/gitlab-org/api/client-go"

func (glb *GitLabClient) CreateTag(projectId string, tagName, ref string) error {
	_, _, err := glb.Client.Tags.CreateTag(projectId, &gitlab.CreateTagOptions{
		TagName: &tagName,
		Ref:     &ref,
	}, nil)
	return err
}
