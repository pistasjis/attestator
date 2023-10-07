package main

import (
	"embed"
	"os"
	"time"
)

type App struct {
	App          string
	IconLocation string
	Publisher    string
}

type AppsJson struct {
	DisplayName string `json:"displayName"`
	Reason      string `json:"reason"`
	Verdict     string `json:"verdict"`
}

type FinalResult struct {
	DisplayName string
	Icon        string
	Reason      string
	Verdict     string
}

//go:embed results_template.html
var Templates embed.FS

var Apps = []App{}
var FinalResults = []FinalResult{}

var CurrentTime = time.Now()
var Path, _ = os.Getwd()

// FIXME: We can't use this on Linux or macOS
var directory_seperator = `\`

// TODO: Figure out if it's possible to make this less "jank"?
var FileLocation = Path + directory_seperator + "attestation_" + CurrentTime.Format("2006-01-02_15.04.05") + ".html"
