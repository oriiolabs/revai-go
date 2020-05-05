package revai

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
)

type CaptionService service

// Caption output for a transcription job
type Caption struct {
	Value string
}

// GetCaptionParams to include in the GetCaptions api call.
type GetCaptionParams struct {
	JobID  string
	Accept string
}

// Get returns the caption output for a transcription job.
// https://www.rev.ai/docs#tag/Captions
func (s *CaptionService) Get(ctx context.Context, params *GetCaptionParams) (*Caption, error) {
	urlPath := "/speechtotext/v1/jobs/" + params.JobID + "/captions"

	accept := params.Accept
	if accept == "" {
		accept = "application/x-subrip"
	}

	acceptHeader := &httpHeader{
		Key:   "Accept",
		Value: accept,
	}

	req, err := s.client.newRequest(http.MethodGet, urlPath, nil, acceptHeader)
	if err != nil {
		return nil, fmt.Errorf("failed creating request %w", err)
	}

	resp, err := s.client.do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, resp.Body); err != nil {
		return nil, err
	}

	caption := &Caption{
		Value: buf.String(),
	}

	return caption, nil
}
