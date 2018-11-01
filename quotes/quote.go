package quotes

import "context"

type Quote struct {
	ID      string `json:"id,omitempty"`
	Author  string `json:"author,omitempty"`
	Quote   string `json:"quote,omitempty"`
	Created int64  `json:"created,omitempty"`
}

type Interactor interface {
	List(context.Context) ([]*Quote, error)
	Put(context.Context, *Quote) (*Quote, error)
}
