# go-airtable-api

Simple Go (Gin) API to be wired to Airtable.

## Prerequisites
- Go 1.21+ installed
- Airtable account and base

## Environment Variables
Set these before running:

```bash
export AIRTABLE_TOKEN=your_pat
export AIRTABLE_BASE_ID=appXXXXXXXXXXXXXX
export AIRTABLE_TABLE=Albums
```

## Run
```bash
go run .
```

## Notes
- Module path: `github.com/omarcarreon/go-airtable-api`.
- Do not commit real secrets; use env vars. A `.env` file is ignored by `.gitignore`.
