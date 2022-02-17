job "be-google" {
  datacenters = ["dc1"]

  type = "batch"

  group "grab-stuff" {
    count = 1
    task "steal-google" {
      driver = "raw_exec"
      config {
        command = "ls"
        args = [
          "-R",
        ]
      }

      artifact {
        source      = "https://www.google.com"
        destination = "local/index.html"
      }

      template {
        data = <<EOH
        echo hello template world
        EOH

        destination = "local/hello-world.sh"
      }
    }
  }
}
