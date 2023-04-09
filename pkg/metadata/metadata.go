package metadata

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Metadata struct {
	TrackCreateDate    string
	TrackModifyDate    string
	ImageWidth         int64
	ImageHeight        int64
	XResolution        int64
	YResolution        int64
	VideoFrameRate     float64
	MatrixStructure    string
	MediaCreateDate    string
	MediaModifyDate    string
	Balance            int64
	AudioBitsPerSample int64
	AudioSampleRate    int64
	Copyright          string
	MinorVersion       int64
	Encoder            string
	Software           string
	Artwork            string
	Rotation           int64
}

/* example data
Track Create Date               : 2023:02:26 04:41:35
Track Modify Date               : 2023:02:26 04:41:35
Image Width                     : 1080
Image Height                    : 1920
X Resolution                    : 72
Y Resolution                    : 72
Video Frame Rate                : 30
Matrix Structure                : 1 0 0 0 1 0 0 0 1
Media Create Date               : 2023:02:26 04:41:35
Media Modify Date               : 2023:02:26 04:41:35
Balance                         : 0
Audio Bits Per Sample           : 16
Audio Sample Rate               : 44100
Copyright                       : 63089a1ab188837503e83ca2d4d0b1a3
Minor Version                   : 512
Encoder                         : Lavf57.71.100
Software                        : {{"publicMode":"1"},{"TEEditor":"2"},{"isFastImport":"0"},{"transType":"2"},{"te_is_reencode":"1"},{"source":""}}
Artwork                         : {"source_type":"vicut","data":{"adsTemplateId":"","appVersion":"7.7.2","capabilityName":"subtitle_recognition","ccTtVid":"","draftInfo":{"adjust":0,"aiMatting":0,"audioToText":69,"chroma":0,"coverTemplateId":"","curveSpeed":0,"faceEffectId":"","filterId":"","gameplayAlgorithm":"","graphs":0,"keyframe":0,"mask":0,"mixMode":0,"motion_blur_cnt":0,"normalSpeed":0,"pip":0,"reverse":0,"slowMotion":0,"soundId":"","stickerId":"","textAnimationId":"","textEffectId":"","textShapeId":"","textTemplateId":"","textToAudio":0,"transitionId":"","videoAnimationId":"","videoEffectId":"","videoMaterialId":"","videoTracking":0},"editSource":"album","editType":"edit","enterFrom":"","exportType":"export","filterId":"","infoStickerId":"","motion_blur_cnt":0,"musicId":"","os":"ios","product":"vicut","provider":"ad_site","region":"VN","resourceTypeApplied":"","slowMotion":"none","stickerId":"","templateId":"","textSpecialEffect":"","transferMethod":"","transitions":"","videoAnimation":"","videoEffectId":"","videoId":"3D322E87-4108-438D-88D9-2A99E1F9E886"}}
Rotation                        : 0
*/

