package constants

var ProductCategory = map[string]struct{}{
	"Beverage":   {},
	"Food":       {},
	"Snack":      {},
	"Condiments": {},
	"Additions":  {},
}

func IsValidProductCategory(cat string) bool {
	_, ok := ProductCategory[cat]
	return ok
}
