// Code generated by protoc-gen-gomeet-service. DO NOT EDIT.
// source: pb/profile.proto
package remotecli

func (c *remoteCli) cmdHelp(args []string) (string, error) {
	h := `HELP :

	┌─ version
	└─ call version service

	┌─ services_status
	└─ call services_status service

	┌─ create <gender [UNKNOW|MALE|FEMALE]> <email [string]> <name [string]> <birthday [string]>
	└─ call create service

	┌─ read <uuid [string]>
	└─ call read service

	┌─ list <page_number [uint32]> <page_size [uint32]> <order [string]> <exclude_soft_deleted [bool]> <soft_deleted_only [bool]> <gender [UNKNOW|MALE|FEMALE]>
	└─ call list service

	┌─ update <uuid [string]> <gender [UNKNOW|MALE|FEMALE]> <email [string]> <name [string]> <birthday [string]> <created_at [string]> <updated_at [string]> <deleted_at [string]>
	└─ call update service

	┌─ soft_delete <uuid [string]>
	└─ call soft_delete service

	┌─ hard_delete <uuid [string]>
	└─ call hard_delete service

	┌─ service_address
	└─ return service address

	┌─ jwt [<token>]
	└─ display current jwt or save none if it's set

	┌─ console_version
	└─ return console version

	┌─ tls_config
	└─ display TLS client configuration

	┌─ help
	└─ display this help
`
	if c.ctxCall == ConsoleCall {
		h += `
	┌─ exit
	└─ exit the console
`
	}

	return h + "\n", nil
}
