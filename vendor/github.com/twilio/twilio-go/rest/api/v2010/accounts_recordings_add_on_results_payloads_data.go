/*
 * This code was generated by
 * ___ _ _ _ _ _    _ ____    ____ ____ _    ____ ____ _  _ ____ ____ ____ ___ __   __
 *  |  | | | | |    | |  | __ |  | |__| | __ | __ |___ |\ | |___ |__/ |__|  | |  | |__/
 *  |  |_|_| | |___ | |__|    |__| |  | |    |__] |___ | \| |___ |  \ |  |  | |__| |  \
 *
 * Twilio - Api
 * This is the public Twilio REST API.
 *
 * NOTE: This class is auto generated by OpenAPI Generator.
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */

package openapi

import (
	"encoding/json"
	"net/url"
	"strings"
)

// Optional parameters for the method 'FetchRecordingAddOnResultPayloadData'
type FetchRecordingAddOnResultPayloadDataParams struct {
	// The SID of the [Account](https://www.twilio.com/docs/iam/api/account) that created the Recording AddOnResult Payload resource to fetch.
	PathAccountSid *string `json:"PathAccountSid,omitempty"`
}

func (params *FetchRecordingAddOnResultPayloadDataParams) SetPathAccountSid(PathAccountSid string) *FetchRecordingAddOnResultPayloadDataParams {
	params.PathAccountSid = &PathAccountSid
	return params
}

// Fetch an instance of a result payload
func (c *ApiService) FetchRecordingAddOnResultPayloadData(ReferenceSid string, AddOnResultSid string, PayloadSid string, params *FetchRecordingAddOnResultPayloadDataParams) (*ApiV2010RecordingAddOnResultPayloadData, error) {
	path := "/2010-04-01/Accounts/{AccountSid}/Recordings/{ReferenceSid}/AddOnResults/{AddOnResultSid}/Payloads/{PayloadSid}/Data.json"
	if params != nil && params.PathAccountSid != nil {
		path = strings.Replace(path, "{"+"AccountSid"+"}", *params.PathAccountSid, -1)
	} else {
		path = strings.Replace(path, "{"+"AccountSid"+"}", c.requestHandler.Client.AccountSid(), -1)
	}
	path = strings.Replace(path, "{"+"ReferenceSid"+"}", ReferenceSid, -1)
	path = strings.Replace(path, "{"+"AddOnResultSid"+"}", AddOnResultSid, -1)
	path = strings.Replace(path, "{"+"PayloadSid"+"}", PayloadSid, -1)

	data := url.Values{}
	headers := map[string]interface{}{
		"Content-Type": "application/x-www-form-urlencoded",
	}

	resp, err := c.requestHandler.Get(c.baseURL+path, data, headers)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	ps := &ApiV2010RecordingAddOnResultPayloadData{}
	if err := json.NewDecoder(resp.Body).Decode(ps); err != nil {
		return nil, err
	}

	return ps, err
}