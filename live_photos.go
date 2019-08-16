package onfido

import (
	"encoding/json"
	"time"
)

// LivePhoto represents a LivePhoto in Onfido API
type LivePhoto struct {
	ID           string     `json:"id,omitempty"`
	CreatedAt    *time.Time `json:"created_at,omitempty"`
	Href         string     `json:"href,omitempty"`
	DownloadHref string     `json:"download_href,omitempty"`
	FileName     string     `json:"file_name,omitempty"`
	FileType     string     `json:"file_type,omitempty"`
	FileSize     int32      `json:"file_size,omitempty"`
}

// LivePhotoIter represents a LivePhoto iterator
type LivePhotoIter struct {
	*iter
}

// LivePhoto returns the current item in the iterator as a LivePhoto.
func (i *LivePhotoIter) LivePhoto() *LivePhoto {
	return i.Current().(*LivePhoto)
}

// ListPhotos retrieves the list of photos for the provided applicant.
// see https://documentation.onfido.com/?shell#live-photos
func (c *Client) ListLivePhotos(applicantID string) *LivePhotoIter {
	return &LivePhotoIter{&iter{
		c:       c,
		nextURL: "/live_photos?applicant_id=" + applicantID,
		handler: func(body []byte) ([]interface{}, error) {
			var r struct {
				LivePhotos []*LivePhoto `json:"live_photos"`
			}

			if err := json.Unmarshal(body, &r); err != nil {
				return nil, err
			}

			values := make([]interface{}, len(r.LivePhotos))
			for i, v := range r.LivePhotos {
				values[i] = v
			}
			return values, nil
		},
	}}
}
