// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Salesforce displays some IDs in the old 15 digit case sensitive format (such as in the url)
// if a user pastes the old format in their config it leads to permadiffs with the ids read from the API
// the conversion source code is posted here
// https://help.salesforce.com/articleView?id=000319308&type=1&mode=1
func normalizeId(id string) string {
	if len(id) != 15 {
		// if the string is empty or already 18 characters, or not a proper id, just return it
		// let it error upstream
		return id
	}
	var addon string
	for block := 0; block < 3; block++ {
		loop := 0
		for position := 0; position < 5; position++ {
			current := id[block*5+position]
			if current >= 'A' && current <= 'Z' {
				loop += 1 << position
			}
		}
		addon += string("ABCDEFGHIJKLMNOPQRSTUVWXYZ012345"[loop])
	}
	return id + addon
}

type NormalizeId struct{}

func (NormalizeId) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.StateValue.IsNull() || req.StateValue.IsUnknown() {
		return
	}
	plan := req.PlanValue.ValueString()
	state := req.StateValue.ValueString()
	if normalizeId(plan) == normalizeId(state) {
		resp.PlanValue = req.StateValue
	} else {
		resp.PlanValue = req.PlanValue
	}
}

type resourceDefaults struct {
	defaults map[string]attr.Value
}

func (r resourceDefaults) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	resp.PlanValue = req.PlanValue
	if req.ConfigValue.IsNull() {
		def, ok := r.defaults[req.Path.String()]
		if ok {
			resp.PlanValue = def.(types.String)
		}
	}
}

type staticComputed struct{}

func (staticComputed) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.StateValue.IsNull() {
		return
	}
	resp.PlanValue = req.StateValue
}

type fixNullToUnknown struct{}

func (fixNullToUnknown) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.StateValue.IsNull() {
		return
	}
	state := req.StateValue
	config := req.ConfigValue
	if config.IsUnknown() && state.IsNull() {
		resp.PlanValue = config
	}
}

type booleanNilIsFalse struct{}

func (booleanNilIsFalse) PlanModifyBool(ctx context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
	resp.PlanValue = req.PlanValue
	if req.ConfigValue.IsNull() {
		resp.PlanValue = types.BoolValue(false)
	}
}
