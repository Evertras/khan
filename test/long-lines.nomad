job "long-lines" {
  datacenters = ["dc1"]

  type = "batch"

  group "long-group" {
    count = 1
    task "long-task-is-long" {
      driver = "docker"
      config {
        image = "alpine:latest"
        args = ["echo", "Hello this is a really long and has a lot of stuff going on but it's a good thing we have wrapping in our logs or you would miss this very special offer buy now for just $9.99 for 99 months glhf it really is amazing how much text can fit onto a reasonably sized monitor hopefully this is actually going to wrap on your monitor if it doesn't that's pretty cool that you have such a ridiculously large monitor but most of us don't but still we really need to handle those edge cases and it would be pretty silly to ask you to make your terminal really thin just because you have an obnoxiously wide monitor that no mortal should wield"]
      }
    }
  }
}
