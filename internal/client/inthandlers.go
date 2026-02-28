package client

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// InterviewHandler represents an interview handler
type InterviewHandler struct {
	ObjectId    string `json:"ObjectId"`
	DisplayName string `json:"DisplayName"`
	DtmfAccessId string `json:"DtmfAccessId"`
	Recipient   string `json:"RecipientDistributionListObjectId"`
}

type interviewHandlersResponse struct {
	Total string             `json:"@total"`
	Items OneOrMany[InterviewHandler] `json:"InterviewHandler"`
}

// InterviewQuestion represents a question in an interview handler
type InterviewQuestion struct {
	QuestionNumber int    `json:"QuestionNumber"`
	Active         string `json:"Active"`
}

type interviewQuestionsResponse struct {
	Total string              `json:"@total"`
	Items OneOrMany[InterviewQuestion] `json:"InterviewQuestion"`
}

// ListInterviewHandlers returns interview handlers
func ListInterviewHandlers(host string, port int, user, pass string, query string, rowsPerPage int) ([]InterviewHandler, error) {
	path := "/handlers/interviewhandlers"
	params := url.Values{}
	if query != "" {
		params.Set("query", query)
	}
	if rowsPerPage > 0 {
		params.Set("rowsPerPage", fmt.Sprintf("%d", rowsPerPage))
	}
	if len(params) > 0 {
		path = path + "?" + params.Encode()
	}

	body, err := Get(host, port, user, pass, path)
	if err != nil {
		return nil, fmt.Errorf("failed to list interview handlers: %w", err)
	}

	var resp interviewHandlersResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse interview handlers response: %w", err)
	}

	return resp.Items, nil
}

// GetInterviewHandler retrieves an interview handler by name or ObjectId
func GetInterviewHandler(host string, port int, user, pass, nameOrID string) (*InterviewHandler, error) {
	if isUUID(nameOrID) {
		return getInterviewHandlerByID(host, port, user, pass, nameOrID)
	}
	return getInterviewHandlerByName(host, port, user, pass, nameOrID)
}

func getInterviewHandlerByID(host string, port int, user, pass, objectID string) (*InterviewHandler, error) {
	body, err := Get(host, port, user, pass, "/handlers/interviewhandlers/"+objectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get interview handler: %w", err)
	}

	var ih InterviewHandler
	if err := json.Unmarshal(body, &ih); err != nil {
		return nil, fmt.Errorf("failed to parse interview handler: %w", err)
	}

	return &ih, nil
}

func getInterviewHandlerByName(host string, port int, user, pass, name string) (*InterviewHandler, error) {
	q := fmt.Sprintf("(displayname is %s)", name)
	handlers, err := ListInterviewHandlers(host, port, user, pass, q, 1)
	if err != nil {
		return nil, err
	}
	if len(handlers) == 0 {
		return nil, fmt.Errorf("interview handler '%s' not found", name)
	}
	return &handlers[0], nil
}

// CreateInterviewHandler creates a new interview handler
func CreateInterviewHandler(host string, port int, user, pass string, fields map[string]interface{}) error {
	_, err := Post(host, port, user, pass, "/handlers/interviewhandlers", fields)
	if err != nil {
		return fmt.Errorf("failed to create interview handler: %w", err)
	}
	return nil
}

// UpdateInterviewHandler updates an interview handler by name or ObjectId
func UpdateInterviewHandler(host string, port int, user, pass, nameOrID string, fields map[string]interface{}) error {
	ih, err := GetInterviewHandler(host, port, user, pass, nameOrID)
	if err != nil {
		return err
	}
	return Put(host, port, user, pass, "/handlers/interviewhandlers/"+ih.ObjectId, fields)
}

// DeleteInterviewHandler deletes an interview handler by name or ObjectId
func DeleteInterviewHandler(host string, port int, user, pass, nameOrID string) error {
	ih, err := GetInterviewHandler(host, port, user, pass, nameOrID)
	if err != nil {
		return err
	}
	return Delete(host, port, user, pass, "/handlers/interviewhandlers/"+ih.ObjectId)
}

// ListInterviewQuestions returns all questions for an interview handler
func ListInterviewQuestions(host string, port int, user, pass, handlerObjectId string) ([]InterviewQuestion, error) {
	path := fmt.Sprintf("/handlers/interviewhandlers/%s/interviewquestions", url.PathEscape(handlerObjectId))
	body, err := Get(host, port, user, pass, path)
	if err != nil {
		return nil, fmt.Errorf("failed to list interview questions: %w", err)
	}

	var resp interviewQuestionsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse interview questions response: %w", err)
	}

	return resp.Items, nil
}

// GetInterviewQuestion retrieves a specific interview question
func GetInterviewQuestion(host string, port int, user, pass, handlerObjectId string, questionNumber int) (*InterviewQuestion, error) {
	path := fmt.Sprintf("/handlers/interviewhandlers/%s/interviewquestions/%d", url.PathEscape(handlerObjectId), questionNumber)
	body, err := Get(host, port, user, pass, path)
	if err != nil {
		return nil, fmt.Errorf("failed to get interview question: %w", err)
	}

	var q InterviewQuestion
	if err := json.Unmarshal(body, &q); err != nil {
		return nil, fmt.Errorf("failed to parse interview question: %w", err)
	}

	return &q, nil
}

// UpdateInterviewQuestion updates an interview question
func UpdateInterviewQuestion(host string, port int, user, pass, handlerObjectId string, questionNumber int, fields map[string]interface{}) error {
	path := fmt.Sprintf("/handlers/interviewhandlers/%s/interviewquestions/%d", url.PathEscape(handlerObjectId), questionNumber)
	if err := Put(host, port, user, pass, path, fields); err != nil {
		return fmt.Errorf("failed to update interview question: %w", err)
	}
	return nil
}
