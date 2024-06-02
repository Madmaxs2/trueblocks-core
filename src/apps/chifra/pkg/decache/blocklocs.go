package decache

import (
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/cache"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/identifiers"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/rpc"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
)

func LocationsFromBlockIds(conn *rpc.Connection, ids []identifiers.Identifier, logs, trace bool) ([]cache.Locator, error) {
	locations := make([]cache.Locator, 0)
	for _, br := range ids {
		blockNums, err := br.ResolveBlocks(conn.Chain)
		if err != nil {
			return nil, err
		}
		for _, bn := range blockNums {
			if logs {
				logGroup := &types.LogGroup{
					BlockNumber:      bn,
					TransactionIndex: base.NOPOSN,
				}
				locations = append(locations, logGroup)

			} else if trace {
				traceGroup := &types.TraceGroup{
					BlockNumber:      bn,
					TransactionIndex: base.NOPOSN,
				}
				locations = append(locations, traceGroup)

			} else {
				block, err := conn.GetBlockHeaderByNumber(bn)
				if err != nil {
					return nil, err
				}
				locations = append(locations, &types.LightBlock{
					BlockNumber: bn,
				})
				receiptGroup := &types.ReceiptGroup{
					BlockNumber:      bn,
					TransactionIndex: base.NOPOSN,
				}
				locations = append(locations, receiptGroup)
				for index := range block.Transactions {
					locations = append(locations, &types.Transaction{
						BlockNumber:      bn,
						TransactionIndex: base.Txnum(index),
					})
					locations = append(locations, &types.TraceGroup{
						BlockNumber:      bn,
						TransactionIndex: base.Txnum(index),
					})
				}
			}
		}
	}
	return locations, nil
}
