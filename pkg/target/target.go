/*
SPDX-FileCopyrightText: 2021 SAP SE or an SAP affiliate company and Gardener contributors

SPDX-License-Identifier: Apache-2.0
*/

package target

import (
	"errors"

	gardencore "github.com/gardener/gardener/pkg/apis/core"
	"sigs.k8s.io/controller-runtime/pkg/client"

	clientgarden "github.com/gardener/gardenctl-v2/internal/client/garden"
)

/*
	Targets are either

	garden -> seed -> shoot
	garden -> project -> shoot
*/

// Target represents the Kubernetes cluster/namespace which should
// be the target for user operations in gardenctl. It works similar
// to the context defined in a kubeconfig.
type Target interface {
	// GardenName returns the currently targeted garden cluster name.
	GardenName() string
	// ProjectName returns the currently targeted project name.
	ProjectName() string
	// SeedName returns the currently targeted seed cluster name.
	SeedName() string
	// ShootName returns the currently targeted shoot cluster name.
	ShootName() string
	// ControlPlane returns true if shoot control plane is targeted.
	ControlPlane() bool
	// WithGardenName returns a copy of the target with the garden name updated.
	// The returned target can be invalid.
	WithGardenName(name string) Target
	// WithProjectName returns a copy of the target with the project name updated.
	// The returned target can be invalid.
	WithProjectName(name string) Target
	// WithSeedName returns a copy of the target with the seed name updated.
	// The returned target can be invalid.
	WithSeedName(name string) Target
	// WithShootName returns a copy of the target with the shoot name updated.
	// The returned target can be invalid.
	WithShootName(name string) Target
	// WithControlPlane returns a copy of the target with the control plane flag updated.
	// The returned target can be invalid.
	WithControlPlane(controlPlane bool) Target
	// Validate checks for semantical correctness of the target, without
	// actually connecting to the targeted clusters.
	Validate() error
	// IsEmpty returns true if all values of the target are empty
	IsEmpty() bool
	// AsListOption returns the target as list option
	AsListOption() client.ListOption

	// DeepCopy returns a deep copy of the target
	DeepCopy() Target
}

type targetImpl struct {
	Garden           string `json:"garden,omitempty"`
	Project          string `json:"project,omitempty"`
	Seed             string `json:"seed,omitempty"`
	Shoot            string `json:"shoot,omitempty"`
	ControlPlaneFlag bool   `json:"controlPlane,omitempty"`
}

var _ Target = &targetImpl{}

// NewTarget returns a new target. This function does not perform any validation,
// so the returned target can be invalid.
func NewTarget(gardenName, projectName, seedName, shootName string) Target {
	return &targetImpl{
		Garden:           gardenName,
		Project:          projectName,
		Seed:             seedName,
		Shoot:            shootName,
		ControlPlaneFlag: false,
	}
}

func newTargetImpl(gardenName, projectName, seedName, shootName string, controlPlane bool) *targetImpl {
	return &targetImpl{Garden: gardenName, Project: projectName, Seed: seedName, Shoot: shootName, ControlPlaneFlag: controlPlane}
}

// Validate checks that the target is not malformed and all required
// steps are configured correctly.
func (t *targetImpl) Validate() error {
	if len(t.Project) > 0 && len(t.Seed) > 0 {
		return errors.New("seed and project must not be configured at the same time")
	}

	return nil
}

// GardenName returns the currently targeted garden cluster name.
func (t *targetImpl) GardenName() string {
	return t.Garden
}

// ProjectName returns the currently targeted project name.
func (t *targetImpl) ProjectName() string {
	return t.Project
}

// SeedName returns the currently targeted seed cluster name.
func (t *targetImpl) SeedName() string {
	return t.Seed
}

// ShootName returns the currently targeted shoot cluster name.
func (t *targetImpl) ShootName() string {
	return t.Shoot
}

// ControlPlane returns true if shoot control plane is targeted.
func (t *targetImpl) ControlPlane() bool {
	return t.ControlPlaneFlag
}

// WithGardenName returns a copy of the target with the garden name updated.
// The returned target can be invalid.
func (t *targetImpl) WithGardenName(name string) Target {
	return newTargetImpl(name, t.Project, t.Seed, t.Shoot, t.ControlPlaneFlag)
}

// WithProjectName returns a copy of the target with the project name updated.
// The returned target can be invalid.
func (t *targetImpl) WithProjectName(name string) Target {
	return newTargetImpl(t.Garden, name, t.Seed, t.Shoot, t.ControlPlaneFlag)
}

// WithSeedName returns a copy of the target with the seed name updated.
// The returned target can be invalid.
func (t *targetImpl) WithSeedName(name string) Target {
	return newTargetImpl(t.Garden, t.Project, name, t.Shoot, t.ControlPlaneFlag)
}

// WithShootName returns a copy of the target with the shoot name updated.
// The returned target can be invalid.
func (t *targetImpl) WithShootName(name string) Target {
	return newTargetImpl(t.Garden, t.Project, t.Seed, name, t.ControlPlaneFlag)
}

// WithControlPlane returns a copy of the target with the control plane flag updated.
// The returned target can be invalid.
func (t *targetImpl) WithControlPlane(controlPlane bool) Target {
	return newTargetImpl(t.Garden, t.Project, t.Seed, t.Shoot, controlPlane)
}

func (t *targetImpl) IsEmpty() bool {
	return t.Garden == "" && t.Project == "" && t.Seed == "" && t.Shoot == ""
}

func (t *targetImpl) AsListOption() client.ListOption {
	opt := clientgarden.ProjectFilter{}

	if t.ShootName() != "" {
		opt["metadata.name"] = t.ShootName()
	}

	if t.ProjectName() != "" {
		opt["project"] = t.ProjectName()
	} else if t.SeedName() != "" {
		opt[gardencore.ShootSeedName] = t.SeedName()
	}

	return opt
}

func (t *targetImpl) DeepCopy() Target {
	return &targetImpl{
		Garden:           t.Garden,
		Project:          t.Project,
		Seed:             t.Seed,
		Shoot:            t.Shoot,
		ControlPlaneFlag: t.ControlPlaneFlag,
	}
}
