package couch

import (
	"encoding/json"
	"fmt"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/structs"
)

// TEMPLATES

// GetAllTemplates returns a list of all the room Templates in the database.
func (c *CouchDB) GetAllTemplates() ([]structs.Template, error) {
	var toReturn []structs.Template
	var query IDPrefixQuery
	query.Selector.ID.GT = "\x00"
	query.Limit = 100

	b, err := json.Marshal(query)
	if err != nil {
		return toReturn, fmt.Errorf("failed to marshal query to get all templates: %s", err)
	}

	var resp templateQueryResponse

	err = c.MakeRequest("POST", fmt.Sprintf("%v/_find", OPTIONS), "application/json", b, &resp)
	if err != nil {
		return toReturn, fmt.Errorf("failed to get all templates: %s", err)
	}

	for _, doc := range resp.Docs {
		if len(doc.Template.UIConfig.Api) > 0 {
			toReturn = append(toReturn, *doc.Template)
		}
	}

	return toReturn, err
}

// GetTemplate returns a template UIConfig.
func (c *CouchDB) GetTemplate(id string) (structs.UIConfig, error) {
	log.L.Info(id)
	template, err := c.getTemplate(id)
	return *template.UIConfig, err
}

func (c *CouchDB) getTemplate(id string) (uiconfig, error) {
	var toReturn uiconfig

	err := c.MakeRequest("GET", fmt.Sprintf("%v/%v", OPTIONS, id), "", nil, &toReturn)
	if err != nil {
		err = fmt.Errorf("failed to get template %s: %s", id, err)
	}

	return toReturn, err
}

// UpdateTemplate sends an updated template to the database.
func (c *CouchDB) UpdateTemplate(id string, newTemp structs.UIConfig) (structs.UIConfig, error) {
	var toReturn structs.UIConfig

	if id == newTemp.ID { // the template ID isn't changing
		// get the rev of the template
		oldTemp, err := c.getTemplate(id)
		if err != nil {
			return toReturn, fmt.Errorf("unable to get template %s to update: %s", id, err)
		}

		// marshal the new template
		b, err := json.Marshal(newTemp)
		if err != nil {
			return toReturn, fmt.Errorf("unable to marshal new template: %s", err)
		}

		// update the template
		err = c.MakeRequest("PUT", fmt.Sprintf("%s/%s?rev=%v", OPTIONS, id, oldTemp.Rev), "application/json", b, &toReturn)
		if err != nil {
			return toReturn, fmt.Errorf("failed to update template %s: %s", id, err)
		}
	} else { // the template ID is changing :|
		// delete the old template
		err := c.deleteTemplate(id)
		if err != nil {
			return toReturn, fmt.Errorf("unable to delete old template %s: %s", id, err)
		}

		// marshal the new template
		b, err := json.Marshal(newTemp)
		if err != nil {
			return toReturn, fmt.Errorf("unable to marshal new template: %s", err)
		}

		// post new template
		var resp CouchUpsertResponse
		err = c.MakeRequest("POST", fmt.Sprintf("%v/%v", OPTIONS, newTemp.ID), "", b, &resp)
		if err != nil {
			if _, ok := err.(*Conflict); ok { // a template with the same ID already exists
				return toReturn, fmt.Errorf("template already exists, please update this template or change IDs. error: %s", err)
			}

			// or an unknown error
			return toReturn, fmt.Errorf("unable to create template %s : %s", id, err)
		}
	}

	return toReturn, nil
}

func (c *CouchDB) deleteTemplate(id string) error {
	// get the template to delete
	template, err := c.getTemplate(id)
	if err != nil {
		return fmt.Errorf("unable to get template %s to delete: %s", id, err)
	}

	// delete the template
	err = c.MakeRequest("DELETE", fmt.Sprintf("%v/%v?rev=%v", OPTIONS, id, template.Rev), "", nil, nil)
	if err != nil {
		return fmt.Errorf("unable to delete template %s: %s", id, err)
	}

	return nil
}

