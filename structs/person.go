package structs

// Person - our authentication struct.
type Person struct {
	ID   string `json:"net-id"`
	Name string `json:"name"`
	// Basic struct {
	// 	Links struct {
	// 		BasicInfo struct {
	// 			Rel    string `json:"rel"`
	// 			Href   string `json:"href"`
	// 			Method string `json:"method"`
	// 		} `json:"basic__info"`
	// 		BasicModify struct {
	// 			Rel    string `json:"rel"`
	// 			Href   string `json:"href"`
	// 			Method string `json:"method"`
	// 		} `json:"basic__modify"`
	// 	} `json:"links"`
	// 	Metadata struct {
	// 		Restricted         bool `json:"restricted"`
	// 		ValidationResponse struct {
	// 			Code    int    `json:"code"`
	// 			Message string `json:"message"`
	// 		} `json:"validation_response"`
	// 		ValidationInformation []string `json:"validation_information"`
	// 	} `json:"metadata"`
	// 	UpdatedDatetime struct {
	// 		Value   time.Time `json:"value"`
	// 		APIType string    `json:"api_type"`
	// 	} `json:"updated_datetime"`
	// 	UpdatedByByuID struct {
	// 		Value       string `json:"value"`
	// 		Description string `json:"description"`
	// 		APIType     string `json:"api_type"`
	// 	} `json:"updated_by_byu_id"`
	// 	CreatedDatetime struct {
	// 		Value   time.Time `json:"value"`
	// 		APIType string    `json:"api_type"`
	// 	} `json:"created_datetime"`
	// 	CreatedByByuID struct {
	// 		Value       string `json:"value"`
	// 		Description string `json:"description"`
	// 		APIType     string `json:"api_type"`
	// 	} `json:"created_by_byu_id"`
	// 	ByuID struct {
	// 		Value   string `json:"value"`
	// 		APIType string `json:"api_type"`
	// 		Key     bool   `json:"key"`
	// 	} `json:"byu_id"`
	// 	PersonID struct {
	// 		Value   string `json:"value"`
	// 		APIType string `json:"api_type"`
	// 	} `json:"person_id"`
	// 	NetID struct {
	// 		Value           string `json:"value"`
	// 		APIType         string `json:"api_type"`
	// 		RelatedResource string `json:"related_resource"`
	// 	} `json:"net_id"`
	// 	PersonalEmailAddress struct {
	// 		Value           string `json:"value"`
	// 		APIType         string `json:"api_type"`
	// 		RelatedResource string `json:"related_resource"`
	// 	} `json:"personal_email_address"`
	// 	PrimaryPhoneNumber struct {
	// 		Value           interface{} `json:"value"`
	// 		APIType         string      `json:"api_type"`
	// 		RelatedResource string      `json:"related_resource"`
	// 	} `json:"primary_phone_number"`
	// 	Deceased struct {
	// 		Value           bool   `json:"value"`
	// 		APIType         string `json:"api_type"`
	// 		RelatedResource string `json:"related_resource"`
	// 	} `json:"deceased"`
	// 	Sex struct {
	// 		Value   string `json:"value"`
	// 		APIType string `json:"api_type"`
	// 	} `json:"sex"`
	// 	FirstName struct {
	// 		Value   string `json:"value"`
	// 		APIType string `json:"api_type"`
	// 	} `json:"first_name"`
	// 	MiddleName struct {
	// 		Value   string `json:"value"`
	// 		APIType string `json:"api_type"`
	// 	} `json:"middle_name"`
	// 	Surname struct {
	// 		Value   string `json:"value"`
	// 		APIType string `json:"api_type"`
	// 	} `json:"surname"`
	// 	Suffix struct {
	// 		Value   string `json:"value"`
	// 		APIType string `json:"api_type"`
	// 	} `json:"suffix"`
	// 	PreferredFirstName struct {
	// 		Value   string `json:"value"`
	// 		APIType string `json:"api_type"`
	// 	} `json:"preferred_first_name"`
	// 	PreferredSurname struct {
	// 		Value   string `json:"value"`
	// 		APIType string `json:"api_type"`
	// 	} `json:"preferred_surname"`
	// 	RestOfName struct {
	// 		Value   string `json:"value"`
	// 		APIType string `json:"api_type"`
	// 	} `json:"rest_of_name"`
	// 	NameLnf struct {
	// 		Value   string `json:"value"`
	// 		APIType string `json:"api_type"`
	// 	} `json:"name_lnf"`
	// 	NameFnf struct {
	// 		Value   string `json:"value"`
	// 		APIType string `json:"api_type"`
	// 	} `json:"name_fnf"`
	// 	PreferredName struct {
	// 		Value   string `json:"value"`
	// 		APIType string `json:"api_type"`
	// 	} `json:"preferred_name"`
	// 	HomeTown struct {
	// 		Value   string `json:"value"`
	// 		APIType string `json:"api_type"`
	// 	} `json:"home_town"`
	// 	HomeStateCode struct {
	// 		Value       string `json:"value"`
	// 		Description string `json:"description"`
	// 		APIType     string `json:"api_type"`
	// 	} `json:"home_state_code"`
	// 	HomeCountryCode struct {
	// 		Value       string `json:"value"`
	// 		Description string `json:"description"`
	// 		APIType     string `json:"api_type"`
	// 	} `json:"home_country_code"`
	// 	MergeInProcess struct {
	// 		Value   bool   `json:"value"`
	// 		APIType string `json:"api_type"`
	// 	} `json:"merge_in_process"`
	// } `json:"basic"`
}
