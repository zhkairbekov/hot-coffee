package models

type TotalSalesResponse struct {
	TotalSales float64 `json:"total_sales"`
}

type PopularItemsResponse struct {
	Items []PopularItem `json:"popular_items"`
}

type PopularItem struct {
	ProductID   string `json:"product_id"`
	Name        string `json:"name"`
	TotalOrders int    `json:"total_orders"`
}
