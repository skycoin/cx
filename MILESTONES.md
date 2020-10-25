# CX Chain Milestones

## v1.0 (Initial release)
- [ ] Simplify chain creation **(Due: 26/10/2020 Monday)**.
    - [x] Simplify cx chain creation to 2 commands.
    - [x] Command for manual connection management.
    - [x] Use base64 encoding for program state to save space.
    - [x] CLI should provide ENV help menu and global flag support.
    - [ ] Implement retry logic for obtaining chain program state. This will avoid failures with tx injection.
    - [ ] Documentation updates and README cleanups.
    - [ ] Remove legacy code.
    - [ ] Fix merge conflicts.
    
- [ ] Integration with CX Tracker **(Due: 02/11/2020 Monday)**.
    - [ ] Implement method to reference and register chain spec.
        - [ ] Subcommand `cxchain-cli register-spec`.
    - [ ] Implement method to auto-register CX nodes.
        - [ ] have CX Tracker serves peer list that gets updated. 3 Sections:
            1. Trusted nodes (publisher nodes).
            2. App nodes (nodes that possess an app address's secret key).
            3. Other nodes (other client nodes).
        - [ ] Nodes should auto-register themselves via the CX Tracker.
        - [ ] Modify chain spec to have CX Tracker URL instead of "peer_list_url".
        
- [ ] Integration with Dmsg **(Due: 09/11/2020 Monday)**.
    - [ ] App nodes should expose an endpoint via dmsg.
    - [ ] Logic for other nodes to connect to app nodes.
    - [ ] To reference an app node: `{chain_pk}.{app_pk}`.

## v1.1 (Backlogged changes)

- [ ] Allow 0-coin transactions (Switch cx transactions to be 0-coin transactions).
- [ ] Fix `MainExpression` representation in transactions (allows for more secure transaction verification).
- [ ] Rework `pex`, `daemon` and associated modules to be able to manually add/remove peers/connections.
- [ ] Improve CLI interface further: add descriptions.