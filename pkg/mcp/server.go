package mcp

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/anthropic/swisseph-mcp/pkg/chart"
	"github.com/anthropic/swisseph-mcp/pkg/geo"
	"github.com/anthropic/swisseph-mcp/pkg/julian"
	"github.com/anthropic/swisseph-mcp/pkg/models"
	"github.com/anthropic/swisseph-mcp/pkg/sweph"
	"github.com/anthropic/swisseph-mcp/pkg/transit"
)

// Server implements the MCP protocol via JSON-RPC over stdio
type Server struct {
	ephePath string
}

// NewServer creates a new MCP server
func NewServer(ephePath string) *Server {
	return &Server{ephePath: ephePath}
}

// JSON-RPC structures
type jsonRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type jsonRPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	Error   *rpcError   `json:"error,omitempty"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// MCP protocol structures
type initializeResult struct {
	ProtocolVersion string     `json:"protocolVersion"`
	Capabilities    capability `json:"capabilities"`
	ServerInfo      serverInfo `json:"serverInfo"`
}

type capability struct {
	Tools *toolsCap `json:"tools,omitempty"`
}

type toolsCap struct {
	ListChanged bool `json:"listChanged"`
}

type serverInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type toolsListResult struct {
	Tools []toolDef `json:"tools"`
}

type toolDef struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	InputSchema json.RawMessage `json:"inputSchema"`
}

type callToolParams struct {
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments"`
}

type callToolResult struct {
	Content []contentItem `json:"content"`
	IsError bool          `json:"isError,omitempty"`
}

type contentItem struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// Run starts the MCP server, reading from stdin and writing to stdout
func (s *Server) Run() error {
	// Initialize Swiss Ephemeris
	absPath, _ := filepath.Abs(s.ephePath)
	sweph.Init(absPath)
	defer sweph.Close()

	decoder := json.NewDecoder(os.Stdin)
	encoder := json.NewEncoder(os.Stdout)

	for {
		var req jsonRPCRequest
		if err := decoder.Decode(&req); err != nil {
			if err == io.EOF {
				return nil
			}
			continue
		}

		resp := s.handleRequest(&req)
		if resp != nil {
			encoder.Encode(resp)
		}
	}
}

func (s *Server) handleRequest(req *jsonRPCRequest) *jsonRPCResponse {
	switch req.Method {
	case "initialize":
		return s.handleInitialize(req)
	case "notifications/initialized":
		return nil // notification, no response
	case "tools/list":
		return s.handleToolsList(req)
	case "tools/call":
		return s.handleToolsCall(req)
	default:
		return &jsonRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error:   &rpcError{Code: -32601, Message: "method not found: " + req.Method},
		}
	}
}

func (s *Server) handleInitialize(req *jsonRPCRequest) *jsonRPCResponse {
	return &jsonRPCResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result: initializeResult{
			ProtocolVersion: "2024-11-05",
			Capabilities: capability{
				Tools: &toolsCap{ListChanged: false},
			},
			ServerInfo: serverInfo{
				Name:    "swisseph-mcp",
				Version: "1.0.0",
			},
		},
	}
}

