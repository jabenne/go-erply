package legacy

import (
	"encoding/json"
	"fmt"

	"github.com/google/go-querystring/query"
)


type CustomerService service

type RewardPoints struct {
	Points float64 `json:"points"`
	CustomerID int
}

type GetRewardPointOpts struct {
	CustomerID int `url:"customerID" json:"customerID"`
}

func (o *GetRewardPointOpts) MarshalJSON() ([]byte, error) {
	bulkStruct := struct {
		GetRewardPointOpts	
		RequestName string `json:"requestName"`
		RequestID int `json:"requestID"`
	}{
		GetRewardPointOpts: *o,
		RequestName: "getCustomerRewardPoints",
		RequestID: o.CustomerID,
	}
	return json.Marshal(bulkStruct)
}

func (s *CustomerService) GetRewardPoints(opts *GetRewardPointOpts) (*RewardPoints, error) {
	var resSucc ErplyResponse[RewardPoints]
	qS, _ := query.Values(opts)
	qS.Add("request", "getCustomerRewardPoints")
	_, err := s.client.R().
		SetQueryString(qS.Encode()).
		SetSuccessResult(&resSucc).
		Post("")

	if err != nil {
		return nil, err 
	}
	
	if len(resSucc.Records) > 0 {
		resSucc.Records[0].CustomerID = opts.CustomerID
		return &resSucc.Records[0], nil
	}

	return nil, fmt.Errorf("No Records Found")
}

func (s *CustomerService) BulkGetRewardPoints(opts []GetRewardPointOpts) ([]RewardPoints, error) {
	var resSucc ErplyBulkResponse[RewardPoints]
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
	var reqRes []RewardPoints
	for _, res := range resSucc.Requests {
		if len(res.Records) > 0 {
			res.Records[0].CustomerID = res.Status.RequestID
			reqRes = append(reqRes, res.Records...)
		}
	}
	
	return reqRes, nil 
}

type SubtractedRewardPoints struct {
	CustomerID int `json:"customerID"`
	RemainingPoints float64 `json:"remainingPoints"`
	SubtractedPoints float64 `json:"subtractedPoints"`
	TransactionID int `json:"transactionID"`
}

type SubtractRewardPointsOpts struct {
	CustomerID int `url:"customerID,omitempty" json:"customerID,omitempty"`
	CampaignID int `url:"campaignID,omitempty" json:"campaignID,omitempty"` 
	WarehouseID int `url:"warehouseID,omitempty" json:"warehouseID,omitempty"` 
	SalepointID int `url:"salepointID,omitempty" json:"salepointID,omitempty"` 
	SalespersonID int `url:"salespersonID,omitempty" json:"salespersonID,omitempty"` 
	Points float64 `url:"points,omitempty" json:"points,omitempty"` 
	SubtractedUnixTime int `url:"subtractedUnixTime,omitempty" json:"subtractedUnixTime,omitempty"` 
	IssuedCouponID int `url:"issuedCouponID,omitempty" json:"issuedCouponID,omitempty"` 
	Description string `url:"description,omitempty" json:"description,omitempty"` 
}

func (o *SubtractRewardPointsOpts) MarshalJSON() ([]byte, error) {
	bulkStruct := struct {
		SubtractRewardPointsOpts	
		RequestName string `json:"requestName"`
		RequestID int `json:"requestID"`
	}{
		SubtractRewardPointsOpts: *o,
		RequestName: "subtractCustomerRewardPoints",
		RequestID: o.CustomerID,
	}
	return json.Marshal(bulkStruct)
}
	
func (s *CustomerService) SubtractRewardPoints(opts *SubtractRewardPointsOpts) (*SubtractedRewardPoints, error) {
	var resSucc ErplyResponse[SubtractedRewardPoints]
	qS, _ := query.Values(opts)
	qS.Add("request", "subtractCustomerRewardPoints")
	_, err := s.client.R().
		SetQueryString(qS.Encode()).
		SetSuccessResult(&resSucc).
		Post("")
	if err != nil {
		return nil, err 
	}

	if resSucc.Status.ResponseStatus == "error" {
		return nil, fmt.Errorf("API Error: %d", resSucc.Status.ErrorCode)
	}
	
	if len(resSucc.Records) > 0 {
		return &resSucc.Records[0], nil
	}

	return nil, fmt.Errorf("No Records Found")
}

func (s *CustomerService) BulkSubtractRewardPoints(opts []SubtractRewardPointsOpts) ([]SubtractedRewardPoints, error) {
	var resSucc ErplyBulkResponse[SubtractedRewardPoints]
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
	var reqRes []SubtractedRewardPoints
	for _, res := range resSucc.Requests {
		reqRes = append(reqRes, res.Records...)
	}
	
	return reqRes, nil 
}

type AddRewardPointsOpts struct {
	CustomerID int `url:"customerID,omitempty" json:"customerID,omitempty"`
	InvoiceID int `url:"invoiceID,omitempty" json:"invoiceID,omitempty"` 
	Points float64 `url:"points,omitempty" json:"points,omitempty"` 
	CreatedUnixTime int `url:"createdUnixTime,omitempty" json:"createdUnixTime,omitempty"` 
	ExpiryUnixTime int `url:"expiryUnixTime,omitempty" json:"expiryUnixTime,omitempty"` 
	EmployeeID int `url:"employeeID,omitempty" json:"employeeID,omitempty"` 
	Description string `url:"description,omitempty" json:"description,omitempty"` 
}

func (o *AddRewardPointsOpts) MarshalJSON() ([]byte, error) {
	bulkStruct := struct {
		AddRewardPointsOpts	
		RequestName string `json:"requestName"`
		RequestID int `json:"requestID"`
	}{
		AddRewardPointsOpts: *o,
		RequestName: "addCustomerRewardPoints",
		RequestID: o.CustomerID,
	}
	return json.Marshal(bulkStruct)
}

func (s *CustomerService) AddRewardPoints(opts *AddRewardPointsOpts) (*RewardPoints, error) {
	var resSucc ErplyResponse[RewardPoints]
	qS, _ := query.Values(opts)
	qS.Add("request", "addCustomerRewardPoints")
	_, err := s.client.R().
		SetQueryString(qS.Encode()).
		SetSuccessResult(&resSucc).
		Post("")
	if err != nil {
		return nil, err 
	}

	if resSucc.Status.ResponseStatus == "error" {
		return nil, fmt.Errorf("API Error: %d", resSucc.Status.ErrorCode)
	}
	
	if len(resSucc.Records) > 0 {
		return &resSucc.Records[0], nil
	}

	return nil, fmt.Errorf("No Records Found")
}

func (s *CustomerService) BulkAddRewardPoints(opts []AddRewardPointsOpts) ([]RewardPoints, error) {
	var resSucc ErplyBulkResponse[RewardPoints]
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
	var reqRes []RewardPoints
	for _, res := range resSucc.Requests {
		if len(res.Records) > 0 {
			reqRes = append(reqRes, res.Records...)
		} 
	}
	
	return reqRes, nil 
}
