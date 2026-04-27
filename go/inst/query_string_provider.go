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
	reset_slave
	change_master_to_master_host
	reset_slave_50603_all
	reset_master
	select_master_pos_wait
	change_master_to_with_params
	change_master_to_master_delay
	set_sql_slave_skip_counter
	change_master_to_get_source_public_key
	select_user_host
	master_host_assign
	master_port_assign
	master_log_file_assign
	master_log_pos_assign
	master_auto_position_assign
	master_ssl_capath_param
	master_ssl_cipher_param
	master_ssl_crl_param
	master_ssl_crlpath_param
	master_ssl_verify_server_cert_param
	master_tls_version_param
	master_tls_ciphersuites_param
	master_public_key_path_param
	master_get_source_public_key_param
	replica_status_ssl_ca_file
	replica_status_ssl_ca_path
	replica_status_ssl_cert
	replica_status_ssl_cipher
	replica_status_ssl_crl_file
	replica_status_ssl_crl_path
	replica_status_ssl_key
	replica_status_ssl_verify_server_cert
	replica_status_tls_version
	replica_status_tls_ciphersuites
	replica_status_public_key_path
	replica_status_get_source_public_key
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

func (qps *QueryStringProvider) reset_slave() string {
	return qps.queries[reset_slave]
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

func (qps *QueryStringProvider) select_user_host() string {
	return qps.queries[select_user_host]
}

func (qps *QueryStringProvider) master_host_assign() string {
	return qps.queries[master_host_assign]
}

func (qps *QueryStringProvider) master_port_assign() string {
	return qps.queries[master_port_assign]
}

func (qps *QueryStringProvider) master_log_file_assign() string {
	return qps.queries[master_log_file_assign]
}

func (qps *QueryStringProvider) master_log_pos_assign() string {
	return qps.queries[master_log_pos_assign]
}

func (qps *QueryStringProvider) master_auto_position_assign() string {
	return qps.queries[master_auto_position_assign]
}

func (qps *QueryStringProvider) master_ssl_capath_param() string {
	return qps.queries[master_ssl_capath_param]
}

func (qps *QueryStringProvider) master_ssl_cipher_param() string {
	return qps.queries[master_ssl_cipher_param]
}

func (qps *QueryStringProvider) master_ssl_crl_param() string {
	return qps.queries[master_ssl_crl_param]
}

func (qps *QueryStringProvider) master_ssl_crlpath_param() string {
	return qps.queries[master_ssl_crlpath_param]
}

func (qps *QueryStringProvider) master_ssl_verify_server_cert_param() string {
	return qps.queries[master_ssl_verify_server_cert_param]
}

func (qps *QueryStringProvider) master_tls_version_param() string {
	return qps.queries[master_tls_version_param]
}

func (qps *QueryStringProvider) master_tls_ciphersuites_param() string {
	return qps.queries[master_tls_ciphersuites_param]
}

func (qps *QueryStringProvider) master_public_key_path_param() string {
	return qps.queries[master_public_key_path_param]
}

func (qps *QueryStringProvider) master_get_source_public_key_param() string {
	return qps.queries[master_get_source_public_key_param]
}

func (qps *QueryStringProvider) replica_status_ssl_ca_file() string {
	return qps.queries[replica_status_ssl_ca_file]
}

func (qps *QueryStringProvider) replica_status_ssl_ca_path() string {
	return qps.queries[replica_status_ssl_ca_path]
}

func (qps *QueryStringProvider) replica_status_ssl_cert() string {
	return qps.queries[replica_status_ssl_cert]
}

func (qps *QueryStringProvider) replica_status_ssl_cipher() string {
	return qps.queries[replica_status_ssl_cipher]
}

func (qps *QueryStringProvider) replica_status_ssl_crl_file() string {
	return qps.queries[replica_status_ssl_crl_file]
}

func (qps *QueryStringProvider) replica_status_ssl_crl_path() string {
	return qps.queries[replica_status_ssl_crl_path]
}

func (qps *QueryStringProvider) replica_status_ssl_key() string {
	return qps.queries[replica_status_ssl_key]
}

func (qps *QueryStringProvider) replica_status_ssl_verify_server_cert() string {
	return qps.queries[replica_status_ssl_verify_server_cert]
}

func (qps *QueryStringProvider) replica_status_tls_version() string {
	return qps.queries[replica_status_tls_version]
}

func (qps *QueryStringProvider) replica_status_tls_ciphersuites() string {
	return qps.queries[replica_status_tls_ciphersuites]
}

func (qps *QueryStringProvider) replica_status_public_key_path() string {
	return qps.queries[replica_status_public_key_path]
}

func (qps *QueryStringProvider) replica_status_get_source_public_key() string {
	return qps.queries[replica_status_get_source_public_key]
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

	stop_slave_io_thread:                   "stop slave io_thread",
	stop_slave_sql_thread:                  "stop slave sql_thread",
	start_slave_sql_thread:                 "start slave sql_thread",
	start_slave_io_thread:                  "start slave io_thread",
	stop_slave:                             "stop slave",
	start_slave_until_master_log:           "start slave until master_log_file=?, master_log_pos=?",
	master_user_param:                      "master_user = ?",
	master_password_param:                  "master_password = ?",
	master_ssl_ca_param:                    "master_ssl_ca = ?",
	master_ssl_cert_param:                  "master_ssl_cert = ?",
	master_ssl:                             "master_ssl",
	master_ssl_key_param:                   "master_ssl_key = ?",
	reset_slave:                            "reset slave",
	change_master_to_master_host:           "change master to master_host='_'",
	reset_slave_50603_all:                  "reset slave /*!50603 all */",
	reset_master:                           "reset master",
	select_master_pos_wait:                 "select master_pos_wait(?, ?)",
	change_master_to_with_params:           "change master to %s",
	change_master_to_master_delay:          "change master to master_delay=%d",
	set_sql_slave_skip_counter:             "set global sql_slave_skip_counter := 1",
	change_master_to_get_source_public_key: "select 1",
	select_user_host:                       "select user, substring_index(host, ':', 1) as slave_hostname from information_schema.processlist where command IN ('Binlog Dump', 'Binlog Dump GTID')",

	master_host_assign:                    "master_host=?",
	master_port_assign:                    "master_port=?",
	master_log_file_assign:                "master_log_file=?",
	master_log_pos_assign:                 "master_log_pos=?",
	master_auto_position_assign:           "master_auto_position=?",
	master_ssl_capath_param:               "master_ssl_capath = ?",
	master_ssl_cipher_param:               "master_ssl_cipher = ?",
	master_ssl_crl_param:                  "master_ssl_crl = ?",
	master_ssl_crlpath_param:              "master_ssl_crlpath = ?",
	master_ssl_verify_server_cert_param:   "master_ssl_verify_server_cert = ?",
	master_tls_version_param:              "master_tls_version = ?",
	master_tls_ciphersuites_param:         "master_tls_ciphersuites = ?",
	master_public_key_path_param:          "master_public_key_path = ?",
	master_get_source_public_key_param:    "get_master_public_key = ?",
	replica_status_ssl_ca_file:            "Master_SSL_CA_File",
	replica_status_ssl_ca_path:            "Master_SSL_CA_Path",
	replica_status_ssl_cert:               "Master_SSL_Cert",
	replica_status_ssl_cipher:             "Master_SSL_Cipher",
	replica_status_ssl_crl_file:           "Master_SSL_CRL_File",
	replica_status_ssl_crl_path:           "Master_SSL_CRL_Path",
	replica_status_ssl_key:                "Master_SSL_Key",
	replica_status_ssl_verify_server_cert: "Master_SSL_Verify_Server_Cert",
	replica_status_tls_version:            "Master_TLS_Version",
	replica_status_tls_ciphersuites:       "Master_TLS_Ciphersuites",
	replica_status_public_key_path:        "Master_public_key_path",
	replica_status_get_source_public_key:  "Get_master_public_key",
}

var queryStrings8014 = func() map[QueryStringKey]string {
	m := make(map[QueryStringKey]string, len(queryStrings80))

	// copy everything from version < 8.0.14
	for k, v := range queryStrings80 {
		m[k] = v
	}

	// change for 8.0.14 and newer 8.0
	m[select_user_host] = "select user, substring_index(host, ':', 1) as slave_hostname from performance_schema.processlist where command IN ('Binlog Dump', 'Binlog Dump GTID')"

	return m

}()

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

	stop_slave_io_thread:                   "stop replica io_thread",
	stop_slave_sql_thread:                  "stop replica sql_thread",
	start_slave_sql_thread:                 "start replica sql_thread",
	start_slave_io_thread:                  "start replica io_thread",
	stop_slave:                             "stop replica",
	start_slave_until_master_log:           "start replica until source_log_file=?, source_log_pos=?",
	master_user_param:                      "source_user = ?",
	master_password_param:                  "source_password = ?",
	master_ssl_ca_param:                    "source_ssl_ca = ?",
	master_ssl_cert_param:                  "source_ssl_cert = ?",
	master_ssl:                             "source_ssl",
	master_ssl_key_param:                   "source_ssl_key = ?",
	reset_slave:                            "reset replica",
	change_master_to_master_host:           "change replication source to source_host='_'",
	reset_slave_50603_all:                  "reset replica /*!50603 all */",
	reset_master:                           "reset binary logs and gtids",
	select_master_pos_wait:                 "select source_pos_wait(?, ?)",
	change_master_to_with_params:           "change replication source to %s",
	change_master_to_master_delay:          "change replication source to source_delay=%d",
	set_sql_slave_skip_counter:             "set global sql_replica_skip_counter := 1",
	change_master_to_get_source_public_key: "change replication source to get_source_public_key=1",
	select_user_host:                       "select user, substring_index(host, ':', 1) as slave_hostname from performance_schema.processlist where command IN ('Binlog Dump', 'Binlog Dump GTID')",

	master_host_assign:                    "source_host=?",
	master_port_assign:                    "source_port=?",
	master_log_file_assign:                "source_log_file=?",
	master_log_pos_assign:                 "source_log_pos=?",
	master_auto_position_assign:           "source_auto_position=?",
	master_ssl_capath_param:               "source_ssl_capath = ?",
	master_ssl_cipher_param:               "source_ssl_cipher = ?",
	master_ssl_crl_param:                  "source_ssl_crl = ?",
	master_ssl_crlpath_param:              "source_ssl_crlpath = ?",
	master_ssl_verify_server_cert_param:   "source_ssl_verify_server_cert = ?",
	master_tls_version_param:              "source_tls_version = ?",
	master_tls_ciphersuites_param:         "source_tls_ciphersuites = ?",
	master_public_key_path_param:          "source_public_key_path = ?",
	master_get_source_public_key_param:    "get_source_public_key = ?",
	replica_status_ssl_ca_file:            "Source_SSL_CA_File",
	replica_status_ssl_ca_path:            "Source_SSL_CA_Path",
	replica_status_ssl_cert:               "Source_SSL_Cert",
	replica_status_ssl_cipher:             "Source_SSL_Cipher",
	replica_status_ssl_crl_file:           "Source_SSL_CRL_File",
	replica_status_ssl_crl_path:           "Source_SSL_CRL_Path",
	replica_status_ssl_key:                "Source_SSL_Key",
	replica_status_ssl_verify_server_cert: "Source_SSL_Verify_Server_Cert",
	replica_status_tls_version:            "Source_TLS_Version",
	replica_status_tls_ciphersuites:       "Source_TLS_Ciphersuites",
	replica_status_public_key_path:        "Source_public_key_path",
	replica_status_get_source_public_key:  "Get_source_public_key",
}

var queryStringProvider80 = QueryStringProvider{
	queries: queryStrings80,
}

var queryStringProvider8014 = QueryStringProvider{
	queries: queryStrings8014,
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
	if version >= "8.0.14" {
		return queryStringProvider8014
	}

	return queryStringProvider80
}
