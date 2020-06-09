package handlers

var (
	// Deploy and Migrate share the same command
	// the api honors the same request/response signature for both
	CommandHandlerMap = map[string]interface{}{
		"deploy":  NewDeployHandler,
		"migrate": NewMigrateHandler,
	}
)
