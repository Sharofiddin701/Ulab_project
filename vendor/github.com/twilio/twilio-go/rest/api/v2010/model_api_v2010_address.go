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

// ApiV2010Address struct for ApiV2010Address
type ApiV2010Address struct {
	// The SID of the [Account](https://www.twilio.com/docs/iam/api/account) that is responsible for the Address resource.
	AccountSid *string `json:"account_sid,omitempty"`
	// The city in which the address is located.
	City *string `json:"city,omitempty"`
	// The name associated with the address.This property has a maximum length of 16 4-byte characters, or 21 3-byte characters.
	CustomerName *string `json:"customer_name,omitempty"`
	// The date and time in GMT that the resource was created specified in [RFC 2822](https://www.ietf.org/rfc/rfc2822.txt) format.
	DateCreated *string `json:"date_created,omitempty"`
	// The date and time in GMT that the resource was last updated specified in [RFC 2822](https://www.ietf.org/rfc/rfc2822.txt) format.
	DateUpdated *string `json:"date_updated,omitempty"`
	// The string that you assigned to describe the resource.
	FriendlyName *string `json:"friendly_name,omitempty"`
	// The ISO country code of the address.
	IsoCountry *string `json:"iso_country,omitempty"`
	// The postal code of the address.
	PostalCode *string `json:"postal_code,omitempty"`
	// The state or region of the address.
	Region *string `json:"region,omitempty"`
	// The unique string that that we created to identify the Address resource.
	Sid *string `json:"sid,omitempty"`
	// The number and street address of the address.
	Street *string `json:"street,omitempty"`
	// The URI of the resource, relative to `https://api.twilio.com`.
	Uri *string `json:"uri,omitempty"`
	// Whether emergency calling has been enabled on this number.
	EmergencyEnabled *bool `json:"emergency_enabled,omitempty"`
	// Whether the address has been validated to comply with local regulation. In countries that require valid addresses, an invalid address will not be accepted. `true` indicates the Address has been validated. `false` indicate the country doesn't require validation or the Address is not valid.
	Validated *bool `json:"validated,omitempty"`
	// Whether the address has been verified to comply with regulation. In countries that require valid addresses, an invalid address will not be accepted. `true` indicates the Address has been verified. `false` indicate the country doesn't require verified or the Address is not valid.
	Verified *bool `json:"verified,omitempty"`
	// The additional number and street address of the address.
	StreetSecondary *string `json:"street_secondary,omitempty"`
}