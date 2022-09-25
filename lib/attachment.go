package lib

import commons "github.com/cryptnode-software/commons/pkg"

type AttachmentType string

const (
	AttachmentTypeNotImplemented AttachmentType = "NOT_IMPLEMENTED"
	AttachmentTypeImage          AttachmentType = "IMAGE"
	AttachmentTypeFile           AttachmentType = "FILE"
)

type Attachment struct {
	Type AttachmentType
	URL  string
	commons.Model
}
