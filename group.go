package box

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

type GroupService struct {
	*Client
}

type Group struct {
	Type       string `json:"type"`
	ID         string `json:"id"`
	Name       string `json:"name"`
	CreatedAt  string `json:"created_at"`
	ModifiedAt string `json:"modified_at"`
}

type GroupCollection struct {
	TotalCount int     `json:"total_count"`
	Entries    []Group `json:"entries"`
	Limit      int     `json:"limit"`
	Offset     int     `json:"offset"`
}

// Docs: https://developers.box.com/docs/#groups-get-all-groups
// TODO(ttacon): test it
func (c *GroupService) Groups(filter string) ([]Group, error) {
	reqBody := map[string]string{}
	if filter != "" {
		reqBody["filter"] = filter
	}
	var data *GroupCollection
	var groups []Group

	for len(groups) != data.TotalCount && data.TotalCount == 0 {
		reqBody["offset"] = strconv.Itoa(data.Offset + data.Limit)

		req, err := c.NewRequest(
			"GET",
			"/groups",
			reqBody,
		)
		if err != nil {
			return nil, err
		}

		_, err = c.Do(req, data)
		if data != nil {
			groups = append(groups, data.Entries...)
		}
	}

	return groups, nil
}

// Docs: https://developers.box.com/docs/#groups-create-a-group
// TODO(ttacon): test it
func (c *GroupService) CreateGroup(name string) (*http.Response, *Group, error) {
	req, err := c.NewRequest(
		"POST",
		"/groups",
		map[string]string{
			"name": name,
		},
	)
	if err != nil {
		return nil, nil, err
	}

	var data Group
	resp, err := c.Do(req, &data)
	return resp, &data, err
}

// Docs: https://developers.box.com/docs/#update-a-group
// TODO(ttacon): test it
func (c *GroupService) UpdateGroup(groupID, name string) (*http.Response, *Group, error) {
	req, err := c.NewRequest(
		"PUT",
		fmt.Sprintf("/groups/%s", groupID),
		map[string]string{
			"name": name,
		},
	)
	if err != nil {
		return nil, nil, err
	}

	var data Group
	resp, err := c.Do(req, &data)
	return resp, &data, err
}

// Docs: https://developers.box.com/docs/#delete-a-group
// TODO(ttacon): test it
func (c *GroupService) DeleteGroup(groupID string) (*http.Response, bool, error) {
	req, err := c.NewRequest(
		"PUT",
		fmt.Sprintf("/groups/%s", groupID),
		nil,
	)
	if err != nil {
		return nil, false, err
	}

	var data Group
	resp, err := c.Do(req, &data)
	return resp, resp.StatusCode == 204, err
}

// Documentation: https://developers.box.com/docs/#groups-get-the-membership-list-for-a-group
func (g *GroupService) ListMembership(groupID string) (*http.Response, *MembershipCollection, error) {
	req, err := g.NewRequest(
		"GET",
		fmt.Sprintf("/groups/%s/memberships", groupID),
		nil,
	)
	if err != nil {
		return nil, nil, err
	}

	var membershipCollection MembershipCollection
	resp, err := g.Do(req, &membershipCollection)
	return resp, &membershipCollection, err
}

type CollectionInfo struct {
	TotalCount int `json:"total_count"`
	Offset     int `json:"offset"`
	Limit      int `json:"limit"`
}

type Membership struct {
	Type  string `json:"type"`
	ID    string `json:"id"`
	User  *User  `json:"user"`
	Group *Group `json:"group"`
	Role  string `json:"role"`
}

type MembershipCollection struct {
	CollectionInfo
	Entries []*Membership `json:"entries"`
}

func (g *GroupService) Membership(membershipEntryID string) (*http.Response, *Membership, error) {
	req, err := g.NewRequest(
		"GET",
		fmt.Sprintf("/group_memberships/%s", membershipEntryID),
		nil,
	)
	if err != nil {
		return nil, nil, err
	}

	var membership Membership
	resp, err := g.Do(req, &membership)
	return resp, &membership, err
}

// Documentation: https://developers.box.com/docs/#groups-add-a-member-to-a-group
func (g *GroupService) AddUserToGroup(uID, gID, role string) (*http.Response, *Membership, error) {
	// try to be nice about errors
	if len(uID) == 0 {
		return nil, nil, errors.New("must provide user ID")
	} else if len(gID) == 0 {
		return nil, nil, errors.New("must provide group ID")
	}

	var toSend = map[string]interface{}{
		"user": map[string]string{
			"id": uID,
		},
		"group": map[string]string{
			"id": gID,
		},
	}
	if len(role) > 0 {
		toSend["role"] = role
	}

	req, err := g.NewRequest(
		"POST",
		"/group_memberships",
		toSend,
	)
	if err != nil {
		return nil, nil, err
	}

	var membership Membership
	resp, err := g.Do(req, &membership)
	return resp, &membership, err
}

// Documentation: https://developers.box.com/docs/#groups-update-a-group-membership
func (g *GroupService) UpdateMembership(membershipID, role string) (*http.Response, *Membership, error) {
	req, err := g.NewRequest(
		"PUT",
		fmt.Sprintf("/group_memberships/%s", membershipID),
		map[string]string{
			"role": role,
		},
	)
	if err != nil {
		return nil, nil, err
	}

	var membership Membership
	resp, err := g.Do(req, &membership)
	return resp, &membership, err
}

// Documentation: https://developers.box.com/docs/#groups-delete-a-group-membership
func (g *GroupService) DeleteMembership(membershipID string) (*http.Response, error) {
	req, err := g.NewRequest(
		"DELETE",
		fmt.Sprintf("/group_memberships/%s", membershipID),
		nil,
	)
	if err != nil {
		return nil, err
	}

	return g.Do(req, nil)
}

// Documentation: https://developers.box.com/docs/#get-all-collaborations-for-a-group
func (g *GroupService) GroupCollaborations(groupID string) (*http.Response, *Collaborations, error) {
	req, err := g.NewRequest(
		"GET",
		fmt.Sprintf("/groups/%s/collaborations", groupID),
		nil,
	)
	if err != nil {
		return nil, nil, err
	}

	var collaborations Collaborations
	resp, err := g.Do(req, &collaborations)
	return resp, &collaborations, err
}
