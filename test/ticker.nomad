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
          "-ec",
          "while true; do sleep 1; date; done",
        ]
      }
    }
  }
}
