package rtcli

import (
	"encoding/json"
	"fmt"
	"time"
)

/*
 *
 */
type TicketCreateOptions struct {
	Subject      string            `json:"Subject,omitempty"`
	Queue        string            `json:"Queue,omitempty"`
	Status       string            `json:"Status,omitempty"`
	Priority     string            `json:"Priority,omitempty"`
	Owner        string            `json:"Owner,omitempty"`
	Requestor    string            `json:"Requestor,omitempty"`
	Content      string            `json:"Content,omitempty"`
	ContentType  string            `json:"ContentType,omitempty"`
	Parent       string            `json:"Parent,omitempty"`
	Due          string            `json:"Due,omitempty"`
	CustomFields map[string]string `json:"CustomFields,omitempty"`
}

type TicketUpdateOptions struct {
	Status       *string           `json:"Status,omitempty"`
	CustomFields map[string]string `json:"CustomFields,omitempty"`
}

/*
 *
 */
type Ticket struct {
	ID              int           `json:"id,omitempty"`
	Subject         string        `json:"Subject,omitempty"`
	Queue           Queue         `json:"Queue,omitempty"`
	Status          string        `json:"Status,omitempty"`
	FinalPriority   string        `json:"FinalPriority,omitempty"`
	Owner           User          `json:"Owner,omitempty"`
	Requestor       []User        `json:"Requestor,omitempty"`
	Created         *time.Time    `json:"Created,omitempty"`
	Cc              []User        `json:"Cc,omitempty"`
	Creator         User          `json:"Creator,omitempty"`
	TimeLeft        string        `json:"TimeLeft,omitempty"`
	TimeEstimated   string        `json:"TimeEstimated,omitempty"`
	AdminCc         []User        `json:"AdminCc,omitempty"`
	Starts          *time.Time    `json:"Starts,omitempty"`
	Started         *time.Time    `json:"Started,omitempty"`
	LastUpdated     *time.Time    `json:"LastUpdated,omitempty"`
	InitialPriority string        `json:"InitialPriority,omitempty"`
	Due             *time.Time    `json:"Due,omitempty"`
	LastUpdatedBy   User          `json:"LastUpdatedBy,omitempty"`
	Priority        string        `json:"Priority,omitempty"`
	Resolved        *time.Time    `json:"Resolved,omitempty"`
	EffectiveID     Item          `json:"EffectiveID,omitempty"`
	CustomFields    []CustomField `json:"CustomFields,omitempty"`
}

// const (
// 	StatusOpen    = "open"
// 	StatusClosed  = "closed"
// 	StatusPending = "pending"
// )

type Item struct {
	ID   string `json:"id,omitempty"`
	Type string `json:"type,omitempty"`
	URL  string `json:"_url,omitempty"`
}

type CustomField struct {
	Item
	Name   string   `json:"name,omitempty"`
	Values []string `json:"values,omitempty"`
}

/*
 *
 */
type ticketCreated struct {
	URL  string `json:"_url,omitempty"`
	Type string `json:"type,omitempty"`
	ID   string `json:"id,omitempty"`
}

/*
 *
 */
func (o *Client) TicketCreate(request *TicketCreateOptions) (*Ticket, error) {
	resp, err := o.doRequest("POST", "ticket", request, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating ticket: %w", err)
	}

	var result ticketCreated
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	return o.TicketGet(result.ID)
}

func (o *Client) TicketGet(ticketID string) (*Ticket, error) {
	params := map[string]string{
		"fields[Queue]":   "Name",
		"fields[Owner]":   "Name,EmailAddress",
		"fields[Creator]": "Name,EmailAddress",
	}
	resp, err := o.doRequest("GET", fmt.Sprintf("ticket/%s", ticketID), nil, params)
	if err != nil {
		return nil, fmt.Errorf("error getting ticket: %w", err)
	}

	var ticket Ticket
	if err := json.Unmarshal(resp, &ticket); err != nil {
		return nil, fmt.Errorf("error parsing ticket: %w", err)
	}
	// Iterate through requestors and fetch additional details
	for i := range ticket.Requestor {
		if ticket.Requestor[i].ID == "" {
			continue
		}

		user, err := o.UserGet(ticket.Requestor[i].ID)
		if err != nil {
			return nil, fmt.Errorf("error getting requestor details: %w", err)
		}
		ticket.Requestor[i].EmailAddress = user.EmailAddress
		ticket.Requestor[i].Name = user.Name
	}
	return &ticket, nil
}

func (o *Client) TicketUpdate(ticketID string, request *TicketUpdateOptions) error {
	_, err := o.doRequest("PUT", fmt.Sprintf("ticket/%s", ticketID), request, nil)
	if err != nil {
		return fmt.Errorf("error updating ticket: %w", err)
	}

	return nil
}

func (o *Client) TicketTake(ticketID string) error {
	_, err := o.doRequest("PUT", fmt.Sprintf("ticket/%s/take", ticketID), nil, nil)
	if err != nil {
		return fmt.Errorf("error take ticket: %w", err)
	}

	return nil
}

func (o *Client) TicketUntake(ticketID string) error {
	_, err := o.doRequest("PUT", fmt.Sprintf("ticket/%s/untake", ticketID), nil, nil)
	if err != nil {
		return fmt.Errorf("error untake ticket: %w", err)
	}

	return nil
}

func (o *Client) TicketSteal(ticketID string) error {
	_, err := o.doRequest("PUT", fmt.Sprintf("ticket/%s/steal", ticketID), nil, nil)
	if err != nil {
		return fmt.Errorf("error steal ticket: %w", err)
	}

	return nil
}

func (o *Client) TicketDelete(ticketID string) error {
	_, err := o.doRequest("DELETE", fmt.Sprintf("ticket/%s", ticketID), nil, nil)
	if err != nil {
		return fmt.Errorf("error steal ticket: %w", err)
	}

	return nil
}
