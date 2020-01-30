package file

import (
	"fmt"
	"testing"

	"github.com/ne2blink/antenna/pkg/storage"
)

func Test_file_AppCGUD(t *testing.T) {
	options := make(map[string]interface{})
	options["path"] = "./db"
	file, err := storage.New("file", options)
	if err != nil {
		t.Errorf("%v", err)
	}
	app := storage.App{}
	app.Name = "001"
	app.SetSecret("")

	id, err := file.CreateApp(app)
	fmt.Println(app)

	app, err = file.GetApp(id)
	if err != nil {
		t.Errorf("%v", err)
	}
	if app.Name != "001" {
		t.Errorf("file storage CreateApp or GetApp")
	}
	fmt.Println(app)

	app.Name = "002"
	file.UpdateApp(app)
	if app.Name != "002" {
		t.Errorf("file storage UpdateApp")
	}
	fmt.Println(app)

	file.DeleteApp(app.ID)
	_, err = file.GetApp(id)
	if err == nil {
		t.Errorf("file storage DeleteApp")
	}
	// t.Errorf("Test")
}
