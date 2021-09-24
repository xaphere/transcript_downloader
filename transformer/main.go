package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

const storageURL = "http://storage:4000/"
const transformURL = "http://transformer:5000/transform"

func main() {
	filepaths := collectStenogramFilePaths()

	for _, p := range filepaths {
		log.WithField("path", p).Info("extracting data")
		result, err := extractVotingFilepathsFromStenogram(p)
		if err != nil {
			log.Error(err)
			continue
		}
		base, _ := filepath.Split(p)
		data, err := readCSVDataFromVotingFile(result.PartyLoc)
		if err != nil {
			log.WithError(err).Error("failed to get party vote")
			continue
		}
		err = storeData(data, filepath.Join(base, "group_vote.csv"))
		if err != nil {
			log.WithError(err).Error("failed to store group data")
		}
		data, err = readCSVDataFromVotingFile(result.IndividualLoc)
		if err != nil {
			log.WithError(err).Error("failed to get individual vote")
			continue
		}
		err = storeData(data, filepath.Join(base, "individual_vote.csv"))
		if err != nil {
			log.WithError(err).Error("failed to store individual data")
		}
		time.Sleep(time.Millisecond * 100)
	}
}

func storeData(data [][]string, location string) error {
	f, err := os.Create(location)
	if err != nil {
		return err
	}
	defer f.Close()
	err = csv.NewWriter(f).WriteAll(data)
	if err != nil {
		return err
	}
	return nil
}

//////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////

func readCSVDataFromVotingFile(fileLoc string) ([][]string, error) {
	fileLoc = strings.ReplaceAll(fileLoc, `\`, `/`)
	fileLoc = storageURL + fileLoc
	data, _ := json.Marshal(struct {
		File string `json:"fileURL"`
	}{File: fileLoc})

	resp, err := http.Post(transformURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Error(string(body))
		return nil, err
	}
	reader := csv.NewReader(resp.Body)
	return reader.ReadAll()
}

//////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////

type VoteMeta struct {
	ID            string
	IndividualLoc string
	PartyLoc      string
}

func extractVotingFilepathsFromStenogram(p string) (*VoteMeta, error) {
	f, err := os.Open(p)
	if err != nil {
		log.WithError(err).WithField("file", p).Errorln("failed to read file")
		return nil, err
	}
	defer f.Close()
	stenogram := Stenogram{}
	err = json.NewDecoder(f).Decode(&stenogram)
	if err != nil {
		log.WithError(err).WithField("file", p).Errorln("failed to decode json file")
		return nil, err
	}
	//p = strings.ReplaceAll(p, `\`, `/`)
	base, fName := filepath.Split(p)

	result := VoteMeta{
		ID: strings.TrimSuffix(fName, ".json"),
	}
	for _, attachment := range stenogram.FileLoc {
		if attachment.Type != XLSType {
			continue
		}
		_, loc := filepath.Split(attachment.Location)
		loc = filepath.Join(base, loc)
		if attachment.Name == "Поименно гласуване" {
			result.IndividualLoc = loc
		}
		if attachment.Name == "Гласуване по парламентарни групи" {
			result.PartyLoc = loc
		}

	}
	if result.PartyLoc == "" {
		return nil, fmt.Errorf("missing party vote file for '%s'", p)
	}
	if result.IndividualLoc == "" {
		return nil, fmt.Errorf("missing individual vote file for '%s'", p)
	}
	return &result, nil
}

func collectStenogramFilePaths() []string {
	stenogramFiles := []string{}
	err := filepath.Walk("./data/", func(p string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(p) != ".json" {
			return nil
		}

		stenogramFiles = append(stenogramFiles, p)
		return nil
	})
	if err != nil {
		log.WithError(err).Error()
	}
	return stenogramFiles
}

type Stenogram struct {
	ID      int        `json:"Pl_Sten_id"`
	Date    string     `json:"Pl_Sten_date"` // format "2015-03-27"
	Title   string     `json:"Pl_Sten_sub"`
	Body    string     `json:"Pl_Sten_body"`
	FileLoc []FileMeta `json:"files"`
}

type FileType string

const (
	PDFType FileType = "pdf"
	XLSType FileType = "xls"
)

type FileMeta struct {
	ID       int      `json:"Pl_StenDid"`
	Name     string   `json:"Pl_StenDname"`
	Location string   `json:"Pl_StenDfile"`
	Type     FileType `json:"Pl_StenDtype"`
}
