package addchapter

import (
	"encoding/json"
	"github.com/Pr3c10us/absolutego/internals/domains/queueport"
	"github.com/Pr3c10us/absolutego/internals/services/book/commands"
)

func Handler(c *queueport.Context) (*queueport.HandlerResult, error) {
	var data commands.AddChapterParameter
	err := json.Unmarshal(c.Data, &data)
	if err != nil {
		return nil, err
	}
	id, err := c.Services.BookServices.AddChapter.Handle(data)
	return &queueport.HandlerResult{
		ChapterId: id,
	}, err
}
