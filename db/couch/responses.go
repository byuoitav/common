package couch

import (
	"github.com/byuoitav/common/state/statedefinition"
	"github.com/byuoitav/common/structs"
)

type building struct {
	Rev string `json:"_rev,omitempty"`
	*structs.Building
}

type buildingQueryResponse struct {
	Docs     []building `json:"docs"`
	Bookmark string     `json:"bookmark"`
	Warning  string     `json:"warning"`
}

type room struct {
	Rev string `json:"_rev,omitempty"`
	*structs.Room
}

type roomQueryResponse struct {
	Docs     []room `json:"docs"`
	Bookmark string `json:"bookmark"`
	Warning  string `json:"warning"`
}

type roomConfiguration struct {
	Rev string `json:"_rev,omitempty"`
	*structs.RoomConfiguration
}

type roomConfigurationQueryResponse struct {
	Docs     []roomConfiguration `json:"docs"`
	Bookmark string              `json:"bookmark"`
	Warning  string              `json:"warning"`
}

type device struct {
	Rev string `json:"_rev,omitempty"`
	*structs.Device
}

type deviceQueryResponse struct {
	Docs     []device `json:"docs"`
	Bookmark string   `json:"bookmark"`
	Warning  string   `json:"warning"`
}

type deviceStateQueryResponse struct {
	Docs     []statedefinition.StaticDevice `json:"docs"`
	Bookmark string                         `json:"bookmark"`
	Warning  string                         `json:"warning"`
}

type deviceType struct {
	Rev string `json:"_rev,omitempty"`
	*structs.DeviceType
}

type deviceTypeQueryResponse struct {
	Docs     []deviceType `json:"docs"`
	Bookmark string       `json:"bookmark"`
	Warning  string       `json:"warning"`
}

type uiconfig struct {
	Rev string `json:"_rev,omitempty"`
	*structs.UIConfig
}

type uiconfigQueryResponse struct {
	Docs     []uiconfig `json:"docs"`
	Bookmark string     `json:"bookmark"`
	Warning  string     `json:"warning"`
}

type icons struct {
	Rev      string   `json:"_rev,omitempty"`
	IconList []string `json:"Icons"`
}

type deviceRoles struct {
	Rev      string         `json:"_rev,omitempty"`
	RoleList []structs.Role `json:"roles"`
}

type roomDesignations struct {
	Rev       string   `json:"_rev,omitempty"`
	DesigList []string `json:"designations"`
}

type closureCodes struct {
	Rev   string   `json:"_rev,omitempty"`
	Codes []string `json:"closure_codes"`
}

type tags struct {
	Rev     string   `json:"_rev,omitempty"`
	TagList []string `json:"tags"`
}

type template struct {
	Rev string `json:"_rev,omitempty"`
	*structs.Template
}

type templateQueryResponse struct {
	Docs     []template `json:"docs"`
	Bookmark string     `json:"bookmark"`
	Warning  string     `json:"warning"`
}

type menu struct {
	Rev   string   `json:"_rev,omitempty"`
	Order []string `json:"order"`
}

type attributeGroup struct {
	Rev string `json:"_rev,omitempty"`
	structs.Group
}

type attributeQueryResponse struct {
	Docs     []attributeGroup `json:"docs"`
	Bookmark string           `json:"bookmark"`
	Warning  string           `json:"warning"`
}

type roomAttachmentResponse struct {
	ID           string                 `json:"_id"`
	Rev          string                 `json:"_rev,omitempty"`
	LinkedImages []string               `json:"linkedImages"`
	Attachments  map[string]interface{} `json:"_attachments"`
}

// type jobs struct {
// Rev string `json:"_rev,omitempty"`
// *structs.Jobs
// }

// type jobsQueryResponse struct {
// 	Docs     []template `json:"docs"`
// 	Bookmark string     `json:"bookmark"`
// 	Warning  string     `json:"warning"`
// }
