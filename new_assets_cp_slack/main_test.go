package main

import (
	"new_assets_cp_slack/utils/globals"
	"testing"
)

func TestSetDefaultValue(t *testing.T) {

	if *globals.Config.ConfigSlack.HttpAttempt != 3 {
		t.Errorf("ConfigSlack !=3")
	}

	if *globals.Config.ConfigCP.HttpAttempt != 3 {
		t.Errorf("ConfigCP !=3")
	}

}
