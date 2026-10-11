package store

// Article is a published article row.
type Article struct {
	ID        string `gorm:"primaryKey"`
	Path      string `gorm:"not null"`
	Title     string `gorm:"not null"`
	Date      string `gorm:"not null;index"` // YYYY-MM-DD
	Summary   string `gorm:"not null"`
	Tags      string `gorm:"not null"` // JSON array
	Category  string `gorm:"not null"` // JSON array of ancestor sidebar titles
	InSidebar bool   `gorm:"not null"`
	HTML      string `gorm:"not null"`
}

// SidebarDoc holds the pruned sidebar tree as JSON in a single row.
type SidebarDoc struct {
	ID   int    `gorm:"primaryKey"`
	JSON string `gorm:"not null"`
}

// Asset is a published file. File is relative to the build's assets folder.
type Asset struct {
	URL         string `gorm:"primaryKey"` // /assets/<hash>/<name>
	Name        string `gorm:"not null"`
	Hash        string `gorm:"not null"`
	Size        int64  `gorm:"not null"`
	Kind        string `gorm:"not null"`
	ContentType string `gorm:"not null"`
	File        string `gorm:"not null"`
}

// ftsSchema is the full-text index. id is stored but not indexed; the trigram
// tokenizer enables substring matching, including CJK text.
const ftsSchema = `CREATE VIRTUAL TABLE IF NOT EXISTS articles_fts USING fts5(
	id UNINDEXED, title, summary, tags, body, tokenize='trigram'
)`
