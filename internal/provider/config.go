package provider

const TEMP_DIR_PATTERN = ".terraform-provider-extended.tmp-*"

type Config struct {
	Input          string
	InputSensitive string
}
