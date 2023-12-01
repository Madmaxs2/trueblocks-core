// Copyright 2021 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.

package transactionsPkg

import (
	"context"
	"fmt"
	"sort"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/articulate"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/identifiers"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/output"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/utils"
)

func (opts *TransactionsOptions) HandleLogs() error {
	chain := opts.Globals.Chain
	testMode := opts.Globals.TestMode
	nErrors := 0

	abiCache := articulate.NewAbiCache(opts.Conn, opts.Articulate)
	emitters := []base.Address{}
	for _, e := range opts.Emitter {
		emitters = append(emitters, base.HexToAddress(e))
	}
	topics := []base.Hash{}
	for _, t := range opts.Topic {
		topics = append(topics, base.HexToHash(t))
	}
	logFilter := types.SimpleLogFilter{
		Emitters: emitters,
		Topics:   topics,
	}

	ctx, cancel := context.WithCancel(context.Background())
	fetchData := func(modelChan chan types.Modeler[types.RawLog], errorChan chan error) {
		if sliceOfMaps, cnt, err := identifiers.AsSliceOfMaps[types.SimpleTransaction](chain, opts.TransactionIds); err != nil {
			errorChan <- err
			cancel()

		} else if cnt == 0 {
			errorChan <- fmt.Errorf("no transactions found")
			cancel()

		} else {
			bar := logger.NewBar(logger.BarOptions{
				Enabled: !testMode && !utils.IsTerminal(),
				Total:   int64(cnt),
			})

			for _, thisMap := range sliceOfMaps {
				thisMap := thisMap
				for app := range thisMap {
					thisMap[app] = new(types.SimpleTransaction)
				}

				iterFunc := func(app types.SimpleAppearance, value *types.SimpleTransaction) error {
					if tx, err := opts.Conn.GetTransactionByAppearance(&app, opts.Traces /* needsTraces */); err != nil {
						return fmt.Errorf("transaction at %s returned an error: %w", app.Orig(), err)

					} else if tx == nil {
						return fmt.Errorf("transaction at %s has no logs", app.Orig())

					} else {
						if opts.Articulate && tx.ArticulatedTx == nil {
							if err = abiCache.ArticulateTransaction(tx); err != nil {
								errorChan <- err // continue even with an error
							}
						}
						*value = *tx
						bar.Tick()
						return nil
					}
				}

				iterErrorChan := make(chan error)
				iterCtx, iterCancel := context.WithCancel(context.Background())
				defer iterCancel()
				go utils.IterateOverMap(iterCtx, iterErrorChan, thisMap, iterFunc)
				for err := range iterErrorChan {
					if !testMode || nErrors == 0 {
						errorChan <- err
						nErrors++
					}
				}

				items := make([]types.SimpleTransaction, 0, len(thisMap))
				for _, tx := range thisMap {
					items = append(items, *tx)
				}
				sort.Slice(items, func(i, j int) bool {
					if items[i].BlockNumber == items[j].BlockNumber {
						return items[i].TransactionIndex < items[j].TransactionIndex
					}
					return items[i].BlockNumber < items[j].BlockNumber
				})

				for _, item := range items {
					item := item
					for _, log := range item.Receipt.Logs {
						log := log
						if logFilter.PassesFilter(&log) {
							modelChan <- &log
						}
					}
				}
			}
			bar.Finish(true)
		}
	}

	extra := map[string]interface{}{
		"logs":      opts.Logs,
		"traces":    opts.Traces,
		"addresses": opts.Uniq,
	}

	return output.StreamMany(ctx, fetchData, opts.Globals.OutputOptsWithExtra(extra))
}
