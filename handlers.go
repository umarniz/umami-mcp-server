package main

import (
	"encoding/json"
	"fmt"
)

func (s *MCPServer) execGetWebsites(args json.RawMessage) (any, *Error) {
	var params struct {
		IncludeTeams bool `json:"includeTeams"`
	}

	if len(args) > 0 {
		_ = json.Unmarshal(args, &params)
	}

	websites, err := s.client.GetWebsites(params.IncludeTeams)
	if err != nil {
		return nil, &Error{Code: -32603, Message: fmt.Sprintf("Failed to get websites: %v", err)}
	}

	data, _ := json.MarshalIndent(websites, "", "  ")
	content := []map[string]string{{
		"type": "text",
		"text": string(data),
	}}

	return map[string]any{"content": content}, nil
}

func (s *MCPServer) execGetStats(args json.RawMessage) (any, *Error) {
	var params struct {
		WebsiteID string `json:"website_id"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}

	if err := json.Unmarshal(args, &params); err != nil {
		return nil, &Error{Code: -32602, Message: "Invalid arguments"}
	}

	if err := validateWebsiteID(params.WebsiteID); err != nil {
		return nil, &Error{Code: -32602, Message: "Invalid website_id"}
	}

	params.StartDate = normalizeDate(params.StartDate)
	params.EndDate = normalizeDate(params.EndDate)

	stats, err := s.client.GetStats(params.WebsiteID, params.StartDate, params.EndDate)
	if err != nil {
		return nil, &Error{Code: -32603, Message: fmt.Sprintf("Failed to get stats: %v", err)}
	}

	data, _ := json.MarshalIndent(stats, "", "  ")
	content := []map[string]string{{
		"type": "text",
		"text": string(data),
	}}

	return map[string]any{"content": content}, nil
}

func (s *MCPServer) execGetPageViews(args json.RawMessage) (any, *Error) {
	var params struct {
		WebsiteID string `json:"website_id"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
		Unit      string `json:"unit"`
	}

	if err := json.Unmarshal(args, &params); err != nil {
		return nil, &Error{Code: -32602, Message: "Invalid arguments"}
	}

	if err := validateWebsiteID(params.WebsiteID); err != nil {
		return nil, &Error{Code: -32602, Message: "Invalid website_id"}
	}

	if params.Unit == "" {
		params.Unit = "day"
	}

	params.StartDate = normalizeDate(params.StartDate)
	params.EndDate = normalizeDate(params.EndDate)

	pageviews, err := s.client.GetPageViews(params.WebsiteID, params.StartDate, params.EndDate, params.Unit)
	if err != nil {
		return nil, &Error{Code: -32603, Message: fmt.Sprintf("Failed to get page views: %v", err)}
	}

	data, _ := json.MarshalIndent(pageviews, "", "  ")
	content := []map[string]string{{
		"type": "text",
		"text": string(data),
	}}

	return map[string]any{"content": content}, nil
}

func (s *MCPServer) execGetMetrics(args json.RawMessage) (any, *Error) {
	var params struct {
		WebsiteID  string `json:"website_id"`
		StartDate  string `json:"start_date"`
		EndDate    string `json:"end_date"`
		MetricType string `json:"metric_type"`
		Limit      int    `json:"limit"`
	}

	if err := json.Unmarshal(args, &params); err != nil {
		return nil, &Error{Code: -32602, Message: "Invalid arguments"}
	}

	if err := validateWebsiteID(params.WebsiteID); err != nil {
		return nil, &Error{Code: -32602, Message: "Invalid website_id"}
	}

	if params.Limit == 0 {
		params.Limit = 10
	}

	params.StartDate = normalizeDate(params.StartDate)
	params.EndDate = normalizeDate(params.EndDate)

	metrics, err := s.client.GetMetrics(
		params.WebsiteID, params.StartDate, params.EndDate, params.MetricType, params.Limit,
	)
	if err != nil {
		return nil, &Error{Code: -32603, Message: fmt.Sprintf("Failed to get metrics: %v", err)}
	}

	data, _ := json.MarshalIndent(metrics, "", "  ")
	content := []map[string]string{{
		"type": "text",
		"text": string(data),
	}}

	return map[string]any{"content": content}, nil
}

