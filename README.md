# zteOnu

A tool to open the factory / telnet mode on ZTE ONU devices via the `webFac`
interface. Type `./zteonu -h` for help.

## Build

```bash
go build -o zteonu .
```

## Usage

```bash
# Old method (default), open telnet on 192.168.1.1:8080
./zteonu -i 192.168.1.1

# New method, derive the SendInfo payload from the interface MAC
./zteonu -i 192.168.1.1 --new

# New method using a specific network interface's MAC
./zteonu -i 192.168.1.1 --new --iface en0

# New method using a custom MAC address for the SendInfo payload
./zteonu -i 192.168.1.1 --new -m 00:07:29:55:35:57

# Also enable permanent telnet (user: root, pass: Zte521)
./zteonu -i 192.168.1.1 --new --telnet
```

## Flags

| Flag       | Short | Default                   | Description                                                                          |
|------------|-------|---------------------------|--------------------------------------------------------------------------------------|
| `--user`   | `-u`  | `telecomadmin`            | factory mode auth username                                                           |
| `--pass`   | `-p`  | `nE7jA%5m`                | factory mode auth password                                                           |
| `--ip`     | `-i`  | `192.168.1.1`             | ONU ip address                                                                       |
| `--port`   |       | `8080`                    | ONU http port                                                                        |
| `--telnet` |       | `false`                   | permanent telnet (user: `root`, pass: `Zte521`)                                      |
| `--tp`     |       | `23`                      | ONU telnet port                                                                      |
| `--new`    |       | `false`                   | use the new method; the `SendInfo` payload is derived from the current interface MAC |
| `--iface`  |       | `""` (first non-loopback) | network interface to read the MAC from                                               |
| `--mac`    | `-m`  | `""`                      | custom client MAC used to derive the `SendInfo` payload (e.g. `00:07:29:55:35:57`); defaults to the interface MAC |

## Notes on `--new`

The new method sends a `SendInfo` payload that encodes the MAC address of a local network interface (see `app/factory`).
The device only authorizes MAC addresses it accepts, so:

- Use `--iface` to choose which interface's MAC is used (defaults to the first non-loopback interface with a valid MAC).
- Use `--mac` to supply a custom MAC directly; this overrides the interface MAC and is what the `SendInfo` payload is derived from.
- The device MAC must be one the device accepts. Historically the device accepted `00:07:29:55:35:57`; supply it via `--mac`
  or spoof the interface MAC (or use a device that accepts the current MAC) so the payload matches what the device expects.

The payload transformation is derived from reverse-engineering the device's verification VM: the 46-byte payload is 12 little-endian `uint16` values, each followed by two zero bytes,
with the MAC XORed into the value bytes (`info=12` = 12 values = 6 MAC bytes × 2).

## Star History

[![Star History Chart](https://api.star-history.com/svg?repos=Septrum101/zteOnu&type=Date)](https://star-history.com/#Septrum101/zteOnu&Date)
