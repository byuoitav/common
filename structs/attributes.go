package structs

// AttributeSet is an object that contains a set of attributes and an identifier for this set
type AttributeSet struct {
	Name       string                 `json:"name"`
	Attributes map[string]interface{} `json:"attributes"`
}
