package main

type ProductFull struct {
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

type ProductSimple struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Image    string `json:"image"`
}

type AllergenInfo struct {
	Section string `json:"section"`
	Text    string `json:"text"`
}

type AllergenFull struct {
	Id    int            `json:"id"`
	Name  string         `json:"name"`
	Image string         `json:"image"`
	Info  []AllergenInfo `json:"info"`
}

type AllergenSimple struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
}

type Res struct {
	Success bool   `json:"success"`
	Msg     string `json:"msg"`
	Data    any    `json:"data"`
}