// ICONS

// GetIcons returns a list of IOConfigurations.
func (c *CouchDB) GetIcons() ([]string, error) {
	i, err := c.getIcons()
	return i.IconList, err
}

func (c *CouchDB) getIcons() (icons, error) {
	var toReturn icons

	err := c.MakeRequest("GET", fmt.Sprintf("%v/%v", OPTIONS, ICONS), "", nil, &toReturn)
	if err != nil {
		err = fmt.Errorf("failed to get icon mapping : %s", err)
	}

	return toReturn, err
}

// UpdateIcons puts an updated list of icons in the database.
func (c *CouchDB) UpdateIcons(iconList []string) ([]string, error) {
	var toReturn []string

	// get the rev of the icon list
	oldList, err := c.getIcons()
	if err != nil {
		return toReturn, fmt.Errorf("unable to get icon list to update: %s", err)
	}

	// marshal the new icon list
	b, err := json.Marshal(iconList)
	if err != nil {
		return toReturn, fmt.Errorf("unable to marshal new icon list: %s", err)
	}

	// update the icon list
	err = c.MakeRequest("PUT", fmt.Sprintf("%s/%s?rev=%v", OPTIONS, ICONS, oldList.Rev), "application/json", b, &toReturn)
	if err != nil {
		return toReturn, fmt.Errorf("failed to update the icon list : %s", err)
	}

	return toReturn, nil
}

// ROLES

// GetDeviceRoles returns a list of the device roles.
func (c *CouchDB) GetDeviceRoles() ([]structs.Role, error) {
	roles, err := c.getDeviceRoles()
	return roles.RoleList, err
}

func (c *CouchDB) getDeviceRoles() (deviceRoles, error) {
	var toReturn deviceRoles

	err := c.MakeRequest("GET", fmt.Sprintf("%v/%v", OPTIONS, ROLES), "", nil, &toReturn)
	if err != nil {
		err = fmt.Errorf("failed to get device roles : %s", err)
	}

	return toReturn, err
}

// UpdateDeviceRoles puts an updated list of device roles in the database.
func (c *CouchDB) UpdateDeviceRoles(roles []structs.Role) ([]structs.Role, error) {
	var toReturn []structs.Role

	// get the rev of the device role list
	oldList, err := c.getDeviceRoles()
	if err != nil {
		return toReturn, fmt.Errorf("unable to get the device role list to update: %s", err)
	}

	// marshal the new device role list
	b, err := json.Marshal(roles)
	if err != nil {
		return toReturn, fmt.Errorf("unable to marshal new device role list: %s", err)
	}

	// update the device role list
	err = c.MakeRequest("PUT", fmt.Sprintf("%s/%s?rev=%v", OPTIONS, ROLES, oldList.Rev), "application/json", b, &toReturn)
	if err != nil {
		return toReturn, fmt.Errorf("failed to update the device role list : %s", err)
	}

	return toReturn, nil
}

// DESIGNATIONS

// GetRoomDesignations returns a list of the room designations.
func (c *CouchDB) GetRoomDesignations() ([]string, error) {
	d, err := c.getRoomDesignations()
	return d.DesigList, err
}

func (c *CouchDB) getRoomDesignations() (roomDesignations, error) {
	var toReturn roomDesignations

	err := c.MakeRequest("GET", fmt.Sprintf("%v/%v", OPTIONS, ROOM_DESIGNATIONS), "", nil, &toReturn)
	if err != nil {
		err = fmt.Errorf("failed to get room designations : %s", err)
	}

	return toReturn, err
}

