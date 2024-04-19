// Copyright 2016, 2024 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
/*
 * Parts of this file were auto generated. Edit only those parts of
 * the code inside of 'EXISTING_CODE' tags.
 */

package types

// EXISTING_CODE
import (
	"encoding/json"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
)

// EXISTING_CODE

type RawAbi struct {
	Address   string   `json:"address"`
	Functions []string `json:"functions"`
	// EXISTING_CODE
	// EXISTING_CODE
}

type SimpleAbi struct {
	Address   base.Address     `json:"address"`
	Functions []SimpleFunction `json:"functions"`
	raw       *RawAbi          `json:"-"`
	// EXISTING_CODE
	// EXISTING_CODE
}

func (s *SimpleAbi) String() string {
	bytes, _ := json.Marshal(s)
	return string(bytes)
}

func (s *SimpleAbi) Raw() *RawAbi {
	return s.raw
}

func (s *SimpleAbi) SetRaw(raw *RawAbi) {
	s.raw = raw
}

func (s *SimpleAbi) Model(chain, format string, verbose bool, extraOptions map[string]any) Model {
	var model = map[string]interface{}{}
	var order = []string{}

	// EXISTING_CODE
	model[s.Address.Hex()] = s.Functions
	order = append(order, s.Address.Hex())
	// EXISTING_CODE

	return Model{
		Data:  model,
		Order: order,
	}
}

// FinishUnmarshal is used by the cache. It may be unused depending on auto-code-gen
func (s *SimpleAbi) FinishUnmarshal() {
	// EXISTING_CODE
	// EXISTING_CODE
}

// EXISTING_CODE
// EXISTING_CODE