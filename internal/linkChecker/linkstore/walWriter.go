package linkstore

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/Negat1v9/link-checker/config"
	"github.com/Negat1v9/link-checker/pkg/logger"
)

type WalOp struct {
	Type        string   `json:"type"` // w - write, d - delete
	LinkGroupID int      `json:"link_group_id"`
	Links       []string `json:"links"`
}

type walWriter struct {
	log                 *logger.Logger
	path                string        // path wal file
	syncN               int           // length on buffer to sync
	syncDuractionSecond time.Duration // time to sync

	f       *os.File // wal file
	buffer  []WalOp
	walOpCh chan WalOp
}

func NewWalWriter(cfg *config.LinkStoreCfg, log *logger.Logger) (*walWriter, error) {
	f, err := os.OpenFile(cfg.WalFilePath, os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}

	w := &walWriter{
		log:                 log,
		path:                cfg.WalFilePath,
		syncN:               cfg.WalFileBufferSize,
		syncDuractionSecond: time.Duration(cfg.WalFileWriteDuration) * time.Second,
		f:                   f,
		buffer:              make([]WalOp, 0, cfg.WalFileBufferSize),
		walOpCh:             make(chan WalOp),
	}

	go w.run()
	return w, nil
}

func (w *walWriter) Append(groupID int, links []string) {
	w.walOpCh <- WalOp{Type: "w", LinkGroupID: groupID, Links: links}
}

func (w *walWriter) GetWalOps() ([]WalOp, error) {
	reader := bufio.NewReader(w.f)
	res := make([]WalOp, 0)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil && err != io.EOF {
			return nil, err
		} else if err == io.EOF {
			break
		}

		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		var op WalOp
		if err = json.Unmarshal(line, &op); err != nil {
			// file is broken
			break
		}

		res = append(res, op)
	}

	return res, nil
}

// clears the WAL data was overwritten in the main storage file
func (w *walWriter) FlushWAL() error {
	f, err := os.OpenFile(w.path, os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	f.Close()
	return err
}

func (w *walWriter) Stop(shutDown context.Context) error {
	close(w.walOpCh)
	for {
		select {
		case <-shutDown.Done():
			return fmt.Errorf("linkstore.walWriter.Stop: %v", shutDown.Err())
		default:
			// write latest updates
			err := w.syncWalFile()

			return err
		}

	}
}

// run walWriter
func (w *walWriter) run() {

	ticker := time.NewTicker(w.syncDuractionSecond)
	for {
		select {
		case op, ok := <-w.walOpCh:
			if !ok {
				return
			}
			w.buffer = append(w.buffer, op)
			if len(w.buffer) >= w.syncN {
				if err := w.syncWalFile(); err != nil {
					w.log.Errorf("walWriter.syncWalFile: %v", err)
				}
				ticker.Reset(w.syncDuractionSecond)
			}
		case <-ticker.C:
			if err := w.syncWalFile(); err != nil {
				w.log.Errorf("walWriter.syncWalFile: %v", err)
			}
			ticker.Reset(w.syncDuractionSecond)
		}
	}
}

func (w *walWriter) syncWalFile() error {
	if len(w.buffer) == 0 {
		return nil
	}

	for _, op := range w.buffer {
		b, _ := json.Marshal(&op)

		if _, err := w.f.Write(append(b, '\n')); err != nil {
			return err
		}
	}
	if err := w.f.Sync(); err != nil {
		return err
	}
	w.buffer = w.buffer[:0]

	return nil
}
