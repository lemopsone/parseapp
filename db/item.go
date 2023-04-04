package db

import (
	"database/sql"
	"github.com/lemopsone/parseapp/models"
)

func (db Database) GetAllItems() (*models.ItemList, error) {
	list := &models.ItemList{}
	rows, err := db.Conn.Query("SELECT * FROM items")
	if err != nil {
		return list, err
	}
	for rows.Next() {
		var item models.Item
		err := rows.Scan(&item.Href,
			&item.CurrentPriceYEN, &item.CurrentPriceRUR,
			&item.BlitzPriceYEN, &item.BlitzPriceRUR,
			&item.TimeLeft,
			&item.TelegramID)
		if err != nil {
			return list, err
		}
		list.Items = append(list.Items, item)
	}
	return list, nil
}

func (db Database) AddItem(item *models.Item) error {
	query := "INSERT INTO items (href, title, current_price_yen, current_price_rur, blitz_price_yen, blitz_price_rur, time_left, telegram_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);"
	if _, err := db.Conn.Exec(query, item.Href, item.Title, item.CurrentPriceYEN, item.CurrentPriceRUR, item.BlitzPriceYEN, item.BlitzPriceRUR, item.TimeLeft, item.TelegramID); err != nil {
		return err
	}
	return nil
}
func (db Database) GetItemByHref(itemHref string) (models.Item, error) {
	item := models.Item{}
	query := "SELECT * FROM items WHERE href = $1;"
	row := db.Conn.QueryRow(query, itemHref)
	switch err := row.Scan(&item.Href, &item.Title, &item.CurrentPriceYEN, &item.CurrentPriceRUR, &item.BlitzPriceYEN, &item.BlitzPriceRUR, &item.TimeLeft, &item.TelegramID); err {
	case sql.ErrNoRows:
		return item, ErrNoMatch
	default:
		return item, err
	}
}

func (db Database) DeleteItem(itemHref string) error {
	query := "DELETE FROM items WHERE href = $1;"
	_, err := db.Conn.Exec(query, itemHref)
	switch err {
	case sql.ErrNoRows:
		return ErrNoMatch
	default:
		return err
	}
}

func (db Database) UpdateItem(itemHref string, itemData models.Item) error {
	query := "UPDATE items SET title=$2, current_price_yen=$3, current_price_rur=$4, blitz_price_yen=$5, blitz_price_rur=$6, time_left=$7, telegram_id=$8 WHERE href = $1;"
	if _, err := db.Conn.Exec(query, itemHref, itemData.Title, itemData.CurrentPriceYEN, itemData.CurrentPriceRUR, itemData.BlitzPriceYEN, itemData.BlitzPriceRUR, itemData.TimeLeft, itemData.TelegramID); err != nil {
		return err
	}
	return nil
}

func (db Database) TruncateTable() error {
	query := "TRUNCATE TABLE items;"
	_, err := db.Conn.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
