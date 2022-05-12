// Copyright 2021 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.

package validate

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/blockRange"
	tslibPkg "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/tslib"
	"github.com/araddon/dateparse"
)

func IsBlockHash(str string) bool {
	if !Is0xPrefixed(str) {
		return false
	}

	if len(str) != 66 {
		return false
	}

	if !IsHex(str) {
		return false
	}

	return true
}

type blockNum = uint32

func IsBlockNumber(str string) (bool, blockNum) {
	base := 10
	source := str

	if Is0xPrefixed(str) {
		base = 16
		source = str[2:]
	}

	value, err := strconv.ParseUint(source, base, 32)
	if err != nil {
		return false, 0
	}

	return true, blockNum(value)
}

func IsBlockNumberList(strs []string) (bool, []blockNum) {
	result := make([]blockNum, len(strs))

	for index, stringValue := range strs {
		check, value := IsBlockNumber(stringValue)
		if !check {
			return false, nil
		}

		result[index] = value
	}

	return true, result
}

func IsDateTimeString(str string) bool {
	bRange, err := blockRange.New(str)
	if err != nil {
		return false
	}
	return bRange.StartType == blockRange.BlockRangeDate
}

func IsBeforeFirstBlock(chain, str string) bool {
	if !IsDateTimeString(str) {
		return false
	}

	// From https://github.com/araddon/dateparse
	time.Local, _ = time.LoadLocation("UTC")
	dt, _ := dateparse.ParseLocal(str) // already validated as a date
	return dt.Before(tslibPkg.DateFromName(chain, "first"))
}

func IsRange(chain, str string) (bool, error) {
	// Disallow "start only ranges" like "1000" or "london"
	if !strings.Contains(str, "-") {
		return false, &InvalidIdentifierLiteralError{
			Value: str,
		}
	}

	bRange, err := blockRange.New(str)

	if err == nil {
		if bRange.Start.Special == "latest" {
			return false, errors.New("cannot start range with 'latest'")
		}

		if bRange.StartType == blockRange.BlockRangeSpecial &&
			!tslibPkg.IsSpecialBlock(chain, bRange.Start.Special) {
			return false, &InvalidIdentifierLiteralError{
				Value: bRange.Start.Special,
			}
		}

		if bRange.EndType == blockRange.BlockRangeSpecial &&
			!tslibPkg.IsSpecialBlock(chain, bRange.End.Special) {
			return false, &InvalidIdentifierLiteralError{
				Value: bRange.End.Special,
			}
		}

		onlyNumbers := bRange.StartType == blockRange.BlockRangeBlockNumber &&
			bRange.EndType == blockRange.BlockRangeBlockNumber

		if onlyNumbers && bRange.Start.Block >= bRange.End.Block {
			return false, errors.New("'stop' must be strictly larger than 'start'")
		}

		return true, nil
	}

	if modifierErr, ok := err.(*blockRange.WrongModifierError); ok {
		return false, errors.New("Input argument appears to be invalid. No such skip marker: " + modifierErr.Token)
	}

	return false, err
}

// Errors returned by ValidateIdentifiers (note: it can also return an
// error passed from IsRange)
var ErrTooManyRanges = errors.New("too many ranges")

type InvalidIdentifierLiteralError struct {
	Value string
	Msg   string
}

func (e *InvalidIdentifierLiteralError) Error() string {
	if len(e.Msg) == 0 {
		e.Msg = "is not a numeral or a special named block."
	}
	return fmt.Sprintf("The given value '%s' %s", e.Value, e.Msg)
}

func IsValidBlockId(chain string, ids []string, validTypes ValidArgumentType) (bool, error) {
	err := ValidateIdentifiers(chain, ids, validTypes, 1)
	return err == nil, err
}