func (s *MCPServer) execGetActive(args json.RawMessage) (any, *Error) {
	var params struct {
		WebsiteID string `json:"website_id"`
	}

	if err := json.Unmarshal(args, &params); err != nil {
		return nil, &Error{Code: -32602, Message: "Invalid arguments"}
	}

	if err := validateWebsiteID(params.WebsiteID); err != nil {
		return nil, &Error{Code: -32602, Message: "Invalid website_id"}
	}

	active, err := s.client.GetActive(params.WebsiteID)
	if err != nil {
		return nil, &Error{Code: -32603, Message: fmt.Sprintf("Failed to get active visitors: %v", err)}
	}

	data, _ := json.MarshalIndent(active, "", "  ")
	content := []map[string]string{{
		"type": "text",
		"text": string(data),
	}}

	return map[string]any{"content": content}, nil
}

func (s *MCPServer) execListReports(args json.RawMessage) (any, *Error) {
	var params struct {
		WebsiteID  string `json:"website_id"`
		ReportType string `json:"report_type"`
		Page       int    `json:"page"`
		PageSize   int    `json:"page_size"`
	}

	if err := json.Unmarshal(args, &params); err != nil {
		return nil, &Error{Code: -32602, Message: "Invalid arguments"}
	}

	if err := validateWebsiteID(params.WebsiteID); err != nil {
		return nil, &Error{Code: -32602, Message: "Invalid website_id"}
	}
	if params.ReportType != "" {
		if err := validateReportType(params.ReportType); err != nil {
			return nil, &Error{Code: -32602, Message: "Invalid report_type"}
		}
	}

	data, err := s.client.ListReports(params.WebsiteID, params.ReportType, params.Page, params.PageSize)
	if err != nil {
		return nil, &Error{Code: -32603, Message: fmt.Sprintf("Failed to list reports: %v", err)}
	}

	return rawJSONContent(data), nil
}

func (s *MCPServer) execGetReport(args json.RawMessage) (any, *Error) {
	var params struct {
		ReportID string `json:"report_id"`
	}

	if err := json.Unmarshal(args, &params); err != nil {
		return nil, &Error{Code: -32602, Message: "Invalid arguments"}
	}

	if err := validateReportID(params.ReportID); err != nil {
		return nil, &Error{Code: -32602, Message: "Invalid report_id"}
	}

	data, err := s.client.GetReport(params.ReportID)
	if err != nil {
		return nil, &Error{Code: -32603, Message: fmt.Sprintf("Failed to get report: %v", err)}
	}

	return rawJSONContent(data), nil
}

func (s *MCPServer) execCreateReport(args json.RawMessage) (any, *Error) {
	report, rpcErr := parseReportRequest(args)
	if rpcErr != nil {
		return nil, rpcErr
	}

	data, err := s.client.CreateReport(report)
	if err != nil {
		return nil, &Error{Code: -32603, Message: fmt.Sprintf("Failed to create report: %v", err)}
	}

	return rawJSONContent(data), nil
}

func (s *MCPServer) execUpdateReport(args json.RawMessage) (any, *Error) {
	var params struct {
		ReportID       string         `json:"report_id"`
		WebsiteID      string         `json:"website_id"`
		WebsiteIDCamel string         `json:"websiteId"`
		Type           string         `json:"type"`
		Name           string         `json:"name"`
		Description    string         `json:"description"`
		Parameters     map[string]any `json:"parameters"`
	}

	if err := json.Unmarshal(args, &params); err != nil {
		return nil, &Error{Code: -32602, Message: "Invalid arguments"}
	}

	if err := validateReportID(params.ReportID); err != nil {
		return nil, &Error{Code: -32602, Message: "Invalid report_id"}
	}

	report, rpcErr := validateReportRequest(ReportRequest{
		WebsiteID:   firstNonEmpty(params.WebsiteID, params.WebsiteIDCamel),
		Type:        params.Type,
		Name:        params.Name,
		Description: params.Description,
		Parameters:  params.Parameters,
	})
	if rpcErr != nil {
		return nil, rpcErr
	}

	data, err := s.client.UpdateReport(params.ReportID, report)
	if err != nil {
		return nil, &Error{Code: -32603, Message: fmt.Sprintf("Failed to update report: %v", err)}
	}

	return rawJSONContent(data), nil
}

func (s *MCPServer) execDeleteReport(args json.RawMessage) (any, *Error) {
	var params struct {
		ReportID string `json:"report_id"`
	}

	if err := json.Unmarshal(args, &params); err != nil {
		return nil, &Error{Code: -32602, Message: "Invalid arguments"}
	}

	if err := validateReportID(params.ReportID); err != nil {
		return nil, &Error{Code: -32602, Message: "Invalid report_id"}
	}

	data, err := s.client.DeleteReport(params.ReportID)
	if err != nil {
		return nil, &Error{Code: -32603, Message: fmt.Sprintf("Failed to delete report: %v", err)}
	}

	if len(data) == 0 {
		data = json.RawMessage(`{"deleted":true}`)
	}
	return rawJSONContent(data), nil
}

