package client

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
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

// ASTPerfmonCounter represents a perfmon counter with description
type ASTPerfmonCounter struct {
	Host        string
	ObjectName  string
	CounterName string
	Description string
	Unit        string
	IsExcluded  bool
}

// GetASTPerfmonCounters retrieves detailed counter list for a specific perfmon object.
// Query: PerfmonListCounter&Host=<host>&Object=<objectName>
func GetASTPerfmonCounters(host, user, pass, objectName string) ([]ASTPerfmonCounter, error) {
	query := fmt.Sprintf("PerfmonListCounter&Host=%s&Object=%s",
		url.QueryEscape(host), url.QueryEscape(objectName))
	body, err := astGet(host, user, pass, query)
	if err != nil {
		return nil, err
	}

	var reply struct {
		ReturnCode string `xml:"ReturnCode,attr"`
		Clusters   []struct {
			Name  string `xml:"Name,attr"`
			Hosts []struct {
				Name    string `xml:"Name,attr"`
				Objects []struct {
					Name     string `xml:"Name,attr"`
					Counters []struct {
						Name        string `xml:"Name,attr"`
						IsExcluded  string `xml:"IsExcluded,attr"`
						Description string `xml:"Description,attr"`
						Unit        string `xml:"Unit,attr"`
					} `xml:"Counter"`
				} `xml:"Object"`
			} `xml:"Host"`
		} `xml:"Cluster"`
	}

	if err := xml.Unmarshal(body, &reply); err != nil {
		return nil, fmt.Errorf("failed to parse perfmon counters: %w", err)
	}

	if reply.ReturnCode != "0" {
		return nil, fmt.Errorf("perfmon counter request returned error code: %s", reply.ReturnCode)
	}

	var results []ASTPerfmonCounter
	for _, cluster := range reply.Clusters {
		for _, h := range cluster.Hosts {
			for _, obj := range h.Objects {
				for _, c := range obj.Counters {
					results = append(results, ASTPerfmonCounter{
						Host:        h.Name,
						ObjectName:  obj.Name,
						CounterName: c.Name,
						Description: c.Description,
						Unit:        c.Unit,
						IsExcluded:  c.IsExcluded == "true",
					})
				}
			}
		}
	}

	return results, nil
}

// ASTPerfmonDataPoint represents a collected real-time perfmon counter value
type ASTPerfmonDataPoint struct {
	Host        string
	ObjectName  string
	Instance    string
	CounterName string
	Value       string
	CStatus     string
}

// GetASTPerfmonData collects real-time counter values for a specific perfmon object.
// Query: PerfmonCollectCounterData&Host=<host>&Object=<objectName>
func GetASTPerfmonData(host, user, pass, objectName string) ([]ASTPerfmonDataPoint, error) {
	query := fmt.Sprintf("PerfmonCollectCounterData&Host=%s&Object=%s",
		url.QueryEscape(host), url.QueryEscape(objectName))
	body, err := astGet(host, user, pass, query)
	if err != nil {
		return nil, err
	}

	var reply struct {
		ReturnCode string `xml:"ReturnCode,attr"`
		Clusters   []struct {
			Name  string `xml:"Name,attr"`
			Hosts []struct {
				Name    string `xml:"Name,attr"`
				Objects []struct {
					Name      string `xml:"Name,attr"`
					Instances []struct {
						Name     string `xml:"Name,attr"`
						Counters []struct {
							Name    string `xml:"Name,attr"`
							Value   string `xml:"Value,attr"`
							CStatus string `xml:"CStatus,attr"`
						} `xml:"Counter"`
					} `xml:"Instance"`
				} `xml:"Object"`
			} `xml:"Host"`
		} `xml:"Cluster"`
	}

	if err := xml.Unmarshal(body, &reply); err != nil {
		return nil, fmt.Errorf("failed to parse perfmon data: %w", err)
	}

	if reply.ReturnCode != "0" {
		return nil, fmt.Errorf("perfmon collect request returned error code: %s", reply.ReturnCode)
	}

	var results []ASTPerfmonDataPoint
	for _, cluster := range reply.Clusters {
		for _, h := range cluster.Hosts {
			for _, obj := range h.Objects {
				for _, inst := range obj.Instances {
					for _, c := range inst.Counters {
						results = append(results, ASTPerfmonDataPoint{
							Host:        h.Name,
							ObjectName:  obj.Name,
							Instance:    inst.Name,
							CounterName: c.Name,
							Value:       c.Value,
							CStatus:     c.CStatus,
						})
					}
				}
			}
		}
	}

	return results, nil
}

