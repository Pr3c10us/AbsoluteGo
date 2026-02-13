package queue

type Queue string

const (
	QueueAddChapter     Queue = "add_chapter"
	QueueGenScript      Queue = "gen_script"
	QueueGenScriptSplit Queue = "gen_script_split"
	QueueGenAudio       Queue = "gen_audio"
	QueueGenVideo       Queue = "gen_video"
	QueueMergeVideo     Queue = "merge_video"
)

type MessageParams struct {
	Queue   Queue
	Message Message
}

type Message struct {
	EventId int64
	Data    []byte
}
