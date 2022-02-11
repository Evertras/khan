job "echo" {
  datacenters = ["dc1"]

  type = "batch"

  group "example" {
    count = 1
    task "uptime" {
      driver = "docker"
      config {
        image = "hello-world"
      }
    }
  }
}
