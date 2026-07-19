package publication

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// FixtureRegistry drives the hermetic offline-registry.mjs test fixture over
// HTTP. It is the adapter the gate exercises: no credential, no network egress
// beyond the given base URL, no claim about public npm's real behavior. base is
// a BENCH_RELEASE_REGISTRY-style URL, e.g. http://127.0.0.1:PORT.
type FixtureRegistry struct {
	Base   string
	Client *http.Client
}

func NewFixtureRegistry(base string) *FixtureRegistry {
	return &FixtureRegistry{Base: base, Client: http.DefaultClient}
}

func (f *FixtureRegistry) client() *http.Client {
	if f.Client != nil {
		return f.Client
	}
	return http.DefaultClient
}

func packageFileName(name, version string) string {
	if strings.HasPrefix(name, "@redbench/") {
		target := strings.TrimPrefix(name, "@redbench/")
		return fmt.Sprintf("redbench-%s-%s.tgz", target, version)
	}
	return fmt.Sprintf("redbench-%s.tgz", version)
}

func encodedPackageName(name string) string {
	if strings.HasPrefix(name, "@") {
		return url.PathEscape(name)
	}
	return name
}

func (f *FixtureRegistry) do(ctx context.Context, method, path string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, f.Base+path, body)
	if err != nil {
		return nil, err
	}
	return f.client().Do(req)
}

func (f *FixtureRegistry) Publish(ctx context.Context, name, version, tag string, tarball []byte) (string, error) {
	file := packageFileName(name, version)
	response, err := f.do(ctx, http.MethodPut, "/upload/"+url.PathEscape(file)+"?tag="+url.QueryEscape(tag), bytes.NewReader(tarball))
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusCreated {
		data, _ := io.ReadAll(response.Body)
		return "", fmt.Errorf("registry publish failed for %s@%s: %d %s", name, version, response.StatusCode, data)
	}
	integrity, live, err := f.Integrity(ctx, name, version)
	if err != nil {
		return "", err
	}
	if !live {
		return "", fmt.Errorf("registry publish did not make %s@%s live", name, version)
	}
	return integrity, nil
}

func (f *FixtureRegistry) StageSubmit(ctx context.Context, name, version string, tarball []byte) (string, error) {
	file := packageFileName(name, version)
	response, err := f.do(ctx, http.MethodPut, "/-/stage/"+url.PathEscape(file), bytes.NewReader(tarball))
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusCreated {
		data, _ := io.ReadAll(response.Body)
		return "", fmt.Errorf("registry stage-submit failed for %s@%s: %d %s", name, version, response.StatusCode, data)
	}
	var payload struct {
		StageID string `json:"stage_id"`
	}
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		return "", fmt.Errorf("registry stage-submit response is malformed: %w", err)
	}
	return payload.StageID, nil
}

func (f *FixtureRegistry) Approve(ctx context.Context, stageID string) error {
	response, err := f.do(ctx, http.MethodPost, "/-/approve/"+url.PathEscape(stageID), nil)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		data, _ := io.ReadAll(response.Body)
		return fmt.Errorf("registry approve failed for stage %s: %d %s", stageID, response.StatusCode, data)
	}
	return nil
}

func (f *FixtureRegistry) Integrity(ctx context.Context, name, version string) (string, bool, error) {
	path := "/-/integrity/" + encodedPackageName(name) + "/" + url.PathEscape(version)
	response, err := f.do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return "", false, err
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusNotFound {
		return "", false, nil
	}
	if response.StatusCode != http.StatusOK {
		data, _ := io.ReadAll(response.Body)
		return "", false, fmt.Errorf("registry integrity query failed for %s@%s: %d %s", name, version, response.StatusCode, data)
	}
	var payload struct {
		Integrity string `json:"integrity"`
	}
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		return "", false, fmt.Errorf("registry integrity response is malformed: %w", err)
	}
	integrity := strings.TrimSpace(payload.Integrity)
	if integrity == "" {
		// The version exists (200 OK) but reports no integrity — a
		// malformed/hostile registry response, not a live/not-live fact.
		return "", false, fmt.Errorf("registry reports %s@%s live with no integrity value", name, version)
	}
	return integrity, true, nil
}

func (f *FixtureRegistry) TagAdd(ctx context.Context, name, tag, version string) error {
	encoded, _ := json.Marshal(version)
	response, err := f.do(ctx, http.MethodPut, "/-/package/"+encodedPackageName(name)+"/dist-tags/"+url.PathEscape(tag), bytes.NewReader(encoded))
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		data, _ := io.ReadAll(response.Body)
		return fmt.Errorf("registry dist-tag add failed for %s %s: %d %s", name, tag, response.StatusCode, data)
	}
	return nil
}

func (f *FixtureRegistry) TagRemove(ctx context.Context, name, tag string) error {
	response, err := f.do(ctx, http.MethodDelete, "/-/package/"+encodedPackageName(name)+"/dist-tags/"+url.PathEscape(tag), nil)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		data, _ := io.ReadAll(response.Body)
		return fmt.Errorf("registry dist-tag remove failed for %s %s: %d %s", name, tag, response.StatusCode, data)
	}
	return nil
}

func (f *FixtureRegistry) Deprecate(ctx context.Context, name, version, message string) error {
	body, _ := json.Marshal(map[string]string{"version": version, "message": message})
	response, err := f.do(ctx, http.MethodPost, "/-/package/"+encodedPackageName(name)+"/deprecate", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		data, _ := io.ReadAll(response.Body)
		return fmt.Errorf("registry deprecate failed for %s@%s: %d %s", name, version, response.StatusCode, data)
	}
	return nil
}
