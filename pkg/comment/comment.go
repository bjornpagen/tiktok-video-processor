package comment

import (
	"fmt"
	"strings"

	"github.com/bjornpagen/goplay/pkg/chrome"
)

type CommentData struct {
	Username     string
	Comment      string
	ProfileImage string
}

func NewCommentData(username, comment string) *CommentData {
	return &CommentData{
		Username: username,
		Comment:  comment,
	}
}

type CommentBuilder struct {
	Chrome *chrome.Browser
}

func NewCommentBuilder() *CommentBuilder {
	return &CommentBuilder{}
}

func (cb *CommentBuilder) Start() error {
	c, err := chrome.New()
	if err != nil {
		return err
	}
	cb.Chrome = c

	err = cb.Chrome.Start()
	if err != nil {
		return err
	}

	err = cb.Chrome.Navigate("https://tokcomment.com")
	if err != nil {
		return err
	}

	return nil
}

func (cb *CommentBuilder) UpdateComment(cd *CommentData) error {
	updateUsername := fmt.Sprintf(`document.getElementById("resultName").innerHTML = "%s"`, cd.Username)
	updateComment := fmt.Sprintf(`document.getElementById("resultComment").innerHTML = "%s"`, cd.Comment)
	_, err := cb.Chrome.Evaluate(strings.Join([]string{updateUsername, updateComment}, ";"))
	if err != nil {
		return err
	}

	return nil
}

func (cb *CommentBuilder) DownloadComment() error {
	_, err := cb.Chrome.Evaluate("onDownloadClick()")
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
