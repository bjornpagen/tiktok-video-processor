package schemas

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

	tiktokdb.UserAwemesStartAwemeIdsVector(builder, len(awemeIds))
	for i := len(awemeIds) - 1; i >= 0; i-- {
		builder.PrependUOffsetT(videoIds[i])
	}
	videoIdsVector := builder.EndVector(len(awemeIds))

	tiktokdb.UserAwemesStart(builder)
	tiktokdb.UserAwemesAddAwemeIds(builder, videoIdsVector)
	userAwemes := tiktokdb.UserAwemesEnd(builder)
	builder.Finish(userAwemes)

	return tiktokdb.GetRootAsUserAwemes(builder.FinishedBytes(), 0)
}

func NewUser(latestUsername string, latestMincursor string, awemes *tiktokdb.UserAwemes) *tiktokdb.User {
	builder := flatbuffers.NewBuilder(0)

	username := builder.CreateString(latestUsername)
	mincursor := builder.CreateString(latestMincursor)

	awemesBuilder := flatbuffers.NewBuilder(0)
	awemesBytes := awemes.Table().Bytes
	awemesBuilder.StartObject(1)
	awemesBuilder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(awemesBuilder.CreateByteVector(awemesBytes)), 0)
	userAwemes := awemesBuilder.EndObject()
	awemesBuilder.Finish(userAwemes)

	builder.StartObject(3)
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(username), 0)
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(mincursor), 0)
	builder.PrependUOffsetTSlot(2, flatbuffers.UOffsetT(builder.CreateByteVector(awemesBuilder.FinishedBytes())), 0)
	user := builder.EndObject()

	builder.Finish(user)

	return tiktokdb.GetRootAsUser(builder.FinishedBytes(), 0)
}
