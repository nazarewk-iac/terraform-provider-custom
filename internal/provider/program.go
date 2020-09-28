package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
)

type Program struct {
	ProgramSpec    []string
	providerConfig *Config
	data           *schema.ResourceData
	context        context.Context
	programKey     string
	tmpDir         string
	env            []string
	files          map[string]string
	perms          map[string]os.FileMode
}

func NewProgramFromResource(ctx context.Context, data *schema.ResourceData, config *Config, programKey string) *Program {
	oldStateV, newStateV := data.GetChange("state")
	p := &Program{
		context:        ctx,
		data:           data,
		providerConfig: config,
		programKey:     programKey,
		ProgramSpec:    data.Get(programKey).([]string),
	}

	// Can't cast interface{} -> []string directly, need to do this manually
	programSpecV := data.Get(programKey).([]interface{})
	p.ProgramSpec = make([]string, len(programSpecV))
	for _, obj := range programSpecV {
		p.ProgramSpec = append(p.ProgramSpec, obj.(string))
	}

	p.files = map[string]string{
		"provider_input":           p.providerConfig.Input,
		"provider_input_sensitive": p.providerConfig.InputSensitive,
		"input":                    p.data.Get("input").(string),
		"input_sensitive":          p.data.Get("input_sensitive").(string),
		"state":                    newStateV.(string),
		"old_state":                oldStateV.(string),
		"id":                       p.data.Id(),
	}

	p.perms = map[string]os.FileMode{
		"state": 0600,
		"id":    0600,
	}

	p.env = os.Environ()
	p.env = append(p.env, fmt.Sprintf("%s=%s", "EXT_DIR", p.tmpDir))

	for name, _ := range p.files {
		currentPath := path.Join(p.tmpDir, name)
		p.env = append(p.env, fmt.Sprintf("EXT_FILE_%s=%s", name, currentPath))
	}

	return p
}

func (p *Program) setupDir() (diags diag.Diagnostics) {
	cwd, err := os.Getwd()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error retrieving current working director"),
			Detail:   err.Error(),
		})
		return
	}
	p.tmpDir, err = ioutil.TempDir(cwd, TEMP_DIR_PATTERN)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error creating temporary directory %s in %s", TEMP_DIR_PATTERN, cwd),
			Detail:   err.Error(),
		})
		return
	}

	for name, content := range p.files {
		currentPath := path.Join(p.tmpDir, name)
		perm, ok := p.perms[name]
		if !ok {
			perm = 0400
		}
		diags = append(diags, p.createFile(currentPath, content, perm)...)
		if diags.HasError() {
			return
		}
	}
	return
}

func (p *Program) executeCommand() (diags diag.Diagnostics) {
	cmd := exec.CommandContext(p.context, p.ProgramSpec[0], p.ProgramSpec[1:]...)
	cmd.Env = p.env
	if err := cmd.Start(); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error when running %s", p.programKey),
			Detail:   err.Error(),
		})
		return
	}
	return
}

func (p *Program) storeNewId() (diags diag.Diagnostics) {
	text, diags := p.readFile("id")
	if diags.HasError() {
		return
	}
	p.data.SetId(text)
	return
}

func (p *Program) setNewState() (diags diag.Diagnostics) {
	text, diags := p.readFile("state")
	if diags.HasError() {
		return
	}

	if err := p.data.Set("state", text); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error when setting \"state\" attribute during %s", p.programKey),
			Detail:   err.Error(),
		})
		return
	}
	return
}

func runProgram(ctx context.Context, data *schema.ResourceData, config *Config, programKey string) (diags diag.Diagnostics) {
	p := NewProgramFromResource(ctx, data, config, programKey)
	diags = append(diags, p.setupDir()...)
	if diags.HasError() {
		return
	}
	defer func() {
		if err := os.RemoveAll(p.tmpDir); err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  fmt.Sprintf("Error when cleaning up temporary directory %s", p.tmpDir),
				Detail:   err.Error(),
			})
		}
	}()

	diags = append(diags, p.executeCommand()...)
	if diags.HasError() {
		return
	}

	diags = append(diags, p.storeNewId()...)
	if diags.HasError() {
		return
	}

	diags = append(diags, p.setNewState()...)
	if diags.HasError() {
		return
	}

	return
}

func (p *Program) readFile(name string) (text string, diags diag.Diagnostics) {
	fullPath := path.Join(p.tmpDir, name)
	content, err := ioutil.ReadFile(fullPath)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error opening file %s", fullPath),
			Detail:   err.Error(),
		})
	}
	text = string(content)
	return
}

func (p *Program) createFile(name string, content string, perm os.FileMode) (diags diag.Diagnostics) {
	fullPath := path.Join(p.tmpDir, name)
	file, err := os.OpenFile(fullPath, os.O_CREATE|os.O_WRONLY, perm)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error creating file %s", fullPath),
			Detail:   err.Error(),
		})
		return
	}

	defer func() {
		if err := file.Close(); err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Error when closing %s", fullPath),
				Detail:   err.Error(),
			})
		}
	}()

	if _, err := file.WriteString(content); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error writing to file %s", fullPath),
			Detail:   err.Error(),
		})
		return
	}
	return
}
