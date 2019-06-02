package shop

type Product struct {
	Result int
}

func CreateProduct(result int) Product {
	return Product{
		Result: result,
	}
}

type Offer struct {
	Product chan Product
}
