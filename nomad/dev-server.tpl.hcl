################################################################################
# Local test Nomad config
#
# The .tpl file is intended as a starting point.  The Makefile will update things
# such as the data dir which require an absolute path which cannot be written
# in git.  To regenerate it, delete the generated file and run 'make nomad-test-server'

# Replace with a local absolute path... handled by Makefile
data_dir  = "DATADIR"

bind_addr = "0.0.0.0"

server {
  enabled          = true
  bootstrap_expect = 1
}

client {
  enabled = true
}

plugin "raw_exec" {
  config {
    enabled = true
  }
}

