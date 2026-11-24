package inst

import (
	"database/sql"
	"strings"
	"testing"
)

func newDemotedSourceTestInstance() *Instance {
	instance := NewInstance()
	instance.Key = InstanceKey{Hostname: "host_1", Port: 3306}
	instance.Version = "8.4.0"
	instance.QSP = GetQueryStringProvider(instance.Version)
	return instance
}

func newTLSPolicyTestInstance(allowTLS bool, verifyServerCert sql.NullBool) *Instance {
	policy := NewInstance()
	policy.Key = InstanceKey{Hostname: "host_2", Port: 3306}
	policy.AllowTLS = allowTLS
	// Paths on the promoted replica must not be copied to the demoted source.
	policy.ReplicationSSLCAFile = "/host_2/ca.pem"
	policy.ReplicationSSLCert = "/host_2/client-cert.pem"
	policy.ReplicationSSLKey = "/host_2/client-key.pem"
	policy.ReplicationSSLVerifyServerCert = verifyServerCert
	return policy
}

func renderTLSFragments(instance *Instance) (string, []interface{}) {
	var params []string
	var args []interface{}
	appendReplicationChangeTLSFragments(instance, &params, &args)
	return strings.Join(params, ", "), args
}

func TestApplyDemotedSourceTLSProfileUsesOwnServerCertificates(t *testing.T) {
	for _, verify := range []bool{true, false} {
		instance := newDemotedSourceTestInstance()
		policy := newTLSPolicyTestInstance(true, sql.NullBool{Valid: true, Bool: verify})

		if !applyDemotedSourceTLSProfile(instance, policy, serverTLSConfiguration{CA: "/host_1/ca.pem", Cert: "/host_1/server-cert.pem", Key: "/host_1/server-key.pem"}) {
			t.Fatalf("expected TLS profile to be applied")
		}
		if !instance.AllowTLS {
			t.Fatalf("expected AllowTLS to be set")
		}
		if instance.ReplicationSSLCAFile != "/host_1/ca.pem" || instance.ReplicationSSLCert != "/host_1/server-cert.pem" || instance.ReplicationSSLKey != "/host_1/server-key.pem" {
			t.Fatalf("expected demoted source's own server certificates, got ca=%q cert=%q key=%q", instance.ReplicationSSLCAFile, instance.ReplicationSSLCert, instance.ReplicationSSLKey)
		}
		if !instance.ReplicationSSLVerifyServerCert.Valid || instance.ReplicationSSLVerifyServerCert.Bool != verify {
			t.Fatalf("expected verify_server_cert=%t copied from policy, got %+v", verify, instance.ReplicationSSLVerifyServerCert)
		}

		fragments, _ := renderTLSFragments(instance)
		for _, expected := range []string{"source_ssl = 1", "source_ssl_ca = ?", "source_ssl_cert = ?", "source_ssl_key = ?", "source_ssl_verify_server_cert = ?"} {
			if !strings.Contains(fragments, expected) {
				t.Fatalf("expected %q in %q", expected, fragments)
			}
		}
		expectedVerify := 0
		if verify {
			expectedVerify = 1
		}
		if value := tlsFragmentValues(instance)["source_ssl_verify_server_cert"]; value != expectedVerify {
			t.Fatalf("expected source_ssl_verify_server_cert=%d, got %v", expectedVerify, value)
		}
	}
}

func TestApplyDemotedSourceTLSProfileSetsVerifyServerCertExplicitly(t *testing.T) {
	instance := newDemotedSourceTestInstance()
	// The policy's flag is unknown: it must still be emitted explicitly, as 0.
	policy := newTLSPolicyTestInstance(true, sql.NullBool{Valid: false})

	applyDemotedSourceTLSProfile(instance, policy, serverTLSConfiguration{CA: "ca.pem", Cert: "server-cert.pem", Key: "server-key.pem"})
	if !instance.ReplicationSSLVerifyServerCert.Valid || instance.ReplicationSSLVerifyServerCert.Bool {
		t.Fatalf("expected explicit verify_server_cert=0, got %+v", instance.ReplicationSSLVerifyServerCert)
	}
	fragments, _ := renderTLSFragments(instance)
	if !strings.Contains(fragments, "source_ssl_verify_server_cert = ?") {
		t.Fatalf("expected verify_server_cert in %q", fragments)
	}
}

