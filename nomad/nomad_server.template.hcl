// golem/nomad/nomad_server_template.hcl
// More details about the config can be found here
// https://www.nomadproject.io/docs/configuration

bind_addr = "{{- .Vars.NOMAD_BIND_ADDRESS -}}"
datacenter = "@golem.NOMAD_DC"
data_dir = "/opt/nomad"
log_level = "INFO"
region = "@golem.NOMAD_REGION"

addresses {
  http = "{{- .Vars.NOMAD_SERVER_ADDRESSES -}}"
}

server {
  bootstrap_expect = @golem.NOMAD_BOOTSTRAP_EXPECT
  enabled = true
  data_dir = "/opt/nomad/"
  encrypt = "@golem.NOMAD_GOSSIP_KEY"
  server_join {
    retry_join = [
{{- range $_, $s := (match "tags" "contains" "@golem.NOMAD_SERVER_TAG") -}}
    {{- if not ($s).PrivateIP -}}
    {{- else -}}
        "{{ ($s).PrivateIP }}",
    {{- end -}}
{{- end -}}
    ]
  }
}

tls {
  http = true
  rpc  = true

  ca_file   = "/etc/nomad.d/nomad-ca.pem"
  cert_file = "/etc/nomad.d/server.pem"
  key_file  = "/etc/nomad.d/server-key.pem"

  verify_server_hostname = true
  verify_https_client    = true
}

// This setup does not include consul
// consul {
//   address = "127.0.0.1:8500"
//   token = "add_actual_token"
//   auto_advertise = true
//   server_auto_join = true
//   client_auto_join = true
// }