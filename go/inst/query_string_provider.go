package inst

import "unicode"

type QueryStringKey int

type QueryStringProvider struct {
	queries map[QueryStringKey]string
}

const (
	log_slave_updates QueryStringKey = iota
	start_slave
	show_slave_status
	slave_io_running
	slave_sql_running
	show_master_status
	master_gtid_info
	master_user
	slave_io_state
	master_log_file
	read_master_log_pos
	relay_master_log_file
	relay_master_log_position
	master_uuid
	master_host
	master_port
	seconds_behind_master
	master_ssl_allowed
	show_slave_hosts

	stop_slave_io_thread
	stop_slave_sql_thread
	start_slave_sql_thread
	start_slave_io_thread
	stop_slave
	start_slave_until_master_log
	master_user_param
	master_password_param
	master_ssl_ca_param
	master_ssl_cert_param
	master_ssl
	master_ssl_key_param
	change_master_to_master_ssl
	reset_slave
	change_master_to_master_host_port
	change_master_to_master_host_port_log_gtid_no
	change_master_to_master_host_port_log_autoposition_no
	change_master_to_master_host_port_autoposition_yes
	change_master_to_master_host_port_log
	change_master_to_master_host
	reset_slave_50603_all
	reset_master
	select_master_pos_wait
	change_master_to_with_params
	change_master_to_master_delay
	set_sql_slave_skip_counter
	change_master_to_get_source_public_key
)

func (qps *QueryStringProvider) log_slave_updates() string {
	return qps.queries[log_slave_updates]
}

func (qps *QueryStringProvider) start_slave() string {
	return qps.queries[start_slave]
}

func (qps *QueryStringProvider) show_slave_status() string {
	return qps.queries[show_slave_status]
}

func (qps *QueryStringProvider) slave_io_running() string {
	return qps.queries[slave_io_running]
}

func (qps *QueryStringProvider) slave_sql_running() string {
	return qps.queries[slave_sql_running]
}

func (qps *QueryStringProvider) master_gtid_info() string {
	return qps.queries[master_gtid_info]
}

func (qps *QueryStringProvider) show_master_status() string {
	return qps.queries[show_master_status]
}

func (qps *QueryStringProvider) master_user() string {
	return qps.queries[master_user]
}

func (qps *QueryStringProvider) slave_io_state() string {
	return qps.queries[slave_io_state]
}

func (qps *QueryStringProvider) master_log_file() string {
	return qps.queries[master_log_file]
}

func (qps *QueryStringProvider) read_master_log_pos() string {
	return qps.queries[read_master_log_pos]
}

func (qps *QueryStringProvider) relay_master_log_file() string {
	return qps.queries[relay_master_log_file]
}

func (qps *QueryStringProvider) relay_master_log_position() string {
	return qps.queries[relay_master_log_position]
}

func (qps *QueryStringProvider) master_uuid() string {
	return qps.queries[master_uuid]
}

func (qps *QueryStringProvider) master_host() string {
	return qps.queries[master_host]
}

func (qps *QueryStringProvider) master_port() string {
	return qps.queries[master_port]
}

func (qps *QueryStringProvider) seconds_behind_master() string {
	return qps.queries[seconds_behind_master]
}

func (qps *QueryStringProvider) master_ssl_allowed() string {
	return qps.queries[master_ssl_allowed]
}

func (qps *QueryStringProvider) show_slave_hosts() string {
	return qps.queries[show_slave_hosts]
}

func (qps *QueryStringProvider) stop_slave_io_thread() string {
	return qps.queries[stop_slave_io_thread]
}

func (qps *QueryStringProvider) stop_slave_sql_thread() string {
	return qps.queries[stop_slave_sql_thread]
}

func (qps *QueryStringProvider) start_slave_sql_thread() string {
	return qps.queries[start_slave_sql_thread]
}

func (qps *QueryStringProvider) start_slave_io_thread() string {
	return qps.queries[start_slave_io_thread]
}

func (qps *QueryStringProvider) stop_slave() string {
	return qps.queries[stop_slave]
}

func (qps *QueryStringProvider) start_slave_until_master_log() string {
	return qps.queries[start_slave_until_master_log]
}

func (qps *QueryStringProvider) master_user_param() string {
	return qps.queries[master_user_param]
}

func (qps *QueryStringProvider) master_password_param() string {
	return qps.queries[master_password_param]
}

func (qps *QueryStringProvider) master_ssl_ca_param() string {
	return qps.queries[master_ssl_ca_param]
}

func (qps *QueryStringProvider) master_ssl_cert_param() string {
	return qps.queries[master_ssl_cert_param]
}

