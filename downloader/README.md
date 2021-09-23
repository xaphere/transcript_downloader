# Transcript downloader

This repo is intended for downloading transcripts for meetings of the Bulgarian national assembly from `https://parliament.bg`.

It contains a script that downloads the transcript and attached files for every Assembly meeting and stores them in a file structure in the `data` folder.
The data is organized by date `year/month/day`.



## Build

You need to have golang installed. 

```shell
go get github.com/xaphere/transcript_downloader
cd $GOPATH/src/github.com/xaphere/transcript_downloader
go build -mod=readonly

./transcript_downloader 
```

## Example of data 

 `data/2010/5/5/702.json`

```json
{
	"Pl_Sten_id": 702,
	"Pl_Sten_date": "2010-05-05",
	"Pl_Sten_sub": "СТОТНО ЗАСЕДАНИЕ\r\nСофия, сряда, 5 май 2010 г.\r\nОткрито в 9,02 ч.\r\n\r\n",
	"Pl_Sten_body": "Председателствал: председателят Цецка Цачева и заместник-председателите Лъчезар ...",
	"files": [
		{
			"Pl_StenDid": 1453,
			"Pl_StenDname": "Гласуване по парламентарни групи",
			"Pl_StenDfile": "/pub/StenD/gv050510.pdf",
			"Pl_StenDtype": "pdf"
		},
		{
			"Pl_StenDid": 1454,
			"Pl_StenDname": "Гласуване по парламентарни групи",
			"Pl_StenDfile": "/pub/StenD/gv050510.xls",
			"Pl_StenDtype": "xls"
		},
		{
			"Pl_StenDid": 1455,
			"Pl_StenDname": "Поименно гласуване",
			"Pl_StenDfile": "/pub/StenD/iv050510.pdf",
			"Pl_StenDtype": "pdf"
		},
		{
			"Pl_StenDid": 1456,
			"Pl_StenDname": "Поименно гласуване",
			"Pl_StenDfile": "/pub/StenD/iv050510.xls",
			"Pl_StenDtype": "xls"
		}
	],
	"video": {
		"Vid": 44,
		"Vidate": "2010-05-05"
	}
}
```