package dataparser

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"sort"
	"text/tabwriter"
)

type user struct {
	Name    string `json:"name"`
	Lines   int    `json:"lines"`
	Commits int    `json:"commits"`
	Files   int    `json:"files"`
}

func getSliceOfAuthors() []user {
	var authorSlice []user

	for name, info := range authors {
		authorSlice = append(authorSlice, user{
			Name:    name,
			Lines:   info.lines,
			Commits: info.commit,
			Files:   info.fileCount,
		})
	}

	sort.Slice(authorSlice, func(i, j int) bool {
		switch info.orderby {
		case "lines":
			return authorSlice[i].Lines > authorSlice[j].Lines ||
				(authorSlice[i].Lines == authorSlice[j].Lines &&
					authorSlice[i].Commits > authorSlice[j].Commits) ||
				(authorSlice[i].Lines == authorSlice[j].Lines &&
					authorSlice[i].Commits == authorSlice[j].Commits &&
					authorSlice[i].Files > authorSlice[j].Files) ||
				(authorSlice[i].Lines == authorSlice[j].Lines &&
					authorSlice[i].Commits == authorSlice[j].Commits &&
					authorSlice[i].Files == authorSlice[j].Files &&
					authorSlice[i].Name < authorSlice[j].Name)
		case "commits":
			return authorSlice[i].Commits > authorSlice[j].Commits ||
				(authorSlice[i].Commits == authorSlice[j].Commits &&
					authorSlice[i].Lines > authorSlice[j].Lines) ||
				(authorSlice[i].Commits == authorSlice[j].Commits &&
					authorSlice[i].Lines == authorSlice[j].Lines &&
					authorSlice[i].Files > authorSlice[j].Files) ||
				(authorSlice[i].Lines == authorSlice[j].Lines &&
					authorSlice[i].Commits == authorSlice[j].Commits &&
					authorSlice[i].Files == authorSlice[j].Files &&
					authorSlice[i].Name < authorSlice[j].Name)
		case "files":
			return authorSlice[i].Files > authorSlice[j].Files ||
				(authorSlice[i].Files == authorSlice[j].Files &&
					authorSlice[i].Lines > authorSlice[j].Lines) ||
				(authorSlice[i].Files == authorSlice[j].Files &&
					authorSlice[i].Lines == authorSlice[j].Lines &&
					authorSlice[i].Commits > authorSlice[j].Commits) ||
				(authorSlice[i].Lines == authorSlice[j].Lines &&
					authorSlice[i].Commits == authorSlice[j].Commits &&
					authorSlice[i].Files == authorSlice[j].Files &&
					authorSlice[i].Name < authorSlice[j].Name)
		}
		return false
	})

	return authorSlice
}

func showDataTabular() (string, error) {
	buf := bytes.Buffer{}
	w := tabwriter.NewWriter(&buf, 0, 0, 1, ' ', 0)
	_, err := fmt.Fprintln(w, "Name\tLines\tCommits\tFiles")
	if err != nil {
		return "", err
	}
	for _, data := range getSliceOfAuthors() {
		_, err := fmt.Fprintf(
			w,
			"%s\t%d\t%d\t%d\n",
			data.Name,
			data.Lines,
			data.Commits,
			data.Files,
		)
		if err != nil {
			return "", err
		}
	}
	w.Flush()
	return buf.String(), nil
}

func showDataCSV() (string, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	err := writer.Write([]string{"Name", "Lines", "Commits", "Files"})
	if err != nil {
		return "", fmt.Errorf("error formatting to CSV: %v", err)
	}
	for _, data := range getSliceOfAuthors() {
		err := writer.Write([]string{
			data.Name,
			fmt.Sprintf("%d", data.Lines),
			fmt.Sprintf("%d", data.Commits),
			fmt.Sprintf("%d", data.Files),
		})
		if err != nil {
			return "", fmt.Errorf("error formatting to CSV: %v", err)
		}
	}
	writer.Flush()

	return buf.String(), nil
}

func showDataJSON() (string, error) {
	var jsonData []byte
	items := getSliceOfAuthors()
	jsonData, err := json.Marshal(items)
	if err != nil {
		return "", fmt.Errorf("error formatting to JSON: %v", err)
	}
	return string(jsonData), nil
}

func showDataJSONLines() (string, error) {
	var jsonData string
	for _, data := range getSliceOfAuthors() {
		jsonItem, err := json.Marshal(data)
		if err != nil {
			return "", fmt.Errorf("error formatting to JSON-lines: %v", err)
		}
		jsonData += string(jsonItem) + "\n"
	}
	return jsonData, nil
}

func ShowData() (string, error) {
	switch info.format {
	case "tabular":
		return showDataTabular()
	case "csv":
		return showDataCSV()
	case "json":
		return showDataJSON()
	case "json-lines":
		return showDataJSONLines()
	}
	return "", nil
}
