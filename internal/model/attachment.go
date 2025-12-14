package model

import "time"

type Attachment struct {
    FileName   string    `json:"file_name" bson:"file_name"`
    URL        string    `json:"url" bson:"url"`
    MimeType   string    `json:"mime_type" bson:"mime_type"`
    Size       int64     `json:"size" bson:"size"`
    UploadedAt time.Time `json:"uploaded_at" bson:"uploaded_at"`
}
