Registering CX chains on the tracker can be done via 2 scenarios:

**1. Run from tracker**
Run the chain with two arguments - config path and config hash (genesis hash). example:

```
cx --publisher \
    --config-path /path/to/cx-config.json \
    --config-hash $GENESIS_SIGNATURE (this is the tracker hash shown on the web app)
```

Config-hash can be removed if the config is already locally downloaded (after first time running).

**2. Run new app**
Create a json config with the following details:

```
{
        "genesisHash": "yourhash",
        "genesisAddress": "youraddress",
        "publicKey": "yourpubkey",
        "secretKey": "yoursecretkey"
}
```

A real value example would be:

```
{
        "genesisHash": "a439e7e0db4277d777a73a2887112b897665fd5dc681943b1b987594fb23432548bd52738a58e8afaf1c85e182cc321c6f95a3b2827f5e9775928a752d4480f601",
        "genesisAddress": "2Z2usqLxYXcepULXU681EGmHzmPK4DDxG8T",
        "publicKey": "03d92088691e72b25a88a649e1afc5d117489bfacb4f0b1f53c85b785d820da0db",
        "secretKey": "a77caae3b5d7bf50e7692789d35d698751ba03122bf5eb1f22cfefd59c752451"
}
```

Run the chain with one argument - config path. example:

```
cx --publisher \
    --config-path /path/to/your/cx-config.json
```