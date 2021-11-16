// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package docs

import (
	v1 "github.com/marmotedu/api/apiserver/v1"
	"github.com/marmotedu/component-base/pkg/core"
	metav1 "github.com/marmotedu/component-base/pkg/meta/v1"
)

// swagger:route POST /secrets secrets createSecretRequest
//
// Create secrets.
//
// Create secrets according to input parameters.
//
//     Security:
//       api_key:
//
//     Responses:
//       default: errResponse
//       200: createSecretResponse

// swagger:route DELETE /secrets/{name} secrets deleteSecretRequest
//
// Delete secret.
//
// Delete secret according to input parameters.
//
//     Security:
//       api_key:
//
//     Responses:
//       default: errResponse
//       200: okResponse

// swagger:route PUT /secrets/{name} secrets updateSecretRequest
//
// Update secret.
//
// Update secret according to input parameters.
//
//     Security:
//       api_key:
//
//     Responses:
//       default: errResponse
//       200: updateSecretResponse

// swagger:route GET /secrets/{name} secrets getSecretRequest
//
// Get details for specified secret.
//
// Get details for specified secret according to input parameters.
//
//     Responses:
//       default: errResponse
//       200: getSecretResponse

// swagger:route GET /secrets secrets listSecretRequest
//
// List secrets.
//
// List secrets.
//
//     Responses:
//       default: errResponse
//       200: listSecretResponse

// List users request.
// swagger:parameters listSecretRequest
type listSecretRequestParamsWrapper struct {
	// in:query
	metav1.ListOptions
}

// List secrets response.
// swagger:response listSecretResponse
type listSecretResponseWrapper struct {
	// in:body
	Body v1.SecretList
}

// Secret response.
// swagger:response createSecretResponse
type createSecretResponseWrapper struct {
	// in:body
	Body v1.Secret
}

// Secret response.
// swagger:response updateSecretResponse
type updateSecretResponseWrapper struct {
	// in:body
	Body v1.Secret
}

// Secret response.
// swagger:response getSecretResponse
type getSecretResponseWrapper struct {
	// in:body
	Body v1.Secret
}

// swagger:parameters createSecretRequest updateSecretRequest
type secretRequestParamsWrapper struct {
	// Secret information.
	// in:body
	Body v1.Secret
}

// swagger:parameters deleteSecretRequest getSecretRequest updateSecretRequest
type secretNameParamsWrapper struct {
	// Secret name.
	// in:path
	Name string `json:"name"`
}

// ErrResponse defines the return messages when an error occurred.
// swagger:response errResponse
type errResponseWrapper struct {
	// in:body
	Body core.ErrResponse
}

// Return nil json object.
// swagger:response okResponse
type okResponseWrapper struct{}