func NewMetadata(videoPath string) (*Metadata, error) {
	rand.Seed(time.Now().UnixNano())

	// Get the creation time of the video file
	videoInfo, err := os.Stat(videoPath)
	if err != nil {
		fmt.Println("Error getting file information:", err)
		return nil, err
	}
	creationTime := videoInfo.ModTime()

	// Convert the creation time to the correct format
	time := ConvertTimeToExifTime(creationTime)

	// Generate a random video ID
	videoID := fmt.Sprintf("%X-%X-%X-%X-%X", rand.Int31(), rand.Int31(), rand.Int31(), rand.Int31(), rand.Int31())

	// Run ffprobe to get metadata
	cmd := exec.Command("ffprobe", "-v", "quiet", "-print_format", "json", "-show_format", "-show_streams", videoPath)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error running ffprobe:", err)
		return nil, err
	}

	var probeOutput map[string]interface{}
	err = json.Unmarshal(output, &probeOutput)
	if err != nil {
		fmt.Println("Error unmarshalling ffprobe output:", err)
		return nil, err
	}

	streams, ok := probeOutput["streams"].([]interface{})
	if !ok {
		fmt.Println("Error: cannot assert probeOutput[\"streams\"] to []interface{}")
		return nil, errors.New("invalid ffprobe output")
	}
	var videoStream, audioStream map[string]interface{}
	for _, stream := range streams {
		streamMap, ok := stream.(map[string]interface{})
		if !ok {
			fmt.Println("Error: cannot assert stream to map[string]interface{}")
			return nil, errors.New("invalid ffprobe output")
		}
		if streamMap["codec_type"].(string) == "video" {
			videoStream = streamMap
		} else if streamMap["codec_type"].(string) == "audio" {
			audioStream = streamMap
		}
	}

	// Extract metadata
	imageWidth := int64(videoStream["width"].(float64))
	imageHeight := int64(videoStream["height"].(float64))
	displayAspectRatio := videoStream["display_aspect_ratio"].(string)
	split := strings.Split(displayAspectRatio, ":")
	numerator, err := strconv.Atoi(split[0])
	if err != nil {
		fmt.Println("Error converting numerator:", err)
		return nil, err
	}
	denominator, err := strconv.Atoi(split[1])
	if err != nil {
		fmt.Println("Error converting denominator:", err)
		return nil, err
	}
	// TODO: this is absolutely not how you calculate the x and y resolution
	xResolution := int64(numerator)
	yResolution := int64(denominator)

	frameRate := videoStream["avg_frame_rate"].(string)
	split = strings.Split(frameRate, "/")
	numerator, err = strconv.Atoi(split[0])
	if err != nil {
		fmt.Println("Error converting numerator:", err)
		return nil, err
	}
	denominator, err = strconv.Atoi(split[1])
	if err != nil {
		fmt.Println("Error converting denominator:", err)
		return nil, err
	}
	videoFrameRate := float64(numerator) / float64(denominator)

	audioBitsPerSample := int64(audioStream["bits_per_sample"].(float64))
	audioSampleRateStr := audioStream["sample_rate"].(string)
	audioSampleRate, err := strconv.ParseInt(audioSampleRateStr, 10, 64)
	if err != nil {
		fmt.Println("Error converting audioSampleRate:", err)
		return nil, err
	}

	// Return the metadata
	return &Metadata{
		TrackCreateDate:    time,
		TrackModifyDate:    time,
		ImageWidth:         imageWidth,
		ImageHeight:        imageHeight,
		XResolution:        xResolution,
		YResolution:        yResolution,
		VideoFrameRate:     videoFrameRate,
		MatrixStructure:    "1 0 0 0 1 0 0 0 1",
		MediaCreateDate:    time,
		MediaModifyDate:    time,
		Balance:            0,
		AudioBitsPerSample: audioBitsPerSample,
		AudioSampleRate:    audioSampleRate,
		Copyright:          "63089a1ab188837503e83ca2d4d0b1a3",
		MinorVersion:       512,
		Encoder:            "Lavf58.76.100",
		Software:           `{{"publicMode":"1"},{"TEEditor":"2"},{"isFastImport":"0"},{"transType":"2"},{"te_is_reencode":"1"},{"source":""}}`,
		Artwork:            fmt.Sprintf(`{"source_type":"vicut","data":{"adsTemplateId":"","appVersion":"7.7.2","capabilityName":"subtitle_recognition","ccTtVid":"","draftInfo":{"adjust":0,"aiMatting":0,"audioToText":0,"chroma":0,"coverTemplateId":"","curveSpeed":0,"faceEffectId":"","filterId":"","gameplayAlgorithm":"","graphs":0,"keyframe":0,"mask":0,"mixMode":0,"motion_blur_cnt":0,"normalSpeed":0,"pip":0,"reverse":0,"slowMotion":0,"soundId":"","stickerId":"","textAnimationId":"","textEffectId":"","textShapeId":"","textTemplateId":"","textToAudio":0,"transitionId":"","videoAnimationId":"","videoEffectId":"","videoMaterialId":"","videoTracking":0},"editSource":"album","editType":"edit","enterFrom":"","exportType":"export","filterId":"","infoStickerId":"","motion_blur_cnt":0,"musicId":"","os":"ios","product":"vicut","provider":"ad_site","region":"US","resourceTypeApplied":"","slowMotion":"none","stickerId":"","templateId":"","textSpecialEffect":"","transferMethod":"","transitions":"","videoAnimation":"","videoEffectId":"","videoId":"%s"}}`, videoID),
		Rotation:           0,
	}, nil
}

// ConvertTimeToExifTime converts a time.Time to an "YYYY:mm:dd HH:MM:SS[.ss][+/-HH:MM|Z]" formatted string
func ConvertTimeToExifTime(t time.Time) string {
	return t.Format("2006:01:02 15:04:05.00-07:00")
}

