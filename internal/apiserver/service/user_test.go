// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package service

import (
	"context"
	"os"
	"testing"

	"github.com/AlekSi/pointer"
	metav1 "github.com/marmotedu/component-base/pkg/meta/v1"

	"github.com/marmotedu/iam/internal/apiserver/store"
	"github.com/marmotedu/iam/internal/apiserver/store/fake"
)

func TestMain(m *testing.M) {
	fakeStore, _ := fake.NewFakeStore()
	store.SetClient(fakeStore)
	os.Exit(m.Run())
}

func TestListUser(t *testing.T) {
	var limit int64 = 3
	opts := metav1.ListOptions{
		Offset: pointer.ToInt64(0),
		Limit:  pointer.ToInt64(limit),
	}

	got, err := ListUser(context.TODO(), opts)
	if err != nil {
		t.Errorf("ListUser() error = %v, wantErr %v", err, nil)
	}

	if got.TotalCount != fake.ResourceCount {
		t.Errorf("ListUser() TotalCount = %v, want %v", got.TotalCount, fake.ResourceCount)
	}

	if len(got.Items) != int(limit) {
		t.Errorf("len(UserListV2.Items)= %v, want %v", len(got.Items), limit)
	}

	if got.Items[0].Name != "user1" {
		t.Errorf("ListUser() User[0]= %v, want %v", got.Items[0].Name, "user1")
	}

	if got.Items[1].Name != "user2" {
		t.Errorf("ListUser() User[1]= %v, want %v", got.Items[1].Name, "user2")
	}

	if got.Items[2].Name != "user3" {
		t.Errorf("ListUser() User[2]= %v, want %v", got.Items[2].Name, "user3")
	}
}

func BenchmarkListUser(b *testing.B) {
	opts := metav1.ListOptions{
		Offset: pointer.ToInt64(0),
		Limit:  pointer.ToInt64(50),
	}

	for i := 0; i < b.N; i++ {
		// _, _ = ListUserBadPerformance(opts)
		_, _ = ListUser(context.TODO(), opts)
	}
}
