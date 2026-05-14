// Copyright (c) 2026 Tencent Inc.
// SPDX-License-Identifier: Apache-2.0

package cubesandbox

import (
	"context"
	"fmt"
	"strconv"
)

type Files struct {
	runner codeRunner
}

func (f *Files) Read(ctx context.Context, path string) (string, error) {
	execution, err := f.runner.RunCode(ctx, "open("+strconv.Quote(path)+").read()", RunCodeOptions{})
	if err != nil {
		return "", err
	}
	if execution.Error != nil {
		return "", fmt.Errorf("Failed to read %s: %s", path, execution.Error.Value)
	}
	return execution.mainText(), nil
}
