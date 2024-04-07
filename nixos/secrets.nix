let
  devKeys = [
    # Foehammer 
    "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIAP0Hp1XEQ0Fx4rpuBuH7YusGU1f6s3Q4Sb1O1U0oKWA foehammer@disroot.org"
  ];

  machineKeys = [
    # Backend
    "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIMhvdHu6oPcTN5UJTeL0AbqXF1Dx3+0CC6D6Zp5Knj3j"
  ];

  keys = devKeys ++ machineKeys;
in
rec {
  "secrets/ts_key.age".publicKeys = keys;
}

