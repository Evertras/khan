job "example" {
  datacenters = ["dc1"]

  type = "service"

  group "cache" {
    count = 1

    network {
      port "db" {
        to = 6379
      }
    }

    service {
      name = "redis-cache"
      tags = ["global", "cache"]
      port = "db"

    }

    restart {
      attempts = 2
      interval = "30m"

      delay = "15s"

      mode = "fail"
    }

    ephemeral_disk {
      size = 300
    }

    task "redis" {
      driver = "docker"

      config {
        image = "redis:3.2"

        ports = ["db"]
      }

      resources {
        cpu    = 500 # 500 MHz
        memory = 256 # 256MB
      }
    }
  }
}
