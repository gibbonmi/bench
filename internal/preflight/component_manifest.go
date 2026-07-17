package preflight

import (
	"encoding/json"
	"fmt"
)

func decodeComponentManifest(data []byte) (componentManifest, error) {
	schema := requirements.ComponentManifest
	root, err := decodeExactObject(data, schema.RootFields)
	if err != nil {
		return componentManifest{}, err
	}
	var manifest componentManifest
	if err := json.Unmarshal(root[schema.RootFields[0]], &manifest.SchemaVersion); err != nil {
		return componentManifest{}, err
	}
	component, err := decodeExactObject(root[schema.RootFields[1]], schema.ComponentFields)
	if err != nil {
		return componentManifest{}, err
	}
	if err := json.Unmarshal(component[schema.ComponentFields[0]], &manifest.Component.Name); err != nil {
		return componentManifest{}, err
	}
	if err := json.Unmarshal(component[schema.ComponentFields[1]], &manifest.Component.Version); err != nil {
		return componentManifest{}, err
	}
	target, err := decodeExactObject(component[schema.ComponentFields[2]], schema.TargetFields)
	if err != nil {
		return componentManifest{}, err
	}
	if err := json.Unmarshal(target[schema.TargetFields[0]], &manifest.Component.Target.OS); err != nil {
		return componentManifest{}, err
	}
	if err := json.Unmarshal(target[schema.TargetFields[1]], &manifest.Component.Target.Arch); err != nil {
		return componentManifest{}, err
	}
	var files []json.RawMessage
	if err := json.Unmarshal(root[schema.RootFields[2]], &files); err != nil {
		return componentManifest{}, err
	}
	manifest.Files = make([]manifestFile, 0, len(files))
	for _, raw := range files {
		object, err := decodeExactObject(raw, schema.FileFields)
		if err != nil {
			return componentManifest{}, err
		}
		var file manifestFile
		values := []any{&file.Path, &file.Mode, &file.Size, &file.SHA256}
		for i, value := range values {
			if err := json.Unmarshal(object[schema.FileFields[i]], value); err != nil {
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
	fields := requirements.ComponentManifest.FileFields
	values := make([]map[string]any, 0, len(files))
	for _, file := range files {
		values = append(values, map[string]any{fields[0]: file.Path, fields[1]: file.Mode, fields[2]: file.Size, fields[3]: file.SHA256})
	}
	return canonicalJSON(values)
}
