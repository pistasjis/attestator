package vars

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

var Templates embed.FS

var Apps = []App{}
var FinalResults = []FinalResult{}

var CurrentTime = time.Now()
var Path, _ = os.Getwd()

// FIXME: We can't use this on Linux or macOS. Maybe we could make cmd.Execute() set the directory separator on program start based on opearing system?
var directory_separator = `\`

// TODO: Figure out if it's possible to make this less "jank"?
var FileLocation = Path + directory_separator + "attestation_" + CurrentTime.Format("2006-01-02_15.04.05") + ".html"

var JsonFileLocation = Path + directory_separator + "attestation_" + CurrentTime.Format("2006-01-02_15.04.05") + ".json"
