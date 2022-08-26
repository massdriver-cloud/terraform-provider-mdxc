/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 * Original repo: https://github.com/hashicorp/terraform-provider-google
 */

package thirdparty

import (
	"fmt"
	"sort"

	"github.com/hashicorp/errwrap"
	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/googleapi"
)

type conditionKey struct {
	Description string
	Expression  string
	Title       string
}

type iamBindingKey struct {
	Role      string
	Condition conditionKey
}

func (k conditionKey) Empty() bool {
	return k == conditionKey{}
}

func (k conditionKey) String() string {
	return fmt.Sprintf("%s/%s/%s", k.Title, k.Description, k.Expression)
}

func AddBinding(bindings []*cloudresourcemanager.Binding, binding *cloudresourcemanager.Binding) []*cloudresourcemanager.Binding {
	return mergeBindings(append(bindings, binding))
}

func RemoveBinding(bindings []*cloudresourcemanager.Binding, binding *cloudresourcemanager.Binding) []*cloudresourcemanager.Binding {
	return subtractFromBindings(bindings, binding)
}

// Removes given role+condition/bound-member pairs from the given Bindings (i.e subtraction).
func subtractFromBindings(bindings []*cloudresourcemanager.Binding, toRemove ...*cloudresourcemanager.Binding) []*cloudresourcemanager.Binding {
	currMap := createIamBindingsMap(bindings)
	toRemoveMap := createIamBindingsMap(toRemove)

	for key, removeSet := range toRemoveMap {
		members, ok := currMap[key]
		if !ok {
			continue
		}
		// Remove all removed members
		for m := range removeSet {
			delete(members, m)
		}
		// Remove role+condition from bindings
		if len(members) == 0 {
			delete(currMap, key)
		}
	}

	return listFromIamBindingMap(currMap)
}

// Flattens a list of Bindings so each role+condition has a single Binding with combined members
func mergeBindings(bindings []*cloudresourcemanager.Binding) []*cloudresourcemanager.Binding {
	bm := createIamBindingsMap(bindings)
	return listFromIamBindingMap(bm)
}

func createIamBindingsMap(bindings []*cloudresourcemanager.Binding) map[iamBindingKey]map[string]struct{} {
	bm := make(map[iamBindingKey]map[string]struct{})
	// Get each binding
	for _, b := range bindings {
		members := make(map[string]struct{})
		key := iamBindingKey{b.Role, conditionKeyFromCondition(b.Condition)}
		// Initialize members map
		if _, ok := bm[key]; ok {
			members = bm[key]
		}
		// Get each member (user/principal) for the binding
		for _, m := range b.Members {
			members[m] = struct{}{}
		}
		if len(members) > 0 {
			bm[key] = members
		} else {
			delete(bm, key)
		}
	}
	return bm
}

// Return list of Bindings for a map of role to member sets
func listFromIamBindingMap(bm map[iamBindingKey]map[string]struct{}) []*cloudresourcemanager.Binding {
	rb := make([]*cloudresourcemanager.Binding, 0, len(bm))
	var keys []iamBindingKey
	for k := range bm {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		keyI := keys[i]
		keyJ := keys[j]
		return fmt.Sprintf("%s%s", keyI.Role, keyI.Condition.String()) < fmt.Sprintf("%s%s", keyJ.Role, keyJ.Condition.String())
	})
	for _, key := range keys {
		members := bm[key]
		if len(members) == 0 {
			continue
		}
		b := &cloudresourcemanager.Binding{
			Role:    key.Role,
			Members: stringSliceFromGolangSet(members),
		}
		if !key.Condition.Empty() {
			b.Condition = &cloudresourcemanager.Expr{
				Description: key.Condition.Description,
				Expression:  key.Condition.Expression,
				Title:       key.Condition.Title,
			}
		}
		rb = append(rb, b)
	}
	return rb
}

func conditionKeyFromCondition(condition *cloudresourcemanager.Expr) conditionKey {
	if condition == nil {
		return conditionKey{}
	}
	return conditionKey{condition.Description, condition.Expression, condition.Title}
}

func stringSliceFromGolangSet(sset map[string]struct{}) []string {
	ls := make([]string, 0, len(sset))
	for s := range sset {
		ls = append(ls, s)
	}
	sort.Strings(ls)

	return ls
}

func IsConflictError(err error) bool {
	if e, ok := err.(*googleapi.Error); ok && (e.Code == 409 || e.Code == 412) {
		return true
	} else if !ok && errwrap.ContainsType(err, &googleapi.Error{}) {
		e := errwrap.GetType(err, &googleapi.Error{}).(*googleapi.Error)
		if e.Code == 409 || e.Code == 412 {
			return true
		}
	}
	return false
}
