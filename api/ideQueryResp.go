package api

type IdeResp struct {
	ElementList        []IdeElement `json:"elementList"`
	Total              int          `json:"total"`
	TotalPages         int          `json:"totalPages"`
	ActualPage         int          `json:"actualPage"`
	ItemsPerPage       int          `json:"itemsPerPage"`
	NumPaginations     int          `json:"numPaginations"`
	HiddenResults      bool         `json:"hiddenResults"`
	Summary            []string     `json:"summary"`
	AlertName          string       `json:"alertName"`
	LowerRangePosition int          `json:"lowerRangePosition"`
	UpperRangePosition int          `json:"upperRangePosition"`
	Paginable          bool         `json:"paginable"`
}

type IdeElement struct {
	PropertyCode      int     `json:"propertyCode" db:"propertycode"`
	Thumbnail         string  `json:"thumbnail,omitempty" db:"thumbnail"`
	ExternalReference string  `json:"externalReference,omitempty" db:"externalreference"`
	NumPhotos         int     `json:"numPhotos" db:"numphotos"`
	Floor             string  `json:"floor,omitempty" db:"onfloor"`
	Price             float64 `json:"price" db:"price"`
	PropertyType      string  `json:"propertyType" db:"propertytype"`
	Operation         string  `json:"operation" db:"operation"`
	Size              float64 `json:"size" db:"hassize"`
	Exterior          bool    `json:"exterior" db:"exterior"`
	Rooms             int     `json:"rooms" db:"rooms"`
	Bathrooms         int     `json:"bathrooms" db:"bathrooms"`
	Address           string  `json:"address" db:"address"`
	Province          string  `json:"province" db:"province"`
	Municipality      string  `json:"municipality" db:"municipality"`
	District          string  `json:"district" db:"district"`
	Country           string  `json:"country" db:"country"`
	Latitude          float64 `json:"latitude" db:"latitude"`
	Longitude         float64 `json:"longitude" db:"longitude"`
	ShowAddress       bool    `json:"showAddress" db:"showaddress"`
	URL               string  `json:"url" db:"url"`
	Distance          string  `json:"distance" db:"distance"`
	Description       string  `json:"description,omitempty" db:"description"`
	HasVideo          bool    `json:"hasVideo" db:"hasvideo"`
	Status            string  `json:"status" db:"status"`
	NewDevelopment    bool    `json:"newDevelopment" db:"newdevelopment"`
	HasLift           bool    `json:"hasLift,omitempty" db:"haslift"`
	PriceByArea       float64 `json:"priceByArea" db:"pricebyarea"`
	SuggestedTexts    struct {
		Subtitle string `json:"subtitle" db:"subtitle"`
		Title    string `json:"title" db:"title"`
	} `json:"suggestedTexts" db:"suggestedtexts"`
	HasPlan           bool `json:"hasPlan" db:"hasplan"`
	Has3DTour         bool `json:"has3DTour" db:"has3dtour"`
	Has360            bool `json:"has360" db:"has360"`
	HasStaging        bool `json:"hasStaging" db:"hasstaging"`
	TopNewDevelopment bool `json:"topNewDevelopment" db:"topnewdevelopment"`
	ParkingSpace      struct {
		HasParkingSpace               bool    `json:"hasParkingSpace" db:"hasparkingspace"`
		IsParkingSpaceIncludedInPrice bool    `json:"isParkingSpaceIncludedInPrice" db:"isparkingspaceincludedinprice"`
		ParkingSpacePrice             float64 `json:"parkingSpacePrice" db:"parkingspaceprice"`
	} `json:"parkingSpace,omitempty" db:"parkingspace"`
	Neighborhood string `json:"neighborhood,omitempty" db:"neighborhood"`
	DetailedType struct {
		Typology    string `json:"typology" db:"typology"`
		SubTypology string `json:"subTypology" db:"subtypology"`
	} `json:"detailedType,omitempty" db:"detailedtype"`
}
