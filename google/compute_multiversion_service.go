package google

import (
	"fmt"

	computeBeta "google.golang.org/api/compute/v0.beta"
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
	v1     *compute.Service
	v0beta *computeBeta.Service
}

func (s *ComputeMultiversionService) WaitOperation(project string, operationName string, scopeType ScopeType, scope string) (*computeBeta.Operation, error) {
	var operation *compute.Operation
	var err error
	switch scopeType {
	case GLOBAL:
		operation, err = s.v1.GlobalOperations.Get(project, operationName).Do()
	case REGION:
		operation, err = s.v1.RegionOperations.Get(project, scope, operationName).Do()
	case ZONE:
		operation, err = s.v1.ZoneOperations.Get(project, scope, operationName).Do()
	default:
		operation, err = nil, fmt.Errorf("Awaited operation with unknown scope. %v %s", scopeType, scope)
	}

	if err != nil {
		return nil, err
	}

	v0BetaOperation := &computeBeta.Operation{}
	err = Convert(operation, v0BetaOperation)
	if err != nil {
		return nil, err
	}

	return v0BetaOperation, nil

}

func (s *ComputeMultiversionService) InsertInstanceGroupManager(project string, zone string, resource *computeBeta.InstanceGroupManager, version ComputeApiVersion) (*computeBeta.Operation, error) {
	op := &computeBeta.Operation{}
	switch version {

	case v0beta:
		v0betaResource := &computeBeta.InstanceGroupManager{}
		err := Convert(resource, v0betaResource)
		if err != nil {
			return nil, err
		}

		v0betaOp, err := s.v0beta.InstanceGroupManagers.Insert(project, zone, v0betaResource).Do()
		if err != nil {
			return nil, err
		}

		err = Convert(v0betaOp, op)
		if err != nil {
			return nil, err
		}

		return op, nil

	case v1:
		v1Resource := &compute.InstanceGroupManager{}
		err := Convert(resource, v1Resource)
		if err != nil {
			return nil, err
		}

		v1Op, err := s.v1.InstanceGroupManagers.Insert(project, zone, v1Resource).Do()
		if err != nil {
			return nil, err
		}

		err = Convert(v1Op, op)
		if err != nil {
			return nil, err
		}

		return op, nil

	}

	return nil, fmt.Errorf("Unknown API version.")
}

func (s *ComputeMultiversionService) GetInstanceGroupManager(project string, zone string, resource string, version ComputeApiVersion) (*computeBeta.InstanceGroupManager, error) {
	res := &computeBeta.InstanceGroupManager{}
	switch version {

	case v0beta:
		r, err := s.v0beta.InstanceGroupManagers.Get(project, zone, resource).Do()
		if err != nil {
			return nil, err
		}

		err = Convert(r, res)
		if err != nil {
			return nil, err
		}

		return res, nil

	case v1:
		r, err := s.v1.InstanceGroupManagers.Get(project, zone, resource).Do()
		if err != nil {
			return nil, err
		}

		err = Convert(r, res)
		if err != nil {
			return nil, err
		}

		return res, nil

	}

	return nil, fmt.Errorf("Unknown API version.")
}

func (s *ComputeMultiversionService) DeleteInstanceGroupManager(project string, zone string, resource string, version ComputeApiVersion) (*computeBeta.Operation, error) {
	op := &computeBeta.Operation{}
	switch version {

	case v0beta:
		v0betaOp, err := s.v0beta.InstanceGroupManagers.Delete(project, zone, resource).Do()
		if err != nil {
			return nil, err
		}

		err = Convert(v0betaOp, op)
		if err != nil {
			return nil, err
		}

		return op, nil

	case v1:
		v1Op, err := s.v1.InstanceGroupManagers.Delete(project, zone, resource).Do()
		if err != nil {
			return nil, err
		}

		err = Convert(v1Op, op)
		if err != nil {
			return nil, err
		}

		return op, nil

	}

	return nil, fmt.Errorf("Unknown API version.")
}

func (s *ComputeMultiversionService) UpdateInstanceGroupManager(project string, zone string, resourceName string, resource *computeBeta.InstanceGroupManager, version ComputeApiVersion) (*computeBeta.Operation, error) {
	op := &computeBeta.Operation{}
	switch version {

	case v0beta:
		v0betaResource := &computeBeta.InstanceGroupManager{}
		err := Convert(resource, v0betaResource)
		if err != nil {
			return nil, err
		}
		v0betaOp, err := s.v0beta.InstanceGroupManagers.Update(project, zone, resourceName, v0betaResource).Do()
		if err != nil {
			return nil, err
		}

		err = Convert(v0betaOp, op)
		if err != nil {
			return nil, err
		}

		return op, nil

	}

	return nil, fmt.Errorf("Unknown API version.")
}

func (s *ComputeMultiversionService) InsertAddress(project string, region string, resource *computeBeta.Address, version ComputeApiVersion) (*computeBeta.Operation, error) {
	op := &computeBeta.Operation{}
	switch version {

	case v0beta:
		v0betaResource := &computeBeta.Address{}
		err := Convert(resource, v0betaResource)
		if err != nil {
			return nil, err
		}

		v0betaOp, err := s.v0beta.Addresses.Insert(project, region, v0betaResource).Do()
		if err != nil {
			return nil, err
		}

		err = Convert(v0betaOp, op)
		if err != nil {
			return nil, err
		}

		return op, nil

	case v1:
		v1Resource := &compute.Address{}
		err := Convert(resource, v1Resource)
		if err != nil {
			return nil, err
		}

		v1Op, err := s.v1.Addresses.Insert(project, region, v1Resource).Do()
		if err != nil {
			return nil, err
		}

		err = Convert(v1Op, op)
		if err != nil {
			return nil, err
		}

		return op, nil

	}

	return nil, fmt.Errorf("Unknown API version.")
}

func (s *ComputeMultiversionService) GetAddress(project string, region string, resource string, version ComputeApiVersion) (*computeBeta.Address, error) {
	res := &computeBeta.Address{}
	switch version {

	case v0beta:
		r, err := s.v0beta.Addresses.Get(project, region, resource).Do()
		if err != nil {
			return nil, err
		}

		err = Convert(r, res)
		if err != nil {
			return nil, err
		}

		return res, nil

	case v1:
		r, err := s.v1.Addresses.Get(project, region, resource).Do()
		if err != nil {
			return nil, err
		}

		err = Convert(r, res)
		if err != nil {
			return nil, err
		}

		return res, nil

	}

	return nil, fmt.Errorf("Unknown API version.")
}

func (s *ComputeMultiversionService) DeleteAddress(project string, region string, resource string, version ComputeApiVersion) (*computeBeta.Operation, error) {
	op := &computeBeta.Operation{}
	switch version {

	case v0beta:
		v0betaOp, err := s.v0beta.Addresses.Delete(project, region, resource).Do()
		if err != nil {
			return nil, err
		}

		err = Convert(v0betaOp, op)
		if err != nil {
			return nil, err
		}

		return op, nil

	case v1:
		v1Op, err := s.v1.Addresses.Delete(project, region, resource).Do()
		if err != nil {
			return nil, err
		}

		err = Convert(v1Op, op)
		if err != nil {
			return nil, err
		}

		return op, nil

	}

	return nil, fmt.Errorf("Unknown API version.")
}
