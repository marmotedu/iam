/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

package log_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/marmotedu/iam/pkg/log"
)

func Test_Options_Validate(t *testing.T) {
	opts := &log.Options{
		Level:            "test",
		Format:           "test",
		EnableColor:      true,
		DisableCaller:    false,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	errs := opts.Validate()
	expected := `[unrecognized level: "test" not a valid log format: "test"]`
	assert.Equal(t, expected, fmt.Sprintf("%s", errs))
}
