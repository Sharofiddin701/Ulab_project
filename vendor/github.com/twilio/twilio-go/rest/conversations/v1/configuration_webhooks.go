/*
 * This code was generated by
 * ___ _ _ _ _ _    _ ____    ____ ____ _    ____ ____ _  _ ____ ____ ____ ___ __   __
 *  |  | | | | |    | |  | __ |  | |__| | __ | __ |___ |\ | |___ |__/ |__|  | |  | |__/
 *  |  |_|_| | |___ | |__|    |__| |  | |    |__] |___ | \| |___ |  \ |  |  | |__| |  \
 *
 * Twilio - Conversations
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
)

//
func (c *ApiService) FetchConfigurationWebhook() (*ConversationsV1ConfigurationWebhook, error) {
	path := "/v1/Configuration/Webhooks"

	data := url.Values{}
	headers := map[string]interface{}{
		"Content-Type": "application/x-www-form-urlencoded",
	}

	resp, err := c.requestHandler.Get(c.baseURL+path, data, headers)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	ps := &ConversationsV1ConfigurationWebhook{}
	if err := json.NewDecoder(resp.Body).Decode(ps); err != nil {
		return nil, err
	}

	return ps, err
}

// Optional parameters for the method 'UpdateConfigurationWebhook'
type UpdateConfigurationWebhookParams struct {
	// The HTTP method to be used when sending a webhook request.
	Method *string `json:"Method,omitempty"`
	// The list of webhook event triggers that are enabled for this Service: `onMessageAdded`, `onMessageUpdated`, `onMessageRemoved`, `onConversationUpdated`, `onConversationRemoved`, `onParticipantAdded`, `onParticipantUpdated`, `onParticipantRemoved`
	Filters *[]string `json:"Filters,omitempty"`
	// The absolute url the pre-event webhook request should be sent to.
	PreWebhookUrl *string `json:"PreWebhookUrl,omitempty"`
	// The absolute url the post-event webhook request should be sent to.
	PostWebhookUrl *string `json:"PostWebhookUrl,omitempty"`
	//
	Target *string `json:"Target,omitempty"`
}

func (params *UpdateConfigurationWebhookParams) SetMethod(Method string) *UpdateConfigurationWebhookParams {
	params.Method = &Method
	return params
}
func (params *UpdateConfigurationWebhookParams) SetFilters(Filters []string) *UpdateConfigurationWebhookParams {
	params.Filters = &Filters
	return params
}
func (params *UpdateConfigurationWebhookParams) SetPreWebhookUrl(PreWebhookUrl string) *UpdateConfigurationWebhookParams {
	params.PreWebhookUrl = &PreWebhookUrl
	return params
}
func (params *UpdateConfigurationWebhookParams) SetPostWebhookUrl(PostWebhookUrl string) *UpdateConfigurationWebhookParams {
	params.PostWebhookUrl = &PostWebhookUrl
	return params
}
func (params *UpdateConfigurationWebhookParams) SetTarget(Target string) *UpdateConfigurationWebhookParams {
	params.Target = &Target
	return params
}

//
func (c *ApiService) UpdateConfigurationWebhook(params *UpdateConfigurationWebhookParams) (*ConversationsV1ConfigurationWebhook, error) {
	path := "/v1/Configuration/Webhooks"

	data := url.Values{}
	headers := map[string]interface{}{
		"Content-Type": "application/x-www-form-urlencoded",
	}

	if params != nil && params.Method != nil {
		data.Set("Method", *params.Method)
	}
	if params != nil && params.Filters != nil {
		for _, item := range *params.Filters {
			data.Add("Filters", item)
		}
	}
	if params != nil && params.PreWebhookUrl != nil {
		data.Set("PreWebhookUrl", *params.PreWebhookUrl)
	}
	if params != nil && params.PostWebhookUrl != nil {
		data.Set("PostWebhookUrl", *params.PostWebhookUrl)
	}
	if params != nil && params.Target != nil {
		data.Set("Target", *params.Target)
	}

	resp, err := c.requestHandler.Post(c.baseURL+path, data, headers)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	ps := &ConversationsV1ConfigurationWebhook{}
	if err := json.NewDecoder(resp.Body).Decode(ps); err != nil {
		return nil, err
	}

	return ps, err
}