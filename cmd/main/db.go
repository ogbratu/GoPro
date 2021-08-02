package main

import "github.com/hashicorp/go-memdb"

func setupDb() *memdb.MemDB {
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"product": {
				Name: "product",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.IntFieldIndex{Field: "Id"},
					},
					"name": {
						Name:    "name",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "Name"},
					},
					"price": {
						Name:    "price",
						Unique:  false,
						Indexer: &memdb.IntFieldIndex{Field: "Price"},
					},
					"stock": {
						Name:    "stock",
						Unique:  false,
						Indexer: &memdb.IntFieldIndex{Field: "Stock"},
					},
				},
			},
			"user": {
				Name: "user",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.IntFieldIndex{Field: "Id"},
					},
					"name": {
						Name:    "name",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "Name"},
					},
					"password": {
						Name:    "password",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "Password"},
					},
					"email": {
						Name:    "email",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "Email"},
					},
				},
			},
			"cart": {
				Name: "cart",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.IntFieldIndex{Field: "Id"},
					},
					"productId": {
						Name:    "productId",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "ProductId"},
					},
					"userId": {
						Name:    "userId",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "UserId"},
					},
					"qty": {
						Name:    "qty",
						Unique:  false,
						Indexer: &memdb.IntFieldIndex{Field: "UserId"},
					},
				},
			},
		},
	}

	// Create a new database
	db, err := memdb.NewMemDB(schema)
	if err != nil {
		panic(err)
	}

	// Create a write transaction
	txn := db.Txn(true)

	// Insert seed data
	products := []*Product{
		{1, "Product1", 100, 2},
		{2, "Product2", 300, 3},
	}
	for _, p := range products {
		if err := txn.Insert("product", p); err != nil {
			panic(err)
		}
	}
	users := []*User{
		{1, "Jon Doe", "1234", "jon.doe@company.com"},
		{2, "Jon Doe 2", "1234", "jon.doe2@company.com"},
	}
	for _, u := range users {
		if err := txn.Insert("user", u); err != nil {
			panic(err)
		}
	}

	// Commit the transaction
	txn.Commit()

	return db
}
