package v1

type GetProjectSubmodulesReq struct {
	ProjectID int64  `json:"project_id" binding:"required,min=1"`
	Ref       string `json:"ref" binding:"omitempty"`
}
