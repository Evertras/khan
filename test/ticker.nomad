job "ticker" {
  datacenters = ["dc1"]

  type = "batch"

  group "tickergroup" {
    count = 1
    task "ticktask" {
      driver = "raw_exec"
      config {
        command = "bash"
        args = [
          "-c",
          "while true; do sleep 1s; date; done",
        ]
      }
    }
  }
}