func (s *Server) handleToolsList(req *jsonRPCRequest) *jsonRPCResponse {
	tools := []toolDef{
		{
			Name:        "geocode",
			Description: "根据地点名称返回地理坐标（经纬度）和时区",
			InputSchema: json.RawMessage(`{
				"type": "object",
				"properties": {
					"location_name": {"type": "string", "description": "地点名称，支持中英文"}
				},
				"required": ["location_name"]
			}`),
		},
		{
			Name:        "datetime_to_jd",
			Description: "将公历日期时间（ISO 8601）转换为儒略日（UT 和 TT）",
			InputSchema: json.RawMessage(`{
				"type": "object",
				"properties": {
					"datetime": {"type": "string", "description": "ISO 8601 格式日期时间"},
					"calendar": {"type": "string", "enum": ["GREGORIAN", "JULIAN"], "default": "GREGORIAN"}
				},
				"required": ["datetime"]
			}`),
		},
		{
			Name:        "jd_to_datetime",
			Description: "将儒略日转换为公历日期时间（ISO 8601）",
			InputSchema: json.RawMessage(`{
				"type": "object",
				"properties": {
					"jd": {"type": "number", "description": "儒略日"},
					"timezone": {"type": "string", "default": "UTC", "description": "目标时区"}
				},
				"required": ["jd"]
			}`),
		},
		{
			Name:        "calc_single_chart",
			Description: "单盘计算：在固定时间点计算天体位置、宫位和相位",
			InputSchema: json.RawMessage(`{
				"type": "object",
				"properties": {
					"latitude": {"type": "number"},
					"longitude": {"type": "number"},
					"jd_ut": {"type": "number"},
					"planets": {"type": "array", "items": {"type": "string"}},
					"house_system": {"type": "string", "default": "PLACIDUS"},
					"orb_config": {"type": "object"}
				},
				"required": ["latitude", "longitude", "jd_ut"]
			}`),
		},
		{
			Name:        "calc_double_chart",
			Description: "双盘计算：计算内外盘各自天体位置及跨盘相位",
			InputSchema: json.RawMessage(`{
				"type": "object",
				"properties": {
					"inner_latitude": {"type": "number"},
					"inner_longitude": {"type": "number"},
					"inner_jd_ut": {"type": "number"},
					"inner_planets": {"type": "array", "items": {"type": "string"}},
					"outer_latitude": {"type": "number"},
					"outer_longitude": {"type": "number"},
					"outer_jd_ut": {"type": "number"},
					"outer_planets": {"type": "array", "items": {"type": "string"}},
					"special_points": {"type": "object"},
					"house_system": {"type": "string", "default": "PLACIDUS"},
					"orb_config": {"type": "object"}
				},
				"required": ["inner_latitude", "inner_longitude", "inner_jd_ut",
					"outer_latitude", "outer_longitude", "outer_jd_ut"]
			}`),
		},
		{
			Name:        "calc_transit",
			Description: "推运计算：在时间范围内搜索行运天体与本命天体之间所有占星事件",
			InputSchema: json.RawMessage(`{
				"type": "object",
				"properties": {
					"natal_latitude": {"type": "number"},
					"natal_longitude": {"type": "number"},
					"natal_jd_ut": {"type": "number"},
					"natal_planets": {"type": "array", "items": {"type": "string"}},
					"transit_latitude": {"type": "number"},
					"transit_longitude": {"type": "number"},
					"start_jd_ut": {"type": "number"},
					"end_jd_ut": {"type": "number"},
					"transit_planets": {"type": "array", "items": {"type": "string"}},
					"special_points": {"type": "object"},
					"event_config": {"type": "object"},
					"house_system": {"type": "string", "default": "PLACIDUS"},
					"orb_config": {"type": "object"}
				},
				"required": ["natal_latitude", "natal_longitude", "natal_jd_ut",
					"transit_latitude", "transit_longitude", "start_jd_ut", "end_jd_ut"]
			}`),
		},
	}

	return &jsonRPCResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  toolsListResult{Tools: tools},
	}
}

func (s *Server) handleToolsCall(req *jsonRPCRequest) *jsonRPCResponse {
	var params callToolParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return errorResponse(req.ID, -32602, "invalid params")
	}

	var result interface{}
	var err error

	switch params.Name {
	case "geocode":
		result, err = s.handleGeocode(params.Arguments)
	case "datetime_to_jd":
		result, err = s.handleDatetimeToJD(params.Arguments)
	case "jd_to_datetime":
		result, err = s.handleJDToDatetime(params.Arguments)
	case "calc_single_chart":
		result, err = s.handleCalcSingleChart(params.Arguments)
	case "calc_double_chart":
		result, err = s.handleCalcDoubleChart(params.Arguments)
	case "calc_transit":
		result, err = s.handleCalcTransit(params.Arguments)
	default:
		return errorResponse(req.ID, -32601, "unknown tool: "+params.Name)
	}

	if err != nil {
		return &jsonRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: callToolResult{
				Content: []contentItem{{Type: "text", Text: fmt.Sprintf("Error: %v", err)}},
				IsError: true,
			},
		}
	}

	jsonBytes, _ := json.MarshalIndent(result, "", "  ")
	return &jsonRPCResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result: callToolResult{
			Content: []contentItem{{Type: "text", Text: string(jsonBytes)}},
		},
	}
}

