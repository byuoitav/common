package structs

// Resource is the base level object returned by the UAPI.
type Resource struct {
	Links    map[string]Link `json:"links,omitempty"`
	Metadata Metadata        `json:"metadata,omitempty"`
	Basic    SubResource     `json:"basic,omitempty"`
	State    SubResource     `json:"av_state,omitempty"`
	Config   SubResource     `json:"av_config,omitempty"`
}

// SubResource is an object that helps to comprise a Resource.
type SubResource struct {
	Links        map[string]Link `json:"links,omitempty"`
	Metadata     Metadata        `json:"metadata,omitempty"`
	Building     Property        `json:"building,omitempty"`
	Room         Property        `json:"room,omitempty"`
	Displays     Property        `json:"displays,omitempty"`
	AudioDevices Property        `json:"audio_devices,omitempty"`
}

// Link contains information about accessing the Resource.
type Link struct {
	Rel    string `json:"rel,omitempty"`
	Href   string `json:"href,omitempty"`
	Method string `json:"method,omitempty"`
}

// Metadata contains high level metadata about the Resource or SubResource
type Metadata struct {
	ValidationResponse    ValidationResponse `json:"validation_response,omitempty"`
	ValidationInformation []string           `json:"validation_information,omitempty"`
	Cache                 Cache              `json:"cache,omitempty"`
	Restricted            *bool              `json:"restricted,omitempty"`
	FieldSetsReturned     []string           `json:"field_sets_returned,omitempty"`
	FieldSetsAvailable    []string           `json:"field_sets_available,omitempty"`
	FieldSetsDefault      []string           `json:"field_sets_default,omitempty"`
}

// ValidationResponse has information about the request.
type ValidationResponse struct {
	Code    *int   `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

// Cache contains a DateTime about the Resource if it was cached.
type Cache struct {
	DateTime string `json:"date_time,omitempty"`
}

// Property is an attribute of a Resource or SubResource.
type Property struct {
	Type            string `json:"api_type,omitempty"`
	Key             bool   `json:"key,omitempty"`
	Value           string `json:"value,omitempty"`
	Description     string `json:"description,omitempty"`
	DisplayLabel    string `json:"display_label,omitempty"`
	Domain          string `json:"domain,omitempty"`
	LongDescription string `json:"long_description,omitempty"`
	RelatedResource string `json:"related_resource,omitempty"`
}
