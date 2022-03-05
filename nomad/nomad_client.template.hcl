// golem/nomad/nomad_client_template.hcl
// More details about the config can be found here
// https://www.nomadproject.io/docs/configuration

bind_addr = "{{ .Vars.NOMAD_BIND_ADDRESS }}"
datacenter = "@golem.NOMAD_DC"
data_dir = "/opt/nomad"
log_level = "INFO"
region = "@golem.NOMAD_REGION"

addresses {
  http = "{{ .Vars.NOMAD_CLIENT_ADDRESSES }}"
}

client {
  enabled = true
  servers = [
{{- range $_, $s := (match "tags" "contains" "@golem.NOMAD_SERVER_TAG") -}}
    {{- if not ($s).PrivateIP -}}
    {{- else -}}
        "{{ ($s).PrivateIP }}",
    {{- end -}}
{{- end -}}
  ]

  host_volume "nomad_host" {
    path      = "/opt/nomad_host"
    read_only = false
  }
  network_interface = "@golem.NOMAD_CLIENT_NETWORK_INTERFACE"
}

tls {
  http = true
  rpc  = true

  ca_file   = "/etc/nomad.d/nomad-ca.pem"
  cert_file = "/etc/nomad.d/client.pem"
  key_file  = "/etc/nomad.d/client-key.pem"

  verify_server_hostname = true
  verify_https_client    = true
}

plugin "docker" {
  config {
    volumes {
      enabled = true
    }
  }
}

// This setup does not include consul
// consul {
//   address = "127.0.0.1:8500"
//   token = ""
//   auto_advertise = true
//   server_auto_join = true
//   client_auto_join = true
// }