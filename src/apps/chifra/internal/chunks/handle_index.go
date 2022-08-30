// Copyright 2021 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.

package chunksPkg

import (
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/cache"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/config"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/index"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
)

func (opts *ChunksOptions) showIndex(ctx *WalkContext, path string, first bool) (bool, error) {
	path = config.ToIndexPath(path)
	header, err := index.ReadChunkHeader(opts.Globals.Chain, path, true)
	if err != nil {
		return false, err
	}

	rng, err := cache.RangeFromFilename(path)
	if err != nil {
		return false, err
	}

	obj := types.SimpleIndex{
		Range:           rng,
		Magic:           header.Magic,
		Hash:            header.Hash,
		AddressCount:    header.AddressCount,
		AppearanceCount: header.AppearanceCount,
		Size:            file.FileSize(path),
	}

	// TODO: Feature - customize display strings
	// opts.Globals.Format = "Magic,Hash,Size,AppearanceCount,AddressCount,Range"
	// opts.Globals.Format = "Range\tAppearanceCount\tAddressCount"
	// TODO: Fix export without arrays
	err = opts.Globals.RenderObject(obj, first)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (opts *ChunksOptions) HandleIndex(blockNums []uint64) error {
	defer opts.Globals.RenderFooter()
	err := opts.Globals.RenderHeader(types.SimpleIndex{}, &opts.Globals.Writer, opts.Globals.Format, opts.Globals.ApiMode, opts.Globals.NoHeader, true)
	if err != nil {
		return err
	}

	ctx := WalkContext{
		VisitFunc: opts.showIndex,
	}

	return opts.WalkIndexFiles(&ctx, cache.Index_Bloom, blockNums)
}
