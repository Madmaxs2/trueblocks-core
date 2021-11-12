package export

/*-------------------------------------------------------------------------------------------
 * qblocks - fast, easily-accessible, fully-decentralized data from blockchains
 * copyright (c) 2016, 2021 TrueBlocks, LLC (http://trueblocks.io)
 *
 * This program is free software: you may redistribute it and/or modify it under the terms
 * of the GNU General Public License as published by the Free Software Foundation, either
 * version 3 of the License, or (at your option) any later version. This program is
 * distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even
 * the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU
 * General Public License for more details. You should have received a copy of the GNU General
 * Public License along with this program. If not, see http://www.gnu.org/licenses/.
 *-------------------------------------------------------------------------------------------*/
/*
 * The file was auto generated with makeClass --gocmds. DO NOT EDIT.
 */

import (
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
)

type ExportOptionsType struct {
	Appearances bool
	Receipts    bool
	Statements  bool
	Logs        bool
	Traces      bool
	Accounting  bool
	Articulate  bool
	Cache       bool
	CacheTraces bool
	Factory     bool
	Count       bool
	FirstRecord uint64
	MaxRecords  uint64
	Relevant    bool
	Emitter     []string
	Topic       []string
	Clean       bool
	Freshen     bool
	Staging     bool
	Unripe      bool
	Load        string
	Reversed    bool
	ByDate      bool
	SummarizeBy string
	SkipDdos    bool
	MaxTraces   uint64
	FirstBlock  uint64
	LastBlock   uint64
}

var Options ExportOptionsType

func (opts *ExportOptionsType) TestLog() {
	logger.Log(logger.Test, "Appearances: ", opts.Appearances)
	logger.Log(logger.Test, "Receipts: ", opts.Receipts)
	logger.Log(logger.Test, "Statements: ", opts.Statements)
	logger.Log(logger.Test, "Logs: ", opts.Logs)
	logger.Log(logger.Test, "Traces: ", opts.Traces)
	logger.Log(logger.Test, "Accounting: ", opts.Accounting)
	logger.Log(logger.Test, "Articulate: ", opts.Articulate)
	logger.Log(logger.Test, "Cache: ", opts.Cache)
	logger.Log(logger.Test, "CacheTraces: ", opts.CacheTraces)
	logger.Log(logger.Test, "Factory: ", opts.Factory)
	logger.Log(logger.Test, "Count: ", opts.Count)
	logger.Log(logger.Test, "FirstRecord: ", opts.FirstRecord)
	logger.Log(logger.Test, "MaxRecords: ", opts.MaxRecords)
	logger.Log(logger.Test, "Relevant: ", opts.Relevant)
	logger.Log(logger.Test, "Emitter: ", opts.Emitter)
	logger.Log(logger.Test, "Topic: ", opts.Topic)
	logger.Log(logger.Test, "Clean: ", opts.Clean)
	logger.Log(logger.Test, "Freshen: ", opts.Freshen)
	logger.Log(logger.Test, "Staging: ", opts.Staging)
	logger.Log(logger.Test, "Unripe: ", opts.Unripe)
	logger.Log(logger.Test, "Load: ", opts.Load)
	logger.Log(logger.Test, "Reversed: ", opts.Reversed)
	logger.Log(logger.Test, "ByDate: ", opts.ByDate)
	logger.Log(logger.Test, "SummarizeBy: ", opts.SummarizeBy)
	logger.Log(logger.Test, "SkipDdos: ", opts.SkipDdos)
	logger.Log(logger.Test, "MaxTraces: ", opts.MaxTraces)
	logger.Log(logger.Test, "FirstBlock: ", opts.FirstBlock)
	logger.Log(logger.Test, "LastBlock: ", opts.LastBlock)
}