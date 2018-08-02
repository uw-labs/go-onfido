package onfido

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"strings"
	"time"
)

// Supported document types and sides
const (
	DocumentTypeUnknown        DocumentType = "unknown"
	DocumentTypePassport       DocumentType = "passport"
	DocumentTypeIDCard         DocumentType = "national_identity_card"
	DocumentTypeDrivingLicense DocumentType = "driving_licence"
	DocumentTypeUKBRP          DocumentType = "uk_biometric_residence_permit"
	DocumentTypeTaxID          DocumentType = "tax_id"
	DocumentTypeVoterID        DocumentType = "voter_id"

	DocumentSideFront DocumentSide = "front"
	DocumentSideBack  DocumentSide = "back"
)

// DocumentType represents a document type (passport, ID, etc)
type DocumentType string

// DocumentSide represents a document side (front, back)
type DocumentSide string

// DocumentRequest represents a document request to Onfido API
type DocumentRequest struct {
	File io.ReadSeeker
	Type DocumentType
	Side DocumentSide
}

// Document represents a document in Onfido API
type Document struct {
	ID           string       `json:"id,omitempty"`
	CreatedAt    *time.Time   `json:"created_at,omitempty"`
	Href         string       `json:"href,omitempty"`
	DownloadHref string       `json:"download_href,omitempty"`
	FileName     string       `json:"file_name,omitempty"`
	FileType     string       `json:"file_type,omitempty"`
	FileSize     int          `json:"file_size,omitempty"`
	Type         DocumentType `json:"type,omitempty"`
	Side         DocumentSide `json:"side,omitempty"`
}

// Documents represents a list of documents from the Onfido API
type Documents struct {
	Documents []*Document `json:"documents"`
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

// createFormFile creates a new form-data header with the provided field name,
// file name, and file content type.
// this is used instead of multipart.Writer.CreateFormFile because Onfido API
// doesn't accept 'application/octet-stream' as content-type.
func createFormFile(writer *multipart.Writer, fieldname string, file io.ReadSeeker) (io.Writer, error) {
	buffer := make([]byte, 512)
	if _, err := file.Read(buffer); err != nil {
		return nil, err
	}
	if _, err := file.Seek(0, 0); err != nil {
		return nil, err
	}
	var filename string
	if f, ok := file.(*os.File); ok {
		filename = f.Name()
	}

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			escapeQuotes(fieldname), escapeQuotes(filename)))
	h.Set("Content-Type", http.DetectContentType(buffer))

	return writer.CreatePart(h)
}

// UploadDocument uploads a document for the provided applicant.
// see https://documentation.onfido.com/?shell#upload-document
func (c *Client) UploadDocument(ctx context.Context, applicantID string, dr DocumentRequest) (*Document, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := createFormFile(writer, "file", dr.File)
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(part, dr.File); err != nil {
		return nil, err
	}
	if err := writer.WriteField("type", string(dr.Type)); err != nil {
		return nil, err
	}
	if err := writer.WriteField("side", string(dr.Side)); err != nil {
		return nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}

	req, err := c.newRequest("POST", "/applicants/"+applicantID+"/documents", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if err != nil {
		return nil, err
	}

	var resp Document
	_, err = c.do(ctx, req, &resp)
	return &resp, err
}

// GetDocument retrieves a single document for the provided applicant by its ID.
// see https://documentation.onfido.com/?shell#retrieve-document
func (c *Client) GetDocument(ctx context.Context, applicantID, id string) (*Document, error) {
	req, err := c.newRequest("GET", "/applicants/"+applicantID+"/documents/"+id, nil)
	if err != nil {
		return nil, err
	}

	var resp Document
	_, err = c.do(ctx, req, &resp)
	return &resp, err
}

// DocumentIter represents a document iterator
type DocumentIter struct {
	*iter
}

// Document returns the current item in the iterator as a Document.
func (i *DocumentIter) Document() *Document {
	return i.Current().(*Document)
}

// ListDocuments retrieves the list of documents for the provided applicant.
// see https://documentation.onfido.com/?shell#list-documents
func (c *Client) ListDocuments(applicantID string) *DocumentIter {
	handler := func(body []byte) ([]interface{}, error) {
		var d Documents
		if err := json.Unmarshal(body, &d); err != nil {
			return nil, err
		}

		values := make([]interface{}, len(d.Documents))
		for i, v := range d.Documents {
			values[i] = v
		}
		return values, nil
	}

	return &DocumentIter{&iter{
		c:       c,
		nextURL: "/applicants/" + applicantID + "/documents",
		handler: handler,
	}}
}
