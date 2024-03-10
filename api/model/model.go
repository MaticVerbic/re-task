package model

// UpdatePackageSizes serves as both the Request and Response struct for UpdatePackageSizes endpoint.
type UpdatePackageSizes struct {
	Sizes []int `json:"sizes"`
}

type CalculateBestPackagesRequest struct {
	Order int `json:"order"`
}

type CalculateBestPackagesResponse struct {
	Packages []int `json:"packages"`
}
