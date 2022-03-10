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
      "usages": ["signing", "key encipherment", "server auth", "client auth", "@golem.HASHI_PATH"]
    }
  }
}
            EOF
        }
        destination = "template.parsed1"
    }
    artifact {
        template {
            path = "./testdata/template.tpl"
        }
        destination = "template.parsed2"
    }
    artifact {
        template {
            path = "https://raw.githubusercontent.com/sudhanshuraheja/golem/main/testdata/template.tpl"
        }
        destination = "template.parsed3"
    }
    commands = [
        "rm template.parsed1",
        "rm template.parsed2",
        "rm template.parsed3",
    ]
}
