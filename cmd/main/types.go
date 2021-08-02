package main

type (
	Product struct {
		Id int `json:"id"`
		Name string `json:"name"`
		Price int `json:"price"`
		Stock int `json:"stock"`
	}
	User struct {
		Id int `json:"id"`
		Name string `json:"name"`
		Password string `json:"password"`
		Email string `json:"email"`
	}

	Cart struct {
		Id int `json:"id"`
		ProductId int `json:"productId"`
		UserId int `json:"userId"`
		Qty int `json:"qty"`
	}

	CartProduct struct {
		ProductId int
		Qty int
	}
)