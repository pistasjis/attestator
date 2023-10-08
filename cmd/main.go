// attestator (pistasjis Attestator) is an open-source application which lets you evaluate how apps on your computer respect your privacy and security.
package cmd

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/pistasjis/attestator/vars"
	"golang.org/x/sys/windows/registry"
)

func removeDuplicates(results []vars.FinalResult) []vars.FinalResult {
	seen := make(map[string]bool)
	uniqueResults := []vars.FinalResult{}

	for _, result := range results {
		key := result.DisplayName + result.Verdict + result.Reason

		// Check if the key has been seen before.
		if _, ok := seen[key]; !ok {
			seen[key] = true
			uniqueResults = append(uniqueResults, result)
		}
	}

	return uniqueResults
}

func RunAttestator(outputasjson bool) {
	fmt.Println("pistasjis Attestator running...")

	// TODO: DRY / make this code look better

	// open 64-bits uninstall registry
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall`, registry.READ)
	if err != nil {
		fmt.Println("Error opening registry key:", err)
		return
	}
	defer key.Close()

	// open 32-bits uninstall registry
	key2, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\WOW6432Node\Microsoft\Windows\CurrentVersion\Uninstall`, registry.READ)
	if err != nil {
		fmt.Println("Error opening registry key:", err)
		return
	}
	defer key.Close()

	// open users installed apps. Wow, this is really starting to annoy me
	key3, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Uninstall`, registry.READ)
	if err != nil {
		fmt.Println("Error opening registry key:", err)
		return
	}
	defer key.Close()

	subkeyNames, err := key.ReadSubKeyNames(-1)
	if err != nil {
		fmt.Println("Error reading subkey names:", err)
		return
	}

	subkey2Names, err := key2.ReadSubKeyNames(-1)
	if err != nil {
		fmt.Println("Error reading subkey names:", err)
		return
	}

	subkey3Names, err := key3.ReadSubKeyNames(-1)
	if err != nil {
		fmt.Println("Error reading subkey names:", err)
		return
	}

	for _, subkeyName := range subkeyNames {
		go addToApps(key, subkeyName)
	}

	for _, subkey2Name := range subkey2Names {
		go addToApps(key2, subkey2Name)
	}

	for _, subkey3Name := range subkey3Names {
		go addToApps(key3, subkey3Name)
	}

	var httpClient = &http.Client{Timeout: 10 * time.Second}

	fmt.Println("\nGetting applications from database")

	file, err := httpClient.Get("https://raw.githubusercontent.com/pistasjis/attestator/main/assets/apps.json")
	if err != nil {
		fmt.Println("Could not get JSON file")
		panic(err)
	}
	if file.StatusCode != 200 {
		fmt.Println("Error getting the JSON file, status code", file.StatusCode)
		os.Exit(1)
	}
	defer file.Body.Close()

	var appsJson []vars.AppsJson
	decoder := json.NewDecoder(file.Body)
	if err := decoder.Decode(&appsJson); err != nil {
		fmt.Println("Can't open json ", err)
	}

	for _, application := range vars.Apps {
		for _, app := range appsJson {
			if strings.Contains(application.App, app.DisplayName) {
				vars.FinalResults = append(vars.FinalResults, vars.FinalResult{Icon: application.IconLocation, DisplayName: app.DisplayName, Verdict: app.Verdict, Reason: app.Reason})
			}
		}
	}

	vars.FinalResults = removeDuplicates(vars.FinalResults)

	if !outputasjson {
		if err := createHTML(); err != nil {
			fmt.Println("Could not create HTML!")
			panic(err)
		}
	} else {
		if err := createJSON(); err != nil {
			fmt.Println("Could not create JSON!")
			panic(err)
		}
	}

	// open in browser. FIXME: once we add cross-platform support we need to add things like xdg-open etc
	if !outputasjson {
		exec.Command("cmd", "/c", "start", vars.FileLocation).Run()
		fmt.Printf("\ndone, saved at %s", vars.FileLocation)
	} else {
		exec.Command("cmd", "/c", "start", vars.JsonFileLocation).Run()
		fmt.Printf("\ndone, saved at %s", vars.FileLocation)
	}

}

// instead of repeating the code for adding each subkey, this function can be used inside of a for loop for each subkey inside of the three keys we have to loop through
func addToApps(key registry.Key, path string) error {
	subkey, err := registry.OpenKey(key, path, registry.READ)
	if err != nil {
		fmt.Println("Error opening subkey: ", err)
		return err
	}
	defer subkey.Close()

	displayName, _, err := subkey.GetStringValue("DisplayName")
	if err != nil {
		return err
	}

	publisher, _, err := subkey.GetStringValue("Publisher")
	if err != nil {
		return err
	}

	iconLocation, _, err := subkey.GetStringValue("DisplayIcon")
	if err != nil {
		return err
	}

	vars.Apps = append(vars.Apps, vars.App{App: displayName, IconLocation: iconLocation, Publisher: publisher})

	return nil
}

// this function outputs the result as HTML.
func createHTML() error {
	output, err := os.Create(vars.FileLocation)
	if err != nil {
		fmt.Println(err)
		panic("Couldn't create html file!")
	}

	tmpl, err := template.ParseFS(vars.Templates, "results_template.html")
	if err != nil {
		fmt.Println(err)
		return err
	}

	templateData := struct {
		Date    string
		Results interface{}
	}{
		Date:    vars.CurrentTime.Format("2006-01-02 15:04:05"),
		Results: vars.FinalResults,
	}

	if err := tmpl.Execute(output, templateData); err != nil {
		panic("couldnt execute template")
	}

	return nil
}

// This function outputs the result as JSON.
func createJSON() error {
	// Idea: Maybe use a struct instead of creating the data using this map? It would look cleaner.
	data := map[string]interface{}{
		"note":    "These are all the apps installed on your system that were in Attestator's database. If you feel like an app is missing, please make an issue on our GitHub (https://github.com/pistasjis/attestator) for the app.",
		"entries": vars.FinalResults,
		"date":    vars.CurrentTime.Unix(),
	}

	json, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Could not marshal JSON!")
		return err
	}

	if err := os.WriteFile(vars.JsonFileLocation, json, 0644); err != nil {
		fmt.Println("Could not create JSON file!")
		return err
	}

	return nil
}
