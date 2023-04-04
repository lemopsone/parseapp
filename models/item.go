package models

import (
	"fmt"
)

type Item struct {
	Href            string
	Title           string
	CurrentPriceYEN string
	CurrentPriceRUR string
	BlitzPriceYEN   string
	BlitzPriceRUR   string
	TimeLeft        string
	TelegramID      string
}

type ItemList struct {
	Items []Item
}

func (i *Item) Bind() error {
	if i.Href == "" {
		return fmt.Errorf("HREF is required")
	}
	return nil
}
func (i *Item) Render() error {
	return nil
}
