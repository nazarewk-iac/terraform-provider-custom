package provider

const TempDirBase = ".terraform/terraform-provider-custom"
const TempDirPattern = "*"

type Config struct {
	Input          string
	InputSensitive string
}
