// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package docs

import (
	v1 "github.com/marmotedu/api/apiserver/v1"
	metav1 "github.com/marmotedu/component-base/pkg/meta/v1"

	"github.com/marmotedu/iam/internal/apiserver/controller/v1/user"
)

// swagger:route POST /users users createUserRequest
//
// Create users.
//
// Create users according to input parameters.
//
//     Security:
//       api_key:
//
//     Responses:
//       default: errResponse
//       200: createUserResponse

// swagger:route DELETE /users/{name} users deleteUserRequest
//
// Delete user.
//
// Delete user according to input parameters.
//
//     Responses:
//       default: errResponse
//       200: okResponse

// swagger:route DELETE /users users deleteUserCollectionRequest
//
// Batch delete user.
//
// Delete users
//
//     Security:
//       api_key:
//
//     Responses:
//       default: errResponse
//       200: okResponse

// swagger:route PUT /users/{name} users updateUserRequest
//
// Update user.
//
// Update user according to input parameters.
//
//     Security:
//       api_key:
//
//     Responses:
//       default: errResponse
//       200: updateUserResponse

// swagger:route PUT /users/{name}/change_password users changePasswordRequest
//
// Change user password.
//
// Change user password.
//
//     Security:
//       api_key:
//
//     Responses:
//       default: errResponse
//       200: okResponse

// swagger:route GET /users/{name} users getUserRequest
//
// Get details for specified user.
//
// Get details for specified user according to input parameters.
//
//     Responses:
//       default: errResponse
//       200: getUserResponse

// swagger:route GET /users users listUserRequest
//
// List users.
//
// List users.
//
//     Responses:
//       default: errResponse
//       200: listUserResponse

// List users request.
// swagger:parameters listUserRequest
type listUserRequestParamsWrapper struct {
	// in:query
	metav1.ListOptions
}

// List users response.
// swagger:response listUserResponse
type listUserResponseWrapper struct {
	// in:body
	Body v1.UserList
}

// User response.
// swagger:response createUserResponse
type createUserResponseWrapper struct {
	// in:body
	Body v1.User
}

// User response.
// swagger:response updateUserResponse
type updateUserResponseWrapper struct {
	// in:body
	Body v1.User
}

// User response.
// swagger:response getUserResponse
type getUserResponseWrapper struct {
	// in:body
	Body v1.User
}

// swagger:parameters createUserRequest updateUserRequest
type userRequestParamsWrapper struct {
	// User information.
	// in:body
	Body v1.User
}

// swagger:parameters deleteUserRequest getUserRequest updateUserRequest
type userNameParamsWrapper struct {
	// User name.
	// in:path
	Name string `json:"name"`
}

// Batch delete users.
// swagger:parameters deleteUserCollectionRequest deletePolicyCollectionRequest
type deleteCollectionRequestParamsWrapper struct {
	// in:query
	Names []string `json:"name"`
}

// Change user password.
// swagger:parameters changePasswordRequest
type changePasswordRequestParamsWrapper struct {
	// The name of user.
	// in:path
	Name string `json:"name"`

	// in:body
	Body user.ChangePasswordRequest
}
