package schema

import (
	"github.com/bjornpagen/tiktok-video-processor/autogen/tiktokdb"
	flatbuffers "github.com/google/flatbuffers/go"
)

func NewAweme(sharelink string) *tiktokdb.Aweme {
	builder := flatbuffers.NewBuilder(0)

	shareLink := builder.CreateString(sharelink)

	tiktokdb.AwemeStart(builder)
	tiktokdb.AwemeAddShareLink(builder, shareLink)
	aweme := tiktokdb.AwemeEnd(builder)
	builder.Finish(aweme)

	return tiktokdb.GetRootAsAweme(builder.FinishedBytes(), 0)
}

func NewUserAwemes(awemeIds []string) *tiktokdb.UserAwemes {
	builder := flatbuffers.NewBuilder(0)

	videoIds := make([]flatbuffers.UOffsetT, len(awemeIds))
	for i, id := range awemeIds {
		videoIds[i] = builder.CreateString(id)
	}

	tiktokdb.UserAwemesStartVideoIdsVector(builder, len(awemeIds))
	for i := len(awemeIds) - 1; i >= 0; i-- {
		builder.PrependUOffsetT(videoIds[i])
	}
	videoIdsVector := builder.EndVector(len(awemeIds))

	tiktokdb.UserAwemesStart(builder)
	tiktokdb.UserAwemesAddVideoIds(builder, videoIdsVector)
	userAwemes := tiktokdb.UserAwemesEnd(builder)
	builder.Finish(userAwemes)

	return tiktokdb.GetRootAsUserAwemes(builder.FinishedBytes(), 0)
}

func NewUser(latestUsername string, latestMincursor string, awemes *tiktokdb.UserAwemes) *tiktokdb.User {
	builder := flatbuffers.NewBuilder(0)

	username := builder.CreateString(latestUsername)
	mincursor := builder.CreateString(latestMincursor)
	awemesBytes := awemes.Table().Bytes

	builder.StartObject(3)
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(username), 0)
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(mincursor), 0)
	builder.PrependUOffsetTSlot(2, flatbuffers.UOffsetT(builder.CreateByteVector(awemesBytes)), 0)
	user := builder.EndObject()

	builder.Finish(user)

	return tiktokdb.GetRootAsUser(builder.FinishedBytes(), 0)
}
