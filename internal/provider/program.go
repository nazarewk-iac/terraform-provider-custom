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

type program struct {
	providerConfig *Config
	data           *schema.ResourceData
	context        context.Context
	name           string
	tmpDir         string
	files          map[string]string
	perms          map[string]os.FileMode
}

func Program(ctx context.Context, data *schema.ResourceData, config *Config) *program {
	oldStateV, newStateV := data.GetChange("state")
	p := &program{
		context:        ctx,
		data:           data,
		providerConfig: config,
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

	return p
}

func (p *program) openDir() (diags diag.Diagnostics) {
	cwd, err := os.Getwd()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error retrieving current working director"),
			Detail:   err.Error(),
		})
		return
	}

	if err := os.MkdirAll(TempDirBase, 0700); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error creating temporary directory parent %s in %s", TempDirBase, cwd),
			Detail:   err.Error(),
		})
		return
	}
	p.tmpDir, err = ioutil.TempDir(TempDirBase, TempDirPattern)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error creating temporary directory %s in %s", TempDirBase, cwd),
			Detail:   err.Error(),
		})
		return
	}

	for name, content := range p.files {
		perm, ok := p.perms[name]
		if !ok {
			perm = 0400
		}
		diags = append(diags, p.createFile(name, content, perm)...)
		if diags.HasError() {
			return
		}
	}
	return
}
func (p *program) prepareEnv() (env []string, diags diag.Diagnostics) {
	env = append(env, os.Environ()...)
	if len(p.tmpDir) == 0 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Cannot prepareEnv() because tmpDir is empty!",
		})
		return
	}

	env = append(env, fmt.Sprintf("%s=%s", "EXT_DIR", p.tmpDir))

	for name, _ := range p.files {
		currentPath := path.Join(p.tmpDir, name)
		env = append(env, fmt.Sprintf("EXT_FILE_%s=%s", name, currentPath))
	}
	return
}

func (p *program) executeCommand(key string) (diags diag.Diagnostics) {
	args := p.getArgs(key)
	cmd := exec.CommandContext(p.context, args[0], args[1:]...)
	env, d := p.prepareEnv()
	diags = append(diags, d...)
	if diags.HasError() {
		return
	}
	cmdRepr := ToString(args)

	cmd.Env = env
	output, err := cmd.CombinedOutput()
	diags = append(diags, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  fmt.Sprintf("Combined output (%d bytes) of %s: %s", len(output), p.name, cmdRepr),
		Detail:   string(output),
	})
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error when running %s", p.name),
			Detail:   fmt.Sprintf("ERROR=%v\nCOMMAND %v\nOUTPUT (%d bytes):\n%v", err.Error(), cmdRepr, len(output), string(output)),
		})
	}
	return
}

func (p *program) storeNewId() (diags diag.Diagnostics) {
	text, diags := p.readFile("id")
	if diags.HasError() {
		return
	}
	p.data.SetId(text)
	return
}

func (p *program) setNewState() (diags diag.Diagnostics) {
	text, diags := p.readFile("state")
	if diags.HasError() {
		return
	}

	if err := p.data.Set("state", text); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error when setting \"state\" attribute during %s", p.name),
			Detail:   err.Error(),
		})
		return
	}
	return
}

func (p *program) getArgs(key string) (spec []string) {
	programSpecV := p.data.Get(key).([]interface{})
	spec = make([]string, len(programSpecV))
	for i, obj := range programSpecV {
		spec[i] = obj.(string)
	}
	return spec
}

func (p *program) closeDir() (diags diag.Diagnostics) {
	if len(p.tmpDir) == 0 {
		return
	}
	if err := os.RemoveAll(p.tmpDir); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  fmt.Sprintf("Error when cleaning up temporary directory %s", p.tmpDir),
			Detail:   err.Error(),
		})
	}
	p.tmpDir = ""
	return
}

func runProgram(ctx context.Context, data *schema.ResourceData, config *Config, name string, commandKey string) (diags diag.Diagnostics) {
	p := Program(ctx, data, config)
	p.name = name
	diags = append(diags, p.openDir()...)
	if diags.HasError() {
		return
	}
	defer func() { diags = append(diags, p.closeDir()...) }()

	diags = append(diags, p.executeCommand(commandKey)...)
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

func (p *program) readFile(name string) (text string, diags diag.Diagnostics) {
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

func (p *program) createFile(name string, content string, perm os.FileMode) (diags diag.Diagnostics) {
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
