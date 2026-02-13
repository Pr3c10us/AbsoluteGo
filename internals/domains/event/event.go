package event

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
	Id        int64
	Status    Status
	Operation Operation
	ChapterId int64
	ScriptId  int64
	VabId     int64
}

type Filter struct {
	Page      int
	Limit     int
	Status    Status
	Operation Operation
}
