// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package docs

import (
	v1 "github.com/marmotedu/api/apiserver/v1"
	metav1 "github.com/marmotedu/component-base/pkg/meta/v1"
)

// swagger:route POST /policies policies createPolicyRequest
//
// Create policies.
//
// Create policies according to input parameters.
//
//     Security:
//       api_key:
//
//     Responses:
//       default: errResponse
//       200: createPolicyResponse

// swagger:route DELETE /policies/{name} policies deletePolicyRequest
//
// Delete policy.
//
// Delete policy according to input parameters.
//
//     Security:
//       api_key:
//
//     Responses:
//       default: errResponse
//       200: okResponse

// swagger:route DELETE /policies policies deletePolicyCollectionRequest
//
// Batch delete policies.
//
// Batch delete policies according to input parameters.
//
//     Security:
//       api_key:
//
//     Responses:
//       default: errResponse
//       200: okResponse

// swagger:route PUT /policies/{name} policies updatePolicyRequest
//
// Update policy.
//
// Update policy according to input parameters.
//
//     Security:
//       api_key:
//
//     Responses:
//       default: errResponse
//       200: updatePolicyResponse

// swagger:route GET /policies/{name} policies getPolicyRequest
//
// Get details for specified policy.
//
// Get details for specified policy according to input parameters.
//
//     Responses:
//       default: errResponse
//       200: getPolicyResponse

// swagger:route GET /policies policies listPolicyRequest
//
// List policies.
//
// List policies.
//
//     Responses:
//       default: errResponse
//       200: listPolicyResponse

// List users request.
// swagger:parameters listPolicyRequest
type listPolicyRequestParamsWrapper struct {
	// in:query
	metav1.ListOptions
}

// List policies response.
// swagger:response listPolicyResponse
type listPolicyResponseWrapper struct {
	// in:body
	Body v1.PolicyList
}

// Policy response.
// swagger:response createPolicyResponse
type createPolicyResponseWrapper struct {
	// in:body
	Body v1.Policy
}

// Policy response.
// swagger:response updatePolicyResponse
type updatePolicyResponseWrapper struct {
	// in:body
	Body v1.Policy
}

// Policy response.
// swagger:response getPolicyResponse
type getPolicyResponseWrapper struct {
	// in:body
	Body v1.Policy
}

// swagger:parameters createPolicyRequest updatePolicyRequest
type policyRequestParamsWrapper struct {
	// Policy information.
	// in:body
	Body v1.Policy
}

// swagger:parameters deletePolicyRequest getPolicyRequest updatePolicyRequest
type policyNameParamsWrapper struct {
	// Policy name.
	// in:path
	Name string `json:"name"`
}
