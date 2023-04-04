package main

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/lemopsone/parseapp/db"
	"github.com/lemopsone/parseapp/models"
	"github.com/lemopsone/parseapp/siteParser"
	"log"
	"os"
	"time"
)

func main() {
	start := time.Now()
	dbUser, dbPassword, dbName :=
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB")
	database, err := db.Initialize(dbUser, dbPassword, dbName)
	if err != nil {
		log.Fatalf("Couldn't set up the database: %v", err)
	}
	defer database.Conn.Close()
	driver, err := postgres.WithInstance(database.Conn, &postgres.Config{})
	if err != nil {
		log.Fatalf("%v", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file:///usr/bin/parseapp/db/migrations",
		"postgres", driver)
	if err != nil {
		log.Fatalf("%v", err)
	}
	err = m.Down()
	if err != nil {
		log.Fatalf("Up and ... %v", err)
	}
	err = m.Up()
	if err != nil {
		log.Fatalf("Up and ... %v", err)
	}
	parsedItems := siteParser.Parse()
	newItems := models.ItemList{}
	if err != nil {
		log.Fatalf("No items found")
	}
	for _, item := range parsedItems.Items {
		_, err := database.GetItemByHref(item.Href) // dbItem
		if err != nil {
			if err == db.ErrNoMatch {
				newItems.Items = append(newItems.Items, item)
			} else {
				log.Fatalf("Error while searching: %v", err)
			}
		} else {
			/* TO-DO: Update the telegram message
			tgID := dbItem.TelegramID
			newTimeLeft = item.TimeLeft
			item.TelegramID = tgID
			*/
		}
		// TO-DO: send all new items to TG
	}
	err = database.TruncateTable()
	if err != nil {
		log.Fatalf("Failed to truncate, %v", err)
	}
	for _, item := range parsedItems.Items {
		err = database.AddItem(&item)
		if err != nil {
			log.Fatalf("Failed to add item, %v", err)
		}
	}
	fmt.Printf("Code executed in %s\n", time.Since(start))
}
