package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	// ErrNotFound is returned when a record is not found
	ErrNotFound = gorm.ErrRecordNotFound
	// ErrDuplicateKey is returned when a duplicate key constraint is violated
	ErrDuplicateKey = errors.New("duplicate key error")
	// ErrNotConnected is returned when the database is not connected
	ErrNotConnected = errors.New("database not connected")
)

// Config holds the database configuration
type Config struct {
	Driver          string        // mysql, postgres, sqlite
	DSN             string        // Data Source Name
	MaxOpenConns    int           // Maximum number of open connections
	MaxIdleConns    int           // Maximum number of idle connections
	ConnMaxLifetime time.Duration // Maximum lifetime of a connection
	ConnMaxIdleTime time.Duration // Maximum idle time of a connection
	LogLevel        logger.LogLevel
}

// DB wraps gorm.DB with additional functionality
type DB struct {
	*gorm.DB
	config *Config
}

// New creates a new database connection
func New(config *Config, dialector gorm.Dialector) (*DB, error) {
	if config == nil {
		return nil, errors.New("config cannot be nil")
	}

	// Set defaults
	if config.MaxOpenConns == 0 {
		config.MaxOpenConns = 100
	}
	if config.MaxIdleConns == 0 {
		config.MaxIdleConns = 10
	}
	if config.ConnMaxLifetime == 0 {
		config.ConnMaxLifetime = time.Hour
	}
	if config.ConnMaxIdleTime == 0 {
		config.ConnMaxIdleTime = 10 * time.Minute
	}

	// GORM config
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(config.LogLevel),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	}

	// Open database connection
	gormDB, err := gorm.Open(dialector, gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying sql.DB
	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(config.ConnMaxIdleTime)

	return &DB{
		DB:     gormDB,
		config: config,
	}, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	if db.DB == nil {
		return ErrNotConnected
	}

	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}

// Ping checks the database connection
func (db *DB) Ping(ctx context.Context) error {
	if db.DB == nil {
		return ErrNotConnected
	}

	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}

	return sqlDB.PingContext(ctx)
}

// Transaction executes a function within a database transaction
func (db *DB) Transaction(fn func(*gorm.DB) error) error {
	return db.DB.Transaction(fn)
}

// WithContext returns a new DB instance with the given context
func (db *DB) WithContext(ctx context.Context) *DB {
	return &DB{
		DB:     db.DB.WithContext(ctx),
		config: db.config,
	}
}

// AutoMigrate runs auto migration for the given models
func (db *DB) AutoMigrate(models ...interface{}) error {
	return db.DB.AutoMigrate(models...)
}

// HealthCheck returns the database health status
func (db *DB) HealthCheck(ctx context.Context) error {
	if db.DB == nil {
		return ErrNotConnected
	}

	sqlDB, err := db.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	// Check connection
	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	// Check stats
	stats := sqlDB.Stats()
	if stats.OpenConnections == 0 {
		return errors.New("no open connections")
	}

	return nil
}

// Stats returns database connection pool statistics
func (db *DB) Stats() (map[string]interface{}, error) {
	if db.DB == nil {
		return nil, ErrNotConnected
	}

	sqlDB, err := db.DB.DB()
	if err != nil {
		return nil, err
	}

	stats := sqlDB.Stats()
	return map[string]interface{}{
		"max_open_connections": stats.MaxOpenConnections,
		"open_connections":     stats.OpenConnections,
		"in_use":               stats.InUse,
		"idle":                 stats.Idle,
		"wait_count":           stats.WaitCount,
		"wait_duration":        stats.WaitDuration.String(),
		"max_idle_closed":      stats.MaxIdleClosed,
		"max_lifetime_closed":  stats.MaxLifetimeClosed,
	}, nil
}
