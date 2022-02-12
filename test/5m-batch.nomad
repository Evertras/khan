job "naptime" {
  datacenters = ["dc1"]

  type = "batch"

  group "nap-group" {
    count = 1
    task "sleep-5m" {
      driver = "docker"
      config {
        image = "alpine:latest"
        args = ["sleep", "5m"]
      }
    }
  }
}
