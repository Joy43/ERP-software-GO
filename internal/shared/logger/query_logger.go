package query_logger

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
	"gorm.io/gorm"
)

// QueryLogger contains configuration for database query logging
type QueryLogger struct {
	slowThreshold time.Duration
	logFile       *os.File
}

//---- QueryContext holds information about a database query execution---
type QueryContext struct {
	Operation   string
	TableName   string
	Query       string
	Args        []interface{}
	StartTime   time.Time
	Duration    time.Duration
	RowsAffected int64
	Error       error         
	Result      interface{}   
}

// ------NewQueryLogger creates a new database query logger -------
func NewQueryLogger(slowThresholdMS int64) *QueryLogger {
	return &QueryLogger{
		slowThreshold: time.Duration(slowThresholdMS) * time.Millisecond,
	}
}

// BeforeQuery logs information before executing a query
func (ql *QueryLogger) BeforeQuery(ctx context.Context, operation, table, query string, args ...interface{}) {
	logEntry := map[string]interface{}{
		"timestamp": time.Now().Format("2006-01-02 15:04:05.000"),
		"level":     "INFO",
		"event":     "QUERY_START",
		"operation": operation,
		"table":     table,
		"query":     query,
		"args":      args,
	}
	
	logMessage := fmt.Sprintf(
		"[QUERY_START] Operation: %s | Table: %s | Query: %s",
		operation, table, query,
	)
	
	log.Println(logMessage)
	ql.logToFile(logEntry)
}

// AfterQuery logs information after executing a query
func (ql *QueryLogger) AfterQuery(qCtx *QueryContext) {
	logEntry := map[string]interface{}{
		"timestamp":     time.Now().Format("2006-01-02 15:04:05.000"),
		"level":         "INFO",
		"event":         "QUERY_SUCCESS",
		"operation":     qCtx.Operation,
		"table":         qCtx.TableName,
		"duration_ms":   qCtx.Duration.Milliseconds(),
		"rows_affected": qCtx.RowsAffected,
	}

	logMessage := fmt.Sprintf(
		"[QUERY_SUCCESS] Operation: %s | Table: %s | Duration: %dms | Rows Affected: %d",
		qCtx.Operation, qCtx.TableName, qCtx.Duration.Milliseconds(), qCtx.RowsAffected,
	)

	log.Println(logMessage)
	ql.logToFile(logEntry)

	// --------- Check for slow query ----------
	if qCtx.Duration > ql.slowThreshold {
		ql.LogSlowQuery(qCtx)
	}
}

//---------- OnQueryError logs errors that occurred during query execution
func (ql *QueryLogger) OnQueryError(qCtx *QueryContext) {
	logEntry := map[string]interface{}{
		"timestamp":     time.Now().Format("2006-01-02 15:04:05.000"),
		"level":         "ERROR",
		"event":         "QUERY_ERROR",
		"operation":     qCtx.Operation,
		"table":         qCtx.TableName,
		"query":         qCtx.Query,
		"args":          qCtx.Args,
		"duration_ms":   qCtx.Duration.Milliseconds(),
		"error":         qCtx.Error.Error(),
	}

	logMessage := fmt.Sprintf(
		"[QUERY_ERROR] Operation: %s | Table: %s | Duration: %dms | Error: %v",
		qCtx.Operation, qCtx.TableName, qCtx.Duration.Milliseconds(), qCtx.Error,
	)

	log.Println(logMessage)
	ql.logToFile(logEntry)
}

// LogSlowQuery logs queries that take longer than the threshold
func (ql *QueryLogger) LogSlowQuery(qCtx *QueryContext) {
	logEntry := map[string]interface{}{
		"timestamp":     time.Now().Format("2006-01-02 15:04:05.000"),
		"level":         "WARN",
		"event":         "SLOW_QUERY",
		"operation":     qCtx.Operation,
		"table":         qCtx.TableName,
		"query":         qCtx.Query,
		"args":          qCtx.Args,
		"duration_ms":   qCtx.Duration.Milliseconds(),
		"threshold_ms":  ql.slowThreshold.Milliseconds(),
	}

	logMessage := fmt.Sprintf(
		"[SLOW_QUERY] Operation: %s | Table: %s | Duration: %dms (Threshold: %dms) | Query: %s",
		qCtx.Operation, qCtx.TableName, qCtx.Duration.Milliseconds(), ql.slowThreshold.Milliseconds(), qCtx.Query,
	)

	log.Println(logMessage)
	ql.logToFile(logEntry)
}

// MonitorDatabaseOperation wraps a database operation with comprehensive logging
func (ql *QueryLogger) MonitorDatabaseOperation(
	ctx context.Context,
	operation string,
	tableName string,
	query string,
	executeFunc func() (interface{}, error),
	args ...interface{},
) (interface{}, error) {
	//------------ Log before executionp--------------
	ql.BeforeQuery(ctx, operation, tableName, query, args...)

	//--------- Record start time--------
	startTime := time.Now()

	//-------- Execute the operation----------
	result, err := executeFunc()

	//------------ Calculate duration---------
	duration := time.Since(startTime)

	// ---------- Create query context-----------
	qCtx := &QueryContext{
		Operation:   operation,
		TableName:   tableName,
		Query:       query,
		Args:        args,
		StartTime:   startTime,
		Duration:    duration,
		Error:       err,
		Result:      result,
	}

	//-------- Log based on result--------
	if err != nil {
		ql.OnQueryError(qCtx)
		return nil, err
	}

	ql.AfterQuery(qCtx)
	return result, nil
}

