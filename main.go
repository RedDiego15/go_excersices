package main

import (
	"fmt"
	"sync"
)

type Product struct {
	ID    int
	Name  string
	Price float64
}

type OrderItem struct {
	ProductID string
	Quantity  int
}

type Order struct {
	ID    int
	Items []OrderItem
}

func calculateOrderTotal(order Order, catalog map[string]Product) (float64, error) {
	total := 0.0
	for _, item := range order.Items {
		productPrice := catalog[item.ProductID].Price
		total += productPrice * float64(item.Quantity)
	}
	return total, nil
}

func processOrder(order Order, catalog map[string]Product, wg *sync.WaitGroup, results chan<- float64) {
	defer wg.Done()
	total, _ := calculateOrderTotal(order, catalog)
	results <- total
}

func main() {
	catalog := map[string]Product{
		"p1": {ID: 1, Name: "keyboard", Price: 50},
		"p2": {ID: 2, Name: "mouse", Price: 25},
		"p3": {ID: 3, Name: "monitor", Price: 300},
	}
	orders := []Order{
		{
			ID: 1,
			Items: []OrderItem{
				{
					ProductID: "p1",
					Quantity:  3,
				},
				{
					ProductID: "p3",
					Quantity:  2,
				},
			},
		},
		{
			ID: 1,
			Items: []OrderItem{
				{
					ProductID: "p1",
					Quantity:  4,
				},
				{
					ProductID: "p3",
					Quantity:  4,
				},
			},
		},
		{
			ID: 1,
			Items: []OrderItem{
				{
					ProductID: "p1",
					Quantity:  10,
				},
				{
					ProductID: "p3",
					Quantity:  10,
				},
			},
		},
	}

	var wg sync.WaitGroup
	results := make(chan float64, 3)
	for _, order := range orders {
		wg.Add(1)
		go processOrder(order, catalog, &wg, results)
	}
	wg.Wait()
	close(results)

	for total := range results {
		fmt.Println("Order total: ", total)
	}

}
