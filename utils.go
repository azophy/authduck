package main

import (
  "os"
)

var (
  BaseUrl = GetEnvOrDefault("BASE_URL", "http://localhost:3000")
)

func GetEnvOrDefault(varName string, defaultValue string) string {
  if os.Getenv(varName) != "" {
    return os.Getenv(varName)
  } else {
    return defaultValue
  }
}
