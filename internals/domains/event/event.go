package event

import "time"

type Status string

const (
	StatusEnqueue    Status = "enqueue"
	StatusProcessing Status = "processing"
	StatusFailed     Status = "failed"
	StatusSuccessful Status = "successful"
	StatusRetry      Status = "retry"
)

func (s Status) IsValid() bool {
	switch s {
	case StatusEnqueue, StatusProcessing, StatusFailed, StatusSuccessful:
		return true
	}
	return false
}

type Operation string

const (
	OpAddChapter     Operation = "add_chapter"
	OpGenScript      Operation = "gen_script"
	OpGenScriptSplit Operation = "gen_script_split"
	OpGenAudio       Operation = "gen_audio"
	OpGenVideo       Operation = "gen_video"
	OpMergeVideo     Operation = "merge_video"
)

func (o Operation) IsValid() bool {
	switch o {
	case OpAddChapter, OpGenScript, OpGenScriptSplit, OpGenAudio, OpGenVideo, OpMergeVideo:
		return true
	}
	return false
}

type Event struct {
	Id          int64
	Status      Status
	Operation   Operation
	Description string
	BookId      int64
	ChapterId   int64
	ScriptId    int64
	VabId       int64
	UpdatedAt   time.Time
}

type Filter struct {
	Page      int       `form:"page" binding:"omitempty,min=1"`
	Limit     int       `form:"limit" binding:"omitempty,min=1,max=100"`
	Status    Status    `form:"status" binding:"omitempty,min=1"`
	Operation Operation `form:"operation" binding:"omitempty,min=1"`
}
