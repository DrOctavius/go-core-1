package autoloader

import (
	"github.com/kyaxcorp/go-core/core/helpers/file"
	"github.com/kyaxcorp/go-core/core/helpers/process/name"
	"github.com/kyaxcorp/go-core/core/helpers/slice"
	// cassandraConfig "github.com/kyaxcorp/go-core/core/clients/db/driver/cassandra/config"

	cockroachConfig "github.com/kyaxcorp/go-core/core/clients/db/driver/cockroach/config"
	mysqlConfig "github.com/kyaxcorp/go-core/core/clients/db/driver/mysql/config"
	websocketClientConfig "github.com/kyaxcorp/go-core/core/clients/websocket/config"
	websocketClientConnection "github.com/kyaxcorp/go-core/core/clients/websocket/connection"
	cfgData "github.com/kyaxcorp/go-core/core/config/data"
	"github.com/kyaxcorp/go-core/core/config/model"
	"github.com/kyaxcorp/go-core/core/helpers/_struct"
	"github.com/kyaxcorp/go-core/core/helpers/err"
	fsPath "github.com/kyaxcorp/go-core/core/helpers/filesystem/path"
	"github.com/kyaxcorp/go-core/core/helpers/hash"
	httpConfig "github.com/kyaxcorp/go-core/core/listeners/http/config"
	websocketServerConfig "github.com/kyaxcorp/go-core/core/listeners/websocket/config"
	loggingConfig "github.com/kyaxcorp/go-core/core/logger/config"
	//brokerConfig "github.com/kyaxcorp/go-core/core/services/broker/config"
	"github.com/spf13/viper"
	"path/filepath"
)

