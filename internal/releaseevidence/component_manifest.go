package releaseevidence

import (
	"encoding/json"
	"fmt"
)

func decodeComponentManifest(data []byte) (componentManifest, error) {
	schema := requirements.ComponentManifest
	rootFields := []string{schema.RootFields.SchemaVersion, schema.RootFields.Component, schema.RootFields.Files}
	componentFields := []string{schema.ComponentFields.Name, schema.ComponentFields.Version, schema.ComponentFields.Target}
	targetFields := []string{schema.TargetFields.OS, schema.TargetFields.Arch}
	fileFields := []string{schema.FileFields.Path, schema.FileFields.Mode, schema.FileFields.Size, schema.FileFields.SHA256}
	root, err := decodeExactObject(data, rootFields)
	if err != nil {
		return componentManifest{}, err
	}
	var manifest componentManifest
	if err := json.Unmarshal(root[schema.RootFields.SchemaVersion], &manifest.SchemaVersion); err != nil {
		return componentManifest{}, err
	}
	component, err := decodeExactObject(root[schema.RootFields.Component], componentFields)
	if err != nil {
		return componentManifest{}, err
	}
	if err := json.Unmarshal(component[schema.ComponentFields.Name], &manifest.Component.Name); err != nil {
		return componentManifest{}, err
	}
	if err := json.Unmarshal(component[schema.ComponentFields.Version], &manifest.Component.Version); err != nil {
		return componentManifest{}, err
	}
	target, err := decodeExactObject(component[schema.ComponentFields.Target], targetFields)
	if err != nil {
		return componentManifest{}, err
	}
	if err := json.Unmarshal(target[schema.TargetFields.OS], &manifest.Component.Target.OS); err != nil {
		return componentManifest{}, err
	}
	if err := json.Unmarshal(target[schema.TargetFields.Arch], &manifest.Component.Target.Arch); err != nil {
		return componentManifest{}, err
	}
	var files []json.RawMessage
	if err := json.Unmarshal(root[schema.RootFields.Files], &files); err != nil {
		return componentManifest{}, err
	}
	manifest.Files = make([]manifestFile, 0, len(files))
	for _, raw := range files {
		object, err := decodeExactObject(raw, fileFields)
		if err != nil {
			return componentManifest{}, err
		}
		var file manifestFile
		values := []struct {
			field string
			value any
		}{{schema.FileFields.Path, &file.Path}, {schema.FileFields.Mode, &file.Mode}, {schema.FileFields.Size, &file.Size}, {schema.FileFields.SHA256, &file.SHA256}}
		for _, item := range values {
			if err := json.Unmarshal(object[item.field], item.value); err != nil {
				return componentManifest{}, err
			}
		}
		manifest.Files = append(manifest.Files, file)
	}
	return manifest, nil
}

func decodeExactObject(data []byte, fields []string) (map[string]json.RawMessage, error) {
	if err := rejectDuplicateJSONKeys(data); err != nil {
		return nil, err
	}
	var object map[string]json.RawMessage
	if err := json.Unmarshal(data, &object); err != nil {
		return nil, err
	}
	if len(object) != len(fields) {
		return nil, fmt.Errorf("component manifest fields do not match schema")
	}
	for _, field := range fields {
		if _, ok := object[field]; !ok {
			return nil, fmt.Errorf("component manifest field %s is missing", field)
		}
	}
	return object, nil
}

func canonicalManifestFiles(files []manifestFile) ([]byte, error) {
	values := make([]map[string]any, 0, len(files))
	for _, file := range files {
		fields := requirements.ComponentManifest.FileFields
		values = append(values, map[string]any{fields.Path: file.Path, fields.Mode: file.Mode, fields.Size: file.Size, fields.SHA256: file.SHA256})
	}
	return canonicalJSON(values)
}
