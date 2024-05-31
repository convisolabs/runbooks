package type_cp

type MetadataResponse struct {
	CurrentPage int `json:"currentPage"`
	TotalPages  int `json:"totalPages"`
}

type AssetsResponse struct {
	Data AssetsDataResponse `json:"data"`
}

type AssetsDataResponse struct {
	Assets AssetsJsonResponse `json:"assets"`
}

type AssetsJsonResponse struct {
	Collection []Asset          `json:"collection"`
	Metadata   MetadataResponse `json:"metadata"`
}

type Asset struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type AssetsByTimeParameters struct {
	CompanyId int
	Page      int
	Limit     int
	Search    AssetsByTimeSearchParameters
}

type AssetsByTimeSearchParameters struct {
	CreatedAt AssetsByTimeCreateAtParameters `json:"createdAt"`
}

type AssetsByTimeCreateAtParameters struct {
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}
