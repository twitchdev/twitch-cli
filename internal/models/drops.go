// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package models

type DropsEntitlementEventSubResponse struct {
	Subscription EventsubSubscription            `json:"subscription"`
	Events       []DropsEntitlementEventSubEvent `json:"events"`
}

type DropsEntitlementEventSubEvent struct {
	ID   string                            `json:"id"`
	Data DropsEntitlementEventSubEventData `json:"data"`
}
type DropsEntitlementEventSubEventData struct {
	EntitlementID  string `json:"entitlement_id"`
	BenefitID      string `json:"benefit_id"`
	CampaignID     string `json:"campaign_id"`
	OrganizationID string `json:"organization_id"`
	CreatedAt      string `json:"created_at"`
	UserID         string `json:"user_id"`
	UserName       string `json:"user_name"`
	UserLogin      string `json:"user_login"`
	CategoryID     string `json:"category_id"`
	CategoryName   string `json:"category_name"`
}
