package erply

type CRMService service

type PatchCustomerOpts struct {}

type BulkPatchCustomerOpts struct {
    ResourceID int `json:"resourceID"`
   	PatchCustomerOpts 
}

type BulkRequest[T any] struct {
    Requests []T `json:"requests"`
}

type BulkResponse struct {
    Results []BulkResult `json:"result"`  
}

type BulkResult struct {
    Message string `json:"message"`
    ResourceID int `json:"resourceID"`
    ResultID int `json:"resultID"`
}

func (s *CRMService)BulkPatchCustomer(opts []BulkPatchCustomerOpts) ([]BulkResult, error) {
    var res BulkResponse
    body := BulkRequest[BulkPatchCustomerOpts] { 
        Requests: opts,
    }
 	_, err := s.client.R().
        SetBody(&body).
        SetSuccessResult(&res).
        SetPathParam("api", "crm").
		Patch("/customers/individuals/bulk")
	if err != nil {
		return nil, err 
	}
   
    return res.Results, err
} 
