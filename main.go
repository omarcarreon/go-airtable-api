package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"

    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
    "github.com/mehanizm/airtable"
)

// album represents data about a record album.
type album struct {
    ID     string  `json:"id"`
    Title  string  `json:"title"`
    Artist string  `json:"artist"`
    Price  float64 `json:"price"`
}

// Store keeps dependencies for data access simple and testable
type Store struct {
    table *airtable.Table
}

func NewStore(table *airtable.Table) *Store { return &Store{table: table} }

func (s *Store) ListAlbums(c *gin.Context) {
    recs, err := s.table.GetRecords().Do()
    if err != nil {
        c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    out := make([]album, 0, len(recs.Records))
    for _, r := range recs.Records {
        f := r.Fields
        out = append(out, album{
            ID:     toString(f["id"]),
            Title:  toString(f["title"]),
            Artist: toString(f["artist"]),
            Price:  toFloat(f["price"]),
        })
    }
    c.IndentedJSON(http.StatusOK, out)
}

func (s *Store) GetAlbumByID(c *gin.Context) {
    id := c.Param("id")
    rec, err := s.table.GetRecord(id)
    if err != nil {
        c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    if rec == nil {
        c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
        return
    }
    f := rec.Fields
    a := album{
        ID:     toString(f["id"]),
        Title:  toString(f["title"]),
        Artist: toString(f["artist"]),
        Price:  toFloat(f["price"]),
    }
    c.IndentedJSON(http.StatusOK, a)
}
func (s *Store) CreateAlbum(c *gin.Context) {
    var a album
    if err := c.BindJSON(&a); err != nil {
        c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    payload := &airtable.Records{
        Records: []*airtable.Record{
            {Fields: map[string]any{
                "id":     a.ID,
                "title":  a.Title,
                "artist": a.Artist,
                "price":  a.Price,
            }},
        },
    }
	recs, err := s.table.AddRecords(payload)
    if err != nil || len(recs.Records) == 0 {
        c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "failed to create album"})
        return
    }
    f := recs.Records[0].Fields
    out := album{
        ID:     toString(f["id"]),
        Title:  toString(f["title"]),
        Artist: toString(f["artist"]),
        Price:  toFloat(f["price"]),
    }
    c.IndentedJSON(http.StatusCreated, out)
}

func main() {
    // Load .env file if present
    if err := godotenv.Load(); err != nil {
        // .env file is optional, continue if not found
    }

    // env config: read and validate vars.
    airtableToken := os.Getenv("AIRTABLE_TOKEN")
    airtableBaseID := os.Getenv("AIRTABLE_BASE_ID")
    airtableTable := os.Getenv("AIRTABLE_TABLE")
    if airtableToken == "" || airtableBaseID == "" || airtableTable == "" {
        log.Fatal("missing required env vars: AIRTABLE_TOKEN, AIRTABLE_BASE_ID, AIRTABLE_TABLE")
    }

    // Airtable client setup
    air := airtable.NewClient(airtableToken)
    tbl := air.GetTable(airtableBaseID, airtableTable)
    store := NewStore(tbl)

    router := gin.Default()
    router.GET("/albums", store.ListAlbums)
	router.GET("/albums/:id", store.GetAlbumByID)
	router.POST("/albums", store.CreateAlbum)

    router.Run("localhost:8080")
}


// removed: postAlbums (replaced by Store.CreateAlbum)

// getAlbumByID now handled by Store.GetAlbumByID

// helpers for Airtable field conversion (keep simple)
func toString(v interface{}) string { return fmt.Sprint(v) }

func toFloat(v interface{}) float64 {
    switch t := v.(type) {
    case float64:
        return t
    case int:
        return float64(t)
    case int64:
        return float64(t)
    case json.Number:
        f, _ := t.Float64()
        return f
    default:
        return 0
    }
}
