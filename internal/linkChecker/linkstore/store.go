package linkstore

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/Negat1v9/link-checker/config"
	"github.com/Negat1v9/link-checker/pkg/logger"
)

type LinkStore struct {
	cfg   *config.LinkStoreCfg
	log   *logger.Logger
	group map[int][]string
	mu    *sync.RWMutex

	backupPeriodSecond time.Duration
	shutDownCh         chan bool
	walWriter          *walWriter
}

func NewLinkStore(cfg *config.LinkStoreCfg, log *logger.Logger) (*LinkStore, error) {
	walWr, err := NewWalWriter(cfg, log)
	if err != nil {
		return nil, err
	}
	l := &LinkStore{
		cfg:   cfg,
		log:   log,
		group: make(map[int][]string),
		mu:    &sync.RWMutex{},

		backupPeriodSecond: time.Duration(cfg.LinkStorageBackupDuration) * time.Second,
		shutDownCh:         make(chan bool),
		walWriter:          walWr,
	}

	go l.run()

	return l, nil
}

// stop store and backup data in disk
func (s *LinkStore) Stop(shutDown context.Context) error {
	s.shutDownCh <- true
	for {
		select {
		case <-shutDown.Done():
			return fmt.Errorf("linkstore.LinkStore.Stop: %v", shutDown.Err())
		default:
			// stop backup store
			if err := s.walWriter.Stop(shutDown); err != nil {
				return err
			}
			s.log.Debugf("walWriter is stopped")
			if err := s.writeBackup(); err != nil {
				return err
			}
			s.log.Debugf("LinkStore is stopped")
			return nil
		}
	}
}

// save link group return links_num
func (s *LinkStore) CreateLinksGroup(links []string) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	groupID := len(s.group) + 1
	s.group[groupID] = links

	// backup
	s.walWriter.Append(groupID, links)

	return groupID
}

// func (s *LinkStore) AddLinkInGroup(groupID int, link string) {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()

// }

// return links in by groupID
func (s *LinkStore) GetLinksByGroup(groupID int) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.group[groupID]
}

// open backup file if exits and load data
func (s *LinkStore) loadStoreBackup() error {
	b, err := os.ReadFile(s.cfg.LinkStoreBackup)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, &s.group)
}

func (s *LinkStore) run() {
	if err := s.loadStoreBackup(); err != nil {
		s.log.Warnf("LinkStore.loadStoreBackup: %v", err)
	}

	if err := s.checkWALFile(); err != nil {
		s.log.Warnf("LinkStore.checkWALFile: %v", err)
	}

	ticker := time.NewTicker(s.backupPeriodSecond)
	for {
		select {
		case <-s.shutDownCh:
			return
		case <-ticker.C:
			if len(s.group) != 0 {
				s.log.Debugf("run LinkStore.writeBackup")
				if err := s.writeBackup(); err != nil {
					s.log.Errorf("LinkStore.writeBackup: %v", err)
				}
			}
			ticker.Reset(s.backupPeriodSecond)
		}
	}
}

func (s *LinkStore) writeBackup() error {
	f, err := os.OpenFile(s.cfg.LinkStoreBackup, os.O_CREATE|os.O_WRONLY|os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	b, err := json.Marshal(&s.group)
	if err != nil {
		return err
	}
	if _, err = f.Write(b); err != nil {

		return err
	}

	if err = f.Sync(); err != nil {
		return err
	}
	return f.Close()
}

// parse wal file and load data call only on start store
func (s *LinkStore) checkWALFile() error {
	// load data from WAL file
	walOperations, err := s.walWriter.GetWalOps()
	if err != nil {
		return err
	}

	// data integrity check
	isUpdated := false
	for _, op := range walOperations {
		// in wal file group exist and in store not -> add it
		if _, ok := s.group[op.LinkGroupID]; !ok {
			if op.Type == "w" {
				// add in
				isUpdated = true
				s.group[op.LinkGroupID] = op.Links
			}
		}
	}
	// save with updated data from WAL if was udpated
	if isUpdated {
		s.log.Infof("unsaved data found in log files")
		if err = s.writeBackup(); err != nil {
			return err
		}
	}

	// clear WAL file
	return s.walWriter.FlushWAL()
}
