job "echo" {
  datacenters = ["dc1"]

  type = "batch"

  group "example" {
    count = 1
    task "before" {
      lifecycle {
        hook = "prestart"
        sidecar = false
      }
      driver = "docker"
      config {
        image = "alpine"
        command = "echo"
        args = ["Hello", "from", "prestart"]
      }
    }
    task "uptime" {
      driver = "docker"
      config {
        image = "hello-world"
      }
    }
  }
}
