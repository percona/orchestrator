/*
   Copyright 2015 Shlomi Noach, courtesy Booking.com

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package inst

import (
	"fmt"

	"github.com/openark/golib/log"
	"github.com/openark/golib/sqlutils"
	"github.com/openark/orchestrator/go/config"
	"github.com/openark/orchestrator/go/db"
)

func IsSQLite() bool {
	return config.Config.IsSQLite()
}

// ReadClusterNameByAlias
func ReadClusterNameByAlias(alias string) (clusterName string, err error) {
	query := `
		select
			cluster_name
		from
			cluster_alias
		where
			alias = ?
			or cluster_name = ?
		`
	err = db.QueryOrchestrator(query, sqlutils.Args(alias, alias), func(m sqlutils.RowMap) error {
		clusterName = m.GetString("cluster_name")
		return nil
	})
	if err != nil {
		return "", err
	}
	if clusterName == "" {
		err = fmt.Errorf("No cluster found for alias %s", alias)
	}
	return clusterName, err
}

// DeduceClusterName attempts to resolve a cluster name given a name or alias.
// If unsuccessful to match by alias, the function returns the same given string
func DeduceClusterName(nameOrAlias string) (clusterName string, err error) {
	if nameOrAlias == "" {
		return "", fmt.Errorf("empty cluster name")
	}
	if name, err := ReadClusterNameByAlias(nameOrAlias); err == nil {
		return name, nil
	}
	return nameOrAlias, nil
}

// ReadAliasByClusterName returns the cluster alias for the given cluster name,
// or the cluster name itself if not explicit alias found
func ReadAliasByClusterName(clusterName string) (alias string, err error) {
	alias = clusterName // default return value
	query := `
		select
			alias
		from
			cluster_alias
		where
			cluster_name = ?
		`
	err = db.QueryOrchestrator(query, sqlutils.Args(clusterName), func(m sqlutils.RowMap) error {
		alias = m.GetString("alias")
		return nil
	})
	return clusterName, err
}

// WriteClusterAlias will write (and override) a single cluster name mapping
func writeClusterAlias(clusterName string, alias string) error {
	writeFunc := func() error {
		_, err := db.ExecOrchestrator(`
			replace into
					cluster_alias (cluster_name, alias, last_registered)
				values
					(?, ?, now())
			`,
			clusterName, alias)
		return log.Errore(err)
	}
	return ExecDBWriteFunc(writeFunc)
}

// writeClusterAliasManualOverride will write (and override) a single cluster name mapping
func writeClusterAliasManualOverride(clusterName string, alias string) error {
	writeFunc := func() error {
		_, err := db.ExecOrchestrator(`
			replace into
					cluster_alias_override (cluster_name, alias)
				values
					(?, ?)
			`,
			clusterName, alias)
		return log.Errore(err)
	}
	return ExecDBWriteFunc(writeFunc)
}

// Original, safe approach, which uses REPLACE INTO
func updateClusterAliasesUsingReplace() error {
	_, err := db.ExecOrchestrator(`
	replace into
			cluster_alias (alias, cluster_name, last_registered)
		select
			suggested_cluster_alias,
				cluster_name,
				now()
			from
			database_instance
			left join database_instance_downtime using (hostname, port)
		  where
			suggested_cluster_alias!=''
				/* exclude newly demoted, downtimed masters */
				and ifnull(
						database_instance_downtime.downtime_active = 1
						and database_instance_downtime.end_timestamp > now()
						and database_instance_downtime.reason = ?
					, 0) = 0
			order by
				ifnull(last_checked <= last_seen, 0) asc,
				read_only desc,
				num_slave_hosts asc
	`, DowntimeLostInRecoveryMessage)

	return err
}

// Optimized approach using INSERT INTO ... ON DUPLICATE KEY UPDATE
// While this approach is much faster and works in most cases, it is not
// guaranteed to be working in every case.
// cluster_alias table has two unique indexes:
// 1. primary on cluster_name column
// 2. alias_uidx on alias colum
//
// The data which is going to be inserted originates from database_instance
// table, in particular the following columns:
// 1. `cluster_name` varchar(128) NOT NULL
// 2. `suggested_cluster_alias` varchar(128) CHARACTER SET ascii COLLATE
//	   ascii_general_ci NOT NULL
//
// So it is possible to end up in the following situation when we use this
// approach:
// create table t1 (a int primary key, b int, unique key (b));
// insert into t1 values (0, 1);
// insert into t1 values (1, 2);
// insert into t1 values (0, 2) on duplicate key update a=0, b=2;
// ERROR 1062 (23000): Duplicate entry '2' for key 't1.b'
func updateClusterAliasesUsingInsert() error {
	_, err := db.ExecOrchestrator(`
	INSERT INTO cluster_alias
	(
		alias,
		cluster_name,
		last_registered
	)
	SELECT    di.suggested_cluster_alias,
		  di.cluster_name,
		  now()
	FROM      database_instance di
	LEFT JOIN database_instance_downtime did
	USING     (hostname, port)
	WHERE     di.suggested_cluster_alias != ''
	/* exclude newly demoted, downtimed masters */
	AND       Ifnull(did.downtime_active = 1
		  AND did.end_timestamp > Now()
		  AND did.reason = ?, 0) = 0
	ORDER BY  Ifnull(di.last_checked <= di.last_seen, 0) ASC,
		  di.read_only DESC,
		  di.num_slave_hosts ASC
	ON DUPLICATE KEY
	UPDATE alias = di.suggested_cluster_alias,
	cluster_name = di.cluster_name,
	last_registered = now()
	`, DowntimeLostInRecoveryMessage)

	return err
}

// UpdateClusterAliases writes down the cluster_alias table based on information
// gained from database_instance
func UpdateClusterAliases() error {
	writeFunc := func() error {
		var err error
		if IsSQLite() {
			// Sql lite backend
			err = updateClusterAliasesUsingReplace()
		} else {
			// MySQL backend (Orchestrator supports only SQLite and MySQL backends)
			// INSERT ON DUPLICATE KEY UPDATE is more performant than REPLACE in MySQL
			err = updateClusterAliasesUsingInsert()
			if err != nil {
				// Fallback to the original, safe implementation
				err = updateClusterAliasesUsingReplace()
			}
		}
		return log.Errore(err)
	}
	if err := ExecDBWriteFunc(writeFunc); err != nil {
		return err
	}
	writeFunc = func() error {
		// Handling the case where no cluster alias exists: we write a dummy alias in the form of the real cluster name.
		_, err := db.ExecOrchestrator(`
			replace into
					cluster_alias (alias, cluster_name, last_registered)
				select
						cluster_name as alias, cluster_name, now()
				  from
				    database_instance
				  group by
				    cluster_name
					having
						sum(suggested_cluster_alias = '') = count(*)
			`)
		return log.Errore(err)
	}
	if err := ExecDBWriteFunc(writeFunc); err != nil {
		return err
	}
	return nil
}

// ForgetLongUnseenClusterAliases will remove entries of cluster_aliases that have long since been last seen.
// This function is compatible with ForgetLongUnseenInstances
func ForgetLongUnseenClusterAliases() error {
	sqlResult, err := db.ExecOrchestrator(`
			delete
				from cluster_alias
			where
			last_registered < NOW() - interval ? hour`,
		config.Config.UnseenInstanceForgetHours,
	)
	if err != nil {
		return log.Errore(err)
	}
	rows, err := sqlResult.RowsAffected()
	if err != nil {
		return log.Errore(err)
	}
	AuditOperation("forget-clustr-aliases", nil, fmt.Sprintf("Forgotten aliases: %d", rows))
	return err
}

// ReplaceAliasClusterName replaces alis mapping of one cluster name onto a new cluster name.
// Used in topology failover/recovery
func ReplaceAliasClusterName(oldClusterName string, newClusterName string) (err error) {
	{
		writeFunc := func() error {
			_, err := db.ExecOrchestrator(`
			update cluster_alias
				set cluster_name = ?
				where cluster_name = ?
			`,
				newClusterName, oldClusterName)
			return log.Errore(err)
		}
		err = ExecDBWriteFunc(writeFunc)
	}
	{
		writeFunc := func() error {
			_, err := db.ExecOrchestrator(`
			update cluster_alias_override
				set cluster_name = ?
				where cluster_name = ?
			`,
				newClusterName, oldClusterName)
			return log.Errore(err)
		}
		if ferr := ExecDBWriteFunc(writeFunc); ferr != nil {
			err = ferr
		}
	}
	return err
}

// ReadUnambiguousSuggestedClusterAliases reads potential master hostname:port who have suggested cluster aliases,
// where no one else shares said suggested cluster alias. Such hostname:port are likely true owners
// of the alias.
func ReadUnambiguousSuggestedClusterAliases() (result map[string]InstanceKey, err error) {
	result = map[string]InstanceKey{}

	query := `
		select
			suggested_cluster_alias,
			min(hostname) as hostname,
			min(port) as port
		from
			database_instance
		where
			suggested_cluster_alias != ''
			and replication_depth=0
		group by
			suggested_cluster_alias
		having
			count(*) = 1
		`
	err = db.QueryOrchestrator(query, sqlutils.Args(), func(m sqlutils.RowMap) error {
		key := InstanceKey{Hostname: m.GetString("hostname"), Port: m.GetInt("port")}
		suggestedAlias := m.GetString("suggested_cluster_alias")
		result[suggestedAlias] = key
		return nil
	})
	return result, err
}