func WriteMetadataToFile(metadata *Metadata, path string) error {
	tmp := ".__temp_data.mp4"

	// Copy the video to a temporary file
	copyCmd := exec.Command("cp", path, tmp)
	err := copyCmd.Run()
	if err != nil {
		return fmt.Errorf("error copying video to temporary file: %w", err)
	}

	// Remove existing metadata, in place
	removeMetadataCmd := exec.Command("exiftool", "-all=", tmp)
	err = removeMetadataCmd.Run()
	if err != nil {
		return fmt.Errorf("error removing metadata: %w", err)
	}

	// Add new metadata to the video, in place
	addMetadataCmd := exec.Command("exiftool",
		"-overwrite_original",
		fmt.Sprintf("-TrackCreateDate=%s", metadata.TrackCreateDate),
		fmt.Sprintf("-TrackModifyDate=%s", metadata.TrackModifyDate),
		fmt.Sprintf("-ImageWidth=%d", metadata.ImageWidth),
		fmt.Sprintf("-ImageHeight=%d", metadata.ImageHeight),
		fmt.Sprintf("-XResolution=%d", metadata.XResolution),
		fmt.Sprintf("-YResolution=%d", metadata.YResolution),
		fmt.Sprintf("-VideoFrameRate=%.2f", metadata.VideoFrameRate),
		fmt.Sprintf("-MatrixStructure=%s", metadata.MatrixStructure),
		fmt.Sprintf("-MediaCreateDate=%s", metadata.MediaCreateDate),
		fmt.Sprintf("-MediaModifyDate=%s", metadata.MediaModifyDate),
		fmt.Sprintf("-Balance=%d", metadata.Balance),
		fmt.Sprintf("-AudioBitsPerSample=%d", metadata.AudioBitsPerSample),
		fmt.Sprintf("-AudioSampleRate=%d", metadata.AudioSampleRate),
		fmt.Sprintf("-Copyright=%s", metadata.Copyright),
		fmt.Sprintf("-MinorVersion=%d", metadata.MinorVersion),
		fmt.Sprintf("-Encoder=%s", metadata.Encoder),
		fmt.Sprintf("-Software=%s", metadata.Software),
		fmt.Sprintf("-Artwork=%s", metadata.Artwork),
		fmt.Sprintf("-Rotation=%d", metadata.Rotation),
		tmp)
	err = addMetadataCmd.Run()
	if err != nil {
		return fmt.Errorf("error adding metadata: %w", err)
	}

	err = os.Rename(tmp, path)
	if err != nil {
		return fmt.Errorf("error renaming temporary output file: %w", err)
	}

	// Clean up temporary files
	os.Remove(tmp)

	return nil
}

func GenerateMetadataAndWriteToFile(path string) error {
	metadata, err := NewMetadata(path)
	if err != nil {
		return fmt.Errorf("error getting metadata: %w", err)
	}

	err = ReencodeVideo(path)
	if err != nil {
		return fmt.Errorf("error reencoding video: %w", err)
	}

	err = WriteMetadataToFile(metadata, path)
	if err != nil {
		return fmt.Errorf("error writing metadata to file: %w", err)
	}

	return nil
}

func ReencodeVideo(path string) error {
	// ffmpeg -i input.mp4 -c:v libx265 -vtag hvc1 -c:a aac -crf 0 -b:v 16M -maxrate 0 -bufsize 16M output.mp4

	// Create a temporary file
	tmp := ".__temp_data.mp4"

	// Reencode the video
	reencodeCmd := exec.Command("ffmpeg",
		"-i", path,
		"-c:v", "libx265",
		"-vtag", "hvc1",
		"-c:a", "aac",
		"-crf", "0",
		"-b:v", "16M",
		"-maxrate", "0",
		"-bufsize", "16M",
		tmp)
	err := reencodeCmd.Run()
	if err != nil {
		return fmt.Errorf("error reencoding video: %w", err)
	}

	// Copy the video to the original file
	copyCmd := exec.Command("cp", tmp, path)
	err = copyCmd.Run()
	if err != nil {
		return fmt.Errorf("error copying video to original file: %w", err)
	}

	// Clean up temporary files
	os.Remove(tmp)

	return nil
}
