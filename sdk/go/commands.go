// Copyright (c) 2026 Tencent Inc.
// SPDX-License-Identifier: Apache-2.0

package cubesandbox

import (
	"context"
	"strconv"
	"strings"
)

type codeRunner interface {
	RunCode(context.Context, string, RunCodeOptions) (*Execution, error)
}

type Commands struct {
	runner codeRunner
}

func (c *Commands) Run(ctx context.Context, cmd string, opts CommandOptions) (*CommandResult, error) {
	code := "import subprocess as _sp\n" +
		"_r = _sp.run(" + strconv.Quote(cmd) + ", shell=True, capture_output=True, text=True)\n" +
		"import sys as _sys\n" +
		"_sys.stdout.write(_r.stdout)\n" +
		"_sys.stderr.write(_r.stderr)\n" +
		"print(_r.returncode)\n"

	var stdoutParts []string
	execution, err := c.runner.RunCode(ctx, code, RunCodeOptions{
		Timeout: opts.Timeout,
		OnStdout: func(message OutputMessage) {
			stdoutParts = append(stdoutParts, message.Text)
		},
	})
	if err != nil {
		return nil, err
	}

	allStdout := strings.Join(stdoutParts, "")
	if allStdout == "" {
		allStdout = strings.Join(execution.Logs.Stdout, "")
	}
	lines := splitLines(allStdout)

	stdout := allStdout
	exitCode := 0
	if len(lines) > 0 && isIntegerLine(strings.TrimSpace(lines[len(lines)-1])) {
		parsed, _ := strconv.Atoi(strings.TrimSpace(lines[len(lines)-1]))
		exitCode = parsed
		stdoutLines := lines[:len(lines)-1]
		stdout = strings.Join(stdoutLines, "\n")
		if len(stdoutLines) > 0 {
			stdout += "\n"
		}
	} else if execution.Error != nil {
		exitCode = 1
	}

	return &CommandResult{
		Stdout:   stdout,
		Stderr:   strings.Join(execution.Logs.Stderr, ""),
		ExitCode: exitCode,
	}, nil
}

func splitLines(s string) []string {
	if s == "" {
		return nil
	}
	normalized := strings.ReplaceAll(s, "\r\n", "\n")
	normalized = strings.ReplaceAll(normalized, "\r", "\n")
	lines := strings.Split(normalized, "\n")
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	return lines
}

func isIntegerLine(s string) bool {
	if s == "" {
		return false
	}
	if s[0] == '-' {
		s = s[1:]
	}
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
