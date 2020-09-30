package provider

const TEMP_DIR_PATTERN = ".terraform-provider-custom.tmp-*"

type Config struct {
	Input          string
	InputSensitive string
}