// GormCallback sets up GORM callbacks for automatic logging
func (ql *QueryLogger) GormCallback(db *gorm.DB) *gorm.DB {
	db.Callback().Create().Before("gorm:create").Register("logger:before_create", func(db *gorm.DB) {
		db.Set("logger_start_time", time.Now())
		ql.BeforeQuery(db.Statement.Context, "CREATE", db.Statement.Table, db.Statement.SQL.String())
	})

	db.Callback().Create().After("gorm:create").Register("logger:after_create", func(db *gorm.DB) {
		startTime, _ := db.Get("logger_start_time")
		duration := time.Since(startTime.(time.Time))
		
		if db.Error != nil {
			ql.OnQueryError(&QueryContext{
				Operation:   "CREATE",
				TableName:   db.Statement.Table,
				Query:       db.Statement.SQL.String(),
				Duration:    duration,
				Error:       db.Error,
			})
		} else {
			ql.AfterQuery(&QueryContext{
				Operation:    "CREATE",
				TableName:    db.Statement.Table,
				Query:        db.Statement.SQL.String(),
				Duration:     duration,
				RowsAffected: db.RowsAffected,
			})
		}
	})

	db.Callback().Query().Before("gorm:query").Register("logger:before_query", func(db *gorm.DB) {
		db.Set("logger_start_time", time.Now())
		ql.BeforeQuery(db.Statement.Context, "READ", db.Statement.Table, db.Statement.SQL.String())
	})

	db.Callback().Query().After("gorm:query").Register("logger:after_query", func(db *gorm.DB) {
		startTime, _ := db.Get("logger_start_time")
		duration := time.Since(startTime.(time.Time))
		
		if db.Error != nil && db.Error != gorm.ErrRecordNotFound {
			ql.OnQueryError(&QueryContext{
				Operation:   "READ",
				TableName:   db.Statement.Table,
				Query:       db.Statement.SQL.String(),
				Duration:    duration,
				Error:       db.Error,
			})
		} else {
			ql.AfterQuery(&QueryContext{
				Operation:    "READ",
				TableName:    db.Statement.Table,
				Query:        db.Statement.SQL.String(),
				Duration:     duration,
				RowsAffected: db.RowsAffected,
			})
		}
	})

	db.Callback().Update().Before("gorm:update").Register("logger:before_update", func(db *gorm.DB) {
		db.Set("logger_start_time", time.Now())
		ql.BeforeQuery(db.Statement.Context, "UPDATE", db.Statement.Table, db.Statement.SQL.String())
	})

	db.Callback().Update().After("gorm:update").Register("logger:after_update", func(db *gorm.DB) {
		startTime, _ := db.Get("logger_start_time")
		duration := time.Since(startTime.(time.Time))
		
		if db.Error != nil {
			ql.OnQueryError(&QueryContext{
				Operation:   "UPDATE",
				TableName:   db.Statement.Table,
				Query:       db.Statement.SQL.String(),
				Duration:    duration,
				Error:       db.Error,
			})
		} else {
			ql.AfterQuery(&QueryContext{
				Operation:    "UPDATE",
				TableName:    db.Statement.Table,
				Query:        db.Statement.SQL.String(),
				Duration:     duration,
				RowsAffected: db.RowsAffected,
			})
		}
	})

	db.Callback().Delete().Before("gorm:delete").Register("logger:before_delete", func(db *gorm.DB) {
		db.Set("logger_start_time", time.Now())
		ql.BeforeQuery(db.Statement.Context, "DELETE", db.Statement.Table, db.Statement.SQL.String())
	})

	db.Callback().Delete().After("gorm:delete").Register("logger:after_delete", func(db *gorm.DB) {
		startTime, _ := db.Get("logger_start_time")
		duration := time.Since(startTime.(time.Time))
		
		if db.Error != nil {
			ql.OnQueryError(&QueryContext{
				Operation:   "DELETE",
				TableName:   db.Statement.Table,
				Query:       db.Statement.SQL.String(),
				Duration:    duration,
				Error:       db.Error,
			})
		} else {
			ql.AfterQuery(&QueryContext{
				Operation:    "DELETE",
				TableName:    db.Statement.Table,
				Query:        db.Statement.SQL.String(),
				Duration:     duration,
				RowsAffected: db.RowsAffected,
			})
		}
	})

	return db
}

// logToFile writes log entry to file (optional)
func (ql *QueryLogger) logToFile(entry map[string]interface{}) {
	if ql.logFile != nil {
		logLine := fmt.Sprintf("%+v\n", entry)
		ql.logFile.WriteString(logLine)
	}
}

// Close closes the log file
func (ql *QueryLogger) Close() error {
	if ql.logFile != nil {
		return ql.logFile.Close()
	}
	return nil
}
