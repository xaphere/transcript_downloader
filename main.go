package main

import (
	"context"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/xaphere/transcript_downloader/parliament"
)

var log = logrus.New()

var storage = parliament.NewFileStorage("data/")
var getter = parliament.NewGetter(http.DefaultClient)

const delay = time.Second * 2

func main() {
	ctx := context.Background()
	// initial data download
	// we try to download the transcript as well as attached files
	errResults := download(ctx, 2010, 2021)
	// some requests fail because there is something broken with the attached files,
	// so we try again without the files
	for id := range errResults {
		downloadStenogram(ctx, id, false)
		time.Sleep(delay)
	}
}

func download(ctx context.Context, startYear, endYear int) map[int]error {
	errorResults := map[int]error{}
	for year := startYear; year <= endYear; year++ {
		for month := 1; month <= 12; month++ {
			meta, err := getter.GetStenogramsForMonth(ctx, year, month)
			if err != nil {
				log.WithError(err).Fatalln("failed to get stenogram metadata")
			}
			for _, m := range meta {
				err = downloadStenogram(ctx, m.ID, true)
				if err != nil {
					errorResults[m.ID] = err
				}
				time.Sleep(delay)
			}
		}
	}
	return errorResults
}

func downloadStenogram(ctx context.Context, stID int, shouldDownloadFiles bool) error {
	st, err := getter.GetStenogram(ctx, stID, shouldDownloadFiles)
	if err != nil {
		log.WithError(err).WithField("stenogram_id", stID).Errorln("failed to get stenogram")
		return err
	}
	log.WithFields(map[string]interface{}{
		"id":    st.ID,
		"title": st.Title,
		"date":  st.Date,
	}).Info()
	err = storage.Store(ctx, *st)
	if err != nil {
		log.WithError(err).WithField("stenogram_id", stID).Errorln("failed to store stenogram")
		return err
	}
	return nil
}
