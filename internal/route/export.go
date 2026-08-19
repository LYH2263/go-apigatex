package route

import "encoding/json"

// ExportJSON 导出路由表 JSON。
func (t *Table) ExportJSON() ([]byte, error) {
	rows := t.List()
	type row struct {
		ID        string            `json:"id"`
		Path      string            `json:"path"`
		Kind      Kind              `json:"kind"`
		Methods   []string          `json:"methods"`
		Upstream  string            `json:"upstream"`
		Headers   map[string]string `json:"headers"`
		StripPath string            `json:"strip_path"`
		Auth      bool              `json:"auth"`
		Priority  int               `json:"priority"`
	}
	out := make([]row, 0, len(rows))
	for _, r := range rows {
		out = append(out, row{
			ID: r.ID, Path: r.Path, Kind: r.Kind, Methods: CloneStrings(r.Methods),
			Upstream: r.Upstream, Headers: CloneHeaders(r.Headers),
			StripPath: r.StripPath, Auth: r.Auth, Priority: r.Priority,
		})
	}
	return json.Marshal(out)
}

// ImportJSON 从 JSON 导入（替换）。
func (t *Table) ImportJSON(raw []byte) error {
	type row struct {
		ID        string            `json:"id"`
		Path      string            `json:"path"`
		Kind      Kind              `json:"kind"`
		Methods   []string          `json:"methods"`
		Upstream  string            `json:"upstream"`
		Headers   map[string]string `json:"headers"`
		StripPath string            `json:"strip_path"`
		Auth      bool              `json:"auth"`
		Priority  int               `json:"priority"`
	}
	var rows []row
	if err := json.Unmarshal(raw, &rows); err != nil {
		return err
	}
	converted := make([]Route, 0, len(rows))
	for _, r := range rows {
		converted = append(converted, Route{
			ID: r.ID, Path: r.Path, Kind: r.Kind, Methods: CloneStrings(r.Methods),
			Upstream: r.Upstream, Headers: CloneHeaders(r.Headers),
			StripPath: r.StripPath, Auth: r.Auth, Priority: r.Priority,
		})
	}
	t.Replace(converted)
	return nil
}