func (qps *QueryStringProvider) master_ssl() string {
	return qps.queries[master_ssl]
}

func (qps *QueryStringProvider) master_ssl_key_param() string {
	return qps.queries[master_ssl_key_param]
}

func (qps *QueryStringProvider) change_master_to_master_ssl() string {
	return qps.queries[change_master_to_master_ssl]
}

func (qps *QueryStringProvider) reset_slave() string {
	return qps.queries[reset_slave]
}

func (qps *QueryStringProvider) change_master_to_master_host_port() string {
	return qps.queries[change_master_to_master_host_port]
}

func (qps *QueryStringProvider) change_master_to_master_host_port_log_gtid_no() string {
	return qps.queries[change_master_to_master_host_port_log_gtid_no]
}

func (qps *QueryStringProvider) change_master_to_master_host_port_log_autoposition_no() string {
	return qps.queries[change_master_to_master_host_port_log_autoposition_no]
}

func (qps *QueryStringProvider) change_master_to_master_host_port_autoposition_yes() string {
	return qps.queries[change_master_to_master_host_port_autoposition_yes]
}

func (qps *QueryStringProvider) change_master_to_master_host_port_log() string {
	return qps.queries[change_master_to_master_host_port_log]
}

func (qps *QueryStringProvider) change_master_to_master_host() string {
	return qps.queries[change_master_to_master_host]
}

func (qps *QueryStringProvider) reset_slave_50603_all() string {
	return qps.queries[reset_slave_50603_all]
}

func (qps *QueryStringProvider) reset_master() string {
	return qps.queries[reset_master]
}

func (qps *QueryStringProvider) select_master_pos_wait() string {
	return qps.queries[select_master_pos_wait]
}

func (qps *QueryStringProvider) change_master_to_with_params() string {
	return qps.queries[change_master_to_with_params]
}

func (qps *QueryStringProvider) change_master_to_master_delay() string {
	return qps.queries[change_master_to_master_delay]
}

func (qps *QueryStringProvider) set_sql_slave_skip_counter() string {
	return qps.queries[set_sql_slave_skip_counter]
}

func (qps *QueryStringProvider) change_master_to_get_source_public_key() string {
	return qps.queries[change_master_to_get_source_public_key]
}

var queryStrings80 = map[QueryStringKey]string{
	log_slave_updates:         "log_slave_updates",
	start_slave:               "start slave",
	show_slave_status:         "show slave status",
	slave_io_running:          "Slave_IO_Running",
	slave_sql_running:         "Slave_SQL_Running",
	show_master_status:        "show master status",
	master_gtid_info:          "select @@global.gtid_mode, @@global.server_uuid, @@global.gtid_executed, @@global.gtid_purged, @@global.master_info_repository = 'TABLE', @@global.binlog_row_image",
	master_user:               "Master_User",
	slave_io_state:            "Slave_IO_State",
	master_log_file:           "Master_Log_File",
	read_master_log_pos:       "Read_Master_Log_Pos",
	relay_master_log_file:     "Relay_Master_Log_File",
	relay_master_log_position: "Exec_Master_Log_Pos",
	master_uuid:               "Master_UUID",
	master_host:               "Master_Host",
	master_port:               "Master_Port",
	seconds_behind_master:     "Seconds_Behind_Master",
	master_ssl_allowed:        "Master_SSL_Allowed",
	show_slave_hosts:          "show slave hosts",

	stop_slave_io_thread:                          "stop slave io_thread",
	stop_slave_sql_thread:                         "stop slave sql_thread",
	start_slave_sql_thread:                        "start slave sql_thread",
	start_slave_io_thread:                         "start slave io_thread",
	stop_slave:                                    "stop slave",
	start_slave_until_master_log:                  "start slave until master_log_file=?, master_log_pos=?",
	master_user_param:                             "master_user = ?",
	master_password_param:                         "master_password = ?",
	master_ssl_ca_param:                           "master_ssl_ca = ?",
	master_ssl_cert_param:                         "master_ssl_cert = ?",
	master_ssl:                                    "master_ssl",
	master_ssl_key_param:                          "master_ssl_key = ?",
	change_master_to_master_ssl:                   "change master to master_ssl=1",
	reset_slave:                                   "reset slave",
	change_master_to_master_host_port:             "change master to master_host=?, master_port=?",
	change_master_to_master_host_port_log_gtid_no: "change master to master_host=?, master_port=?, master_log_file=?, master_log_pos=?, master_use_gtid=no",
	change_master_to_master_host_port_log_autoposition_no: "change master to master_host=?, master_port=?, master_log_file=?, master_log_pos=?, master_auto_position=0",
	change_master_to_master_host_port_autoposition_yes:    "change master to master_host=?, master_port=?, master_auto_position=1",
	change_master_to_master_host_port_log:                 "change master to master_host=?, master_port=?, master_log_file=?, master_log_pos=?",
	change_master_to_master_host:                          "change master to master_host='_'",
	reset_slave_50603_all:                                 "reset slave /*!50603 all */",
	reset_master:                                          "reset master",
	select_master_pos_wait:                                "select master_pos_wait(?, ?)",
	change_master_to_with_params:                          "change master to %s",
	change_master_to_master_delay:                         "change master to master_delay=%d",
	set_sql_slave_skip_counter:                            "set global sql_slave_skip_counter := 1",
	change_master_to_get_source_public_key:                "select 1",
}

