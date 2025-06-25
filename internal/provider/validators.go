// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type notEmptyString struct{}

func (notEmptyString) Description(ctx context.Context) string {
	return "Ensures the string is not empty."
}

func (notEmptyString) MarkdownDescription(ctx context.Context) string {
	return "Ensures the string is not empty."
}

func (notEmptyString) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsUnknown() {
		return
	}
	if req.ConfigValue.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Empty String",
			"Value must not be empty.",
		)
	}
}

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func isEmailValid(e string) bool {
	if len(e) < 3 || len(e) > 254 {
		return false
	}
	return emailRegex.MatchString(e)
}

type email struct{}

func (email) Description(ctx context.Context) string {
	return "Ensures the string is a valid email address."
}

func (email) MarkdownDescription(ctx context.Context) string {
	return "Ensures the string is a valid email address."
}

func (email) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsUnknown() {
		return
	}
	if !isEmailValid(req.ConfigValue.ValueString()) {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid email address",
			"Value must be a valid email address.",
		)
	}
}

type stringInSlice struct {
	slice    []string
	optional bool
}

func (s stringInSlice) Description(ctx context.Context) string {
	return fmt.Sprintf("Ensures the string is one of: [%s]", strings.Join(s.slice, ", "))
}

func (s stringInSlice) MarkdownDescription(ctx context.Context) string {
	return s.Description(ctx)
}

func (s stringInSlice) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsUnknown() {
		return
	}
	if s.optional && req.ConfigValue.IsNull() {
		return
	}
	for _, item := range s.slice {
		if req.ConfigValue.ValueString() == item {
			return
		}
	}
	resp.Diagnostics.AddAttributeError(
		req.Path,
		"Invalid string",
		fmt.Sprintf("String must be one of: [%s]", strings.Join(s.slice, ", ")),
	)
}
