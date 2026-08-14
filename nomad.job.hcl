job "phalanx-proxy" {
  datacenters = ["dc1"]
  type        = "service"

  group "proxy" {
    count = 2

    network {
      mode = "bridge"
      port "proxy" {
        to = 8443
      }
    }

    service {
      name = "phalanx"
      port = "proxy"

      tags = ["proxy", "tls", "ingress"]

      check {
        name     = "alive"
        type     = "tcp"
        interval = "10s"
        timeout  = "2s"
      }

      # Connect integration for Consul service discovery could be defined here
      # connect {
      #   sidecar_service {}
      # }
    }

    task "phalanx" {
      driver = "docker"

      config {
        image = "phalanx:latest"
        ports = ["proxy"]

        # To load testing certs and config, we can mount a volume or use templates
        # volumes = [
        #   "local/config.yaml:/root/config.yaml",
        #   "local/certs:/root/certs"
        # ]
      }

      # Setup config and certificates from Consul KV or Vault in a real environment
      template {
        data = <<EOH
listen_addr: ":8443"
backends:
  "service1.local": "service1.service.consul:8080"
  "service2.local": "service2.service.consul:8080"
tls:
  cert_file: "/root/certs/server.crt"
  key_file: "/root/certs/server.key"
  ca_cert_file: "/root/certs/ca.crt"
EOH
        destination = "local/config.yaml"
      }

      resources {
        cpu    = 256 # MHz
        memory = 128 # MB
      }
    }
  }
}
