package monitor

// Copyright 2021 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/config"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/index"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/validate"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// Header is the header of the Monitor file. Note that it's the same width as an index.AppearanceRecord
// therefor one should not change its size
type Header struct {
	Magic       uint16 `json:"-"`
	Unused      bool   `json:"-"`
	Deleted     bool   `json:"deleted,omitempty"`
	LastScanned uint32 `json:"lastScanned,omitempty"`
}

// Monitor carries information about a Monitor file and its header
type Monitor struct {
	Address  common.Address `json:"address"`
	Staged   bool           `json:"-"`
	Chain    string         `json:"-"`
	ReadFp   *os.File       `json:"-"`
	TestMode bool           `json:"-"`
	Header
}

const (
	Ext = ".mon.bin"
)

// NewMonitor returns a Monitor (but has not yet read in the AppearanceRecords). If 'create' is
// sent, create the Monitor if it does not already exist
func NewMonitor(chain, addr string, create, testMode bool) Monitor {
	mon := new(Monitor)
	mon.Header = Header{Magic: file.SmallMagicNumber}
	mon.Address = common.HexToAddress(addr)
	mon.Chain = chain
	mon.TestMode = testMode
	mon.Reload(create)
	return *mon
}

// NewStagedMonitor returns a Monitor whose path is in the 'staging' folder
func NewStagedMonitor(chain, addr string, testMode bool) (Monitor, error) {
	mon := Monitor{
		Header:   Header{Magic: file.SmallMagicNumber},
		Address:  common.HexToAddress(addr),
		Chain:    chain,
		TestMode: testMode,
	}

	// Note, we are not yet staged, so Path returns the production folder
	prodPath := mon.Path()
	mon.Staged = true
	// stagedPath := mon.Path()

	// either copy the existing monitor or create a new one
	if file.FileExists(prodPath) {
	} else {
		err := mon.WriteMonHeader(false, 0)
		if err != nil {
			return mon, err
		}
	}
	return mon, nil
}

// String implements the Stringer interface
func (mon Monitor) String() string {
	if mon.Deleted {
		return fmt.Sprintf("%s\t%d\t%d\t%d\t%t", hexutil.Encode(mon.Address.Bytes()), mon.Count(), file.FileSize(mon.Path()), mon.LastScanned, mon.Deleted)
	}
	return fmt.Sprintf("%s\t%d\t%d\t%d", hexutil.Encode(mon.Address.Bytes()), mon.Count(), file.FileSize(mon.Path()), mon.LastScanned)
}

type SimpleMonitor struct {
	Address     string `json:"address"`
	NRecords    int    `json:"nRecords"`
	FileSize    int64  `json:"fileSize"`
	LastScanned uint32 `json:"lastScanned"`
}

func NewSimpleMonitor(mon Monitor) SimpleMonitor {
	return SimpleMonitor{
		Address:     mon.GetAddrStr(),
		NRecords:    int(mon.Count()),
		FileSize:    file.FileSize(mon.Path()),
		LastScanned: mon.Header.LastScanned,
	}
}

// ToJSON returns a JSON object from a Monitor
func (mon Monitor) ToJSON() string {
	sm := NewSimpleMonitor(mon)
	bytes, err := json.Marshal(sm)
	if err != nil {
		return ""
	}
	return string(bytes)
}

// Path returns the path to the Monitor file
func (mon *Monitor) Path() (path string) {
	if mon.Staged {
		path = config.GetPathToCache(mon.Chain) + "monitors/staging/" + strings.ToLower(mon.Address.Hex()) + Ext
	} else {
		path = config.GetPathToCache(mon.Chain) + "monitors/" + strings.ToLower(mon.Address.Hex()) + Ext
	}
	if mon.TestMode {
		path = strings.Replace(path, config.GetPathToCache(mon.Chain), config.GetPathToChainConfig(mon.Chain)+"mocked/", -1)
	}
	return
}

// Reload loads information about the monitor such as the file's size and record count
func (mon *Monitor) Reload(create bool) (uint32, error) {
	if create && !file.FileExists(mon.Path()) {
		// Make sure the file exists since we've been told to monitor it
		err := mon.WriteMonHeader(false, 0)
		if err != nil {
			return 0, err
		}
	}
	return mon.Count(), nil
}

