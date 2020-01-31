package file

import (
	"fmt"
	"os"
	"testing"

	"github.com/ne2blink/antenna/pkg/storage"
)

func Test_file_AppCGUD(t *testing.T) {
	options := make(map[string]interface{})
	path := "./test.db"
	options["path"] = path
	file, err := storage.New("file", options)
	if err != nil {
		t.Errorf("%v", err)
	}
	defer os.Remove(path)
	// defer file.Close()
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
		t.Errorf("file storage CreateApp or GetApp error")
	}
	fmt.Println(app)

	app.Name = "002"
	file.UpdateApp(app)
	if app.Name != "002" {
		t.Errorf("file storage UpdateApp error")
	}
	fmt.Println(app)

	file.DeleteApp(app.ID)
	_, err = file.GetApp(id)
	if err == nil {
		t.Errorf("file storage DeleteApp error")
	}

	// t.Errorf("Test")
}

func Test_file_Subscribe(t *testing.T) {
	options := make(map[string]interface{})
	path := "./test.db"
	options["path"] = path
	file, err := storage.New("file", options)
	if err != nil {
		t.Errorf("%v", err)
	}
	defer os.Remove(path)
	// defer file.Close()
	app := storage.App{}
	app.Name = "001"
	app.SetSecret("")

	id, err := file.CreateApp(app)
	a, err := file.ListSubscribedApps(123)
	fmt.Println(a)
	file.Subscribe(123, id)
	a, err = file.ListSubscribedApps(123)
	fmt.Println(a)
	aa, err := file.GetApp(id)
	fmt.Println(aa)
	aaa, err := file.ListSubscribers(id)
	fmt.Println(aaa)

	// t.Errorf("Test")
}