func errorResponse(id interface{}, code int, msg string) *jsonRPCResponse {
	return &jsonRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error:   &rpcError{Code: code, Message: msg},
	}
}

// === Tool handlers ===

func (s *Server) handleGeocode(args json.RawMessage) (interface{}, error) {
	var input struct {
		LocationName string `json:"location_name"`
	}
	if err := json.Unmarshal(args, &input); err != nil {
		return nil, err
	}
	return geo.Geocode(input.LocationName)
}

func (s *Server) handleDatetimeToJD(args json.RawMessage) (interface{}, error) {
	var input struct {
		Datetime string             `json:"datetime"`
		Calendar models.CalendarType `json:"calendar"`
	}
	if err := json.Unmarshal(args, &input); err != nil {
		return nil, err
	}
	if input.Calendar == "" {
		input.Calendar = models.CalendarGregorian
	}
	return julian.DateTimeToJD(input.Datetime, input.Calendar)
}

func (s *Server) handleJDToDatetime(args json.RawMessage) (interface{}, error) {
	var input struct {
		JD       float64 `json:"jd"`
		Timezone string  `json:"timezone"`
	}
	if err := json.Unmarshal(args, &input); err != nil {
		return nil, err
	}
	if input.Timezone == "" {
		input.Timezone = "UTC"
	}
	dt, err := julian.JDToDateTime(input.JD, input.Timezone)
	if err != nil {
		return nil, err
	}
	return map[string]string{"datetime": dt}, nil
}

func (s *Server) handleCalcSingleChart(args json.RawMessage) (interface{}, error) {
	var input struct {
		Latitude    float64          `json:"latitude"`
		Longitude   float64          `json:"longitude"`
		JDUT        float64          `json:"jd_ut"`
		Planets     []models.PlanetID `json:"planets"`
		HouseSystem models.HouseSystem `json:"house_system"`
		OrbConfig   *models.OrbConfig  `json:"orb_config"`
	}
	if err := json.Unmarshal(args, &input); err != nil {
		return nil, err
	}

	if len(input.Planets) == 0 {
		input.Planets = []models.PlanetID{
			models.PlanetSun, models.PlanetMoon, models.PlanetMercury,
			models.PlanetVenus, models.PlanetMars, models.PlanetJupiter,
			models.PlanetSaturn, models.PlanetUranus, models.PlanetNeptune,
			models.PlanetPluto,
		}
	}
	if input.HouseSystem == "" {
		input.HouseSystem = models.HousePlacidus
	}
	orbs := models.DefaultOrbConfig()
	if input.OrbConfig != nil {
		orbs = *input.OrbConfig
	}

	return chart.CalcSingleChart(input.Latitude, input.Longitude, input.JDUT, input.Planets, orbs, input.HouseSystem)
}