func SaveConfigFromMemory(cfg Config) error {
	c := viper.New()

	// Create the map!
	//MainConfig.ClientsStatus.MySQL.Connections = make(map[string]mysql.Config)
	//MainConfig.Listeners.Http.Instances = make(map[string]http.Config)

	// This is the default config!

	// Set default settings

	var _err error

	//---------------------------------------------------------------------------------\\
	//---------------------------\\    MYSQL CLIENT    //------------------------------\\
	//---------------------------------------------------------------------------------\\

	// Mysql Default Config Instance
	if _, ok := cfgData.MainConfig.Clients.MySQL.Instances["default"]; !ok {
		dbInstance, _err := mysqlConfig.SetDefaults(nil)
		if _err != nil {
			panic(_err)
		}

		// If the map is empty... we need to create it
		if cfgData.MainConfig.Clients.MySQL.Instances == nil {
			// Creating  the map, where we will set afterwards the default object
			cfgData.MainConfig.Clients.MySQL.Instances = make(map[string]mysqlConfig.Config)
		}
		cfgData.MainConfig.Clients.MySQL.Instances["default"] = *dbInstance
	}

	// Looping through instances and setting defaults if they are missing
	for connectionName, dbInstance := range cfgData.MainConfig.Clients.MySQL.Instances {
		_, _err = mysqlConfig.SetDefaults(&dbInstance)
		if _err != nil {
			panic(_err)
		}
		cfgData.MainConfig.Clients.MySQL.Instances[connectionName] = dbInstance
	}

	//---------------------------------------------------------------------------------\\
	//---------------------------\\    MYSQL CLIENT    //------------------------------\\
	//---------------------------------------------------------------------------------\\

	//

	//---------------------------------------------------------------------------------\\
	//-----------------------\\    COCKROACHDB CLIENT    //----------------------------\\
	//---------------------------------------------------------------------------------\\

	// Cockroach Default Config Instance
	if _, ok := cfgData.MainConfig.Clients.Cockroach.Instances["default"]; !ok {
		dbInstance, _err := cockroachConfig.SetDefaults(nil)
		if _err != nil {
			panic(_err)
		}

		// If the map is empty... we need to create it
		if cfgData.MainConfig.Clients.Cockroach.Instances == nil {
			// Creating  the map, where we will set afterwards the default object
			cfgData.MainConfig.Clients.Cockroach.Instances = make(map[string]cockroachConfig.Config)
		}
		cfgData.MainConfig.Clients.Cockroach.Instances["default"] = *dbInstance
	}

	// Looping through instances and setting defaults if they are missing
	for connectionName, dbInstance := range cfgData.MainConfig.Clients.Cockroach.Instances {
		_, _err = cockroachConfig.SetDefaults(&dbInstance)
		if _err != nil {
			panic(_err)
		}
		cfgData.MainConfig.Clients.Cockroach.Instances[connectionName] = dbInstance
	}
	//---------------------------------------------------------------------------------\\
	//-----------------------\\    COCKROACHDB CLIENT    //----------------------------\\
	//---------------------------------------------------------------------------------\\

	//

	//

	//---------------------------------------------------------------------------------\\
	//-----------------------\\    CASSANDRA CLIENT    //------------------------------\\
	//---------------------------------------------------------------------------------\\

	// Cassandra Default Config Instance
	/*if _, ok := cfgData.MainConfig.ClientsStatus.Cassandra.Instances["default"]; !ok {
		_cassandra := &cassandraConfig.Config{
			Hosts: []cassandraConfig.Host{
				cassandraConfig.Host{
					Destination: "",
					Port:        0,
				},
			},
		}
		if _err := _struct.SetDefaultValues(_cassandra); _err != nil {
			panic(_err)
		}
		// If the map is empty... we need to create it
		if cfgData.MainConfig.ClientsStatus.Cassandra.Instances == nil {
			// Creating  the map, where we will set afterwards the default object
			cfgData.MainConfig.ClientsStatus.Cassandra.Instances = make(map[string]cassandraConfig.Config)
		}
		cfgData.MainConfig.ClientsStatus.Cassandra.Instances["default"] = *_cassandra
	}*/

	//---------------------------------------------------------------------------------\\
	//-----------------------\\    CASSANDRA CLIENT    //------------------------------\\
	//---------------------------------------------------------------------------------\\

	//

	//---------------------------------------------------------------------------------\\
	//--------------------------\\    HTTP SERVER    //--------------------------------\\
	//---------------------------------------------------------------------------------\\

	// Http Default Config Instance
	if _, ok := cfgData.MainConfig.Listeners.Http.Instances["default"]; !ok {
		_http := &httpConfig.Config{}
		if _err := _struct.SetDefaultValues(_http); _err != nil {
			panic(_err)
		}
		// If the map is empty... we need to create it
		if cfgData.MainConfig.Listeners.Http.Instances == nil {
			// Creating  the map, where we will set afterwards the default object
			cfgData.MainConfig.Listeners.Http.Instances = make(map[string]httpConfig.Config)
		}
		cfgData.MainConfig.Listeners.Http.Instances["default"] = *_http
	}

	// Loop through all connections and set the default values
	for instanceName, httpInstance := range cfgData.MainConfig.Listeners.Http.Instances {

		// Logger config
		// Setting default values for logger
		if _err := _struct.SetDefaultValues(&httpInstance.Logger); _err != nil {
			panic(_err)
		}

		if _err := _struct.SetDefaultValues(&httpInstance); _err != nil {
			panic(_err)
		}

		if len(httpInstance.ListeningAddresses) == 0 {
			// No listening addresses, let's add one
			httpInstance.ListeningAddresses = []string{
				// the + symbol is the rule that checks if the port is busy, and if it is, it will
				// do +1 until finds a free port!
				"0.0.0.0:8080+",
			}
		}

		if len(httpInstance.ListeningAddressesSSL) == 0 {
			// No listening addresses, let's add one
			httpInstance.ListeningAddressesSSL = []string{
				// the + symbol is the rule that checks if the port is busy, and if it is, it will
				// do +1 until finds a free port!
				"0.0.0.0:8443+",
			}
		}

		// Set the logger to websocket config
		cfgData.MainConfig.Listeners.Http.Instances[instanceName] = httpInstance
	}

	//---------------------------------------------------------------------------------\\
	//--------------------------\\    HTTP SERVER    //--------------------------------\\
	//---------------------------------------------------------------------------------\\

	//

	//---------------------------------------------------------------------------------\\
	//-----------------------\\    WEBSOCKET SERVER    //------------------------------\\
	//---------------------------------------------------------------------------------\\

	// Http Default Config Instance
	if _, ok := cfgData.MainConfig.Listeners.WebSocket.Instances["default"]; !ok {
		_websocketServer := &websocketServerConfig.Config{}
		if _err := _struct.SetDefaultValues(_websocketServer); _err != nil {
			panic(_err)
		}
		// If the map is empty... we need to create it
		if cfgData.MainConfig.Listeners.WebSocket.Instances == nil {
			// Creating  the map, where we will set afterwards the default object
			cfgData.MainConfig.Listeners.WebSocket.Instances = make(map[string]websocketServerConfig.Config)
		}
		cfgData.MainConfig.Listeners.WebSocket.Instances["default"] = *_websocketServer
	}

	// Loop through all connections and set the default values
	for instanceName, wsInstance := range cfgData.MainConfig.Listeners.WebSocket.Instances {

		// Logger config
		// Setting default values for logger
		if _err := _struct.SetDefaultValues(&wsInstance.Logger); _err != nil {
			panic(_err)
		}

		if _err := _struct.SetDefaultValues(&wsInstance); _err != nil {
			panic(_err)
		}

		if len(wsInstance.ListeningAddresses) == 0 {
			// No listening addresses, let's add one
			wsInstance.ListeningAddresses = []string{
				// the + symbol is the rule that checks if the port is busy, and if it is, it will
				// do +1 until finds a free port!
				"0.0.0.0:8080+",
			}
		}

		if len(wsInstance.ListeningAddressesSSL) == 0 {
			// No listening addresses, let's add one
			wsInstance.ListeningAddressesSSL = []string{
				// the + symbol is the rule that checks if the port is busy, and if it is, it will
				// do +1 until finds a free port!
				"0.0.0.0:8443+",
			}
		}

		// Set the logger to websocket config
		cfgData.MainConfig.Listeners.WebSocket.Instances[instanceName] = wsInstance
	}

	//---------------------------------------------------------------------------------\\
	//-----------------------\\    WEBSOCKET SERVER    //------------------------------\\
	//---------------------------------------------------------------------------------\\

	//

	//---------------------------------------------------------------------------------\\
	//---------------------\\    WEBSOCKET CLIENT    //--------------------------------\\
	//---------------------------------------------------------------------------------\\

	// WebSocket Client Default Config
	if _, ok := cfgData.MainConfig.Clients.WebSocket.Instances["default"]; !ok {
		// Create the default config for websocket
		_webSocketClient := &websocketClientConfig.Config{}
		if _err := _struct.SetDefaultValues(_webSocketClient); _err != nil {
			panic(_err)
		}

		// If the map is empty... we need to create it
		if cfgData.MainConfig.Clients.WebSocket.Instances == nil {
			// Creating  the map, where we will set afterwards the default object
			cfgData.MainConfig.Clients.WebSocket.Instances = make(map[string]websocketClientConfig.Config)
		}

		// Set finally the default config
		cfgData.MainConfig.Clients.WebSocket.Instances["default"] = *_webSocketClient
	}

	// Loop through all connections and set the default values
	for instanceName, wsInstance := range cfgData.MainConfig.Clients.WebSocket.Instances {

		// If the default connection doesn't exist, create it!
		if exists, _ := slice.IndexExists(wsInstance.Connections, 0); !exists {
			// Create a default connection config for websocket
			wsInstance.Connections = append(wsInstance.Connections, &websocketClientConnection.Connection{})
		}

		// Loop through other connections and check if hey are ok!
		for connIndex, conn := range wsInstance.Connections {
			// Set the standard values for the object
			if _err := _struct.SetDefaultValues(conn); _err != nil {
				panic(_err)
			}
			// Set the authentication options
			if _err := _struct.SetDefaultValues(&conn.AuthOptions); _err != nil {
				panic(_err)
			}

			wsInstance.Connections[connIndex] = conn
		}

		// Logger config
		// Setting default values for logger
		if _err := _struct.SetDefaultValues(&wsInstance.Logger); _err != nil {
			panic(_err)
		}

		// Reconnect Options
		if _err := _struct.SetDefaultValues(&wsInstance.Reconnect); _err != nil {
			panic(_err)
		}

		if _err := _struct.SetDefaultValues(&wsInstance); _err != nil {
			panic(_err)
		}

		// Set the logger to websocket config
		cfgData.MainConfig.Clients.WebSocket.Instances[instanceName] = wsInstance
	}
	//---------------------------------------------------------------------------------\\
	//---------------------\\    WEBSOCKET CLIENT    //--------------------------------\\
	//---------------------------------------------------------------------------------\\

	//

	//---------------------------------------------------------------------------------\\
	//---------------------\\    LOGGING CHANNELS    //--------------------------------\\
	//---------------------------------------------------------------------------------\\

	// Logging Default config Channel
	if _, ok := cfgData.MainConfig.Logging.Channels["default"]; !ok {
		_logging := &loggingConfig.Config{}
		if _err := _struct.SetDefaultValues(_logging); _err != nil {
			panic(_err)
		}
		// If the map is empty... we need to create it
		if cfgData.MainConfig.Logging.Channels == nil {
			// Creating  the map, where we will set afterwards the default object
			cfgData.MainConfig.Logging.Channels = make(map[string]loggingConfig.Config)
		}
		cfgData.MainConfig.Logging.Channels["default"] = *_logging
	}

	// Check if there are additional channels defined from the main app
	if cfg.AdditionalLoggingChannels != nil {
		for channelName, channel := range cfg.AdditionalLoggingChannels {
			// Set the channel to existing map

			// Check if the channel already exists in the config
			if _, ok := cfgData.MainConfig.Logging.Channels[channelName]; !ok {
				// if it doesn't exist, get the default object
				if _err := _struct.SetDefaultValues(&channel); _err != nil {
					panic(_err)
				}
				// Set back
				cfgData.MainConfig.Logging.Channels[channelName] = channel
			} else {
				// if exists do nothing...
			}

		}
	}

	// Let's set the default values  for all channels
	for channelName, channel := range cfgData.MainConfig.Logging.Channels {
		if _err := _struct.SetDefaultValues(&channel); _err != nil {
			panic(_err)
		}
		// Set back
		cfgData.MainConfig.Logging.Channels[channelName] = channel
	}
	//---------------------------------------------------------------------------------\\
	//---------------------\\    LOGGING CHANNELS    //--------------------------------\\
	//---------------------------------------------------------------------------------\\

	//

	// Setting the main config
	c.Set("main", cfgData.MainConfig)
	// Setting the custom config
	c.Set("custom", cfg.CustomConfig)

	// TODO: save config only by comparing if it's different!
	// If it's diff, then overwrite it!
	configPath := GetConfigFilePath()
	if configPath == "" {
		return err.New(0, "config path is empty")
	}

	configTmpPath := GetConfigTmpFilePath()
	if configTmpPath == "" {
		return err.New(0, "config tmp path is empty")
	}
	// Save the temporary config file
	_err = c.WriteConfigAs(configTmpPath)
	if _err != nil {
		// log.Println("Failed to generate config!")
		return _err
	}

	// Compare the 2 configs
	tmpConfigHash, _err := hash.FileSha256(configTmpPath)
	// Delete the tmp config
	file.Delete(configTmpPath)
	if _err != nil {
		return _err
	}

	realConfigHash, _err := hash.FileSha256(configPath)
	if _err != nil {
		return _err
	}
	// log.Println(realConfigHash, tmpConfigHash)

	// Compare the 2 configs
	if tmpConfigHash == realConfigHash {
		// It's the same configuration!
		// log.Println("Same config!!! skipping save...")
		return nil
	}

	// Save the real config file
	_err = c.WriteConfigAs(configPath)
	if _err != nil {
		// log.Println("Failed to generate config!")
		return _err
	}
	return nil
}

