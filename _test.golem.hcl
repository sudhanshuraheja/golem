vars = {
    APP = "golem"
}

recipe "test-artifact-template" "local" {
    artifact {
        template {
            data = <<EOF
{
  "signing": {
    "default": {
      "expiry": "87600h",
      "usages": ["signing", "key encipherment", "server auth", "client auth", "{{.Vars.HASHI_PATH}}"]
    }
  }
}
            EOF
        }
        destination = "template.parsed1"
    }
    artifact {
        template {
            path = "./template.tpl"
        }
        destination = "template.parsed2"
    }
    artifact {
        template {
            path = "https://raw.githubusercontent.com/sudhanshuraheja/golem/main/template.tpl"
        }
        destination = "template.parsed3"
    }
}
