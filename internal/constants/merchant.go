package constants

var MerchantCategorySet = map[string]struct{}{
	"SmallRestaurant":       {},
	"MediumRestaurant":      {},
	"LargeRestaurant":       {},
	"MerchandiseRestaurant": {},
	"BoothKiosk":            {},
	"ConvenienceStore":      {},
}

func IsValidMerchantCategory(cat string) bool {
	_, ok := MerchantCategorySet[cat]
	return ok
}
