package shop

import (
	"fmt"
)

type Shop struct {
	Products []Product
	shopSize int
}

func RunShop(shop Shop, addProduct chan Product, sellProduct chan Offer, getProducts chan bool) {
	for {
		select {
		case newProduct := <-newSaleGuard(hasSpace(shop), addProduct):
			shop.Products = append(shop.Products, newProduct)
		case sale := <-saleGuard(len(shop.Products) > 0, sellProduct):
			sale.Product <- shop.Products[0]
			shop.Products = shop.Products[1:]
		case <-getProducts:
			for _, product := range shop.Products {
				fmt.Println(product)
			}
		}
	}
}

func newSaleGuard(condition bool, channel chan Product) chan Product {
	if condition {
		return channel
	}
	return nil
}

func saleGuard(condition bool, channel <-chan Offer) <-chan Offer {
	if condition {
		return channel
	}
	return nil
}

func CreateShop(shopSize int) Shop {
	return Shop{make([]Product, 0, shopSize), shopSize}
}

func hasSpace(shop Shop) bool {
	return len(shop.Products) < shop.shopSize
}
