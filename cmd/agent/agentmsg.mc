;// Header
MessageIdTypedef=DWORD

LanguageNames=(
    English=0x409:MSG00409
)


;// Messages
MessageId=1
SymbolicName=MSG_WARN_REGCONFIG_FAILED
Severity=Warning
Language=English
Failed to import config items from registry.  The error was %1
.

MessageId=2
SymbolicName=MSG_WARN_CONFIGUPGRADE_FAILED
Severity=Warning
Language=English
Failed to upgrade configuration.  The error was %1.
.

MessageId=3
SymbolicName=MSG_SERVICE_STARTED
Severity=Informational
Language=English
The %1 service has started.
.

MessageId=4
SymbolicName=MSG_SERVICE_STOPPED
Severity=Informational
Language=English
The %1 service has stopped.
.

MessageId=5
SymbolicName=MSG_UNKNOWN_CONTROL_REQUEST
Severity=Warning
Language=English
Unexpected control request %1
.

MessageId=6
SymbolicName=MSG_SERVICE_STOPPING
Severity=Informational
Language=English
Received stop command, shutting down
.

MessageId=7
SymbolicName=MSG_SERVICE_STARTED
Severity=Informational
Language=English
starting the %1 service.
.

MessageId=8
SymbolicName=MSG_SERVICE_FAILED
Severity=Error
Language=English
The Service failed: %1
.

MessageId=8
SymbolicName=MSG_SERVICE_FAILED
Severity=Warning
Language=English
The Service received an unexpected control request: %1
.

MessageId=9
SymbolicName=MSG_UNEXPECTED_CONTROL_REQUEST
Severity=Error
Language=English
Unexpected message from the service control manager: %1
.

MessageId=10
SymbolicName=MSG_RECEIVED_STOP_COMMAND
Severity=Informational
Language=English
The service %1 received the stop signal, shutting down.
.

MessageId=11
SymbolicName=MSG_AGENT_START_FAILURE
Severity=Error
Language=English
The service failed to start. Error %1
.

MessageId=12
SymbolicName=MSG_RECEIVED_STOP_SVC_COMMAND
Severity=Informational
Language=English
The service %1 received the stop command from the service control manager, shutting down.
.

MessageId=13
SymbolicName=MSG_RECEIVED_STOP_SHUTDOWN
Severity=Informational
Language=English
The service %1 received the shutdown command from the service control manager, shutting down.
.

MessageId=14
SymbolicName=MSG_AGENT_SHUTDOWN_STARTING
Severity=Informational
Language=English
The service %1 service is initiating shutdown.
.

MessageId=15
SymbolicName=MSG_WARNING_PROGRAMDATA_ERROR
Severity=Warning
Language=English
Unable to determine the location of Program Data using the default value %1.
.

MessageId=16
SymbolicName=MSG_INFO_FUNCTION_ENTER
Severity=Informational
Language=English
%1 EnterFunction.
.

MessageId=17
SymbolicName=MSG_INFO_IMPORT_NOKEY
Severity=Informational
Language=English
ImportRegistryConfig() -- registry key doesn't exist
.

MessageId=18
SymbolicName=MSG_INFO_IMPORT_REG_UNEXPECTED
Severity=Informational
Language=English
%1 : Unexpected error reading registry
.

MessageId=19
SymbolicName=MSG_SETUP_CONFIG_SECRETS_ERR
Severity=Error
Language=English
Error setting up secrets config %1
.

MessageId=20
SymbolicName=MSG_API_KEY_PRESENT
Severity=Informational
Language=English
API key present in config file
.

MessageId=21
SymbolicName=MSG_INFO_STARING_PROCESSING
Severity=Informational
Language=English
Parsing install variables in registry
.

MessageId=22
SymbolicName=MSG_INFO_DONE_PARSING
Severity=Informational
Language=English
Done parsing variables from registry %1
.

MessageId=23
SymbolicName=MSG_ERROR_YAML_PARSE_FAILED
Severity=Error
Language=English
Failed to parse yaml %1
.

MessageId=24
SymbolicName=MSG_ERROR_FAIL_WRITE_YAML
Severity=Error
Language=English
Failed to write yaml config to file %1
.

