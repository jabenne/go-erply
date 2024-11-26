package legacy

import (
	"encoding/json"
	"fmt"

	"github.com/google/go-querystring/query"
)

type CouponService service

type SaveIssuedCouponOpts struct {
	IssuedCouponID            int    `url:"issuedCouponID,omitempty" json:"issuedCouponID,omitempty"`
	CouponID                  int    `url:"couponID,omitempty" json:"couponID,omitempty"`
	UniqueIdentifier          string `url:"uniqueIdentifier,omitempty" json:"uniqueIdentifier,omitempty"`
	InvoiceID                 int    `url:"invoiceID,omitempty" json:"invoiceID,omitempty"`
	CustomerID                int 	 `url:"customerID,omitempty" json:"customerID,omitempty"`
	WarehouseID               int    `url:"warehouseID,omitempty" json:"warehouseID,omitempty"`
	PointOfSaleID             int    `url:"pointOfSaleID,omitempty" json:"pointOfSaleID,omitempty"`
	EmployeeID                int    `url:"employeeID,omitempty" json:"employeeID,omitempty"`
	Timestamp                 int    `url:"timestamp,omitempty" json:"timestamp,omitempty"`
	ExpiryDate                string `url:"expiryDate,omitempty" json:"expiryDate,omitempty"`
	IsPrintedAutomatically    int    `url:"isPrintedAutomatically,omitempty" json:"isPrintedAutomatically,omitempty"`
	DoNotSubtractRewardPoints int    `url:"doNotSubtractRewardPoints,omitempty" json:"doNotSubtractRewardPoints,omitempty"`
	RedeemedInvoiceID         int    `url:"redeemedInvoiceID,omitempty" json:"redeemedInvoiceID,omitempty"`
	RedeemedCustomerID        int    `url:"redeemedCustomerID,omitempty" json:"redeemedCustomerID,omitempty"`
	RedeemedTimestamp         int    `url:"redeemedTimestamp,omitempty" json:"redeemedTimestamp,omitempty"`
	RedeemedWarehouseID       int    `url:"redeemedWarehouseID,omitempty" json:"redeemedWarehouseID,omitempty"`
	RedeemedPointOfSaleID     int    `url:"redeemedPointOfSaleID,omitempty" json:"redeemedPointOfSaleID,omitempty"`
	RedeemedEmployeeID        int    `url:"redeemedEmployeeID,omitempty" json:"redeemedEmployeeID,omitempty"`
}

func (o *SaveIssuedCouponOpts) MarshalJSON() ([]byte, error) {
	bulkStruct := struct {
		SaveIssuedCouponOpts
		RequestName string `json:"requestName"`
	}{
		SaveIssuedCouponOpts: *o,
		RequestName: "saveIssuedCoupon",
	}

	return json.Marshal(bulkStruct)
}

type IssuedCoupon struct {
	IssuedCouponID             int `json:"issuedCouponID"`
	CouponID                   int `json:"couponID"`
	CouponCode                 string `json:"couponCode"`
	CouponDescription          string `json:"couponDescription"`
	CouponName                 string `json:"couponName"`
	UniqueIdentifier           string `json:"uniqueIdentifier"`
	PrintingCostInRewardPoints int    `json:"printingCostInRewardPoints"`
}

func (s *CouponService) SaveIssuedCoupon(opts *SaveIssuedCouponOpts) (*IssuedCoupon, error) {
	var resSucc ErplyResponse[IssuedCoupon]
	qS, _ := query.Values(opts)
	qS.Add("request", "saveIssuedCoupon")
	_, err := s.client.R().
		SetQueryString(qS.Encode()).
		SetSuccessResult(&resSucc).
		Post("")
	if err != nil {
		return nil, err
	}

	if len(resSucc.Records) > 0 {
		return &resSucc.Records[0], nil
	}

	if resSucc.Status.ResponseStatus == "error" {
		return nil, fmt.Errorf("API Error: %d", resSucc.Status.ErrorCode)
	}

	return nil, fmt.Errorf("No Records Found")
}

func (s *CouponService) BulkSaveIssuedCoupon(opts []SaveIssuedCouponOpts) ([]IssuedCoupon, error) {
	var resSucc ErplyBulkResponse[IssuedCoupon]
	reqString, err := json.Marshal(opts)
	if err != nil {
		return nil, err
	}

	_, err = s.client.R().
		SetQueryParam("requests", string(reqString)).
		SetSuccessResult(&resSucc).
		Post("")

	if err != nil {
		return nil, err
	}
	var reqRes []IssuedCoupon
	for _, res := range resSucc.Requests {
		reqRes = append(reqRes, res.Records...)
	}
	
	return reqRes, nil 
}
