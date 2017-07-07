package google

import (
	"fmt"
	"google.golang.org/api/compute/v1"
)

type ScopeType uint8

const (
	GLOBAL ScopeType = iota
	REGION
	ZONE
)

// A ComputeService that delegates requests to the appropriate API
// Takes in and returns models of the highest API level supported.
type ComputeMultiversionService struct {
	v1 *compute.Service
}

func (s *ComputeMultiversionService) WaitOperation(project string, operationName string, scopeType ScopeType, scope string) (*compute.Operation, error) {
	switch scopeType {
	case GLOBAL:
		return s.v1.GlobalOperations.Get(project, operationName).Do()
	case REGION:
		return s.v1.RegionOperations.Get(project, scope, operationName).Do()
	case ZONE:
		return s.v1.ZoneOperations.Get(project, scope, operationName).Do()
	}

	return nil, fmt.Errorf("Awaited operation with unknown scope. %v %s", scopeType, scope)
}

func (s *ComputeMultiversionService) InsertInstanceGroupManager(project string, zone string, manager *compute.InstanceGroupManager, version ComputeApiVersion) (*compute.Operation, error) {
	op := &compute.Operation{}
	switch version {
	case v1:
		v1Manager := &compute.InstanceGroupManager{}
		err := convert(manager, v1Manager)
		if err != nil {
			return nil, err
		}

		v1Op, err := s.v1.InstanceGroupManagers.Insert(project, zone, v1Manager).Do()
		err = convert(v1Op, op)
		if err != nil {
			return nil, err
		}

		return op, nil
	}

	return nil, fmt.Errorf("Unknown API version.")
}