MessageId=25
SymbolicName=MSG_INFO_DONE_REGISTRY
Severity=Informational
Language=English
Done writing registry config
.

MessageId=26
SymbolicName=MSG_INFO_CONF_NOT_EXIST
Severity=Informational
Language=English
Conf file doesn't exist, returning %1
.

MessageId=27
SymbolicName=MSG_INFO_NO_API_KEY
Severity=Informational
Language=English
Conf file has no api key
.

MessageId=28
SymbolicName=MSG_INFO_IMPORTING_A5_CONF
Severity=Informational
Language=English
Importing Agent5 config file %1
.

MessageId=29
SymbolicName=MSG_INFO_DONE_IMPORTING_A5_CONF
Severity=Informational
Language=English
Done importing A5 config
.

MessageId=30
SymbolicName=MSG_INFO_STARTING_AGENT
Severity=Informational
Language=English
Starting agent (pvs error %1)
.

MessageId=31
SymbolicName=MSG_INFO_STARTAGENT_START
Severity=Informational
Language=English
StartAgent() starting
.

MessageId=32
SymbolicName=MSG_INFO_STARTAGENT_LOGGING
Severity=Informational
Language=English
StartAgent() setting up logging
.

MessageId=33
SymbolicName=MSG_INFO_LOGGING_FAILED
Severity=Error
Language=English
StartAgent() failed to set up logging %1
.

MessageId=34
SymbolicName=MSG_INFO_STARTING_AGENT
Severity=Informational
Language=English
Starting Datadog Agent v%1.
.

MessageId=35
SymbolicName=MSG_INFO_STARTING_HEALTH
Severity=Informational
Language=English
Starting health server
.

MessageId=36
SymbolicName=MSG_ERROR_STARTING_HEALTH
Severity=Error
Language=English
Error starting health server %v
.

MessageId=37
SymbolicName=MSG_ERROR_STARTED_HEALTH
Severity=Informational
Language=English
Started health check on port %1
.

MessageId=38
SymbolicName=MSG_INFO_HOSTNAME
Severity=Informational
Language=English
Computed hostname to be %1.
.

MessageId=39
SymbolicName=MSG_INFO_STARTING_APISERVER
Severity=Informational
Language=English
Starting API server
.

MessageId=40
SymbolicName=MSG_ERROR_APISERVER_FAILED
Severity=Error
Language=English
Starting API server failed %1.
.

MessageId=41
SymbolicName=MSG_STATSD_STARTED
Severity=Informational
Language=English
Statsd started
.

MessageId=42
SymbolicName=MSG_LOGS_STARTING
Severity=Informational
Language=English
Starting logs
.

MessageId=43
SymbolicName=MSG_STARTING_DEPENDENT_SERVICES
Severity=Informational
Language=English
Starting dependent services
.

MessageId=44
SymbolicName=MSG_DONE_STARTING_AGENT
Severity=Informational
Language=English
Done starting agent
.

MessageId=45
SymbolicName=MSG_STARTING_SERVICES
Severity=Informational
Language=English
Starting dependent service %1.
.

MessageId=46
SymbolicName=MSG_FAIL_CONNECT_SCM
Severity=Error
Language=English
Failed to connect to scm %1.
.

MessageId=47
SymbolicName=MSG_FAIL_OPENSERVICE
Severity=Error
Language=English
Failed to Open service %1.
.

MessageId=48
SymbolicName=MSG_FAIL_STARTSERVICE
Severity=Error
Language=English
Failed to start service %1.
.

MessageId=49
SymbolicName=MSG_ENTER_EXECUTE
Severity=Informational
Language=English
Entered service callback Execute (%1)
.

MessageId=50
SymbolicName=MSG_CALLING_RUN
Severity=Informational
Language=English
Entering run method.
.

MessageId=51
SymbolicName=MSG_RUN_RETURN
Severity=Informational
Language=English
Run method returned %s
.

MessageId=52
SymbolicName=MSG_DEFER_FUNCTION_EXIT
Severity=Informational
Language=English
Function %s has executed the defer().
.
