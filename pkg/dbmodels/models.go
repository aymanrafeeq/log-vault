package models

import (
	"context"
	"fmt"
	"log"
	"logGen/model"
	"os"
	"regexp"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type queryComponent struct {
	key      string
	value    []string
	operator string
}

type LogLevel struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"uniqueIndex;size:100;not null"`
}

type LogComponent struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"uniqueIndex;size:200;not null"`
}

type LogHost struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"uniqueIndex;size:200;not null"`
}

// Main entry
type Entry struct {
	gorm.Model
	TimeStamp time.Time
	// keep string columns for backward-compat/query compatibility
	Level     string
	Component string
	Host      string
	RequestId string
	Message   string

	// foreign keys to normalized tables (nullable, cascade on update, set null on delete)
	LogLevelID     *uint         `gorm:"index"`
	LogLevel       *LogLevel     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:LogLevelID"`
	LogComponentID *uint         `gorm:"index"`
	LogComponent   *LogComponent `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:LogComponentID"`
	LogHostID      *uint         `gorm:"index"`
	LogHost        *LogHost      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:LogHostID"`
}

func (l Entry) String() string {
	if l.TimeStamp.IsZero() {
		return "Empty"
	} else {
		return fmt.Sprintf("%s | %s | %s | %s | %s | %s", l.TimeStamp, l.Level, l.Component, l.Host, l.RequestId, l.Message)

	}
}

func CreateDB(dbUrl string) (*gorm.DB, error) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level (Silent, Error, Warn, Info)
			IgnoreRecordNotFoundError: false,       // Ignore ErrRecordNotFound error
			Colorful:                  true,        // Enable color
		},
	)

	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{Logger: newLogger})
	if err != nil {
		return nil, fmt.Errorf("couldn't open database %v", err)
	}
	return db, nil
}

func InitDb(db *gorm.DB) error {
	// Migrate lookup tables first, then Entry
	if err := db.AutoMigrate(&LogLevel{}, &LogComponent{}, &LogHost{}, &Entry{}); err != nil {
		return err
	}
	return nil
}

// helper: find or create LogLevel by name
func findOrCreateLogLevel(ctx context.Context, db *gorm.DB, name string) (*LogLevel, error) {
	if strings.TrimSpace(name) == "" {
		return nil, nil
	}
	var lvl LogLevel
	tx := db.WithContext(ctx).Where(LogLevel{Name: name}).FirstOrCreate(&lvl)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &lvl, nil
}

func findOrCreateLogComponent(ctx context.Context, db *gorm.DB, name string) (*LogComponent, error) {
	if strings.TrimSpace(name) == "" {
		return nil, nil
	}
	var comp LogComponent
	tx := db.WithContext(ctx).Where(LogComponent{Name: name}).FirstOrCreate(&comp)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &comp, nil
}

func findOrCreateLogHost(ctx context.Context, db *gorm.DB, name string) (*LogHost, error) {
	if strings.TrimSpace(name) == "" {
		return nil, nil
	}
	var h LogHost
	tx := db.WithContext(ctx).Where(LogHost{Name: name}).FirstOrCreate(&h)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &h, nil
}

func AddEntry(db *gorm.DB, e model.LogEntry) error {

	ctx := context.Background()

	lvl, err := findOrCreateLogLevel(ctx, db, string(e.Level))
	if err != nil {
		return fmt.Errorf("find/create log level: %w", err)
	}
	comp, err := findOrCreateLogComponent(ctx, db, e.Component)
	if err != nil {
		return fmt.Errorf("find/create log component: %w", err)
	}
	host, err := findOrCreateLogHost(ctx, db, e.Host)
	if err != nil {
		return fmt.Errorf("find/create log host: %w", err)
	}

	x := Entry{
		TimeStamp: e.Time,
		Level:     string(e.Level),
		Component: e.Component,
		Host:      e.Host,
		RequestId: e.ReqID,
		Message:   e.Msg,
	}

	// attach IDs if found
	if lvl != nil {
		x.LogLevelID = &lvl.ID
	}
	if comp != nil {
		x.LogComponentID = &comp.ID
	}
	if host != nil {
		x.LogHostID = &host.ID
	}

	// Save entry
	if err := db.WithContext(ctx).Create(&x).Error; err != nil {
		return err
	}
	return nil
}

func parseQuery(parts []string) ([]queryComponent, error) {

	var ret []queryComponent
	// parts := strings.Fields(query)

	pattern := `^(?P<key>[^\s=!<>]+)\s*(?P<operator>=|!=|>=|<=|>|<)\s*(?P<value>.+)$`

	r, _ := regexp.Compile(pattern)
	for _, part := range parts {
		matches := r.FindStringSubmatch(part)
		if matches == nil {
			return nil, fmt.Errorf("invalid condition: %s", part)
		}

		val := strings.Split(matches[r.SubexpIndex("value")], ",")
		cond := queryComponent{
			key:      matches[r.SubexpIndex("key")],
			operator: matches[r.SubexpIndex("operator")],
			value:    val,
		}
		ret = append(ret, cond)
	}
	return ret, nil

}

func Query(db *gorm.DB, queryList []string) ([]Entry, error) {
	var ret []Entry

	// Parse the query string
	parsed, err := parseQuery(queryList)
	if err != nil {
		return nil, err
	}

	fmt.Println("Parsed conditions:", parsed)

	q := db
	for _, c := range parsed {
		if len(c.value) == 1 {
			// single value
			fmt.Printf("Applying condition: %s %s %s\n", c.key, c.operator, c.value[0])
			q = q.Where(fmt.Sprintf("%s %s ?", c.key, c.operator), c.value[0])
		} else {
			//multiple values and operator is !=
			if c.operator == "!=" {
				fmt.Printf("Applying NOT IN condition: %s NOT IN %v\n", c.key, c.value)
				q = q.Where(fmt.Sprintf("%s NOT IN ?", c.key), c.value)
			} else {
				// multi value and operator is =
				fmt.Printf("Applying IN condition: %s IN %v\n", c.key, c.value)
				q = q.Where(fmt.Sprintf("%s IN ?", c.key), c.value)
			}

		}
	}

	// Execute final query
	if err := q.Find(&ret).Error; err != nil {
		return nil, err
	}

	return ret, nil
}
