// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package v1

import (
	"context"

	v1 "github.com/marmotedu/api/apiserver/v1"
	metav1 "github.com/marmotedu/component-base/pkg/meta/v1"
	"github.com/marmotedu/errors"

	"github.com/marmotedu/iam/internal/apiserver/store"
	"github.com/marmotedu/iam/internal/pkg/code"
)

// PolicySrv defines functions used to handle policy request.
type PolicySrv interface {
	Create(ctx context.Context, policy *v1.Policy, opts metav1.CreateOptions) error
	Update(ctx context.Context, policy *v1.Policy, opts metav1.UpdateOptions) error
	Delete(ctx context.Context, username string, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, username string, names []string, opts metav1.DeleteOptions) error
	Get(ctx context.Context, username string, name string, opts metav1.GetOptions) (*v1.Policy, error)
	List(ctx context.Context, username string, opts metav1.ListOptions) (*v1.PolicyList, error)
}

type policyService struct {
	store store.Factory
}

var _ PolicySrv = (*policyService)(nil)

func newPolicies(srv *service) *policyService {
	return &policyService{store: srv.store}
}

func (s *policyService) Create(ctx context.Context, policy *v1.Policy, opts metav1.CreateOptions) error {
	if err := s.store.Policies().Create(ctx, policy, opts); err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}

	return nil
}

func (s *policyService) Update(ctx context.Context, policy *v1.Policy, opts metav1.UpdateOptions) error {
	// Save changed fields.
	if err := s.store.Policies().Update(ctx, policy, opts); err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}

	return nil
}

func (s *policyService) Delete(ctx context.Context, username, name string, opts metav1.DeleteOptions) error {
	if err := s.store.Policies().Delete(ctx, username, name, opts); err != nil {
		return err
	}

	return nil
}

func (s *policyService) DeleteCollection(
	ctx context.Context,
	username string,
	names []string,
	opts metav1.DeleteOptions,
) error {
	if err := s.store.Policies().DeleteCollection(ctx, username, names, opts); err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}

	return nil
}

func (s *policyService) Get(ctx context.Context, username, name string, opts metav1.GetOptions) (*v1.Policy, error) {
	policy, err := s.store.Policies().Get(ctx, username, name, opts)
	if err != nil {
		return nil, err
	}

	return policy, nil
}

func (s *policyService) List(ctx context.Context, username string, opts metav1.ListOptions) (*v1.PolicyList, error) {
	policies, err := s.store.Policies().List(ctx, username, opts)
	if err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}

	return policies, nil
}
