package generatevideos

import (
	"encoding/json"
	"github.com/Pr3c10us/absolutego/internals/domains/queueport"
	commands2 "github.com/Pr3c10us/absolutego/internals/services/script/commands"
)

func Handler(c *queueport.Context) (*queueport.HandlerResult, error) {
	var data commands2.GenerateVideosParameter
	err := json.Unmarshal(c.Data, &data)
	if err != nil {
		return nil, err
	}
	scriptId, err := c.Services.ScriptServices.GenerateVideos.Handle(data)
	return &queueport.HandlerResult{
		ScriptId: scriptId,
	}, err
}