func (s *Server) handleCalcDoubleChart(args json.RawMessage) (interface{}, error) {
	var input struct {
		InnerLatitude  float64              `json:"inner_latitude"`
		InnerLongitude float64              `json:"inner_longitude"`
		InnerJDUT      float64              `json:"inner_jd_ut"`
		InnerPlanets   []models.PlanetID    `json:"inner_planets"`
		OuterLatitude  float64              `json:"outer_latitude"`
		OuterLongitude float64              `json:"outer_longitude"`
		OuterJDUT      float64              `json:"outer_jd_ut"`
		OuterPlanets   []models.PlanetID    `json:"outer_planets"`
		SpecialPoints  *models.SpecialPointsConfig `json:"special_points"`
		HouseSystem    models.HouseSystem    `json:"house_system"`
		OrbConfig      *models.OrbConfig     `json:"orb_config"`
	}
	if err := json.Unmarshal(args, &input); err != nil {
		return nil, err
	}

	defaultPlanets := []models.PlanetID{
		models.PlanetSun, models.PlanetMoon, models.PlanetMercury,
		models.PlanetVenus, models.PlanetMars, models.PlanetJupiter,
		models.PlanetSaturn, models.PlanetUranus, models.PlanetNeptune,
		models.PlanetPluto,
	}
	if len(input.InnerPlanets) == 0 {
		input.InnerPlanets = defaultPlanets
	}
	if len(input.OuterPlanets) == 0 {
		input.OuterPlanets = defaultPlanets
	}
	if input.HouseSystem == "" {
		input.HouseSystem = models.HousePlacidus
	}
	orbs := models.DefaultOrbConfig()
	if input.OrbConfig != nil {
		orbs = *input.OrbConfig
	}

	innerChart, outerChart, crossAspects, err := chart.CalcDoubleChart(
		input.InnerLatitude, input.InnerLongitude, input.InnerJDUT, input.InnerPlanets,
		input.OuterLatitude, input.OuterLongitude, input.OuterJDUT, input.OuterPlanets,
		input.SpecialPoints, orbs, input.HouseSystem,
	)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"inner_chart":   innerChart,
		"outer_chart":   outerChart,
		"cross_aspects": crossAspects,
	}, nil
}

func (s *Server) handleCalcTransit(args json.RawMessage) (interface{}, error) {
	var input struct {
		NatalLatitude  float64                     `json:"natal_latitude"`
		NatalLongitude float64                     `json:"natal_longitude"`
		NatalJDUT      float64                     `json:"natal_jd_ut"`
		NatalPlanets   []models.PlanetID           `json:"natal_planets"`
		TransitLatitude  float64                   `json:"transit_latitude"`
		TransitLongitude float64                   `json:"transit_longitude"`
		StartJDUT      float64                     `json:"start_jd_ut"`
		EndJDUT        float64                     `json:"end_jd_ut"`
		TransitPlanets []models.PlanetID           `json:"transit_planets"`
		SpecialPoints  *models.SpecialPointsConfig `json:"special_points"`
		EventConfig    *models.EventConfig         `json:"event_config"`
		HouseSystem    models.HouseSystem          `json:"house_system"`
		OrbConfig      *models.OrbConfig           `json:"orb_config"`
	}
	if err := json.Unmarshal(args, &input); err != nil {
		return nil, err
	}

	defaultPlanets := []models.PlanetID{
		models.PlanetSun, models.PlanetMoon, models.PlanetMercury,
		models.PlanetVenus, models.PlanetMars, models.PlanetJupiter,
		models.PlanetSaturn, models.PlanetUranus, models.PlanetNeptune,
		models.PlanetPluto,
	}
	if len(input.NatalPlanets) == 0 {
		input.NatalPlanets = defaultPlanets
	}
	if len(input.TransitPlanets) == 0 {
		input.TransitPlanets = defaultPlanets
	}
	if input.HouseSystem == "" {
		input.HouseSystem = models.HousePlacidus
	}
	orbs := models.DefaultOrbConfig()
	if input.OrbConfig != nil {
		orbs = *input.OrbConfig
	}
	eventCfg := models.DefaultEventConfig()
	if input.EventConfig != nil {
		eventCfg = *input.EventConfig
	}

	events, err := transit.CalcTransitEvents(transit.TransitCalcInput{
		NatalLat:       input.NatalLatitude,
		NatalLon:       input.NatalLongitude,
		NatalJD:        input.NatalJDUT,
		NatalPlanets:   input.NatalPlanets,
		TransitLat:     input.TransitLatitude,
		TransitLon:     input.TransitLongitude,
		StartJD:        input.StartJDUT,
		EndJD:          input.EndJDUT,
		TransitPlanets: input.TransitPlanets,
		SpecialPoints:  input.SpecialPoints,
		EventConfig:    eventCfg,
		OrbConfig:      orbs,
		HouseSystem:    input.HouseSystem,
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"events": events,
	}, nil
}
