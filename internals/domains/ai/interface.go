package ai

type Interface interface {
	UploadFiles(files []File) ([]UploadedFile, error)
	GenerateText(prompt string, useFastModel bool, uploadedFiles []UploadedFile) (*Response, error)
	GenerateAudio(text string, voice Voice) (*Response, error)
}