func (s *MCPServer) execRunReport(args json.RawMessage) (any, *Error) {
	var params struct {
		ReportType string         `json:"report_type"`
		Payload    map[string]any `json:"payload"`
	}

	if err := json.Unmarshal(args, &params); err != nil {
		return nil, &Error{Code: -32602, Message: "Invalid arguments"}
	}

	if err := validateReportType(params.ReportType); err != nil {
		return nil, &Error{Code: -32602, Message: "Invalid report_type"}
	}
	if params.Payload == nil {
		return nil, &Error{Code: -32602, Message: "payload is required"}
	}
	if params.ReportType == reportTypeFunnel {
		if rpcErr := normalizeFunnelPayload(params.Payload); rpcErr != nil {
			return nil, rpcErr
		}
	}
	if websiteID, ok := params.Payload["websiteId"].(string); ok {
		if err := validateWebsiteID(websiteID); err != nil {
			return nil, &Error{Code: -32602, Message: "Invalid payload.websiteId"}
		}
	}

	data, err := s.client.RunReport(params.ReportType, params.Payload)
	if err != nil {
		return nil, &Error{Code: -32603, Message: fmt.Sprintf("Failed to run report: %v", err)}
	}

	return rawJSONContent(data), nil
}

func parseReportRequest(args json.RawMessage) (ReportRequest, *Error) {
	var params struct {
		WebsiteID      string         `json:"website_id"`
		WebsiteIDCamel string         `json:"websiteId"`
		Type           string         `json:"type"`
		Name           string         `json:"name"`
		Description    string         `json:"description"`
		Parameters     map[string]any `json:"parameters"`
	}
	if err := json.Unmarshal(args, &params); err != nil {
		return ReportRequest{}, &Error{Code: -32602, Message: "Invalid arguments"}
	}
	return validateReportRequest(ReportRequest{
		WebsiteID:   firstNonEmpty(params.WebsiteID, params.WebsiteIDCamel),
		Type:        params.Type,
		Name:        params.Name,
		Description: params.Description,
		Parameters:  params.Parameters,
	})
}

func validateReportRequest(params ReportRequest) (ReportRequest, *Error) {
	if err := validateWebsiteID(params.WebsiteID); err != nil {
		return ReportRequest{}, &Error{Code: -32602, Message: "Invalid websiteId"}
	}
	if err := validateReportType(params.Type); err != nil {
		return ReportRequest{}, &Error{Code: -32602, Message: "Invalid type"}
	}
	if params.Name == "" {
		return ReportRequest{}, &Error{Code: -32602, Message: "name is required"}
	}
	if params.Parameters == nil {
		return ReportRequest{}, &Error{Code: -32602, Message: "parameters is required"}
	}
	if params.Type == reportTypeFunnel {
		if rpcErr := normalizeFunnelPayload(params.Parameters); rpcErr != nil {
			return ReportRequest{}, rpcErr
		}
	}
	return params, nil
}

func normalizeFunnelPayload(payload map[string]any) *Error {
	if rpcErr := normalizeFunnelSteps(payload); rpcErr != nil {
		return rpcErr
	}

	parameters, ok := payload["parameters"].(map[string]any)
	if !ok {
		return nil
	}
	return normalizeFunnelSteps(parameters)
}

func normalizeFunnelSteps(payload map[string]any) *Error {
	steps, ok := payload["steps"]
	if !ok {
		return nil
	}

	stepList, ok := steps.([]any)
	if !ok {
		return &Error{Code: -32602, Message: "funnel parameters.steps must be an array"}
	}

	for i, step := range stepList {
		stepMap, ok := step.(map[string]any)
		if !ok {
			return &Error{Code: -32602, Message: fmt.Sprintf("funnel step %d must be an object", i)}
		}

		stepType, ok := stepMap["type"].(string)
		if !ok || stepType == "" {
			return &Error{Code: -32602, Message: fmt.Sprintf("funnel step %d type is required", i)}
		}

		switch stepType {
		case metricTypePath, "event":
		case "url", "URL", "page", "pathname":
			stepMap["type"] = metricTypePath
		default:
			return &Error{
				Code:    -32602,
				Message: fmt.Sprintf("funnel step %d type must be \"path\" or \"event\"", i),
			}
		}
	}

	return nil
}

func rawJSONContent(data []byte) map[string]any {
	var out any
	if err := json.Unmarshal(data, &out); err != nil {
		out = string(data)
	}
	formatted, _ := json.MarshalIndent(out, "", "  ")

	content := []map[string]string{{
		"type": "text",
		"text": string(formatted),
	}}
	return map[string]any{"content": content}
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
