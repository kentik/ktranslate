package fetch

// #include <sys/file.h>
import "C"

import (
	"github.com/kentik/golog/logger"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	MAX_DOWNLOAD_TIME    = 60 * time.Second
	MAX_AGE_TO_KEEP_FILE = 24 * time.Hour
	MMAP_SUFFIX          = ".mmap"
	COUNT_SUFFIX         = ".count"
	FETCH_LOCK_FILE      = "/tmp/ch_fetch_file.lock"
)

func FetchFromURL(url string, urlBackup string, file string, log *logger.Logger) error {

	// Make a lock file
	lock, err := os.Create(FETCH_LOCK_FILE)
	if err != nil {
		return err
	}
	defer lock.Close()
	fd := lock.Fd()
	returnCode := C.flock(C.int(fd), C.LOCK_EX|C.LOCK_NB)
	if returnCode < 0 {
		// Another process is trying to dowload. Sleep for download time, and then move on.
		log.Warnf("", "Could not get lock. Waiting %v and then continuing", MAX_DOWNLOAD_TIME)
		time.Sleep(MAX_DOWNLOAD_TIME)
	} else {
		log.Infof("", "Lock aquired on %s", file)
		defer func() {
			returnCode = C.flock(C.int(fd), C.LOCK_UN)
			if returnCode < 0 {
				log.Errorf("", "Flock unlock returned bad status: %d", returnCode)
			}
		}()
	}

	// Only fetch file once/hr
	if info, err := os.Stat(file); err == nil {
		if info.ModTime().Add(MAX_AGE_TO_KEEP_FILE).After(time.Now()) {
			log.Infof("", "File %s too new to download", file)
			// Nohting to do
			return nil
		}
	}

	// If we get here, actually download the file
	response, err := http.Get(url)
	if err != nil || response.StatusCode != 200 {
		log.Errorf("", "Error downloading file: %v", err)

		// Try getting from backup url
		if strings.HasPrefix(urlBackup, "http") {
			log.Infof("", "Trying backup url %s", urlBackup)
			return FetchFromURL(urlBackup, "", file, log)
		} else {
			return err
		}
	}

	defer response.Body.Close()

	//open a file for writing
	fh, err := os.Create(file)
	if err != nil {
		return err
	}
	defer fh.Close()

	// Use io.Copy to just dump the response body to the file. This supports huge files
	_, err = io.Copy(fh, response.Body)
	if err != nil {
		return err
	}

	// rm the mmaped version of this file, if it exists.
	log.Infof("", "Removing old version of mmap: %s", file+MMAP_SUFFIX)
	os.Remove(file + MMAP_SUFFIX)

	return nil
}
