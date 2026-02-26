package client

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
)

const pawsBaseURL = "https://%s/platform-services/services/%s"

// wsaAnonymous is the WS-Addressing anonymous reply endpoint for synchronous requests.
const wsaAnonymous = "http://www.w3.org/2005/08/addressing/anonymous"

// ClusterNodeStatus holds PAWS cluster node information.
type ClusterNodeStatus struct {
	Hostname string
	Address  string
	Status   string
	DBRole   string
	Role     int
	Type     string
}

// ReplicationStatus holds PAWS cluster replication health information.
type ReplicationStatus struct {
	OK bool
}

// DRSResult holds the result of initiating a DRS backup.
type DRSResult struct {
	Result  string
	Message string
}

// DRSStatus holds the current DRS backup/restore operation status.
type DRSStatus struct {
	Status  string
	Message string
}

// pawsCall sends a SOAP 1.1 request to a PAWS service endpoint and returns the raw body.
func pawsCall(host, user, pass, service, operation, bodyXML string) ([]byte, error) {
	svcURL := fmt.Sprintf(pawsBaseURL, host, service)
	msgID := "uuid:" + uuid.New().String()
	wsaAction := "urn:" + operation

	envelope := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/"
                  xmlns:vos="http://services.api.platform.vos.cisco.com"
                  xmlns:wsa="http://www.w3.org/2005/08/addressing">
  <soapenv:Header>
    <wsa:To>%s</wsa:To>
    <wsa:ReplyTo>
      <wsa:Address>%s</wsa:Address>
    </wsa:ReplyTo>
    <wsa:Action>%s</wsa:Action>
    <wsa:MessageID>%s</wsa:MessageID>
  </soapenv:Header>
  <soapenv:Body>
    %s
  </soapenv:Body>
</soapenv:Envelope>`, svcURL, wsaAnonymous, wsaAction, msgID, bodyXML)

	if os.Getenv("CUPI_DEBUG") != "" {
		fmt.Fprintf(os.Stderr, "=== PAWS Request [%s/%s] ===\n%s\n", service, operation, envelope)
	}

	req, err := http.NewRequest("POST", svcURL, bytes.NewReader([]byte(envelope)))
	if err != nil {
		return nil, fmt.Errorf("failed to create PAWS request: %w", err)
	}

	req.SetBasicAuth(user, pass)
	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	req.Header.Set("SOAPAction", "")

	httpClient := NewHTTPClient()
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("PAWS call failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read PAWS response: %w", err)
	}

	if os.Getenv("CUPI_DEBUG") != "" {
		fmt.Fprintf(os.Stderr, "=== PAWS Response [%s/%s] HTTP %d ===\n%s\n", service, operation, resp.StatusCode, string(body))
	}

	if resp.StatusCode != http.StatusOK {
		return body, fmt.Errorf("PAWS HTTP %d: %s", resp.StatusCode, pawsTruncate(string(body), 500))
	}

	return body, nil
}

// GetClusterStatus returns node status for all cluster nodes via PAWS ClusterNodesService.
func GetClusterStatus(host, user, pass string) ([]ClusterNodeStatus, error) {
	bodyXML := `<vos:getClusterStatus/>`

	respBody, err := pawsCall(host, user, pass, "ClusterNodesService", "getClusterStatus", bodyXML)
	if err != nil {
		return nil, err
	}

	var response struct {
		Nodes []struct {
			Hostname string `xml:"hostname"`
			Address  string `xml:"address"`
			Status   string `xml:"status"`
			DBRole   string `xml:"dbRole"`
			Role     int    `xml:"role"`
			Type     string `xml:"type"`
		} `xml:"Body>getClusterStatusResponse>return>clusterNodeStatus"`
	}

	if err := xml.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to parse PAWS cluster status response: %w", err)
	}

	var nodes []ClusterNodeStatus
	for _, n := range response.Nodes {
		nodes = append(nodes, ClusterNodeStatus{
			Hostname: n.Hostname,
			Address:  n.Address,
			Status:   n.Status,
			DBRole:   n.DBRole,
			Role:     n.Role,
			Type:     n.Type,
		})
	}

	return nodes, nil
}

// GetReplicationStatus returns cluster replication health via PAWS ClusterNodesService.
func GetReplicationStatus(host, user, pass string) (*ReplicationStatus, error) {
	bodyXML := `<vos:isClusterReplicationOK/>`

	respBody, err := pawsCall(host, user, pass, "ClusterNodesService", "isClusterReplicationOK", bodyXML)
	if err != nil {
		return nil, err
	}

	var response struct {
		OK string `xml:"Body>isClusterReplicationOKResponse>return>clusterReplicationStatusOK"`
	}

	if err := xml.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to parse PAWS replication status response: %w", err)
	}

	return &ReplicationStatus{
		OK: response.OK == "true",
	}, nil
}

// InitiateDRSBackup starts a DRS backup to an SFTP server via PAWS DataExportService.
func InitiateDRSBackup(host, user, pass, sftpServer string, sftpPort int, sftpUser, sftpPass, sftpDir string) (*DRSResult, error) {
	bodyXML := fmt.Sprintf(`<vos:dataExport>
      <vos:args0>%s</vos:args0>
      <vos:args1>%d</vos:args1>
      <vos:args2>%s</vos:args2>
      <vos:args3>%s</vos:args3>
      <vos:args4>%s</vos:args4>
    </vos:dataExport>`,
		escapeXML(sftpServer),
		sftpPort,
		escapeXML(sftpUser),
		escapeXML(sftpPass),
		escapeXML(sftpDir),
	)

	respBody, err := pawsCall(host, user, pass, "DataExportService", "dataExport", bodyXML)
	if err != nil {
		return nil, err
	}

	var response struct {
		Return struct {
			Result  string `xml:"result"`
			Message string `xml:"dataExportResult"`
		} `xml:"Body>dataExportResponse>return"`
	}

	if err := xml.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to parse PAWS DRS backup response: %w", err)
	}

	return &DRSResult{
		Result:  response.Return.Result,
		Message: response.Return.Message,
	}, nil
}

// GetDRSStatus returns the current DRS backup/restore operation status via PAWS DataExportStatusService.
func GetDRSStatus(host, user, pass string) (*DRSStatus, error) {
	bodyXML := `<vos:dataExportStatus/>`

	respBody, err := pawsCall(host, user, pass, "DataExportStatusService", "dataExportStatus", bodyXML)
	if err != nil {
		return nil, err
	}

	var response struct {
		Return struct {
			Status  string `xml:"dataExportStatus"`
			Message string `xml:"result"`
		} `xml:"Body>dataExportStatusResponse>return"`
	}

	if err := xml.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to parse PAWS DRS status response: %w", err)
	}

	return &DRSStatus{
		Status:  response.Return.Status,
		Message: response.Return.Message,
	}, nil
}

func pawsTruncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}

func escapeXML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&apos;")
	return s
}