func TestApplyDemotedSourceTLSProfileNoTLSInTopology(t *testing.T) {
	instance := newDemotedSourceTestInstance()
	policy := newTLSPolicyTestInstance(false, sql.NullBool{Valid: true, Bool: false})

	if applyDemotedSourceTLSProfile(instance, policy, serverTLSConfiguration{CA: "ca.pem", Cert: "server-cert.pem", Key: "server-key.pem"}) {
		t.Fatalf("expected no TLS profile when the promoted replica did not use TLS")
	}
	if fragments, _ := renderTLSFragments(instance); fragments != "" {
		t.Fatalf("expected no TLS fragments, got %q", fragments)
	}
}

func TestApplyDemotedSourceTLSProfileNilPolicy(t *testing.T) {
	instance := newDemotedSourceTestInstance()
	if applyDemotedSourceTLSProfile(instance, nil, serverTLSConfiguration{CA: "ca.pem", Cert: "server-cert.pem", Key: "server-key.pem"}) {
		t.Fatalf("expected no TLS profile without a policy")
	}
	if instance.AllowTLS {
		t.Fatalf("expected AllowTLS to stay false")
	}
}

func TestApplyDemotedSourceTLSProfileKeepsExistingReplicaSettings(t *testing.T) {
	instance := newDemotedSourceTestInstance()
	instance.AllowTLS = true
	instance.ReplicationSSLCAFile = "/replica/ca.pem"
	instance.ReplicationSSLVerifyServerCert = sql.NullBool{Valid: true, Bool: false}
	policy := newTLSPolicyTestInstance(true, sql.NullBool{Valid: true, Bool: true})

	if applyDemotedSourceTLSProfile(instance, policy, serverTLSConfiguration{CA: "/host_1/ca.pem"}) {
		t.Fatalf("expected existing replica TLS settings to be preserved")
	}
	if instance.ReplicationSSLCAFile != "/replica/ca.pem" || instance.ReplicationSSLVerifyServerCert.Bool {
		t.Fatalf("existing replica TLS settings were changed: %+v", instance)
	}
}

// tlsFragmentValues maps each "option = ?" fragment to its argument.
func tlsFragmentValues(instance *Instance) map[string]interface{} {
	var params []string
	var args []interface{}
	appendReplicationChangeTLSFragments(instance, &params, &args)
	values := map[string]interface{}{}
	i := 0
	for _, param := range params {
		if name, found := strings.CutSuffix(param, " = ?"); found {
			values[name] = args[i]
			i++
		}
	}
	return values
}

func TestApplyDemotedSourceTLSProfileIncompleteServerConfiguration(t *testing.T) {
	instance := newDemotedSourceTestInstance()
	policy := newTLSPolicyTestInstance(true, sql.NullBool{Valid: true, Bool: false})

	applyDemotedSourceTLSProfile(instance, policy, serverTLSConfiguration{CA: "/host_1/ca.pem"})
	fragments, _ := renderTLSFragments(instance)
	if !strings.Contains(fragments, "source_ssl = 1") {
		t.Fatalf("expected TLS with available settings, got %q", fragments)
	}
	values := tlsFragmentValues(instance)
	for option, expected := range map[string]string{"source_ssl_ca": "/host_1/ca.pem", "source_ssl_cert": "", "source_ssl_key": ""} {
		if value, found := values[option]; !found || value != expected {
			t.Fatalf("expected %s = %q to be set explicitly, got %q (set: %t) in %q", option, expected, value, found, fragments)
		}
	}
}

// staleTLSPathsTestInstance returns an instance whose channel exists with SOURCE_SSL=0 but still
// reports the certificate paths it was configured with, including a CRL that no longer exists.
func staleTLSPathsTestInstance() *Instance {
	instance := newDemotedSourceTestInstance()
	instance.AllowTLS = false
	instance.ReplicationSSLCAFile = "/stale/ca.pem"
	instance.ReplicationSSLCAPath = "/stale/ca-dir"
	instance.ReplicationSSLCert = "/stale/cert.pem"
	instance.ReplicationSSLKey = "/stale/key.pem"
	instance.ReplicationSSLCRLFile = "/nonexistent/old-crl.pem"
	instance.ReplicationSSLCRLPath = "/stale/crl-dir"
	instance.ReplicationSSLCipher = "ECDHE-RSA-AES128-GCM-SHA256"
	instance.ReplicationTLSVersion = "TLSv1.2"
	instance.ReplicationTLSCiphersuites = "TLS_AES_128_GCM_SHA256"
	return instance
}

