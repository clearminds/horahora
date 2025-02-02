package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/go-redsync/redsync"
	"github.com/jmoiron/sqlx"
)

type VideoDLRequest struct {
	Redsync     *redsync.Redsync
	Db          *sqlx.DB
	VideoID     string // Foreign ID
	ID          int    // Domestic ID
	DownloaddID int
	URL         string
	ParentURL   string
	mut         *redsync.Mutex
}

func (v *VideoDLRequest) SetDownloadSucceeded() error {
	sql := "UPDATE videos SET dlStatus = 1 WHERE id = $1"
	_, err := v.Db.Exec(sql, v.ID)
	return err
}

func (v *VideoDLRequest) SetDownloadFailed() error {
	sql := "UPDATE videos SET dlStatus = 2 WHERE id = $1"
	_, err := v.Db.Exec(sql, v.ID)
	return err
}

// have to pass a transaction for this one because it needs to be atomic with the scheduler query
func (v *VideoDLRequest) SetDownloadInProgress(tx *sql.Tx) error {
	sql := "UPDATE videos SET dlStatus = 3 WHERE id = $1"
	_, err := tx.Exec(sql, v.ID)
	return err
}

func (v *VideoDLRequest) AcquireLockForVideo() error {
	v.mut = v.Redsync.NewMutex(v.VideoID, redsync.SetExpiry(time.Minute*10))
	return v.mut.Lock()
}

func (v *VideoDLRequest) ReleaseLockForVideo() error {
	_, err := v.mut.Unlock()
	return err
}

type event string

const (
	Scheduled  event = "Video %s from %s has been scheduled for download"
	Error      event = "Video %s from %s could not be downloaded, failed with an error. "
	Downloaded event = "Video %s from %s has been downloaded successfully, and uploaded to videoservice"
)

func (v *VideoDLRequest) RecordEvent(inpEvent event, additionalErrorMsg string) error {
	website, err := GetWebsiteFromURL(v.ParentURL)
	if err != nil {
		return err
	}

	formattedMsg := fmt.Sprintf(string(inpEvent), v.VideoID, website)

	if additionalErrorMsg != "" {
		formattedMsg += fmt.Sprintf("\n\nError message: %s", additionalErrorMsg)
	}

	sql := "insert into archival_events (video_url, download_id, parent_url, event_message, event_time) VALUES ($1, $2, $3, $4, Now())"
	_, err = v.Db.Exec(sql, v.URL, v.DownloaddID, v.ParentURL, formattedMsg)
	return err
}
