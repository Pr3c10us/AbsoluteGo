package generatescript

import (
	"encoding/json"
	"github.com/Pr3c10us/absolutego/internals/domains/queueport"
	commands2 "github.com/Pr3c10us/absolutego/internals/services/script/commands"
)

func Handler(c *queueport.Context) error {
	var data commands2.GenerateScriptParameters
	err := json.Unmarshal(c.Data, &data)
	if err != nil {
		return err
	}
	err = c.Services.ScriptServices.GenerateScript.Handle(data)
	return err
}
