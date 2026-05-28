package bolt

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	htmltmpl "html/template"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	texttmpl "text/template"

	appconfig "traveler/pkg/config" /**/
)

type engine struct {
	templateFuncs map[string]any
	templateDir   string
}

type engineConfig struct {
	TemplateDir string `json:"template_dir" mapstructure:"template_dir"`
}

// Format defines template rendering mode.
type Format string

const (
	FormatJSON Format = "json"
	FormatText Format = "text"
	FormatHTML Format = "html"
)

// Request represents a typed render request.
//
// Use Template directly for inline templates, or TemplateName to resolve from file.
// When both are provided, Template is used.
// If Format is empty, json is used.
// ConfigFile is optional and can provide default template_dir.
// TemplateDir overrides template_dir from ConfigFile when both are set.
type Request[T any] struct {
	TemplateName string
	Template     string
	Format       Format
	ConfigFile   string
	TemplateDir  string
	Payload      T
	Model        any
}

type renderRequest interface {
	templateName() string
	templateInline() string
	format() Format
	configFile() string
	templateDir() string
	payload() any
	model() any
}

func (r Request[T]) templateName() string   { return r.TemplateName }
func (r Request[T]) templateInline() string { return r.Template }
func (r Request[T]) format() Format         { return r.Format }
func (r Request[T]) configFile() string     { return r.ConfigFile }
func (r Request[T]) templateDir() string    { return r.TemplateDir }
func (r Request[T]) payload() any           { return r.Payload }
func (r Request[T]) model() any             { return r.Model }

func newEngine() *engine {
	return &engine{
		templateFuncs: map[string]any{
			"toJSON": toJSON,
		},
	}
}

func buildEngineFromRequest(request renderRequest) (*engine, error) {
	e := newEngine()

	if request.configFile() != "" {
		cfg, err := loadConfigFile(request.configFile())
		if err != nil {
			return nil, err
		}
		e.templateDir = cfg.TemplateDir
	}

	if request.templateDir() != "" {
		e.templateDir = request.templateDir()
	}

	if e.templateDir == "" {
		e.templateDir = "."
	}

	return e, nil
}

func loadConfigFile(path string) (engineConfig, error) {
	resolvedPath, err := resolveConfigPath(path)
	if err != nil {
		return engineConfig{}, err
	}

	var cfg engineConfig
	if err := appconfig.LoadInto(resolvedPath, &cfg); err != nil {
		return engineConfig{}, fmt.Errorf("load engine config: %w", err)
	}
	return cfg, nil
}

func resolveConfigPath(path string) (string, error) {
	if path == "" {
		return "", errors.New("config path is required")
	}
	if filepath.IsAbs(path) {
		return path, nil
	}

	return filepath.Abs(path)
}

// RenderAs renders a request into the requested generic output type.
//
// Output type behavior:
// - string: returns rendered bytes as string
// - []byte: returns rendered bytes
// - any other type: requires json format and unmarshals into that type
func RenderAs[TOut any](request renderRequest) (TOut, error) {
	var out TOut

	e, err := buildEngineFromRequest(request)
	if err != nil {
		return out, err
	}

	format, err := resolveFormat(request.format())
	if err != nil {
		return out, err
	}

	templateBody, err := e.resolveTemplate(request.templateInline(), request.templateName())
	if err != nil {
		return out, err
	}

	ctx, err := buildContext(request.model(), request.payload())
	if err != nil {
		return out, err
	}

	rendered, err := e.renderByFormat(format, templateBody, ctx)
	if err != nil {
		return out, err
	}

	return decodeOutput[TOut](format, rendered)
}

// Render renders and returns string output.
func Render(request renderRequest) (string, error) {
	return RenderAs[string](request)
}

func resolveFormat(format Format) (Format, error) {
	if format == "" {
		return FormatJSON, nil
	}
	switch format {
	case FormatJSON, FormatText, FormatHTML:
		return format, nil
	default:
		return "", fmt.Errorf("invalid format %q; use json, text, or html", format)
	}
}

func decodeOutput[TOut any](format Format, rendered []byte) (TOut, error) {
	var out TOut
	outType := reflect.TypeOf((*TOut)(nil)).Elem()

	switch outType.Kind() {
	case reflect.String:
		reflect.ValueOf(&out).Elem().SetString(string(rendered))
		return out, nil
	case reflect.Slice:
		if outType.Elem().Kind() == reflect.Uint8 {
			reflect.ValueOf(&out).Elem().SetBytes(rendered)
			return out, nil
		}
		fallthrough
	default:
		if format != FormatJSON {
			return out, fmt.Errorf("format %q requires string or []byte output type", format)
		}
		if err := json.Unmarshal(rendered, &out); err != nil {
			return out, fmt.Errorf("unmarshal rendered json: %w", err)
		}
		return out, nil
	}
}

func (e *engine) renderByFormat(format Format, tmpl string, ctx map[string]any) ([]byte, error) {
	switch format {
	case FormatJSON:
		rendered, err := e.renderTextTemplate(tmpl, ctx)
		if err != nil {
			return nil, err
		}
		compact := bytes.NewBuffer(nil)
		if err := json.Compact(compact, rendered); err != nil {
			return nil, fmt.Errorf("template result is not valid JSON: %w", err)
		}
		return compact.Bytes(), nil
	case FormatText:
		return e.renderTextTemplate(tmpl, ctx)
	case FormatHTML:
		return e.renderHTMLTemplate(tmpl, ctx)
	default:
		return nil, fmt.Errorf("unsupported format %q", format)
	}
}

