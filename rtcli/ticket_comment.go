package rtcli

import (
	"fmt"
)

/*
 *
 */
type TicketCorrespondOptions struct {
	Content     string `json:"Content,omitempty"`
	ContentType string `json:"ContentType,omitempty"`
}

type TicketCommentOptions struct {
	Content     string `json:"Content,omitempty"`
	ContentType string `json:"ContentType,omitempty"`
}

/*
 *
 */
func (o *Client) TicketComment(ticketID string, comment *TicketCommentOptions) error {
	_, err := o.doRequest("POST", fmt.Sprintf("ticket/%s/comment", ticketID), comment, nil)
	if err != nil {
		return fmt.Errorf("error commenting ticket: %w", err)
	}
	return nil
}

func (o *Client) TicketCorrespond(ticketID string, comment *TicketCorrespondOptions) error {
	_, err := o.doRequest("POST", fmt.Sprintf("ticket/%s/correspond", ticketID), comment, nil)
	if err != nil {
		return fmt.Errorf("error commenting ticket: %w", err)
	}
	return nil
}