// ASTService represents a VOS service entry
type ASTService struct {
	ServiceName   string
	ServiceStatus string
	StartupType   string
	ReasonCode    string
	NodeName      string
}

// GetASTServiceList retrieves the list of services on the CUC node.
// Query: GetServiceList&NodeName=<host>
func GetASTServiceList(host, user, pass string) ([]ASTService, error) {
	query := fmt.Sprintf("GetServiceList&NodeName=%s", url.QueryEscape(host))
	body, err := astGet(host, user, pass, query)
	if err != nil {
		return nil, err
	}

	var reply struct {
		ReturnCode  string `xml:"ReturnCode,attr"`
		ServiceList []struct {
			ServiceName   string `xml:"ServiceName,attr"`
			ServiceStatus string `xml:"ServiceStatus,attr"`
			StartupType   string `xml:"StartupType,attr"`
			ReasonCode    string `xml:"ReasonCode,attr"`
			NodeName      string `xml:"NodeName,attr"`
		} `xml:"ServiceInfoList>service"`
	}

	if err := xml.Unmarshal(body, &reply); err != nil {
		return nil, fmt.Errorf("failed to parse service list: %w", err)
	}

	if reply.ReturnCode != "0" {
		return nil, fmt.Errorf("service list request returned error code: %s", reply.ReturnCode)
	}

	var results []ASTService
	for _, s := range reply.ServiceList {
		results = append(results, ASTService{
			ServiceName:   s.ServiceName,
			ServiceStatus: s.ServiceStatus,
			StartupType:   s.StartupType,
			ReasonCode:    s.ReasonCode,
			NodeName:      s.NodeName,
		})
	}

	return results, nil
}

// DoASTServiceAction performs a Start, Stop, or Restart action on a named VOS service.
// Query: DoServiceAction&NodeName=<host>&ServiceName=<name>&Action=<action>
func DoASTServiceAction(host, user, pass, serviceName, action string) error {
	query := fmt.Sprintf("DoServiceAction&NodeName=%s&ServiceName=%s&Action=%s",
		url.QueryEscape(host), url.QueryEscape(serviceName), url.QueryEscape(action))
	_, err := astGet(host, user, pass, query)
	return err
}

// ASTAlertDetail represents detailed information for a single alert
type ASTAlertDetail struct {
	AlertID           string
	DisplayName       string
	Group             string
	Description       string
	IsEnabled         bool
	IsTriggered       bool
	IsWithinSafeRange bool
	ThresholdType     string
	ThresholdValue    string
	Severity          string
}

// GetASTAlertDetail retrieves detailed information for a specific alert by ID.
// Query: GetAlertDetail&AlertID=<id>
func GetASTAlertDetail(host, user, pass, alertID string) (*ASTAlertDetail, error) {
	query := fmt.Sprintf("GetAlertDetail&AlertID=%s", url.QueryEscape(alertID))
	body, err := astGet(host, user, pass, query)
	if err != nil {
		return nil, err
	}

	var reply struct {
		ReturnCode string `xml:"ReturnCode,attr"`
		Detail     struct {
			AlertID           string `xml:"AlertID,attr"`
			DisplayName       string `xml:"DisplayName,attr"`
			Group             string `xml:"Group,attr"`
			Description       string `xml:"Description,attr"`
			IsEnabled         string `xml:"IsEnabled,attr"`
			IsTriggered       string `xml:"IsTriggered,attr"`
			IsWithinSafeRange string `xml:"IsWithinSafeRange,attr"`
			ThresholdType     string `xml:"ThresholdType,attr"`
			ThresholdValue    string `xml:"ThresholdValue,attr"`
			Severity          string `xml:"Severity,attr"`
		} `xml:"AlertDetail"`
	}

	if err := xml.Unmarshal(body, &reply); err != nil {
		return nil, fmt.Errorf("failed to parse alert detail: %w", err)
	}

	if reply.ReturnCode != "0" {
		return nil, fmt.Errorf("alert detail request returned error code: %s", reply.ReturnCode)
	}

	d := reply.Detail
	return &ASTAlertDetail{
		AlertID:           d.AlertID,
		DisplayName:       d.DisplayName,
		Group:             d.Group,
		Description:       d.Description,
		IsEnabled:         d.IsEnabled == "1",
		IsTriggered:       d.IsTriggered == "1",
		IsWithinSafeRange: d.IsWithinSafeRange == "1",
		ThresholdType:     d.ThresholdType,
		ThresholdValue:    d.ThresholdValue,
		Severity:          d.Severity,
	}, nil
}
