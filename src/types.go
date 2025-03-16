package main

type Product struct {
	Id            int     `json:"id"`
	Name          string  `json:"name"`
	Price         float32 `json:"price"`
	Category      string  `json:"category"`
	Image         string  `json:"image"`
	AmountSold    int     `json:"amountSold"`
	AmountInStock int     `json:"amountInStock"`
	HasAllergens  bool    `json:"hasAllergens"`
	Rating        float32 `json:"rating"`
}

type Res struct {
	Success bool   `json:"success"`
	Msg     string `json:"msg"`
	Data    any    `json:"data"`
}
