package script

import (
	"github.com/gopherjs/gopherjs/js"
	"honnef.co/go/js/dom"
	"honnef.co/go/js/xhr"
)

var document = dom.GetWindow().Document()

// hasUpdatesAvailable returns true if there's at least one remaining update.
func hasUpdatesAvailable() bool {
	updates := document.GetElementsByClassName("go-package-update")
	for _, update := range updates {
		if len(update.GetElementsByClassName("disabled")) == 0 {
			return true
		}
	}
	return false
}

// UpdateGoPackage updates Go packages specified by importPathPattern.
func UpdateGoPackage(importPathPattern string) {
	var go_package = document.GetElementByID(importPathPattern)
	var go_package_button = go_package.GetElementsByClassName("update-button")[0].(*dom.HTMLAnchorElement)

	go_package_button.SetTextContent("Updating...")
	go_package_button.AddEventListener("click", false, func(event dom.Event) { event.PreventDefault() })
	go_package_button.SetTabIndex(-1)
	go_package_button.Class().Add("disabled")

	go func() {
		req := xhr.NewRequest("POST", "http://localhost:7043/-/update")
		req.SetRequestHeader("Content-Type", "application/x-www-form-urlencoded")
		err := req.Send("import_path_pattern=" + importPathPattern)
		if err != nil {
			println(err.Error())
			return
		}

		// Hide the "Updating..." label.
		go_package_button.Style().SetProperty("display", "none", "")

		// Show "No Updates Available" if there are no remaining updates.
		if !hasUpdatesAvailable() {
			document.GetElementByID("no_updates").(dom.HTMLElement).Style().SetProperty("display", "none", "")
		}

		// Move this Go package to "Installed Updates" list.
		installed_updates := document.GetElementByID("installed_updates").(dom.HTMLElement)
		installed_updates.Style().SetProperty("display", "", "")
		installed_updates.ParentNode().InsertBefore(go_package, installed_updates.NextSibling()) // Insert after.
	}() //gopherjs:blocking
}

func main() {
	js.Global.Set("update_go_package", UpdateGoPackage)
}
