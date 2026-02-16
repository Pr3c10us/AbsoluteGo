package generatevideos

import (
	"encoding/binary"
	"github.com/Pr3c10us/absolutego/internals/domains/queueport"
)

func Handler(c *queueport.Context) (*queueport.HandlerResult, error) {
	scriptId := int64(binary.BigEndian.Uint64(c.Data))
	scriptId, err := c.Services.ScriptServices.GenerateVideos.Handle(scriptId)
	return &queueport.HandlerResult{
		ScriptId: scriptId,
	}, err
}