func (e *engine) resolveTemplate(inlineTemplate string, templateName string) (string, error) {
	if inlineTemplate != "" {
		return inlineTemplate, nil
	}
	if templateName == "" {
		return "", errors.New("template or template name is required")
	}
	if filepath.IsAbs(templateName) {
		return "", errors.New("template_name must be a relative path")
	}
	path := filepath.Join(e.templateDir, templateName)
	path = filepath.Clean(path)

	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read template file %q: %w", path, err)
	}
	return string(data), nil
}

func buildContext(model any, payload any) (map[string]any, error) {
	ctx := map[string]any{}
	merged := map[string]any{}

	if model != nil {
		modelAny, modelMap, err := normalizePayload(model)
		if err != nil {
			return nil, fmt.Errorf("normalize model: %w", err)
		}
		if len(modelMap) > 0 {
			ctx["model"] = modelMap
		} else {
			ctx["model"] = modelAny
		}
		for k, v := range modelMap {
			merged[k] = v
		}
	}

	if payload != nil {
		payloadAny, payloadMap, err := normalizePayload(payload)
		if err != nil {
			return nil, fmt.Errorf("normalize payload: %w", err)
		}
		ctx["json"] = payloadAny
		for k, v := range payloadMap {
			merged[k] = v
		}
	}

	for k, v := range merged {
		ctx[k] = v
	}
	return ctx, nil
}

func normalizePayload(payload any) (any, map[string]any, error) {
	switch p := payload.(type) {
	case map[string]any:
		return p, cloneMap(p), nil
	case []byte:
		return decodeJSONBytes(p)
	case string:
		return decodeJSONBytes([]byte(p))
	default:
		if structMap, ok, err := toStructMap(payload); ok {
			if err != nil {
				return nil, nil, err
			}
			return payload, structMap, nil
		}

		b, err := json.Marshal(p)
		if err != nil {
			return nil, nil, fmt.Errorf("marshal payload: %w", err)
		}
		return decodeJSONBytes(b)
	}
}

func decodeJSONBytes(data []byte) (any, map[string]any, error) {
	var anyValue any
	if err := json.Unmarshal(data, &anyValue); err != nil {
		return nil, nil, fmt.Errorf("unmarshal payload as json: %w", err)
	}

	if obj, ok := anyValue.(map[string]any); ok {
		return anyValue, cloneMap(obj), nil
	}
	return anyValue, map[string]any{}, nil
}

func cloneMap(in map[string]any) map[string]any {
	out := make(map[string]any, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

func toJSON(v any) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (e *engine) renderTextTemplate(tmpl string, ctx map[string]any) ([]byte, error) {
	parsed, err := texttmpl.New("bolt").Funcs(texttmpl.FuncMap(e.templateFuncs)).Option("missingkey=error").Parse(tmpl)
	if err != nil {
		return nil, fmt.Errorf("parse template: %w", err)
	}

	var out bytes.Buffer
	if err := parsed.Execute(&out, ctx); err != nil {
		return nil, fmt.Errorf("execute template: %w", err)
	}
	return out.Bytes(), nil
}

func (e *engine) renderHTMLTemplate(tmpl string, ctx map[string]any) ([]byte, error) {
	parsed, err := htmltmpl.New("bolt").Funcs(htmltmpl.FuncMap(e.templateFuncs)).Option("missingkey=error").Parse(tmpl)
	if err != nil {
		return nil, fmt.Errorf("parse template: %w", err)
	}

	var out bytes.Buffer
	if err := parsed.Execute(&out, ctx); err != nil {
		return nil, fmt.Errorf("execute template: %w", err)
	}
	return out.Bytes(), nil
}

func toStructMap(payload any) (map[string]any, bool, error) {
	value := reflect.ValueOf(payload)
	if !value.IsValid() {
		return nil, false, nil
	}

	for value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return nil, true, errors.New("nil struct pointer payload")
		}
		value = value.Elem()
	}

	if value.Kind() != reflect.Struct {
		return nil, false, nil
	}

	mapped, err := valueToTemplateValue(value)
	if err != nil {
		return nil, true, err
	}

	asMap, ok := mapped.(map[string]any)
	if !ok {
		return nil, true, errors.New("struct payload could not be converted to object")
	}
	return asMap, true, nil
}

func valueToTemplateValue(v reflect.Value) (any, error) {
	for v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return nil, nil
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Struct:
		t := v.Type()
		out := map[string]any{}
		for i := 0; i < v.NumField(); i++ {
			field := t.Field(i)
			if field.PkgPath != "" {
				continue
			}

			name := field.Name
			tag := field.Tag.Get("json")
			if tag != "" {
				parts := strings.Split(tag, ",")
				if parts[0] == "-" {
					continue
				}
				if parts[0] != "" {
					name = parts[0]
				}
			}

			fieldValue, err := valueToTemplateValue(v.Field(i))
			if err != nil {
				return nil, err
			}
			out[name] = fieldValue
		}
		return out, nil
	case reflect.Map:
		if v.Type().Key().Kind() != reflect.String {
			return nil, errors.New("map keys must be strings")
		}
		out := map[string]any{}
		iter := v.MapRange()
		for iter.Next() {
			val, err := valueToTemplateValue(iter.Value())
			if err != nil {
				return nil, err
			}
			out[iter.Key().String()] = val
		}
		return out, nil
	case reflect.Slice, reflect.Array:
		out := make([]any, v.Len())
		for i := 0; i < v.Len(); i++ {
			val, err := valueToTemplateValue(v.Index(i))
			if err != nil {
				return nil, err
			}
			out[i] = val
		}
		return out, nil
	default:
		return v.Interface(), nil
	}
}
