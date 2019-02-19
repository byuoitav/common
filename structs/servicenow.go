package structs

type ServiceNowLinkValue struct {
	Link  string `json:"link"`
	Value string `json:"value"`
}

type ResolutionCategories struct {
	Result []Category `json:"result"`
}

type Category struct {
	UAction string `json:"u_action"`
	Hint    string `json:"u_attribute_help_text"`
}

type QueriedUsers struct {
	Result []Users `json:"result"`
}

type Users struct {
	NetID string `json:"user_name,omitempty"`
	Email string `json:"email,omitempty"`
	Phone string `json:"home_phone,omitempty"`
}
