package structs

// BulkUpdateResponse - a collection of responses when making bulk changes to the database.
type BulkUpdateResponse struct {
	ID      string `json:"_id"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}
