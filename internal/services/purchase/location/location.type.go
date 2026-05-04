package location

//--------- LocationType ENUM -------------
type LocationType string

const (
	LocationWarehouse LocationType = "warehouse"
	LocationStore     LocationType = "store"
	LocationShowroom  LocationType = "showroom"
	LocationOutlet    LocationType = "outlet"
)