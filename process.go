package main

import (
	"encoding/json"
	"io"
	"os"
	"slices"
	"strings"
)

type Request struct {
	Page         uint    `json:"page"`
	ItemsPerPage uint    `json:"items_per_page"`
	Filter       *Filter `json:"filter"`
}

type Filter struct {
	NameContains string   `json:"contains"`
	Tags         []string `json:"tags"`
}

type BugEntry struct {
	Image     string   `json:"image"`
	Name      string   `json:"name"`
	LatinName string   `json:"latin"`
	Tags      []string `json:"tags"`
}

func parse(data []byte) (Request, error) {
	var request Request
	if err := json.Unmarshal(data, &request); err != nil {
		return request, err
	}
	if request.Filter != nil {
		for i, tag := range request.Filter.Tags {
			request.Filter.Tags[i] = strings.ToLower(tag)
		}
		if strings.TrimSpace(request.Filter.NameContains) != "" {
			request.Filter.NameContains = strings.TrimSpace(strings.ToLower(request.Filter.NameContains))
		} else {
			request.Filter.NameContains = ""
		}

	}
	return request, nil
}

func process(req Request) (string, error) {
	metadata, err := readMetadata()
	if err != nil {
		return "", err
	}

	res := make([]BugEntry, 0)

	for _, bugEntry := range metadata {
		if filter(&bugEntry, req.Filter) {
			res = append(res, bugEntry)
		}
	}

	output := struct {
		TotalCount uint       `json:"total_count"`
		Items      []BugEntry `json:"items"`
	}{}

	output.Items = res
	output.TotalCount = uint(len(res))

	data, err := json.Marshal(output)
	return string(data), err
}

func filter(entry *BugEntry, filter *Filter) bool {
	if filter == nil {
		return true
	}

	if filter.NameContains != "" {
		if !strings.Contains(strings.ToLower(entry.Name), filter.NameContains) && !strings.Contains(strings.ToLower(entry.LatinName), filter.NameContains) {
			return false
		}
	}

	if len(filter.Tags) > 0 {
		res := false
		for _, tag := range entry.Tags {
			if slices.Contains(filter.Tags, strings.ToLower(tag)) {
				res = true
				break
			}
		}
		return res
	}

	return true
}

func readMetadata() ([]BugEntry, error) {
	var metadata []BugEntry
	metadataFile, err := os.Open(metadataPath)
	if err != nil {
		return metadata, err
	}
	defer metadataFile.Close()

	metadataRaw, err := io.ReadAll(metadataFile)
	if err != nil {
		return metadata, err
	}

	if err := json.Unmarshal(metadataRaw, &metadata); err != nil {
		return metadata, err
	}

	return metadata, nil
}
