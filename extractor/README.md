# Extractor

This script will walk through downloaded assembly `data`, find `individual_vote.csv` and `group_vote.csv` and try to
create `aggregate_vote.json`. It extracts the reason for the vote from `group_vote.csv` and matches is with the number
of the vote in `individual_vote.csv`.

The aggregated data is stored in `aggregate_vote.json` in the same folder as the source data.

## Mapping

Individual votes contains some cryptic data for how the MPs voted.

From my understanding the symbols mean:

|symbol | meaning    |
| ----- | ---------- |
|  +    | For        |
|  =    | Abstain    |
|  -    | Against    |
|  0    | NoVote     |
|  П    | Here       |
|  О    | Absent     |
|  Р    | Registered |

Where `Here`, `Absent`, `Registered` are used when the assembly counts the MPs to see if there is quorum.

### example for individual votes

```
0,,,,,1,2,3,4,5,,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24
1,АЛЕКСАНДЪР  РУМЕНОВ НЕНКОВ,,334.0,ГЕРБ,П,+,+,-,+,,+,+,+,+,+,+,0,+,+,+,+,0,0,0,0,0,0,0,0
2,АЛЕКСАНДЪР ВЛАДИМИРОВ РАДОСЛАВОВ,,333.0,КБ,П,+,+,+,+,,+,+,+,+,=,+,+,+,=,+,+,+,0,0,0,0,0,-,+
3,АЛЕКСАНДЪР СТОЙЧЕВ СТОЙКОВ,,602.0,ГЕРБ,П,+,+,-,+,,+,+,+,+,+,+,+,+,+,+,+,+,+,+,0,0,+,+,+
4,АЛИОСМАН  ИБРАИМ ИМАМОВ,,336.0,ДПС,П,+,+,+,0,,+,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
5,АНАСТАС ВАСИЛЕВ АНАСТАСОВ,,337.0,ГЕРБ,П,+,+,-,+,,+,+,+,+,+,+,+,+,+,+,+,+,+,+,+,0,0,0,0
6,АНАТОЛИЙ ВЕЛИКОВ ЙОРДАНОВ,,338.0,ГЕРБ,П,+,+,-,+,,+,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
```

### example for aggregate_vote.json

```json
[
	{
		"number": 334,
		"name": "АЛЕКСАНДЪР  РУМЕНОВ НЕНКОВ",
		"party": "ГЕРБ",
		"Votes": [
			{
				"label": "Програмата за работата на НС 13 - 15 януари 2010 г.",
				"vote": "for",
				"date": "2010-01-13 09:11:00 +0000 UTC"
			},
			{
				"label": "Програмата за работата на НС 13 - 15 януари 2010 г.-процедура за пряко предаване-Кр.Велчев",
				"vote": "for",
				"date": "2010-01-13 09:19:00 +0000 UTC"
			}
		]
	},
	{
		"number": 333,
		"name": "АЛЕКСАНДЪР ВЛАДИМИРОВ РАДОСЛАВОВ",
		"party": "КБ",
		"Votes": [
			{
				"label": "Програмата за работата на НС 13 - 15 януари 2010 г.",
				"vote": "for",
				"date": "2010-01-13 09:11:00 +0000 UTC"
			},
			{
				"label": "Програмата за работата на НС 13 - 15 януари 2010 г.-процедура за пряко предаване-Кр.Велчев",
				"vote": "for",
				"date": "2010-01-13 09:19:00 +0000 UTC"
			}
		]
	},
	{
		"number": 602,
		"name": "АЛЕКСАНДЪР СТОЙЧЕВ СТОЙКОВ",
		"party": "ГЕРБ",
		"Votes": [
			{
				"label": "Програмата за работата на НС 13 - 15 януари 2010 г.",
				"vote": "for",
				"date": "2010-01-13 09:11:00 +0000 UTC"
			},
			{
				"label": "Програмата за работата на НС 13 - 15 януари 2010 г.-процедура за пряко предаване-Кр.Велчев",
				"vote": "for",
				"date": "2010-01-13 09:19:00 +0000 UTC"
			}
		]
	},
	{
		"number": 336,
		"name": "АЛИОСМАН  ИБРАИМ ИМАМОВ",
		"party": "ДПС",
		"Votes": [
			{
				"label": "Програмата за работата на НС 13 - 15 януари 2010 г.",
				"vote": "for",
				"date": "2010-01-13 09:11:00 +0000 UTC"
			},
			{
				"label": "Програмата за работата на НС 13 - 15 януари 2010 г.-процедура за пряко предаване-Кр.Велчев",
				"vote": "for",
				"date": "2010-01-13 09:19:00 +0000 UTC"
			}
		]
	}
]
```