var queryStrings84 = map[QueryStringKey]string{
	log_slave_updates:         "log_replica_updates",
	start_slave:               "start replica",
	show_slave_status:         "show replica status",
	slave_io_running:          "Replica_IO_Running",
	slave_sql_running:         "Replica_SQL_Running",
	show_master_status:        "show binary log status",
	master_gtid_info:          "select @@global.gtid_mode, @@global.server_uuid, @@global.gtid_executed, @@global.gtid_purged, 1, @@global.binlog_row_image",
	master_user:               "Source_User",
	slave_io_state:            "Replica_IO_State",
	master_log_file:           "Source_Log_File",
	read_master_log_pos:       "Read_Source_Log_Pos",
	relay_master_log_file:     "Relay_Source_Log_File",
	relay_master_log_position: "Exec_Source_Log_Pos",
	master_uuid:               "Source_UUID",
	master_host:               "Source_Host",
	master_port:               "Source_Port",
	seconds_behind_master:     "Seconds_Behind_Source",
	master_ssl_allowed:        "Source_SSL_Allowed",
	show_slave_hosts:          "show replicas",

	stop_slave_io_thread:                          "stop replica io_thread",
	stop_slave_sql_thread:                         "stop replica sql_thread",
	start_slave_sql_thread:                        "start replica sql_thread",
	start_slave_io_thread:                         "start replica io_thread",
	stop_slave:                                    "stop replica",
	start_slave_until_master_log:                  "start replica until source_log_file=?, source_log_pos=?",
	master_user_param:                             "source_user = ?",
	master_password_param:                         "source_password = ?",
	master_ssl_ca_param:                           "source_ssl_ca = ?",
	master_ssl_cert_param:                         "source_ssl_cert = ?",
	master_ssl:                                    "source_ssl",
	master_ssl_key_param:                          "source_ssl_key = ?",
	change_master_to_master_ssl:                   "change replication source to source_ssl=1",
	reset_slave:                                   "reset replica",
	change_master_to_master_host_port:             "change replication source to source_host=?, source_port=?",
	change_master_to_master_host_port_log_gtid_no: "change replication source to source_host=?, source_port=?, source_log_file=?, source_log_pos=?, source_use_gtid=no",
	change_master_to_master_host_port_log_autoposition_no: "change replication source to source_host=?, source_port=?, source_log_file=?, source_log_pos=?, source_auto_position=0",
	change_master_to_master_host_port_autoposition_yes:    "change replication source to source_host=?, source_port=?, source_auto_position=1",
	change_master_to_master_host_port_log:                 "change replication source to source_host=?, source_port=?, source_log_file=?, source_log_pos=?",
	change_master_to_master_host:                          "change replication source to source_host='_'",
	reset_slave_50603_all:                                 "reset replica /*!50603 all */",
	reset_master:                                          "reset binary logs and gtids",
	select_master_pos_wait:                                "select source_pos_wait(?, ?)",
	change_master_to_with_params:                          "change replication source to %s",
	change_master_to_master_delay:                         "change replication source to source_delay=%d",
	set_sql_slave_skip_counter:                            "set global sql_replica_skip_counter := 1",
	change_master_to_get_source_public_key:                "change replication source to get_source_public_key=1",
}

var queryStringProvider80 = QueryStringProvider{
	queries: queryStrings80,
}

var queryStringProvider84 = QueryStringProvider{
	queries: queryStrings84,
}

func GetQueryStringProvider(version string) QueryStringProvider {
	// always fall back to 8.0 provider
	if len(version) == 0 || !unicode.IsDigit(rune(version[0])) {
		return queryStringProvider80
	}

	if version >= "8.4" {
		return queryStringProvider84
	}

	return queryStringProvider80
}
