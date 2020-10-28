// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package analytics

import "testing"

func TestShouldFilter(t *testing.T) {
	record := AnalyticsRecord{
		Username: "colin",
	}

	// test skip_usernames
	filter := AnalyticsFilters{
		SkippedUsernames: []string{"colin"},
	}
	shouldFilter := filter.ShouldFilter(record)
	if shouldFilter == false {
		t.Fatal("filter should be filtering the record")
	}

	// test different usernames
	filter = AnalyticsFilters{
		Usernames: []string{"james"},
	}
	shouldFilter = filter.ShouldFilter(record)
	if shouldFilter == false {
		t.Fatal("filter should be filtering the record")
	}

	// test no filter
	filter = AnalyticsFilters{}
	shouldFilter = filter.ShouldFilter(record)
	if shouldFilter == true {
		t.Fatal("filter should not be filtering the record")
	}
}

func TestHasFilter(t *testing.T) {
	filter := AnalyticsFilters{}

	hasFilter := filter.HasFilter()
	if hasFilter == true {
		t.Fatal("Has filter should be false.")
	}

	filter = AnalyticsFilters{
		Usernames: []string{"colin"},
	}
	hasFilter = filter.HasFilter()
	if hasFilter == false {
		t.Fatal("HasFilter should be true.")
	}
}
