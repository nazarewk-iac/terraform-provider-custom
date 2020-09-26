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

func resourceExtended() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceExtendedCreate,
		ReadContext:   resourceExtendedRead,
		UpdateContext: resourceExtendedUpdate,
		DeleteContext: resourceExtendedDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"program_create": {
				Type:     schema.TypeList,
				Elem:     schema.TypeString,
				MinItems: 1,
			},
			"program_read": {
				Type:     schema.TypeList,
				Elem:     schema.TypeString,
				MinItems: 1,
			},
			"program_update": {
				Type:     schema.TypeList,
				Elem:     schema.TypeString,
				Optional: true,
			},
			"program_delete": {
				Type:     schema.TypeList,
				Elem:     schema.TypeString,
				MinItems: 1,
			},
			"input": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"input_sensitive": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"state": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Default:  "",
			},
		},
	}
}

func resourceExtendedCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	config := meta.(*Config)
	diags = append(diags, runProgram(ctx, data, config, "program_create")...)
	return
}

func resourceExtendedRead(ctx context.Context, data *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	config := meta.(*Config)
	diags = append(diags, runProgram(ctx, data, config, "program_read")...)
	return
}

func resourceExtendedUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	config := meta.(*Config)
	diags = append(diags, runProgram(ctx, data, config, "program_update")...)
	return
}

func resourceExtendedDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	config := meta.(*Config)
	diags = append(diags, runProgram(ctx, data, config, "program_delete")...)
	return
}

func runProgram(ctx context.Context, data *schema.ResourceData, config *Config, programKey string) (diags diag.Diagnostics) {
	programSpec := data.Get(programKey).([]string)
	cwd, err := os.Getwd()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error retrieving current working director"),
			Detail:   err.Error(),
		})
		return
	}
	dir, err := ioutil.TempDir(cwd, TEMP_DIR_PATTERN)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error creating temporary directory %s in %s", TEMP_DIR_PATTERN, cwd),
			Detail:   err.Error(),
		})
		return
	}

	defer func() {
		if err := os.RemoveAll(dir); err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  fmt.Sprintf("Error when cleaning up temporary directory %s", dir),
				Detail:   err.Error(),
			})
		}
	}()

	cmd := exec.CommandContext(ctx, programSpec[0], programSpec[1:]...)
	env := os.Environ()
	env = append(env, fmt.Sprintf("%s=%s", "EXT_DIR", dir))

	oldStateV, newStateV := data.GetChange("state")
	files := map[string]string{
		"provider_input":           config.Input,
		"provider_input_sensitive": config.InputSensitive,
		"input":                    data.Get("input").(string),
		"input_sensitive":          data.Get("input_sensitive").(string),
		"state":                    newStateV.(string),
		"old_state":                oldStateV.(string),
		"id":                       data.Id(),
	}
	perms := map[string]os.FileMode{
		"state": 0600,
		"id":    0600,
	}
	for name, content := range files {
		currentPath := path.Join(dir, name)
		perm, ok := perms[name]
		if !ok {
			perm = 0400
		}
		diags = append(diags, createFile(currentPath, content, perm)...)
		if len(diags) > 0 {
			return
		}
		env = append(env, fmt.Sprintf("EXT_FILE_%s=%s", name, currentPath))
	}
	cmd.Env = env
	if err := cmd.Start(); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error when running %s", programKey),
			Detail:   err.Error(),
		})
		return
	}
	text, d := readFile(path.Join(dir, "id"))
	diags = append(diags, d...)
	if len(diags) > 0 {
		return
	}
	data.SetId(text)

	return
}

func readFile(path string) (text string, diags diag.Diagnostics) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error opening file %s", path),
			Detail:   err.Error(),
		})
	}
	text = string(content)
	return
}

func createFile(path string, content string, perm os.FileMode) (diags diag.Diagnostics) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, perm)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error creating file %s", path),
			Detail:   err.Error(),
		})
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Error when closing %s: %#v", path),
				Detail:   err.Error(),
			})
		}
	}()
	if _, err := file.WriteString(content); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error writing to file %s", path),
			Detail:   err.Error(),
		})
		return
	}
	return
}
