package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

func checkout(w http.ResponseWriter, _ *http.Request){

	log.Printf("Started checkout for user %d\n", getLoggedInUser())
	txn := Db.Txn(false)
	defer txn.Abort()

	// Read cart for logged in user
	it, err := txn.Get("cart", "userId")
	if err != nil {
		panic(err)
	}

	var Checkout []Cart
	for obj := it.Next(); obj != nil; obj = it.Next() {
		cart := obj.(Cart)
		if cart.UserId == getLoggedInUser() {
			Checkout = append(Checkout, cart)
		}
	}

	// Before sending request to payment service lock products in stock
	log.Printf("Sending request to payment service user %d\n", getLoggedInUser())
	txn = Db.Txn(true)
	for _, cart := range Checkout {
		var id interface{} = cart.ProductId
		raw, err := txn.First("product", "id", id)
		if err != nil {
			panic(err)
		}
		raw.(*Product).Stock = raw.(*Product).Stock - cart.Qty
		if err := txn.Insert("product", raw); err != nil {
			panic(err)
		}
	}
	txn.Commit()

	response, err := pay()
	if err != nil {
		fmt.Println(err)
	}

	// In case of payment failure revert unlock products in stock
	log.Printf("Got response code %d from payment service user %d\n", response.StatusCode, getLoggedInUser())
	if response.StatusCode != http.StatusOK {
		log.Printf("Got error from payment service, rollingback stock reservation for user %d\n", getLoggedInUser())
		txn = Db.Txn(true)
		for _, cart := range Checkout {
			var id interface{} = cart.ProductId
			raw, err := txn.First("product", "id", id)
			if err != nil {
				panic(err)
			}
			raw.(*Product).Stock = raw.(*Product).Stock + cart.Qty
			if err := txn.Insert("product", raw); err != nil {
				panic(err)
			}
		}
		txn.Commit()
	} else {
		// checkout ok, remove all products in user's cart and update stocks
		txn = Db.Txn(true)
		it, err := txn.Get("cart", "id")
		if err != nil {
			panic(err)
		}

		var toUpdate []CartProduct
		for obj := it.Next(); obj != nil; obj = it.Next() {
			if obj.(Cart).UserId == Users[len(Users) - 1].Id {

				// update stocks
				toUpdate = append(toUpdate, CartProduct{obj.(Cart).ProductId, obj.(Cart).Qty})

				// remove cart entries
				toDelete := obj.(Cart)
				txn.Delete("cart", toDelete)
			}
		}
		txn.Commit()
		log.Printf("Successfully checkout out products for user %d\n", getLoggedInUser())

		// update stocks
		txn = Db.Txn(true)
		for _, cartProd := range toUpdate {
			raw, err := txn.First("product", "id", cartProd.ProductId)
			if err != nil {
				panic(err)
			}
			raw.(*Product).Stock = raw.(*Product).Stock - cartProd.Qty
			if err := txn.Insert("product", raw); err != nil {
				panic(err)
			}
		}
		txn.Commit()
		log.Printf("Successfully updated stocks after checkout for user out %d\n", getLoggedInUser())

	}

	if err != nil {
		fmt.Print(err.Error())
	}

	json.NewEncoder(w).Encode(Checkout)
}

func addToCart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["productId"]
	productId, err := strconv.Atoi(key)

	if err != nil {
		panic(err)
	}

	log.Printf("Adding product %d to cart for user %d\n", productId, getLoggedInUser())

	// Find cart by logged in userId and productId
	cartId := 0
	cart := Cart{cartId, productId, Users[len(Users) - 1].Id, 0}
	txn := Db.Txn(false)
	defer txn.Abort()

	it, err := txn.Get("cart", "id")
	if err != nil {
		panic(err)
	}

	for obj := it.Next(); obj != nil; obj = it.Next() {
		if obj.(Cart).UserId == Users[len(Users) - 1].Id && obj.(Cart).ProductId == productId {
			cart = obj.(Cart)
			cartId = cart.ProductId
		}
	}

	// Validate stock
	qty := cart.Qty
	raw, err := txn.First("product", "id", productId)
	if err != nil {
		panic(err)
	}
	if raw.(*Product).Stock < qty + 1 {
		log.Printf("Failure adding product %d to cart for user %d. Quantity exceeds stock\n", productId, getLoggedInUser())
		respond(w, http.StatusBadRequest, "Product quantity exceeds stock")
		return
	}

	// Add product to cart
	txn = Db.Txn(true)
	cart.Id = cartId + 1
	cart.Qty = cart.Qty + 1
	if err := txn.Insert("cart", cart); err != nil {
		panic(err)
	}

	txn.Commit()
	log.Printf("Added product %d to cart for user %d\n", productId, getLoggedInUser())
}

func respond(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}