var certificatePathOptions = []string{"source_ssl_ca", "source_ssl_capath", "source_ssl_cert", "source_ssl_key", "source_ssl_crl", "source_ssl_crlpath"}

// Regression test: MySQL keeps options that CHANGE REPLICATION SOURCE TO does not mention, so when
// the server's ssl_* variables are empty, every stale certificate path must be cleared explicitly.
func TestApplyDemotedSourceTLSProfileClearsStalePathsOfDisabledChannel(t *testing.T) {
	instance := staleTLSPathsTestInstance()
	policy := newTLSPolicyTestInstance(true, sql.NullBool{Valid: true, Bool: true})

	if !applyDemotedSourceTLSProfile(instance, policy, serverTLSConfiguration{}) {
		t.Fatalf("expected TLS profile to be applied to a channel with TLS disabled")
	}
	values := tlsFragmentValues(instance)
	for _, option := range certificatePathOptions {
		if value, found := values[option]; !found || value != "" {
			t.Fatalf("expected %s = '' to clear the stale path, got %q (set: %t)", option, value, found)
		}
	}
}

// Regression test: a stale SOURCE_SSL_CRL prevents the connection even when the server certificates
// are valid. With valid server certificates and no server CRL, the CRL options must be cleared.
func TestApplyDemotedSourceTLSProfileClearsStaleCRL(t *testing.T) {
	instance := staleTLSPathsTestInstance()
	policy := newTLSPolicyTestInstance(true, sql.NullBool{Valid: true, Bool: true})

	applyDemotedSourceTLSProfile(instance, policy, serverTLSConfiguration{CA: "/host_1/ca.pem", Cert: "/host_1/server-cert.pem", Key: "/host_1/server-key.pem"})
	values := tlsFragmentValues(instance)
	expected := map[string]string{
		"source_ssl_ca":      "/host_1/ca.pem",
		"source_ssl_capath":  "",
		"source_ssl_cert":    "/host_1/server-cert.pem",
		"source_ssl_key":     "/host_1/server-key.pem",
		"source_ssl_crl":     "",
		"source_ssl_crlpath": "",
	}
	for option, want := range expected {
		if value, found := values[option]; !found || value != want {
			t.Fatalf("expected %s = %q, got %q (set: %t)", option, want, value, found)
		}
	}
}

// Server variables that are set replace the stale paths.
func TestApplyDemotedSourceTLSProfileReplacesStalePathsWithServerConfiguration(t *testing.T) {
	instance := staleTLSPathsTestInstance()
	policy := newTLSPolicyTestInstance(true, sql.NullBool{Valid: true, Bool: true})

	applyDemotedSourceTLSProfile(instance, policy, serverTLSConfiguration{
		CAPath:  "/host_1/ca-dir",
		Key:     "/host_1/server-key.pem",
		CRL:     "/host_1/crl.pem",
		CRLPath: "/host_1/crl-dir",
	})
	values := tlsFragmentValues(instance)
	expected := map[string]string{
		"source_ssl_ca":      "",
		"source_ssl_capath":  "/host_1/ca-dir",
		"source_ssl_cert":    "",
		"source_ssl_key":     "/host_1/server-key.pem",
		"source_ssl_crl":     "/host_1/crl.pem",
		"source_ssl_crlpath": "/host_1/crl-dir",
	}
	for option, want := range expected {
		if value, found := values[option]; !found || value != want {
			t.Fatalf("expected %s = %q, got %q (set: %t)", option, want, value, found)
		}
	}
}

