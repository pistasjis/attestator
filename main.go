// attestator (pistasjis Attestator) is an open-source application which lets you evaluate how apps on your computer respect your privacy and security.
package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"golang.org/x/sys/windows/registry"
)

func removeDuplicates(results []FinalResult) []FinalResult {
	seen := make(map[string]bool)
	uniqueResults := []FinalResult{}

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

func main() {
	if runtime.GOOS != "windows" {
		panic("pistasjis Attestator cannot be run on " + runtime.GOOS)
	}

	fmt.Println("pistasjis Attestator running...")

	if runtime.GOARCH == "arm64" {
		fmt.Println("The arm64 architecture is not exactly supported (I can't test it), so please report your findings on GitHub at https://github.com/pistasjis/attestator.")
	}

	if runtime.GOARCH == "386" && runtime.GOOS == "windows" {
		fmt.Println("The binary for the 386 (x86) architecture cannot get 64-bit apps, please keep that in mind.")
	}

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

	file, err := os.Open("apps.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var appsJson []AppsJson
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&appsJson); err != nil {
		fmt.Println("Can't open json ", err)
	}

	output, err := os.Create(FileLocation)
	if err != nil {
		fmt.Println(err)
		panic("Couldn't create html file!")
	}

	tmpl, err := template.ParseFS(Templates, "results_template.html")
	if err != nil {
		fmt.Println(err)
		panic("Could not parse template!")
	}

	for _, application := range Apps {
		for _, app := range appsJson {
			if strings.Contains(application.App, app.DisplayName) {
				FinalResults = append(FinalResults, FinalResult{Icon: application.IconLocation, DisplayName: app.DisplayName, Verdict: app.Verdict, Reason: app.Reason})
			}
		}
	}

	FinalResults = removeDuplicates(FinalResults)

	templateData := struct {
		Date    string
		Results interface{}
	}{
		Date:    CurrentTime.Format("2006-01-02 15:04:05"),
		Results: FinalResults,
	}

	if err := tmpl.Execute(output, templateData); err != nil {
		panic("couldnt execute template")
	}

	// open in browser
	exec.Command("cmd", "/c", "start", FileLocation).Run()

	fmt.Printf("\ndone, saved at %s", FileLocation)
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

	Apps = append(Apps, App{App: displayName, IconLocation: iconLocation, Publisher: publisher})

	return nil
}
