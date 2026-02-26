package client

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

const astURLPattern = "https://%s/ast/Astisapi.dll"

// astGet makes a GET request to /ast/Astisapi.dll?{query}
func astGet(host, user, pass, query string) ([]byte, error) {
	fullURL := fmt.Sprintf("%s?%s", fmt.Sprintf(astURLPattern, host), query)

	if os.Getenv("CUPI_DEBUG") != "" {
		fmt.Fprintf(os.Stderr, "=== AST GET Request ===\n%s\n", fullURL)
	}

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(user, pass)

	httpClient := NewHTTPClient()
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("AST call failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if os.Getenv("CUPI_DEBUG") != "" {
		fmt.Fprintf(os.Stderr, "=== AST Response ===\n%s\n", string(body))
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("AST HTTP %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// ASTPartition represents disk partition information
type ASTPartition struct {
	Node           string
	Name           string
	PercentageUsed int
	TotalMbytes    int
	UsedMbytes     int
}

// GetASTDiskInfo retrieves disk partition information
func GetASTDiskInfo(host, user, pass string) ([]ASTPartition, error) {
	body, err := astGet(host, user, pass, "GetPreCannedInfo&Items=getPartitionInfoRequest")
	if err != nil {
		return nil, err
	}

	var preCanned struct {
		ReturnCode string `xml:"ReturnCode,attr"`
		Reply      struct {
			ReturnCode string `xml:"ReturnCode,attr"`
			Hosts      []struct {
				Name       string `xml:"Name,attr"`
				ReturnCode string `xml:"ReturnCode,attr"`
				Partitions []struct {
					Name           string `xml:"Name,attr"`
					PercentageUsed string `xml:"PercentageUsed,attr"`
					TotalMbytes    string `xml:"TotalMbytes,attr"`
					UsedMbytes     string `xml:"UsedMbytes,attr"`
				} `xml:"Partition"`
			} `xml:"Host"`
		} `xml:"getPartitionInfoReply"`
	}

	if err := xml.Unmarshal(body, &preCanned); err != nil {
		return nil, fmt.Errorf("failed to parse disk info: %w", err)
	}

	if preCanned.ReturnCode != "0" {
		return nil, fmt.Errorf("AST returned error code: %s", preCanned.ReturnCode)
	}

	if preCanned.Reply.ReturnCode != "0" {
		return nil, fmt.Errorf("disk info request returned error code: %s", preCanned.Reply.ReturnCode)
	}

	var results []ASTPartition
	for _, h := range preCanned.Reply.Hosts {
		for _, partition := range h.Partitions {
			pct, _ := strconv.Atoi(partition.PercentageUsed)
			total, _ := strconv.Atoi(partition.TotalMbytes)
			used, _ := strconv.Atoi(partition.UsedMbytes)

			results = append(results, ASTPartition{
				Node:           h.Name,
				Name:           partition.Name,
				PercentageUsed: pct,
				TotalMbytes:    total,
				UsedMbytes:     used,
			})
		}
	}

	return results, nil
}

// ASTTftpInfo represents TFTP information
type ASTTftpInfo struct {
	Node          string
	TotalRequests int
	Aborted       int
}

// GetASTTftpInfo retrieves TFTP information
func GetASTTftpInfo(host, user, pass string) ([]ASTTftpInfo, error) {
	body, err := astGet(host, user, pass, "GetPreCannedInfo&Items=getTftpInfoRequest")
	if err != nil {
		return nil, err
	}

	var preCanned struct {
		ReturnCode string `xml:"ReturnCode,attr"`
		Reply      struct {
			ReturnCode string `xml:"ReturnCode,attr"`
			Nodes      []struct {
				Name                     string `xml:"Name,attr"`
				TotalTftpRequests        string `xml:"TotalTftpRequests,attr"`
				TotalTftpRequestsAborted string `xml:"TotalTftpRequestsAborted,attr"`
			} `xml:"TftpNode"`
		} `xml:"getTftpInfoReply"`
	}

	if err := xml.Unmarshal(body, &preCanned); err != nil {
		return nil, fmt.Errorf("failed to parse TFTP info: %w", err)
	}

	if preCanned.ReturnCode != "0" {
		return nil, fmt.Errorf("TFTP info request returned error code: %s", preCanned.ReturnCode)
	}

	var results []ASTTftpInfo
	for _, node := range preCanned.Reply.Nodes {
		totalReqs, _ := strconv.Atoi(node.TotalTftpRequests)
		aborted, _ := strconv.Atoi(node.TotalTftpRequestsAborted)

		results = append(results, ASTTftpInfo{
			Node:          node.Name,
			TotalRequests: totalReqs,
			Aborted:       aborted,
		})
	}

	return results, nil
}

// ASTHeartbeat represents heartbeat information
type ASTHeartbeat struct {
	Type string
	Node string
	Rate int
}

// GetASTHeartbeat retrieves heartbeat information
func GetASTHeartbeat(host, user, pass string) ([]ASTHeartbeat, error) {
	body, err := astGet(host, user, pass, "GetPreCannedInfo&Items=getHeartbeatInfoRequest")
	if err != nil {
		return nil, err
	}

	var preCanned struct {
		ReturnCode string `xml:"ReturnCode,attr"`
		Reply      struct {
			ReturnCode string `xml:"ReturnCode,attr"`
			CMs        struct {
				Nodes []struct {
					NodeName      string `xml:"NodeName,attr"`
					HeartbeatRate string `xml:"HeartbeatRate,attr"`
				} `xml:"Node"`
			} `xml:"CMs"`
			TFTPs struct {
				Nodes []struct {
					NodeName      string `xml:"NodeName,attr"`
					HeartbeatRate string `xml:"HeartbeatRate,attr"`
				} `xml:"Node"`
			} `xml:"TFTPs"`
		} `xml:"getHeartbeatInfoReply"`
	}

	if err := xml.Unmarshal(body, &preCanned); err != nil {
		return nil, fmt.Errorf("failed to parse heartbeat info: %w", err)
	}

	if preCanned.ReturnCode != "0" {
		return nil, fmt.Errorf("heartbeat info request returned error code: %s", preCanned.ReturnCode)
	}

	var results []ASTHeartbeat

	for _, node := range preCanned.Reply.CMs.Nodes {
		rate, _ := strconv.Atoi(node.HeartbeatRate)
		results = append(results, ASTHeartbeat{
			Type: "CM",
			Node: node.NodeName,
			Rate: rate,
		})
	}

	for _, node := range preCanned.Reply.TFTPs.Nodes {
		rate, _ := strconv.Atoi(node.HeartbeatRate)
		results = append(results, ASTHeartbeat{
			Type: "TFTP",
			Node: node.NodeName,
			Rate: rate,
		})
	}

	return results, nil
}

// ASTAlert represents alert information
type ASTAlert struct {
	AlertID           string
	DisplayName       string
	Group             string
	IsTriggered       bool
	IsEnabled         bool
	IsWithinSafeRange bool
	Timestamp         string
}

// GetASTAlerts retrieves alert information
func GetASTAlerts(host, user, pass string, triggeredOnly bool) ([]ASTAlert, error) {
	body, err := astGet(host, user, pass, "GetAlertSummaryList")
	if err != nil {
		return nil, err
	}

	var reply struct {
		ReturnCode string `xml:"ReturnCode,attr"`
		Alerts     []struct {
			AlertID           string `xml:"AlertID,attr"`
			DisplayName       string `xml:"DisplayName,attr"`
			IsEnabled         string `xml:"IsEnabled,attr"`
			Group             string `xml:"Group,attr"`
			IsWithinSafeRange string `xml:"IsWithinSafeRange,attr"`
			IsTriggered       string `xml:"IsTriggered,attr"`
			TimeStamp         string `xml:"TimeStamp,attr"`
		} `xml:"AlertSummaryList>AlertSummary"`
	}

	if err := xml.Unmarshal(body, &reply); err != nil {
		return nil, fmt.Errorf("failed to parse alerts: %w", err)
	}

	var results []ASTAlert
	for _, alert := range reply.Alerts {
		isTriggered := alert.IsTriggered == "1"

		if triggeredOnly && !isTriggered {
			continue
		}

		isEnabled := alert.IsEnabled == "1"
		isWithinSafeRange := alert.IsWithinSafeRange == "1"

		results = append(results, ASTAlert{
			AlertID:           alert.AlertID,
			DisplayName:       alert.DisplayName,
			Group:             alert.Group,
			IsTriggered:       isTriggered,
			IsEnabled:         isEnabled,
			IsWithinSafeRange: isWithinSafeRange,
			Timestamp:         alert.TimeStamp,
		})
	}

	return results, nil
}

// ASTPerfmonObject represents perfmon object information
type ASTPerfmonObject struct {
	Host         string
	ObjectName   string
	HasInstances bool
	CounterCount int
}

// GetASTPerfmonObjects retrieves perfmon object list
func GetASTPerfmonObjects(host, user, pass string) ([]ASTPerfmonObject, error) {
	body, err := astGet(host, user, pass, "PerfmonListObject")
	if err != nil {
		return nil, err
	}

	var reply struct {
		ReturnCode string `xml:"ReturnCode,attr"`
		Clusters   []struct {
			Name  string `xml:"Name,attr"`
			Hosts []struct {
				Name    string `xml:"Name,attr"`
				Type    string `xml:"Type,attr"`
				Objects []struct {
					Name         string `xml:"Name,attr"`
					HasInstances string `xml:"HasInstances,attr"`
					Counters     []struct {
						Name string `xml:"Name,attr"`
					} `xml:"Counter"`
				} `xml:"Object"`
			} `xml:"Host"`
		} `xml:"Cluster"`
	}

	if err := xml.Unmarshal(body, &reply); err != nil {
		return nil, fmt.Errorf("failed to parse perfmon objects: %w", err)
	}

	if reply.ReturnCode != "0" {
		return nil, fmt.Errorf("perfmon request returned error code: %s", reply.ReturnCode)
	}

	var results []ASTPerfmonObject
	for _, cluster := range reply.Clusters {
		for _, h := range cluster.Hosts {
			for _, obj := range h.Objects {
				hasInst := obj.HasInstances == "true"
				counterCount := len(obj.Counters)

				results = append(results, ASTPerfmonObject{
					Host:         h.Name,
					ObjectName:   obj.Name,
					HasInstances: hasInst,
					CounterCount: counterCount,
				})
			}
		}
	}

	return results, nil
}