// Regression test: a channel that exists with SOURCE_SSL=0 keeps its cipher and TLS version
// restrictions. For example SOURCE_TLS_VERSION='TLSv1.2' stops replication from a source that only
// accepts TLSv1.3. The restrictions come from the policy, and SOURCE_TLS_CIPHERSUITES, which SHOW
// REPLICA STATUS does not report, is reset to NULL (MySQL's default).
func TestApplyDemotedSourceTLSProfileReplacesStaleTLSRestrictions(t *testing.T) {
	instance := staleTLSPathsTestInstance()
	policy := newTLSPolicyTestInstance(true, sql.NullBool{Valid: true, Bool: true})

	applyDemotedSourceTLSProfile(instance, policy, serverTLSConfiguration{CA: "/host_1/ca.pem", Cert: "/host_1/server-cert.pem", Key: "/host_1/server-key.pem"})
	values := tlsFragmentValues(instance)
	for option, want := range map[string]interface{}{"source_ssl_cipher": "", "source_tls_version": "", "source_tls_ciphersuites": nil} {
		if value, found := values[option]; !found || value != want {
			t.Fatalf("expected %s = %v to replace the stale restriction, got %v (set: %t)", option, want, value, found)
		}
	}

	// The policy's own restrictions are copied.
	instance = staleTLSPathsTestInstance()
	policy.ReplicationSSLCipher = "ECDHE-RSA-AES256-GCM-SHA384"
	policy.ReplicationTLSVersion = "TLSv1.3"
	applyDemotedSourceTLSProfile(instance, policy, serverTLSConfiguration{CA: "/host_1/ca.pem", Cert: "/host_1/server-cert.pem", Key: "/host_1/server-key.pem"})
	values = tlsFragmentValues(instance)
	for option, want := range map[string]interface{}{"source_ssl_cipher": "ECDHE-RSA-AES256-GCM-SHA384", "source_tls_version": "TLSv1.3", "source_tls_ciphersuites": nil} {
		if value, found := values[option]; !found || value != want {
			t.Fatalf("expected %s = %v copied from the policy, got %v (set: %t)", option, want, value, found)
		}
	}
}

// SOURCE_TLS_VERSION and SOURCE_TLS_CIPHERSUITES are set only where CHANGE REPLICATION SOURCE TO
// accepts them; an unknown option would make the statement fail.
func TestReplicationChangeTLSFragmentsReplaceOnlySupportedOptions(t *testing.T) {
	policy := newTLSPolicyTestInstance(true, sql.NullBool{Valid: true, Bool: true})
	for _, c := range []struct {
		version         string
		tlsVersion      bool
		tlsCiphersuites bool
	}{
		{"8.4.8", true, true},
		{"8.0.36-28", true, true},
		{"8.0.19", true, true},
		{"8.0.18", true, false},
		{"5.7.44-log", true, false},
		{"5.7.9", false, false},
		{"10.6.12-MariaDB", false, false},
		{"", false, false},
	} {
		instance := newDemotedSourceTestInstance()
		instance.Version = c.version
		instance.QSP = GetQueryStringProvider(c.version)
		applyDemotedSourceTLSProfile(instance, policy, serverTLSConfiguration{CA: "ca.pem", Cert: "cert.pem", Key: "key.pem"})
		fragments, _ := renderTLSFragments(instance)
		if got := strings.Contains(fragments, "tls_version"); got != c.tlsVersion {
			t.Fatalf("version %q: expected TLS version option %t, got %q", c.version, c.tlsVersion, fragments)
		}
		if got := strings.Contains(fragments, "tls_ciphersuites"); got != c.tlsCiphersuites {
			t.Fatalf("version %q: expected TLS ciphersuites option %t, got %q", c.version, c.tlsCiphersuites, fragments)
		}
		for _, option := range []string{"ssl_ca", "ssl_capath", "ssl_cert", "ssl_key", "ssl_crl", "ssl_crlpath", "ssl_cipher"} {
			if !strings.Contains(fragments, "_"+option+" = ?") {
				t.Fatalf("version %q: expected %s in %q", c.version, option, fragments)
			}
		}
	}
}

func TestVersionAtLeast(t *testing.T) {
	for _, c := range []struct {
		version, minVersion string
		atLeast             bool
	}{
		{"8.0.19", "8.0.19", true},
		{"8.0.36-28", "8.0.19", true},
		{"8.0.18", "8.0.19", false},
		{"8.4.0", "8.0.19", true},
		{"5.7.10-log", "5.7.10", true},
		{"5.7.9", "5.7.10", false},
		{"10.6.12-MariaDB", "8.0.19", true},
		{"8.0", "8.0.19", false},
		{"", "5.7.10", false},
		{"unknown", "5.7.10", false},
	} {
		if got := versionAtLeast(c.version, c.minVersion); got != c.atLeast {
			t.Fatalf("versionAtLeast(%q, %q) = %t, expected %t", c.version, c.minVersion, got, c.atLeast)
		}
	}
}

