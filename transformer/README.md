# Transformer

This repo contains a system that will walk through downloaded assembly `data`. It will try to find Vote data typically stored in `xls` files.
It will convert those `xls` files to `csv` and store them in the same folder as the original vote data. 
The produced names are `individual_vote.csv` and `group_vote.csv`

- `individual_vote` contains the information on how every single MP has voted.
- `group_vote` contains aggregate information on how parties have voted. Also the description of the vote itself.

## Run 

```shell
docker compose up -d
docker logs -f extractor
docker compose down
```

## Setup

This system requires `data` folder with already downloaded assembly information to be in the correct place. See `docker-compose.yml` for more info.

It starts tree services `storage`, `transformer` and `extractor`.

- `transformer` is a service that receives a post request with fileURL from which to download a file and try to transform it in csv format.
- `storage` is a simple file server written in go. we need it to provide a local url for the xls files stored in `data` folder
- `extractor` is a script that craws `data` folder looking for `json` information for any stenogram. From the json it extracts the individual and group `xls` files and sends them to the `transformer`. In the end it stores the produced `csv` files in the same location as the `json`