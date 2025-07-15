package models

import "app/pkg/custom_errors"

type Product struct {
	ID                             int     `json:"id" db:"id"`
	ProductCode                    string  `json:"product_code" db:"product_code"`
	Description                    string  `json:"description" db:"description"`
	Width                          float64 `json:"width" db:"width"`
	Height                         float64 `json:"height" db:"height"`
	Length                         float64 `json:"length" db:"length"`
	NetWeight                      float64 `json:"net_weight" db:"net_weight"`
	ExpirationRate                 int     `json:"expiration_rate" db:"expiration_rate"`
	RecommendedFreezingTemperature float64 `json:"recommended_freezing_temperature" db:"recommended_freezing_temperature"`
	FreezingRate                   int     `json:"freezing_rate" db:"freezing_rate"`
	ProductTypeId                  int     `json:"product_type_id" db:"product_type_id"`
	SellerId                       *int    `json:"seller_id" db:"seller_id"`
}

type ProductRequest struct {
	Id                             int      `json:"id"`
	ProductCode                    *string  `json:"product_code"`
	Description                    *string  `json:"description"`
	Width                          *float64 `json:"width"`
	Height                         *float64 `json:"height"`
	Length                         *float64 `json:"length"`
	NetWeight                      *float64 `json:"net_weight"`
	ExpirationRate                 *int     `json:"expiration_rate"`
	RecommendedFreezingTemperature *float64 `json:"recommended_freezing_temperature"`
	FreezingRate                   *int     `json:"freezing_rate"`
	ProductTypeId                  *int     `json:"product_type_id"`
	SellerId                       *int     `json:"seller_id"`
}

type ProductPatchRequest struct {
	Id                             int      `json:"id"`
	ProductCode                    *string  `json:"product_code"`
	Description                    *string  `json:"description"`
	Width                          *float64 `json:"width"`
	Height                         *float64 `json:"height"`
	Length                         *float64 `json:"length"`
	NetWeight                      *float64 `json:"net_weight"`
	ExpirationRate                 *int     `json:"expiration_rate"`
	RecommendedFreezingTemperature *float64 `json:"recommended_freezing_temperature"`
	FreezingRate                   *int     `json:"freezing_rate"`
	ProductTypeId                  *int     `json:"product_type_id"`
	SellerId                       *int     `json:"seller_id"`
}

func (c ProductRequest) Verify() error {
	if c.ProductCode == nil {
		return &custom_errors.MandatoryArgMissingErr{Argument: "product_code"}
	}

	if c.Description == nil {
		return &custom_errors.MandatoryArgMissingErr{Argument: "description"}
	}

	if c.Width == nil {
		return &custom_errors.MandatoryArgMissingErr{Argument: "width"}
	}

	if c.Height == nil {
		return &custom_errors.MandatoryArgMissingErr{Argument: "height"}
	}

	if c.Length == nil {
		return &custom_errors.MandatoryArgMissingErr{Argument: "length"}
	}

	if c.NetWeight == nil {
		return &custom_errors.MandatoryArgMissingErr{Argument: "net_weight"}
	}

	if c.ExpirationRate == nil {
		return &custom_errors.MandatoryArgMissingErr{Argument: "expiration_rate"}
	}

	if c.RecommendedFreezingTemperature == nil {
		return &custom_errors.MandatoryArgMissingErr{Argument: "recommended_freezing_temperature"}
	}

	if c.FreezingRate == nil {
		return &custom_errors.MandatoryArgMissingErr{Argument: "freezing_rate"}
	}
	if c.ProductTypeId == nil {
		return &custom_errors.MandatoryArgMissingErr{Argument: "product_type_id"}
	}

	return nil
}

func (c ProductPatchRequest) Verify() error {
	return nil
}
