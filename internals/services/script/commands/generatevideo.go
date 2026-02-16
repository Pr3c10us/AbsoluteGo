package commands

import (
	"errors"
	"github.com/Pr3c10us/absolutego/internals/domains/book"
	"github.com/Pr3c10us/absolutego/internals/domains/script"
	"github.com/Pr3c10us/absolutego/internals/domains/storage"
	"github.com/Pr3c10us/absolutego/packages/appError"
	"github.com/Pr3c10us/absolutego/packages/configs"
	"github.com/Pr3c10us/absolutego/packages/utils"
	"os"
	"path/filepath"
)

type GenerateVideo struct {
	bookImplementation    book.Interface
	storageImplementation storage.Interface
	environmentVariable   *configs.EnvironmentVariables
	scriptImplementation  script.Interface
}

func (service *GenerateVideo) Handle(id int64) (int64, error) {
	split, err := service.scriptImplementation.GetSplit(id)
	if err != nil {
		return 0, err
	}
	if split == nil {
		return 0, errors.New("split does not exist")
	}
	if split.AudioURL == nil || split.AudioDuration == nil {
		return 0, errors.New("generate split audio first")
	}

	scr, err := service.scriptImplementation.GetScript(split.ScriptId)
	if err != nil {
		return 0, err
	}
	if scr == nil {
		return 0, appError.BadRequest(errors.New("script does not exist"))
	}

	chapters, err := service.bookImplementation.GetChapters(scr.BookId, scr.Chapters)
	if err != nil {
		return 0, err
	}
	if len(chapters) < 1 {
		return 0, appError.BadRequest(errors.New("script does not exist"))
	}

	file, err := utils.DownloadPage(*split.Panel.URL)
	if err != nil {
		return 0, err
	}
	defer os.Remove(file)

	audio, err := utils.DownloadPage(*split.AudioURL)
	if err != nil {
		return 0, err
	}
	defer os.Remove(audio)

	blurImage, err := utils.DownloadPage(chapters[0].BlurURL)
	if err != nil {
		return 0, err
	}
	defer os.Remove(blurImage)

	vidData := utils.VideoData{
		Panel:    file,
		Duration: *split.AudioDuration,
		Effect:   *split.Effect,
	}

	videoDir, err := utils.GetDirectory("tmp")
	if err != nil {
		return 0, err
	}
	defer os.RemoveAll(videoDir)

	slideshowPath := filepath.Join(videoDir, "slideshow.mp4")

	err = utils.CreateVideoFromImages([]utils.VideoData{
		vidData,
	}, slideshowPath, &utils.CreateVideoOptions{
		FPS:             30,
		Width:           1080,
		Height:          1920,
		BackgroundImage: blurImage,
		HWAccel:         service.environmentVariable.HardwareAccelerator,
	})
	if err != nil {
		return 0, err
	}

	videoPath := filepath.Join(videoDir, "video.mp4")
	err = utils.MergeAudioToVideo(slideshowPath, audio, videoPath, &utils.MergeAudioOptions{
		AudioFade: true,
		Loop:      true,
		Volume:    1,
	})
	if err != nil {
		return 0, err
	}

	osFile, err := os.Open(videoPath)
	if err != nil {
		return 0, err
	}
	defer osFile.Close()

	url, err := service.storageImplementation.UploadFile(service.environmentVariable.Buckets.VideoBucket, osFile)
	if err != nil {
		return 0, err
	}

	err = service.scriptImplementation.UpdateSplit(split.Id, &script.Split{
		VideoURL: &url,
	})
	if err != nil {
		return 0, err
	}

	return scr.Id, nil
}

func NewGenerateVideo(
	bookImplementation book.Interface,
	storageImplementation storage.Interface,
	environmentVariable *configs.EnvironmentVariables,
	scriptImplementation script.Interface,
) *GenerateVideo {
	return &GenerateVideo{
		bookImplementation,
		storageImplementation,
		environmentVariable,
		scriptImplementation,
	}
}
