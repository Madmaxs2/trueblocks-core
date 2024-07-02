/* eslint object-curly-newline: ["error", "never"] */
/* eslint max-len: ["error", 160] */
/*
 * Copyright 2016, 2024 The TrueBlocks Authors. All rights reserved.
 * Use of this source code is governed by a license that can
 * be found in the LICENSE file.
 *
 * This file was auto generated. DO NOT EDIT.
 */

import * as ApiCallers from '../lib/api_callers';
import { address, float64, Message, Monitor, MonitorClean, uint64 } from '../types';

export function getMonitors(
  parameters?: {
    addrs?: address[],
    delete?: boolean,
    undelete?: boolean,
    remove?: boolean,
    clean?: boolean,
    list?: boolean,
    watch?: boolean,
    watchlist?: string,
    commands?: string,
    batchSize?: uint64,
    runCount?: uint64,
    sleep?: float64,
    fmt?: string,
    chain: string,
    noHeader?: boolean,
    cache?: boolean,
    decache?: boolean,
  },
  options?: RequestInit,
) {
  return ApiCallers.fetch<Message[] | MonitorClean[] | Monitor[]>(
    { endpoint: '/monitors', method: 'get', parameters, options },
  );
}
