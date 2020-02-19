package common

import "time"

//Configurations contains all the config data needed for the service
type Configurations struct {
	Server               ServerConfigurations
	AWS                  AWSConfigurations
	Version              string
	FMS                  FMSConfigurations
	HttpTimeoutInSeconds time.Duration
}

// ServerConfigurations exported
type ServerConfigurations struct {
	Host string
	Port int
}

type FMSConfigurations struct {
	URL string
	FetchVehiclesEndPoint string
}

// AWSConfigurations exported
type AWSConfigurations struct {
	Region string
	Tables TableConfigurations
}

// TableConfigurations exported
type TableConfigurations struct {
	Vehicle TableStructure
}

// TableStructure exported
type TableStructure struct {
	TableName string
	HashKey   string
}
