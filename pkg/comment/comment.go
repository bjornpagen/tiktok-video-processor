package comment

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"

	"github.com/bjornpagen/goplay/pkg/chrome"
)

type CommentData struct {
	Username  string
	Comment   string
	ImagePath string
}

func NewCommentData(username, comment, profileimagepath string) *CommentData {
	return &CommentData{
		Username:  username,
		Comment:   comment,
		ImagePath: profileimagepath,
	}
}

type CommentBuilder struct {
	c *chrome.Browser
}

func NewCommentBuilder() *CommentBuilder {
	return &CommentBuilder{}
}

func (cb *CommentBuilder) Start() error {
	c, err := chrome.New()
	if err != nil {
		return err
	}
	cb.c = c

	err = cb.c.Start()
	if err != nil {
		return err
	}

	err = cb.c.Navigate("https://tokcomment.com")
	if err != nil {
		return err
	}

	return nil
}

func (cb *CommentBuilder) UpdateComment(cd *CommentData) error {
	err := cb.UpdateUsername(cd.Username)
	if err != nil {
		return err
	}

	err = cb.UpdateText(cd.Comment)
	if err != nil {
		return err
	}

	if cd.ImagePath != "" {
		err = cb.UpdatePicture(cd.ImagePath)
		if err != nil {
			return err
		}
	}

	return nil
}

func (cb *CommentBuilder) UpdateUsername(username string) error {
	upperText := fmt.Sprintf(`Reply to %s's comment`, username)
	updateUsername := fmt.Sprintf(`document.getElementById("resultName").innerHTML = "%s"`, upperText)
	_, err := cb.c.Evaluate(updateUsername)
	if err != nil {
		return err
	}

	return nil
}

func (cb *CommentBuilder) UpdateText(comment string) error {
	updateComment := fmt.Sprintf(`document.getElementById("resultComment").innerHTML = "%s"`, comment)
	_, err := cb.c.Evaluate(updateComment)
	if err != nil {
		return err
	}

	return nil
}

func (cb *CommentBuilder) UpdatePicture(imagePath string) error {
	// Read the image file
	imageFile, err := os.Open(imagePath)
	if err != nil {
		return err
	}
	defer imageFile.Close()

	// Decode the image to detect the format
	img, format, err := image.Decode(imageFile)
	if err != nil {
		return err
	}

	if format != "jpeg" && format != "png" {
		return fmt.Errorf("unsupported image format: %s", format)
	}

	// Encode the image as base64
	buf := new(bytes.Buffer)
	if format == "jpeg" {
		err = jpeg.Encode(buf, img, nil)
	} else {
		err = png.Encode(buf, img)
	}
	if err != nil {
		return err
	}
	base64Image := base64.StdEncoding.EncodeToString(buf.Bytes())

	// Prepare the JavaScript to update the image src attribute
	mimeType := fmt.Sprintf("image/%s", format)
	srcData := fmt.Sprintf("data:%s;base64,%s", mimeType, base64Image)
	updateImageSrc := fmt.Sprintf(`document.getElementById("resultImage").src = "%s"`, srcData)

	// Evaluate the JavaScript
	_, err = cb.c.Evaluate(updateImageSrc)
	if err != nil {
		return err
	}

	return nil
}

func (cb *CommentBuilder) DownloadComment() error {
	_, err := cb.c.Evaluate("onDownloadClick()")
	if err != nil {
		return err
	}

	return nil
}

// func GenerateCommentHTML(data CommentData) (string, error) {
// 	commentHTML := `<div id="fullResult" class="w-fit pb-5">
// 	<div class="relative bg-white w-[200px] pl-[8px] pt-[9px] pr-[6px] pb-[12px] rounded-bl-none h-fit mx-auto rounded-[5px] flex flex-col
// 		after:content-[''] after:rounded-bl-[5px] after:absolute after:bottom-[-9px] after:left-0 after:w-0 after:h-0 after:border-[11px] after:border-t-white after:border-r-transparent after:border-b-0 after:border-l-0">
// 		<p id="resultName" class="text-[#8b8b8b] ml-[31px] text-[10px] font-proximaNovaSemiBold leading-[0.65rem]">
// 			Reply to {{.Name}}'s comment
// 		</p>
// 		<div class="flex flex-row leading-4">
// 			<img id="resultImage" class="w-[25px] aspect-square h-[25px] inline rounded-full mr-[6px]" src="data:image/png;base64,{{.ProfileImage}}" alt="">
// 			<p id="resultComment" class="break-words text-black text-[12px] font-proximaNovaBold mt-[2px]">
// 				{{.Comment}}
// 			</p>
// 		</div>
// 	</div>
// </div>`

// 	tmpl, err := template.New("comment").Parse(commentHTML)
// 	if err != nil {
// 		return "", err
// 	}

// 	var result strings.Builder
// 	err = tmpl.Execute(&result, data)
// 	if err != nil {
// 		return "", err
// 	}

// 	return result.String(), nil
// }
