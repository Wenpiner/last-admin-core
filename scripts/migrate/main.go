package main

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	esql "entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/schema"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/lib/pq"
	"github.com/wenpiner/last-admin-core/rpc/ent"
)

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	DBType       string
	Host         string
	Port         int
	Username     string
	Password     string
	DatabaseName string
	SSLMode      string
}

// DefaultConfigs returns default database configurations
func DefaultConfigs() map[string]DatabaseConfig {
	return map[string]DatabaseConfig{
		"postgres": {
			DBType:       "postgres",
			Host:         "127.0.0.1",
			Port:         5432,
			Username:     "postgres",
			Password:     "postgres",
			DatabaseName: "last_admin",
			SSLMode:      "disable",
		},
		"mysql": {
			DBType:       "mysql",
			Host:         "127.0.0.1",
			Port:         3306,
			Username:     "root",
			Password:     "root",
			DatabaseName: "last_admin",
			SSLMode:      "",
		},
		"sqlite3": {
			DBType:       "sqlite3",
			Host:         "",
			Port:         0,
			Username:     "",
			Password:     "",
			DatabaseName: "last_admin.db",
			SSLMode:      "",
		},
	}
}

// GetDSN returns the DSN string based on database type
func (c *DatabaseConfig) GetDSN() string {
	switch c.DBType {
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=True", c.Username, c.Password, c.Host, c.Port, c.DatabaseName)
	case "postgres":
		return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s", c.Username, c.Password, c.Host, c.Port, c.DatabaseName, c.SSLMode)
	case "sqlite3":
		return fmt.Sprintf("file:%s?_busy_timeout=100000&_fk=1", c.DatabaseName)
	default:
		return ""
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\n========================================")
	fmt.Println("   数据库迁移工具 - Database Migration")
	fmt.Println("========================================\n")

	// Step 1: Select database type
	dbType := selectDatabaseType(reader)

	// Step 2: Get default config and allow customization
	defaults := DefaultConfigs()
	config := defaults[dbType]

	fmt.Printf("\n请输入数据库配置信息 (按 Enter 使用默认值):\n\n")

	// Customize config based on database type
	if dbType != "sqlite3" {
		config.Host = promptInput(reader, "数据库主机地址", config.Host)
		config.Port = promptIntInput(reader, "数据库端口", config.Port)
		config.Username = promptInput(reader, "数据库用户名", config.Username)
		config.Password = promptInput(reader, "数据库密码", config.Password)
		config.DatabaseName = promptInput(reader, "数据库名称", config.DatabaseName)

		if dbType == "postgres" {
			config.SSLMode = promptInput(reader, "SSL 模式 (disable/require/prefer)", config.SSLMode)
		}
	} else {
		config.DatabaseName = promptInput(reader, "数据库文件路径", config.DatabaseName)
	}

	// Step 3: Confirm configuration
	fmt.Println("\n========================================")
	fmt.Println("   确认数据库配置")
	fmt.Println("========================================")
	printConfig(config)

	confirm := promptYesNo(reader, "是否继续执行迁移?")
	if !confirm {
		fmt.Println("已取消迁移操作")
		return
	}

	// Step 4: Execute migration
	fmt.Println("\n正在执行数据库迁移...")
	if err := performMigration(config); err != nil {
		fmt.Printf("❌ 迁移失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ 数据库迁移成功!")
}

func selectDatabaseType(reader *bufio.Reader) string {
	fmt.Println("请选择数据库类型:")
	fmt.Println("1. PostgreSQL (推荐)")
	fmt.Println("2. MySQL")
	fmt.Println("3. SQLite3")
	fmt.Print("\n请输入选项 (1-3) [默认: 1]: ")

	input := strings.TrimSpace(readLine(reader))
	if input == "" {
		input = "1"
	}

	switch input {
	case "1":
		return "postgres"
	case "2":
		return "mysql"
	case "3":
		return "sqlite3"
	default:
		fmt.Println("无效的选项，使用默认值 PostgreSQL")
		return "postgres"
	}
}

func promptInput(reader *bufio.Reader, prompt, defaultValue string) string {
	fmt.Printf("%s [%s]: ", prompt, defaultValue)
	input := strings.TrimSpace(readLine(reader))
	if input == "" {
		return defaultValue
	}
	return input
}

func promptIntInput(reader *bufio.Reader, prompt string, defaultValue int) int {
	fmt.Printf("%s [%d]: ", prompt, defaultValue)
	input := strings.TrimSpace(readLine(reader))
	if input == "" {
		return defaultValue
	}

	var value int
	_, err := fmt.Sscanf(input, "%d", &value)
	if err != nil {
		fmt.Printf("无效的数字，使用默认值 %d\n", defaultValue)
		return defaultValue
	}
	return value
}

func promptYesNo(reader *bufio.Reader, prompt string) bool {
	fmt.Printf("%s (y/n) [默认: n]: ", prompt)
	input := strings.TrimSpace(strings.ToLower(readLine(reader)))
	return input == "y" || input == "yes"
}

func readLine(reader *bufio.Reader) string {
	line, _ := reader.ReadString('\n')
	return strings.TrimSpace(line)
}

func printConfig(config DatabaseConfig) {
	fmt.Printf("数据库类型: %s\n", config.DBType)
	if config.DBType != "sqlite3" {
		fmt.Printf("主机地址: %s\n", config.Host)
		fmt.Printf("端口: %d\n", config.Port)
		fmt.Printf("用户名: %s\n", config.Username)
		fmt.Printf("密码: %s\n", strings.Repeat("*", len(config.Password)))
		fmt.Printf("数据库名: %s\n", config.DatabaseName)
		if config.DBType == "postgres" {
			fmt.Printf("SSL 模式: %s\n", config.SSLMode)
		}
	} else {
		fmt.Printf("数据库文件: %s\n", config.DatabaseName)
	}
	fmt.Println("========================================")
}

func performMigration(config DatabaseConfig) error {
	// Open database connection
	db, err := sql.Open(config.DBType, config.GetDSN())
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	fmt.Println("✓ 数据库连接成功")

	// Set connection pool settings
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(10)

	// Create ent driver
	drv := esql.OpenDB(config.DBType, db)
	defer drv.Close()

	// Create ent client
	client := ent.NewClient(ent.Driver(drv))
	defer client.Close()

	// Run migration
	fmt.Println("✓ 正在创建/更新数据库表...")
	if err := client.Schema.Create(ctx,
		schema.WithForeignKeys(false),
		schema.WithDropColumn(true),
		schema.WithDropIndex(true),
	); err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	fmt.Println("✓ 数据库表创建/更新完成")
	return nil
}

