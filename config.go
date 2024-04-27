package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Connection struct {
	Name    string `json:"name"`
	DSN     string `json:"dsn"`
	Driver  string `json:"driver"`
	Default bool   `json:"default"`
}

type Config struct {
	Limit          int          `json:"limit"`
	TimeoutSeconds int          `jon:"timeout"`
	Connections    []Connection `json:"connections"`
}

var (
	QueryName      string
	QueryFilePath  string
	ConfigFilePath string
	ConnectionName string
	Limit          int
)

func init() {
	flag.StringVar(&QueryName, "query", "", "the name of the query to execute")
	flag.StringVar(&QueryFilePath, "file", "", "the sqlc compatible file with named queries")
	flag.StringVar(&ConnectionName, "conn", "", "the connection to use")
	flag.StringVar(&ConfigFilePath, "config", "", "the config file path")
	flag.IntVar(&Limit, "limit", 50, "the query result limit")
	flag.Parse()
}

func GetConfigFilePath() string {
	dir, _ := os.UserHomeDir()
	if ConfigFilePath == "" {
		ConfigFilePath = filepath.Join(dir, ".sqlk.json")
	}

	return ConfigFilePath
}

func GetConnectionName() string {
	if ConnectionName == "" {
		ConnectionName = "default"
	}

	return ConnectionName
}

func GetLimit() int {
	if Limit == 0 {
		Limit = 50
	}

	return Limit
}

func GetQueryFilePath(conn Connection) string {
	dir, _ := os.UserHomeDir()
	if QueryFilePath == "" {
		QueryFilePath = filepath.Join(dir, fmt.Sprintf(".sqlk.%s.sql", GetConnectionName()))
	}

	return QueryFilePath
}

func GetQuery(queryFilePath string) (string, error) {
	buffer, err := os.ReadFile(queryFilePath)
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile(`name:\s+(\w+)`)
	queries := strings.Split(string(buffer), ";")

	for i, match := range re.FindAllStringSubmatch(string(buffer), -1) {
		if strings.ToLower(match[1]) == QueryName {
			return strings.TrimSpace(queries[i]), nil
		}
	}

	return "", nil
}