// Credentials that turn TLS on for a channel that had it off replace its stale TLS options too.
func TestApplyCredentialsTLSProfileClearsStaleOptionsWhenEnablingTLS(t *testing.T) {
	instance := staleTLSPathsTestInstance()
	instance.ReplicationSSLVerifyServerCert = sql.NullBool{Valid: true, Bool: true}

	applyCredentialsTLSProfile(instance, &ReplicationCredentials{User: "repl", SSLCaCert: "/creds/ca.pem"})
	if !instance.AllowTLS {
		t.Fatalf("expected credentials with TLS material to enable TLS")
	}
	values := tlsFragmentValues(instance)
	expected := map[string]interface{}{
		"source_ssl_ca":                 "/creds/ca.pem",
		"source_ssl_capath":             "",
		"source_ssl_cert":               "",
		"source_ssl_key":                "",
		"source_ssl_crl":                "",
		"source_ssl_crlpath":            "",
		"source_ssl_cipher":             "",
		"source_tls_version":            "",
		"source_tls_ciphersuites":       nil,
		"source_ssl_verify_server_cert": 1,
	}
	for option, want := range expected {
		if value, found := values[option]; !found || value != want {
			t.Fatalf("expected %s = %v, got %v (set: %t)", option, want, value, found)
		}
	}
}

// Credentials for a channel that already uses TLS replace only what they supply.
func TestApplyCredentialsTLSProfileKeepsExistingTLSSettings(t *testing.T) {
	instance := newDemotedSourceTestInstance()
	instance.AllowTLS = true
	instance.ReplicationSSLCAFile = "/replica/ca.pem"
	instance.ReplicationSSLCRLFile = "/replica/crl.pem"
	instance.ReplicationTLSVersion = "TLSv1.3"

	applyCredentialsTLSProfile(instance, &ReplicationCredentials{User: "repl", SSLCert: "/creds/cert.pem"})
	values := tlsFragmentValues(instance)
	for option, want := range map[string]interface{}{"source_ssl_ca": "/replica/ca.pem", "source_ssl_crl": "/replica/crl.pem", "source_tls_version": "TLSv1.3", "source_ssl_cert": "/creds/cert.pem"} {
		if values[option] != want {
			t.Fatalf("expected %s = %v, got %+v", option, want, values)
		}
	}
	for _, option := range []string{"source_ssl_key", "source_ssl_capath", "source_tls_ciphersuites"} {
		if _, found := values[option]; found {
			t.Fatalf("expected %s to be left to MySQL, got %+v", option, values)
		}
	}

	// Credentials without TLS material change nothing.
	instance = staleTLSPathsTestInstance()
	applyCredentialsTLSProfile(instance, &ReplicationCredentials{User: "repl", Password: "secret"})
	if instance.AllowTLS || instance.ReplicationSSLCRLFile == "" {
		t.Fatalf("expected credentials without TLS material to leave TLS settings alone, got %+v", instance)
	}
}

func TestServerTLSConfigurationComplete(t *testing.T) {
	for _, c := range []struct {
		server   serverTLSConfiguration
		complete bool
	}{
		{serverTLSConfiguration{CA: "ca.pem", Cert: "cert.pem", Key: "key.pem"}, true},
		{serverTLSConfiguration{CAPath: "ca-dir", Cert: "cert.pem", Key: "key.pem"}, true},
		{serverTLSConfiguration{Cert: "cert.pem", Key: "key.pem", CRL: "crl.pem"}, false},
		{serverTLSConfiguration{CA: "ca.pem", Key: "key.pem"}, false},
		{serverTLSConfiguration{CA: "ca.pem", Cert: "cert.pem"}, false},
	} {
		if c.server.complete() != c.complete {
			t.Fatalf("expected complete()=%t for %s", c.complete, c.server)
		}
	}
}

