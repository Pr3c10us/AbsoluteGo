package createvab

import (
	"encoding/json"
	"github.com/Pr3c10us/absolutego/internals/domains/queueport"
	"github.com/Pr3c10us/absolutego/internals/services/vab/commands"
)

func Handler(c *queueport.Context) (*queueport.HandlerResult, error) {
	var data commands.CreateVABParameter
	err := json.Unmarshal(c.Data, &data)
	if err != nil {
		return nil, err
	}
	id, err := c.Services.VABServices.CreateVAB.Handle(data)
	return &queueport.HandlerResult{
		VabId: id,
	}, err
}
