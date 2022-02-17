job "many-echos" {
  datacenters = ["dc1"]

  type = "batch"

  group "maingroup" {
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
    task "hello-from-docker" {
      driver = "docker"
      config {
        image = "hello-world"
      }
    }
  }

  group "another" {
    count = 1
    task "before-another" {
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
      resources {
        # Reserve a Doyota Mius to drive around with for the duration of this task
        device "doyota/car/mius" {}
      }
    }
    task "this-is-familiar" {
      driver = "docker"
      config {
        image = "hello-world"
      }
    }
  }
}