// Without the policy, empty paths are still omitted, so ChangeMasterTo keeps whatever MySQL holds.
func TestReplicationChangeTLSFragmentsOmitEmptyPathsWithoutPolicy(t *testing.T) {
	instance := newDemotedSourceTestInstance()
	instance.AllowTLS = true
	instance.ReplicationSSLCAFile = "/replica/ca.pem"
	values := tlsFragmentValues(instance)
	if values["source_ssl_ca"] != "/replica/ca.pem" {
		t.Fatalf("expected source_ssl_ca to be kept, got %+v", values)
	}
	for _, option := range append(certificatePathOptions[1:], "source_ssl_cipher", "source_tls_version", "source_tls_ciphersuites") {
		if _, found := values[option]; found {
			t.Fatalf("expected empty %s to be omitted, got %+v", option, values)
		}
	}
}

func newTestReplicationCredentials() *ReplicationCredentials {
	return &ReplicationCredentials{
		User:      "repl",
		Password:  "secret",
		SSLCaCert: "/host_2/ca.pem",
		SSLCert:   "/host_2/client-cert.pem",
		SSLKey:    "/host_2/client-key.pem",
	}
}

func TestReplicationCredentialsForTLSPolicySeedsOnlyUserAndPassword(t *testing.T) {
	instance := newDemotedSourceTestInstance()
	policy := newTLSPolicyTestInstance(true, sql.NullBool{Valid: true, Bool: true})
	creds := newTestReplicationCredentials()

	seeded := ReplicationCredentialsForTLSPolicy(creds, instance, policy)
	if seeded.User != "repl" || seeded.Password != "secret" {
		t.Fatalf("expected user and password to be kept, got %+v", seeded)
	}
	if seeded.SSLCaCert != "" || seeded.SSLCert != "" || seeded.SSLKey != "" {
		t.Fatalf("expected TLS fields to be dropped, got %+v", seeded)
	}
	if creds.SSLCaCert == "" || creds.SSLCert == "" || creds.SSLKey == "" {
		t.Fatalf("expected the original credentials to be left unchanged, got %+v", creds)
	}

	// The seeded credentials must not enable TLS on the instance, so the policy still applies.
	if !applyDemotedSourceTLSProfile(instance, policy, serverTLSConfiguration{CA: "/host_1/ca.pem", Cert: "/host_1/server-cert.pem", Key: "/host_1/server-key.pem"}) {
		t.Fatalf("expected TLS profile to be applied after seeding credentials")
	}
	if instance.ReplicationSSLCert != "/host_1/server-cert.pem" || !instance.ReplicationSSLVerifyServerCert.Bool {
		t.Fatalf("expected own certificates and copied verify flag, got cert=%q verify=%+v", instance.ReplicationSSLCert, instance.ReplicationSSLVerifyServerCert)
	}
}

func TestReplicationCredentialsForTLSPolicyKeepsCredentialsWhenPolicyDoesNotApply(t *testing.T) {
	// The replica's channel does not use TLS: credentials, including TLS fields, are used as before.
	creds := newTestReplicationCredentials()
	if seeded := ReplicationCredentialsForTLSPolicy(creds, newDemotedSourceTestInstance(), newTLSPolicyTestInstance(false, sql.NullBool{})); seeded != creds {
		t.Fatalf("expected credentials unchanged when the policy's channel does not use TLS, got %+v", seeded)
	}

	// The instance already has TLS settings of its own: they are preserved, credentials unchanged.
	instance := newDemotedSourceTestInstance()
	instance.AllowTLS = true
	instance.ReplicationSSLCAFile = "/host_1/replication-ca.pem"
	if seeded := ReplicationCredentialsForTLSPolicy(creds, instance, newTLSPolicyTestInstance(true, sql.NullBool{Valid: true, Bool: true})); seeded != creds {
		t.Fatalf("expected credentials unchanged when the instance has TLS settings, got %+v", seeded)
	}

	// No policy: ChangeMasterTo behaviour.
	if seeded := ReplicationCredentialsForTLSPolicy(creds, newDemotedSourceTestInstance(), nil); seeded != creds {
		t.Fatalf("expected credentials unchanged without a policy, got %+v", seeded)
	}
	if seeded := ReplicationCredentialsForTLSPolicy(nil, newDemotedSourceTestInstance(), newTLSPolicyTestInstance(true, sql.NullBool{})); seeded != nil {
		t.Fatalf("expected nil credentials to stay nil, got %+v", seeded)
	}
}