// UpdateRoomDesignations puts an updated list of room designations in the database.
func (c *CouchDB) UpdateRoomDesignations(desigs []string) ([]string, error) {
	var toReturn []string

	// get the rev of the room designation list
	oldList, err := c.getRoomDesignations()
	if err != nil {
		return toReturn, fmt.Errorf("unable to get the room designation list to update: %s", err)
	}

	// marshal the new room designation list
	b, err := json.Marshal(desigs)
	if err != nil {
		return toReturn, fmt.Errorf("unable to marshal new room designation list: %s", err)
	}

	// update the room designation list
	err = c.MakeRequest("PUT", fmt.Sprintf("%s/%s?rev=%v", OPTIONS, ROOM_DESIGNATIONS, oldList.Rev), "application/json", b, &toReturn)
	if err != nil {
		return toReturn, fmt.Errorf("failed to update the room designation list : %s", err)
	}

	return toReturn, nil
}

// CLOSURE CODES

// GetClosureCodes returns a list of the possible closure codes for ServiceNow.
func (c *CouchDB) GetClosureCodes() ([]string, error) {
	codes, err := c.getClosureCodes()
	return codes.Codes, err
}

func (c *CouchDB) getClosureCodes() (closureCodes, error) {
	var toReturn closureCodes

	err := c.MakeRequest("GET", fmt.Sprintf("%v/%v", OPTIONS, CLOSURE_CODES), "", nil, &toReturn)
	if err != nil {
		err = fmt.Errorf("failed to get closure codes : %s", err)
	}

	return toReturn, err
}

// UpdateClosureCodes puts an updated list of closure codes in the database.
func (c *CouchDB) UpdateClosureCodes(codes []string) ([]string, error) {
	var toReturn []string

	// get the rev of the closure code list
	oldList, err := c.getClosureCodes()
	if err != nil {
		return toReturn, fmt.Errorf("unable to get the closure code list to update: %s", err)
	}

	// marshal the new closure code list
	b, err := json.Marshal(codes)
	if err != nil {
		return toReturn, fmt.Errorf("unable to marshal new closure code list: %s", err)
	}

	// update the closure codes list
	err = c.MakeRequest("PUT", fmt.Sprintf("%s/%s?rev=%v", OPTIONS, CLOSURE_CODES, oldList.Rev), "application/json", b, &toReturn)
	if err != nil {
		return toReturn, fmt.Errorf("failed to update the closure code list : %s", err)
	}

	return toReturn, nil
}

// TAGS

// GetTags returns a list of the device roles.
func (c *CouchDB) GetTags() ([]string, error) {
	t, err := c.getTags()
	return t.TagList, err
}

func (c *CouchDB) getTags() (tags, error) {
	var toReturn tags

	err := c.MakeRequest("GET", fmt.Sprintf("%v/%v", OPTIONS, TAGS), "", nil, &toReturn)
	if err != nil {
		err = fmt.Errorf("failed to get device roles : %s", err)
	}

	return toReturn, err
}

// UpdateTags puts an updated list of tagsin the database.
func (c *CouchDB) UpdateTags(newTags []string) ([]string, error) {
	var toReturn []string

	// get the rev of the tag list
	oldList, err := c.getTags()
	if err != nil {
		return toReturn, fmt.Errorf("unable to get the tag list to update: %s", err)
	}

	// marshal the new tag list
	b, err := json.Marshal(newTags)
	if err != nil {
		return toReturn, fmt.Errorf("unable to marshal new tag list: %s", err)
	}

	// update the tag list
	err = c.MakeRequest("PUT", fmt.Sprintf("%s/%s?rev=%v", OPTIONS, TAGS, oldList.Rev), "application/json", b, &toReturn)
	if err != nil {
		return toReturn, fmt.Errorf("failed to update the tag list : %s", err)
	}

	return toReturn, nil
}

// GetMenuTree returns a list of attribute sets from the database
func (c *CouchDB) GetMenuTree() ([]string, error) {
	i, err := c.getMenuTree()
	return i.Order, err
}

func (c *CouchDB) getMenuTree() (menu, error) {
	var toReturn menu

	err := c.MakeRequest("GET", fmt.Sprintf("%v/%v", OPTIONS, MENUTREE), "", nil, &toReturn)
	if err != nil {
		err = fmt.Errorf("failed to get menu tree order : %s", err)
	}

	return toReturn, err
}
