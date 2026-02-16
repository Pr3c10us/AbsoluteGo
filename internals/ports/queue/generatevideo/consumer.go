package generatevideo

import (
	"encoding/binary"
	"github.com/Pr3c10us/absolutego/internals/domains/queueport"
)

func Handler(c *queueport.Context) (*queueport.HandlerResult, error) {
	splitId := int64(binary.BigEndian.Uint64(c.Data))
	scriptId, err := c.Services.ScriptServices.GenerateVideo.Handle(splitId)
	return &queueport.HandlerResult{
		ScriptId: scriptId,
	}, err
}
