package company

import (
	. "wspolbiezne/company/shop"
	"fmt"
	"time"
)

type Customer struct {
	name string
	shop chan Offer
	loud bool
}

func CreateCustomer(name string, shop chan Offer, loud bool) Customer {
	return Customer{name, shop, loud}
}

func StartCustomer(customer Customer) {
	for {
		var product = make(chan Product)
		offer := Offer{Product: product}
		customer.shop <- offer

		boughtItem := <-offer.Product

		if customer.loud {
			fmt.Printf("%s customer bought \t%v\n", customer.name, boughtItem.Result)
		}
		time.Sleep(newOfferInterval)
	}
}
