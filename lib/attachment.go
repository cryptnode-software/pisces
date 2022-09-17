package lib

type AttachmentType string

const (
	AttachmentTypeNotImplemented AttachmentType = "NOT_IMPLEMENTED"
	AttachmentTypeImage          AttachmentType = "IMAGE"
	AttachmentTypeFile           AttachmentType = "FILE"
)

type Attachment struct {
	Type AttachmentType
	URL  string
	Model
}