func (mon *Monitor) Count() uint32 {
	if file.FileSize(mon.Path()) == 0 {
		return 0
	}
	s := uint32(file.FileSize(mon.Path()))
	w := uint32(index.AppRecordWidth)
	n := uint32(s / w)
	return n - 1
}

// GetAddrStr returns the Monitor's address as a string
func (mon *Monitor) GetAddrStr() string {
	return strings.ToLower(mon.Address.Hex())
}

// Close closes an open Monitor if it's open, does nothing otherwise
func (mon *Monitor) Close() {
	if mon.ReadFp != nil {
		mon.ReadFp.Close()
		mon.ReadFp = nil
	}
}

// ReadMonHeader reads the monitor's header and returns without closing the file
func (mon *Monitor) ReadMonHeader() (err error) {
	if mon.ReadFp == nil {
		mon.ReadFp, err = os.OpenFile(mon.Path(), os.O_RDONLY, 0644)
		if err != nil {
			return
		}
	}
	if file.FileSize(mon.Path()) > 0 {
		return binary.Read(mon.ReadFp, binary.LittleEndian, &mon.Header)
	}
	return
}

// ReadAppearanceAt returns the appearance at the one-based index.
func (mon *Monitor) ReadAppearanceAt(idx uint32, app *index.AppearanceRecord) (err error) {
	if idx == 0 || idx > mon.Count() {
		// the file contains a header on record wide, so a one-based index eases caller code
		err = errors.New(fmt.Sprintf("index out of range in ReadAppearanceAt[%d]", idx))
		return
	}

	if mon.ReadFp == nil {
		path := mon.Path()
		mon.ReadFp, err = os.OpenFile(path, os.O_RDONLY, 0644)
		if err != nil {
			return
		}
	}

	// This index is one based because we have to skip over the header
	byteIndex := int64(idx) * index.AppRecordWidth
	_, err = mon.ReadFp.Seek(byteIndex, io.SeekStart)
	if err != nil {
		return
	}

	err = binary.Read(mon.ReadFp, binary.LittleEndian, &app.BlockNumber)
	if err != nil {
		return
	}
	err = binary.Read(mon.ReadFp, binary.LittleEndian, &app.TransactionId)
	return
}

// ReadAppearances returns appearances starting at the first appearance in the file.
func (mon *Monitor) ReadAppearances(apps *[]index.AppearanceRecord) (err error) {
	if mon.ReadFp == nil {
		path := mon.Path()
		mon.ReadFp, err = os.OpenFile(path, os.O_RDONLY, 0644)
		if err != nil {
			return
		}
	}

	// Seek past the header to get to the first record
	_, err = mon.ReadFp.Seek(index.AppRecordWidth, io.SeekStart)
	if err != nil {
		return
	}

	err = binary.Read(mon.ReadFp, binary.LittleEndian, apps)
	if err != nil {
		return
	}
	return
}

