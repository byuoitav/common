package structs

type BulkUpdateResponse struct {
	ID      string `json:"_id"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}
