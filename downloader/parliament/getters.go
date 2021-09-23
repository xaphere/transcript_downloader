package parliament

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/multierr"
)

const location = "https://www.parliament.bg"
const api_location = location + "/api/v1/"

type Stenogram struct {
	ID      int               `json:"Pl_Sten_id"`
	Date    string            `json:"Pl_Sten_date"` // format "2015-03-27"
	Title   string            `json:"Pl_Sten_sub"`
	Body    string            `json:"Pl_Sten_body"`
	FileLoc []FileMeta        `json:"files"`
	Video   VideoMeta         `json:"video"`
	Files   map[string][]byte `json:"-"`
}

type VideoMeta struct {
	ID   int    `json:"Vid"`
	Date string `json:"Vidate"` // format "2015-03-27"
}

type FileType string

const (
	PDFType = "pdf"
	XLSType = "xls"
)

type FileMeta struct {
	ID       int      `json:"Pl_StenDid"`
	Name     string   `json:"Pl_StenDname"`
	Location string   `json:"Pl_StenDfile"`
	Type     FileType `json:"Pl_StenDtype"`
}

type StenogramMeta struct {
	ID    int    `json:"t_id"`
	Label int    `json:"t_label"`
	Date  string `json:"t_date"` // format "2015-03-27"
}

type Getter struct {
	client *http.Client
}

func NewGetter(client *http.Client) *Getter {
	return &Getter{
		client: client,
	}
}

func (g *Getter) GetStenogramsForMonth(ctx context.Context, year, month int) ([]StenogramMeta, error) {
	if !(1 <= month && month <= 12) {
		return nil, fmt.Errorf("invalid month %d", month)
	}
	if !(2010 <= year && year <= 2021) {
		return nil, fmt.Errorf("invalid year %d [2010;2021]", year)
	}

	loc := api_location + fmt.Sprintf("archive-period/bg/Pl_StenV/%d/%d/0/0", year, month)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, loc, nil)
	if err != nil {
		return nil, err
	}
	resp, err := g.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request '%s' failed with %d", req.URL.String(), resp.StatusCode)
	}
	var result []StenogramMeta
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (g *Getter) GetStenogram(ctx context.Context, stID int, shouldDownloadFiles bool) (*Stenogram, error) {
	loc := api_location + fmt.Sprintf("pl-sten/%d", stID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, loc, nil)
	if err != nil {
		return nil, err
	}

	resp, err := g.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request '%s' failed with %d", req.URL.String(), resp.StatusCode)
	}
	var result Stenogram
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	result.Files = map[string][]byte{}
	if !shouldDownloadFiles {
		return &result, nil
	}
	var errs error
	for _, fm := range result.FileLoc {
		data, err := g.GetPublicFile(ctx, fm.Location)
		if err != nil {
			errs = multierr.Append(errs, err)
			continue
		}
		result.Files[fm.Location] = data
	}
	return &result, errs
}

func (g *Getter) GetPublicFile(ctx context.Context, fileURL string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, location+fileURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := g.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request '%s' failed with %d", req.URL.String(), resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}