// WriteMonHeader reads the monitor's header
func (mon *Monitor) WriteMonHeader(deleted bool, lastScanned uint32) (err error) {
	f, err := os.OpenFile(mon.Path(), os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	defer f.Close()

	mon.Deleted = deleted
	if lastScanned > mon.LastScanned {
		mon.LastScanned = lastScanned
		// TODO: BOGUS1
		// fmt.Fprintf(os.Stdout, "%sWritzy--> %s%s\n", colors.Green, mon, colors.Off)
	}

	f.Seek(0, io.SeekStart)
	err = binary.Write(f, binary.LittleEndian, mon.Header)
	return
}

// WriteAppendApps appends appearances to the end of the file, updates the header with
// lastScanned (if later) and returns the number of records written. Note that we should
// be writing to a temporary file.
func (mon *Monitor) WriteAppendApps(lastScanned uint32, apps *[]index.AppearanceRecord) error {
	if !mon.Staged {
		log.Fatal("Trying to write to a non-staged file. Should not happen.")

	} else if mon == nil {
		log.Fatal("Trying to write from a nil monitor. Should not happen.")
	}

	err := mon.WriteMonHeader(mon.Deleted, lastScanned)
	if err != nil {
		return err
	}

	if apps != nil {
		if len(*apps) > 0 {
			_, err := mon.WriteAppearances(*apps)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// WriteAppearances writes appearances to a Monitor
func (mon *Monitor) WriteAppearances(apps []index.AppearanceRecord) (count int, err error) {
	mode := os.O_WRONLY | os.O_CREATE
	if file.FileExists(mon.Path()) {
		// log.Println("Appending to existing monitor", mon.GetAddrStr())
		mode = os.O_WRONLY | os.O_APPEND
	}

	path := mon.Path()
	f, err := os.OpenFile(path, mode, 0644)
	if err != nil {
		return
	}

	f.Seek(index.AppRecordWidth, io.SeekStart)

	b := make([]byte, 4, 4)
	for _, app := range apps {
		binary.LittleEndian.PutUint32(b, app.BlockNumber)
		_, err = f.Write(b)
		if err != nil {
			f.Close()
			return
		}
		binary.LittleEndian.PutUint32(b, app.TransactionId)
		_, err = f.Write(b)
		if err != nil {
			f.Close()
			return
		}
	}

	f.Close() // do not defer this, we need to close it so the fileSize is right
	mon.Reload(false /* create */)
	count = int(mon.Count())

	return
}

// IsDeleted returns true if the monitor has been deleted but not removed
func (mon *Monitor) IsDeleted() bool {
	mon.ReadMonHeader()
	return mon.Header.Deleted
}

// Delete marks the file's delete flag, but does not physically remove the file
func (mon *Monitor) Delete() (prev bool) {
	prev = mon.Deleted
	mon.WriteMonHeader(true, mon.LastScanned)
	mon.Deleted = true
	return
}

// UnDelete unmarks the file's delete flag
func (mon *Monitor) UnDelete() (prev bool) {
	prev = mon.Deleted
	mon.WriteMonHeader(false, mon.LastScanned)
	mon.Deleted = false
	return
}

// Remove removes a previously deleted file, does nothing if the file is not deleted
func (mon *Monitor) Remove() (bool, error) {
	if !mon.IsDeleted() {
		return false, errors.New("cannot remove a monitor that is not deleted")
	}
	file.Remove(mon.Path())
	return !file.FileExists(mon.Path()), nil
}

func addressFromPath(path string) (string, error) {
	_, fileName := filepath.Split(path)
	if len(fileName) == 0 || !strings.HasPrefix(fileName, "0x") || !strings.HasSuffix(fileName, ".mon.bin") {
		return "", errors.New("path does is not a valid monitor filename")
	}
	parts := strings.Split(fileName, ".")
	return strings.ToLower(parts[0]), nil
}

// SentinalAddr is a marker to signify the end of the monitor list produced by ListMonitors
var SentinalAddr = common.HexToAddress("0xdeaddeaddeaddeaddeaddeaddeaddeaddeaddead")

// ListMonitors puts a list of Monitors into the monitorChannel. The list of monitors is built from
// a file called addresses.csv in the current folder or, if not present, from existing monitors
func ListMonitors(chain, folder string, monitorChan chan<- Monitor) {
	defer func() {
		monitorChan <- Monitor{Address: SentinalAddr}
	}()

	info, err := os.Stat("./addresses.csv")
	if err == nil {
		// If the shorthand file exists in the current folder, use it...
		lines := file.AsciiFileToLines(info.Name())
		fmt.Println("Found ", len(lines), " addresses to monitor in ./addresses.csv")
		for _, line := range lines {
			if !strings.HasPrefix(line, "#") {
				parts := strings.Split(line, ",")
				if len(parts) > 0 && validate.IsValidAddress(parts[0]) && !validate.IsZeroAddress(parts[0]) {
					monitorChan <- NewMonitor(chain, parts[0], true /* create */, false)
				}
			}
		}
		return
	}

	// ...otherwise freshen all existing monitors
	pp := config.GetPathToCache(chain) + folder
	filepath.Walk(pp, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			addr, _ := addressFromPath(path)
			if len(addr) > 0 {
				monitorChan <- NewMonitor(chain, addr, true /* create */, false)
			}
		}
		return nil
	})
}

// TODO: This should be non-interuptable
// MoveToProduction moves a previously staged monitor to the monitors folder.
func (mon *Monitor) MoveToProduction() error {
	if !mon.Staged {
		return errors.New("trying to move monitor that is not staged")
	}

	oldPath := mon.Path()
	mon.Staged = false
	return os.Rename(oldPath, mon.Path())
}
