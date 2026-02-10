package ai

type Interface interface {
	UploadFiles(files []File) ([]UploadedFile, error)
	GenerateText(prompt string, useFastModel bool, uploadedFiles []UploadedFile) (*Response, error)
	GenerateAudioLive(text string, voice Voice) (*Response, error)
}
