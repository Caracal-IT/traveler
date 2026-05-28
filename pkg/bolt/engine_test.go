package bolt

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRender_JSONStringIn_JSONStringOut(t *testing.T) {
	out, err := Render(Request[string]{
		Template: `{"message":"Hello {{.name}}","count":{{.count}}}`,
		Payload:  `{"name":"Bolt","count":2}`,
	})
	require.NoError(t, err)
	require.JSONEq(t, `{"message":"Hello Bolt","count":2}`, out)
}

func TestRenderAs_ObjectIn_ObjectOut(t *testing.T) {
	type Input struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	type Output struct {
		Greeting string `json:"greeting"`
		Adult    bool   `json:"adult"`
	}

	out, err := RenderAs[Output](Request[Input]{
		Template: `{"greeting":"Hi {{.name}}","adult":{{if ge .age 18}}true{{else}}false{{end}}}`,
		Payload:  Input{Name: "Ada", Age: 31},
	})
	require.NoError(t, err)
	require.Equal(t, Output{Greeting: "Hi Ada", Adult: true}, out)
}

func TestRenderAs_MixedModelAndPayload(t *testing.T) {
	type Model struct {
		Tenant string `json:"tenant"`
	}
	type Output struct {
		Source string `json:"source"`
		Ref    string `json:"ref"`
	}

	out, err := RenderAs[Output](Request[string]{
		Template: `{"source":"{{.model.tenant}}","ref":"{{.json.user.id}}-{{.tenant}}"}`,
		Model:    Model{Tenant: "caracal"},
		Payload:  `{"user":{"id":"u-99"},"tenant":"override"}`,
	})
	require.NoError(t, err)
	require.Equal(t, Output{Source: "caracal", Ref: "u-99-override"}, out)
}

func TestRender_TextFormat(t *testing.T) {
	out, err := Render(Request[string]{
		Template: `User {{.name}} has role {{.role}}`,
		Format:   FormatText,
		Payload:  `{"name":"Bolt","role":"admin"}`,
	})
	require.NoError(t, err)
	require.Equal(t, "User Bolt has role admin", out)
}

func TestRender_HTMLFormatEscapes(t *testing.T) {
	out, err := Render(Request[string]{
		Template: `<p>{{.name}}</p>`,
		Format:   FormatHTML,
		Payload:  `{"name":"<b>Bolt</b>"}`,
	})
	require.NoError(t, err)
	require.Equal(t, `<p>&lt;b&gt;Bolt&lt;/b&gt;</p>`, out)
}

func TestRenderAs_TextFormat_ObjectOutputRejected(t *testing.T) {
	type Output struct {
		Name string `json:"name"`
	}

	_, err := RenderAs[Output](Request[string]{
		Template: `hello {{.name}}`,
		Format:   FormatText,
		Payload:  `{"name":"bolt"}`,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "requires string or []byte output")
}

func TestRender_RejectsAbsoluteTemplateDir(t *testing.T) {
	_, err := Render(Request[map[string]any]{
		TemplateName: "profile.json.tmpl",
		TemplateDir:  t.TempDir(),
		Payload:      map[string]any{"name": "Bolt"},
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "template_dir must be a relative path")
}

func TestRender_RejectsAbsoluteTemplateName(t *testing.T) {
	absoluteTemplateName := filepath.Join(t.TempDir(), "profile.json.tmpl")
	_, err := Render(Request[map[string]any]{
		TemplateName: absoluteTemplateName,
		Payload:      map[string]any{"name": "Bolt"},
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "template_name must be a relative path")
}

func TestRender_TemplateNameWithTemplateDir(t *testing.T) {
	dir, err := os.MkdirTemp(".", "bolt-test-*")
	require.NoError(t, err)
	defer func() { _ = os.RemoveAll(dir) }()

	templatePath := filepath.Join(dir, "welcome.json.tmpl")
	err = os.WriteFile(templatePath, []byte(`{"msg":"Hello {{.name}}"}`), 0o644)
	require.NoError(t, err)

	out, err := Render(Request[string]{
		TemplateName: "welcome.json.tmpl",
		TemplateDir:  dir,
		Payload:      `{"name":"Bolt"}`,
	})
	require.NoError(t, err)
	require.JSONEq(t, `{"msg":"Hello Bolt"}`, out)
}

func TestRender_ConfigFile_AppliesTemplateDir(t *testing.T) {
	baseDir, err := os.MkdirTemp(".", "bolt-config-*")
	require.NoError(t, err)
	defer func() { _ = os.RemoveAll(baseDir) }()

	templatesDir := filepath.Join(baseDir, "templates")
	err = os.MkdirAll(templatesDir, 0o755)
	require.NoError(t, err)

	templatePath := filepath.Join(templatesDir, "welcome.json.tmpl")
	err = os.WriteFile(templatePath, []byte(`{"msg":"Hello {{.name}}"}`), 0o644)
	require.NoError(t, err)

	configPath := filepath.Join(baseDir, "engine.json")
	err = os.WriteFile(configPath, []byte(`{"template_dir":"`+templatesDir+`"}`), 0o644)
	require.NoError(t, err)

	out, err := Render(Request[string]{
		ConfigFile:   configPath,
		TemplateName: "welcome.json.tmpl",
		Payload:      `{"name":"Bolt"}`,
	})
	require.NoError(t, err)
	require.JSONEq(t, `{"msg":"Hello Bolt"}`, out)
}

func TestRender_ConfigFile_InvalidConfig(t *testing.T) {
	configFile, err := os.CreateTemp(".", "bolt-engine-invalid-*.json")
	require.NoError(t, err)
	defer func() {
		_ = configFile.Close()
		_ = os.Remove(configFile.Name())
	}()

	_, err = configFile.WriteString(`{invalid-json`)
	require.NoError(t, err)

	_, err = Render(Request[string]{
		ConfigFile:   configFile.Name(),
		TemplateName: "welcome.json.tmpl",
		Payload:      `{"name":"Bolt"}`,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "load engine config")
}

func TestRender_ConfigFile_RelativePathUsesCurrentDirectory(t *testing.T) {
	originalWD, err := os.Getwd()
	require.NoError(t, err)

	baseDir := t.TempDir()
	err = os.Chdir(baseDir)
	require.NoError(t, err)
	defer func() { _ = os.Chdir(originalWD) }()

	err = os.MkdirAll("templates", 0o755)
	require.NoError(t, err)

	err = os.WriteFile("templates/welcome.json.tmpl", []byte(`{"msg":"Hello {{.name}}"}`), 0o644)
	require.NoError(t, err)
	err = os.WriteFile("engine.json", []byte(`{"template_dir":"templates"}`), 0o644)
	require.NoError(t, err)

	out, err := Render(Request[string]{
		ConfigFile:   "engine.json",
		TemplateName: "welcome.json.tmpl",
		Payload:      `{"name":"Bolt"}`,
	})
	require.NoError(t, err)
	require.JSONEq(t, `{"msg":"Hello Bolt"}`, out)
}

func BenchmarkRenderJSON(b *testing.B) {
	request := Request[map[string]any]{
		Template: `{"id":"{{.id}}","name":"{{.name}}","meta":{{toJSON .meta}},"ok":true}`,
		Payload: map[string]any{
			"id":   "x-1",
			"name": "bolt",
			"meta": map[string]any{"region": "za", "version": 1},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := RenderAs[[]byte](request)
		if err != nil {
			b.Fatal(err)
		}
	}
}
