package main

import (
	"fmt"
	"log"
	"strings"
)

type LogDetail struct {
	Pairs []LogPair
}

type LogPair struct {
	Key   string
	Value interface{}
}

func init() {
	log.SetFlags(0)
}

func LogOperation(operation string, details []LogDetail) {
	log.Printf("Operation: %s\n%s", operation, LogDetailsWithJustification(details))
}

func LogDetailsWithJustification(details []LogDetail) string {
	maxWidth := getMaxWidth(details)

	var logMessage strings.Builder

	for _, detail := range details {
		for i, pair := range detail.Pairs {
			keyValue := fmt.Sprintf("%s: %v", pair.Key, pair.Value)
			if i == 0 {
				logMessage.WriteString(fmt.Sprintf("%-*s", maxWidth, keyValue))
			} else {
				logMessage.WriteString(fmt.Sprintf(", %-*s", maxWidth, keyValue))
			}
		}
		logMessage.WriteString("\n")
	}
	return logMessage.String()
}

func getMaxWidth(details []LogDetail) int {
	maxWidth := 0
	for _, detail := range details {
		for _, pair := range detail.Pairs {
			pairLength := len(pair.Key) + len(fmt.Sprintf("%v", pair.Value))
			if pairLength > maxWidth {
				maxWidth = pairLength
			}
		}
	}
	return maxWidth
}
