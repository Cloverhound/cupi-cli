package client

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	dimeServiceURL = "https://%s:8443/logcollectionservice/services/DimeGetFileService"
	dimeSoapAction = "http://schemas.cisco.com/ast/soap/action/#LogCollectionPort#GetOneFile"
	dimeRecordVer  = 0x01
)

// GetFile downloads a log file from a CUC node via the DIME log collection service.
// filePath must be an absolute path on the CUC node. If it does not begin with /var/log,
// /var/log/active/ is prepended automatically.
// Returns the raw file bytes from the attachment.
func GetFile(host, user, pass, filePath string) ([]byte, error) {
	if !strings.HasPrefix(filePath, "/var/log") {
		filePath = "/var/log/active/" + strings.TrimPrefix(filePath, "/")
	}

	envelope := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/"
                   xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
                   xmlns:xsd="http://www.w3.org/2001/XMLSchema"
                   xmlns:SOAP-ENC="http://schemas.xmlsoap.org/soap/encoding/">
  <SOAP-ENV:Body SOAP-ENV:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
    <ns:GetOneFile xmlns:ns="http://schemas.cisco.com/ast/soap/">
      <FileName xsi:type="xsd:string">%s</FileName>
    </ns:GetOneFile>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`, filePath)

	svcURL := fmt.Sprintf(dimeServiceURL, host)

	if os.Getenv("CUPI_DEBUG") != "" {
		fmt.Fprintf(os.Stderr, "=== DIME Request ===\nURL: %s\n%s\n", svcURL, envelope)
	}

	req, err := http.NewRequest("POST", svcURL, bytes.NewReader([]byte(envelope)))
	if err != nil {
		return nil, fmt.Errorf("failed to create DIME request: %w", err)
	}

	req.SetBasicAuth(user, pass)
	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	req.Header.Set("SOAPAction", dimeSoapAction)

	httpClient := NewHTTPClient()
	httpClient.Timeout = 120 * time.Second

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("DIME request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read DIME response: %w", err)
	}

	ct := resp.Header.Get("Content-Type")
	if os.Getenv("CUPI_DEBUG") != "" {
		fmt.Fprintf(os.Stderr, "=== DIME Response HTTP %d CT=%q size=%d ===\n", resp.StatusCode, ct, len(body))
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("DIME error (HTTP %d): %s", resp.StatusCode, dimeExtractFault(body))
	}

	mediaType, params, _ := mime.ParseMediaType(ct)

	switch {
	case strings.HasPrefix(mediaType, "multipart/"):
		boundary := params["boundary"]
		if boundary == "" {
			return nil, fmt.Errorf("DIME: multipart/related response missing boundary parameter in Content-Type: %s", ct)
		}
		return dimeParseMultipart(body, boundary)

	case strings.Contains(ct, "application/dime"):
		return dimeParseAttachment(body)

	default:
		return nil, fmt.Errorf("DIME: unexpected response Content-Type %q (HTTP %d)", ct, resp.StatusCode)
	}
}

// dimeParseMultipart extracts the file payload from a MIME multipart/related body.
func dimeParseMultipart(body []byte, boundary string) ([]byte, error) {
	mr := multipart.NewReader(bytes.NewReader(body), boundary)

	partIndex := 0
	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("DIME: failed to read multipart section %d: %w", partIndex+1, err)
		}

		partIndex++
		if partIndex == 2 {
			data, err := io.ReadAll(part)
			if err != nil {
				return nil, fmt.Errorf("DIME: failed to read file attachment: %w", err)
			}
			return data, nil
		}
		if _, err := io.Copy(io.Discard, part); err != nil {
			return nil, fmt.Errorf("DIME: failed to skip SOAP part: %w", err)
		}
	}

	return nil, fmt.Errorf("DIME: no file attachment found in multipart response (got %d parts)", partIndex)
}

// dimeParseAttachment parses a raw DIME binary body and extracts the file payload.
func dimeParseAttachment(data []byte) ([]byte, error) {
	r := bytes.NewReader(data)
	recordIndex := 0
	var fileChunks []byte
	collecting := false

	for {
		var hdr [12]byte
		if _, err := io.ReadFull(r, hdr[:]); err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				break
			}
			return nil, fmt.Errorf("DIME: failed to read record header: %w", err)
		}

		version := (hdr[0] >> 3) & 0x1F
		if version != dimeRecordVer {
			return nil, fmt.Errorf("DIME: unexpected record version %d (expected 1)", version)
		}

		me := (hdr[0] >> 1) & 0x01
		cf := hdr[0] & 0x01

		optionsLen := int(binary.BigEndian.Uint16(hdr[2:4]))
		idLen := int(binary.BigEndian.Uint16(hdr[4:6]))
		typeLen := int(binary.BigEndian.Uint16(hdr[6:8]))
		dataLen := int(binary.BigEndian.Uint32(hdr[8:12]))

		skipLen := dimePadTo4(optionsLen) + dimePadTo4(idLen) + dimePadTo4(typeLen)
		if skipLen > 0 {
			if _, err := io.ReadFull(r, make([]byte, skipLen)); err != nil {
				return nil, fmt.Errorf("DIME: failed to skip record metadata: %w", err)
			}
		}

		paddedLen := dimePadTo4(dataLen)
		payload := make([]byte, paddedLen)
		if paddedLen > 0 {
			if _, err := io.ReadFull(r, payload); err != nil {
				return nil, fmt.Errorf("DIME: failed to read record data: %w", err)
			}
		}
		payload = payload[:dataLen]

		if recordIndex >= 1 || collecting {
			collecting = true
			fileChunks = append(fileChunks, payload...)
			if cf == 0 {
				return fileChunks, nil
			}
		}

		recordIndex++

		if me == 1 {
			break
		}
	}

	if collecting {
		return fileChunks, nil
	}
	return nil, fmt.Errorf("DIME: no file attachment found in response")
}

// dimePadTo4 returns n rounded up to the nearest 4-byte boundary.
func dimePadTo4(n int) int {
	if n%4 == 0 {
		return n
	}
	return n + (4 - n%4)
}

// dimeExtractFault extracts the <faultstring> from a SOAP fault XML body.
func dimeExtractFault(body []byte) string {
	s := string(body)
	start := strings.Index(s, "<faultstring>")
	if start >= 0 {
		end := strings.Index(s[start:], "</faultstring>")
		if end >= 0 {
			return s[start+len("<faultstring>") : start+end]
		}
	}
	if len(s) > 500 {
		return s[:500] + "..."
	}
	return s
}
