package generateaudios

import (
	"encoding/json"
	"github.com/Pr3c10us/absolutego/internals/domains/queueport"
	commands2 "github.com/Pr3c10us/absolutego/internals/services/script/commands"
)

func Handler(c *queueport.Context) (*queueport.HandlerResult, error) {
	var data commands2.GenerateAudiosParameters
	err := json.Unmarshal(c.Data, &data)
	if err != nil {
		return nil, err
	}
	id, err := c.Services.ScriptServices.GenerateAudios.Handle(data)
	return &queueport.HandlerResult{
		ScriptId: id,
	}, err
}
