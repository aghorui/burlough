package project

import (
	"log"
	"testing"

	"github.com/aghorui/burlough/blog"
	"github.com/aghorui/burlough/util"
)


func TestProjectCreation(t *testing.T) {
	var state ProjectState
	var ulog []UpdateLog

	// tempdir, err :+ os.MkdirTemp(os.TempDir(), "burloughtemp*")

	state, err  := Load("/tmp/input")

	if err != nil {
		if err == ErrNoConfigFileFound {
			state, ulog, err = Init("/tmp/input", blog.ConfigFileParams{
				RenderPath: "../output",
				UseFileTimestampAsCreationDate: true,
			}, true)

			if err != nil {
				t.Fatalf("Error: %v", err)
			}

			log.Println("State Initialized.")

		} else {
			t.Fatalf("Error: %v", err)

		}
	} else {
		log.Println("State loaded from file.")
	}

	util.DumpVar(state, ulog, err)

	state.WriteConfig()

	util.DumpFile("/tmp/input/burlough.json")

	err = state.Render()

	if err != nil {
		t.Fatalf("Error: %v", err)
	}
}