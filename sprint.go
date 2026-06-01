package jira

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hendrikiseke1979/jira/jiradata"
)

// GetActiveSprint finds the active sprint for the given project by first
// finding the scrum board for the project, then fetching its active sprint.
func GetActiveSprint(ua HttpClient, endpoint, project string) (*jiradata.Sprint, error) {
	// Step 1: find the scrum board for this project
	uri, err := url.Parse(URLJoin(endpoint, "rest/agile/1.0/board"))
	if err != nil {
		return nil, err
	}
	params := url.Values{}
	params.Set("projectKeyOrId", project)
	params.Set("type", "scrum")
	uri.RawQuery = params.Encode()

	resp, err := ua.GetJSON(uri.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, responseError(resp)
	}

	boardList := &jiradata.BoardList{}
	if err := json.NewDecoder(resp.Body).Decode(boardList); err != nil {
		return nil, err
	}
	if len(boardList.Values) == 0 {
		return nil, fmt.Errorf("no scrum board found for project %q", project)
	}
	boardID := boardList.Values[0].ID

	// Step 2: get the active sprint for that board
	sprintURI, err := url.Parse(URLJoin(endpoint, "rest/agile/1.0/board", fmt.Sprintf("%d", boardID), "sprint"))
	if err != nil {
		return nil, err
	}
	sprintParams := url.Values{}
	sprintParams.Set("state", "active")
	sprintURI.RawQuery = sprintParams.Encode()

	sprintResp, err := ua.GetJSON(sprintURI.String())
	if err != nil {
		return nil, err
	}
	defer sprintResp.Body.Close()

	if sprintResp.StatusCode != 200 {
		return nil, responseError(sprintResp)
	}

	sprintList := &jiradata.SprintList{}
	if err := json.NewDecoder(sprintResp.Body).Decode(sprintList); err != nil {
		return nil, err
	}
	for _, sprint := range sprintList.Values {
		if sprint.OriginBoardID == boardID {
			return sprint, nil
		}
	}
	return nil, fmt.Errorf("no active sprint found for project %q", project)
}

// SprintAddIssues adds the given issues to the sprint with the given ID.
// It reuses the EpicIssues type since the request body is identical.
func SprintAddIssues(ua HttpClient, endpoint string, sprintID int, issues *jiradata.EpicIssues) error {
	encoded, err := json.Marshal(issues)
	if err != nil {
		return err
	}

	uri := URLJoin(endpoint, "rest/agile/1.0/sprint", fmt.Sprintf("%d", sprintID), "issue")
	resp, err := ua.Post(uri, "application/json", bytes.NewBuffer(encoded))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 204 {
		return nil
	}
	return responseError(resp)
}