func GetConfigPath() string {
	if globalConfigPath != "" {
		return globalConfigPath
	}

	// Get the config full path
	configFilePath := GetConfigFilePath()
	// Get only the base dir without file
	return filepath.Dir(configFilePath) + filepath.FromSlash("/")
}

func GetConfigFilePath() string {
	//path := GetConfigPath()
	// TODO: we should add here from the arguments introduced for the process!
	// There are 3 variants:
	// The default one is from the root directory
	// From the argument list from the process
	// From the OS Default config path: /etc/.... or Windows somewhere...

	path := fsPath.Root()

	if path != "" {
		path += GetConfigFullFileName()
	}
	return path
}

//GetConfigFileName
// We should have the same config name if the app name is not changed (doesn't matter if it's on windows or Linux)
// we should remove the file extension!
func GetConfigFileName() string {
	return model.ConfigFileName + "_" + hash.MD5(GetCleanAppFileName())
}

func GetCleanAppFileName() string {
	return name.GetCurrentProcessCleanExecName()
}

func GetConfigFileType() string {
	return model.ConfigFileType
}

func GetConfigFullFileName() string {
	return GetConfigFileName() + "." + GetConfigFileType()
}

func GetTmpConfigFileName() string {
	return model.ConfigTmpFileName + "_" + hash.MD5(GetCleanAppFileName()) + "." + model.ConfigFileType
}

// GetConfigTmpFilePath -> this is the temporary file path of the config... it's only for comparation purpose
func GetConfigTmpFilePath() string {
	path := GetConfigPath()
	if path != "" {
		path += GetTmpConfigFileName()
	}
	return path
}

func GetCustomConfigFilePath() string {
	// TODO: try reading from environment values or input arguments the path of the config!
	path := GetConfigPath()
	if path != "" {
		path = path + model.CustomConfigFileName + "." + model.CustomConfigFileType
	}
	return path
}

func IsConfigExists() bool {
	path := GetConfigFilePath()
	if path == "" {
		return false
	}
	return file.Exists(path)
}

func IsCustomConfigExists() bool {
	path := GetCustomConfigFilePath()
	if path == "" {
		return false
	}
	return file.Exists(path)
}

func IsConfigValid() bool {
	// TODO: load through viper and check if loaded!
	return true
}
