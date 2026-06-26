package meta

import sharedv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/shared/v1"

type Cursor struct {
	NextCursor string `json:"next_cursor,omitempty"`
	PrevCursor string `json:"prev_cursor,omitempty"`
	Limit      int32  `json:"limit,omitempty"`
}

func FromProto(proto *sharedv1.CursorMeta) *Cursor {
	if proto == nil {
		return nil
	}
	return &Cursor{
		NextCursor: proto.GetNextCursor(),
		PrevCursor: proto.GetPrevCursor(),
		Limit:      proto.GetLimit(),
	}
}
