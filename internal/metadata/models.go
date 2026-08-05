package metadata

import "github.com/google/uuid"

type Config struct {
	Batches []Batch `json:"BATCH"`
}

type Batch struct {
	BatchID     uuid.UUID
	BatchCode   string      `json:"BATCH_CODE"`
	Description string      `json:"BATCH_DESCRIPTION"`
	IsActive    string      `json:"BATCH_IS_ACTIVE"`
	Parameters  []Parameter `json:"BATCH_PARAMETERS"`
	Layers      []Layer     `json:"LAYERS"`
}

type Layer struct {
	LayerID     uuid.UUID
	LayerCode   string   `json:"LAYER_CODE"`
	Description string   `json:"DESCRIPTION"`
	IsActive    string   `json:"IS_ACTIVE"`
	Priority    int      `json:"PRIORITY"`
	Modules     []Module `json:"MODULES"`
}

type Module struct {
	ModuleID      uuid.UUID
	ModuleCode    string      `json:"MODULE_CODE"`
	Description   string      `json:"DESCRIPTION"`
	IsActive      string      `json:"IS_ACTIVE"`
	IsInitialFull string      `json:"IS_INITIAL_FULL"`
	Script        string      `json:"MODULE_SCRIPT"`
	Priority      int         `json:"PRIORITY"`
	Parameters    []Parameter `json:"PARAMETERS"`
}

type Parameter struct {
	ParameterID    uuid.UUID
	ParameterCode  string `json:"PARAMETER_CODE"`
	ParameterValue string `json:"PARAMETER_VALUE"`
	Description    string `json:"DESCRIPTION"`
	IsActive       string `json:"IS_ACTIVE"`
}
