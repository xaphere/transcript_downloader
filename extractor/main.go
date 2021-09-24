package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

const dataFolder = "../data/"

const individualFileName = "individual_vote.csv"
const groupFileName = "group_vote.csv"
const aggregateFileName = "aggregate_vote.json"

type MP struct {
	Number int    `json:"number"`
	Name   string `json:"name"`
	Party  string `json:"party"`
	Votes  []MPVote
}
type MPVote struct {
	Label string   `json:"label"`
	Vote  VoteType `json:"vote"`
	Date  string   `json:"date"`
}

func main() {
	filepaths := collectStenogramFilePaths()

	for _, p := range filepaths {
		ivData, err := readCSV(filepath.Join(p, individualFileName))
		if err != nil {
			log.WithError(err).Error("failed to read individual vote data")
			continue
		}
		gvData, err := readCSV(filepath.Join(p, groupFileName))
		if err != nil {
			log.WithField("path", p).WithError(err).Error("failed to read group vote data")
			continue
		}
		iv, err := ExtractIndividualVoteDataFromCSV(ivData)
		if err != nil {
			log.WithField("path", p).WithError(err).Error("failed to extract individual votes")
			continue
		}
		gv, err := ExtractVoteDataFromCSV(gvData)
		if err != nil {
			log.WithField("path", p).WithError(err).Error("failed to extract group votes")
			continue
		}
		aggregate := []MP{}
		for _, i := range iv {
			mp := MP{
				Number: i.Number,
				Name:   i.Name,
				Party:  i.Party,
				Votes:  []MPVote{},
			}
			for _, g := range gv {
				mp.Votes = append(mp.Votes, MPVote{
					Label: g.Title,
					Vote:  i.Votes[g.ID],
					Date:  g.Date.String(),
				})
			}
			aggregate = append(aggregate, mp)
		}
		f, _ := os.Create(filepath.Join(p, aggregateFileName))
		enc := json.NewEncoder(f)
		enc.SetIndent("", "\t")
		enc.Encode(aggregate)
		f.Close()
	}
}

func readCSV(location string) ([][]string, error) {
	f, err := os.Open(location)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	result, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("%s is empty", location)
	}
	return result, nil
}

func collectStenogramFilePaths() []string {
	stenogramFiles := []string{}
	err := filepath.Walk(dataFolder, func(p string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(p) != ".json" {
			return nil
		}

		stenogramFiles = append(stenogramFiles, filepath.Dir(p))
		return nil
	})
	if err != nil {
		log.WithError(err).Error()
	}
	return stenogramFiles
}
