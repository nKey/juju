// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package tools

import (
	"fmt"

	"launchpad.net/juju-core/environs/simplestreams"
	"launchpad.net/juju-core/version"
)

// ToolsMetadataLookupParams is used to query metadata for matching tools.
type ToolsMetadataLookupParams struct {
	simplestreams.MetadataLookupParams
	Version string
}

// ValidateToolsMetadata attempts to load tools metadata for the specified cloud attributes and returns
// any tools versions found, or an error if the metadata could not be loaded.
func ValidateToolsMetadata(params *ToolsMetadataLookupParams) ([]string, error) {
	if params.Region == "" {
		return nil, fmt.Errorf("required parameter region not specified")
	}
	if params.Endpoint == "" {
		return nil, fmt.Errorf("required parameter endpoint not specified")
	}
	if len(params.Architectures) == 0 {
		return nil, fmt.Errorf("required parameter arches not specified")
	}
	if len(params.BaseURLs) == 0 {
		return nil, fmt.Errorf("required parameter baseURLs not specified")
	}
	if params.Version == "" {
		params.Version = version.CurrentNumber().String()
	}
	toolsConstraint := NewToolsConstraint(params.Version, simplestreams.LookupParams{
		CloudSpec: simplestreams.CloudSpec{params.Region, params.Endpoint},
		Series:    params.Series,
		Arches:    params.Architectures,
	})
	matchingTools, err := Fetch(params.BaseURLs, simplestreams.DefaultIndexPath, toolsConstraint, false)
	if err != nil {
		return nil, err
	}
	if len(matchingTools) == 0 {
		return nil, fmt.Errorf("no matching tools found for constraint %+v", toolsConstraint)
	}
	versions := make([]string, len(matchingTools))
	for i, tm := range matchingTools {
		vers := version.Binary{version.MustParse(tm.Version), tm.Release, tm.Arch}
		versions[i] = vers.String()
	}
	return versions, nil
}