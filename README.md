# telegraf-exec-nut

This is a simple tool to extract [Network UPS Tools](https://networkupstools.org/) (NUT) output and output
[Influx line protocol](https://docs.influxdata.com/influxdb/cloud/reference/syntax/line-protocol/);
it is designed to be used with a
[telegraf exec plugin](https://github.com/influxdata/telegraf/tree/master/plugins/inputs/exec).

This parses the output of the NUT TCP daemon and has been developed against Ubuntu 21.10.

## Interactive Run Example

The compiled tool can be run interactively.

```bash
./telegraf-exec-nut -help

Usage of telegraf-exec-nut:
  -host string
        NUT host (default "localhost")
  -password string
        NUT password
  -port int
        NUT port (default 3493)
  -username string
        NUT username
```

## Telegraf Run Example

This is a sample telegraf exec input that assumes the binary has been installed
to `/usr/local/bin/telegraf-exec-nut`:

```toml
[[inputs.exec]]                                                                 
  commands = ["/usr/local/bin/telegraf-exec-nut -username upsmon -password hunter2"]
  timeout = "5s"                                                                
  data_format = "influx"      
```

Then in InfluxDB, the `nut` measurement will have these tags:

```
battery.type
device.mfr
device.model
device.serial
device.type
```

And these fields:

```
battery.charge
battery.charge.low
battery.runtime
input.frequency
input.transfer.high
input.transfer.low
input.voltage
output.frequency
output.frequency.nominal
output.voltage
output.voltage.nominal
ups.beeper.status
ups.delay.shutdown
ups.delay.start
ups.firmware
ups.load
ups.power
ups.power.nominal
ups.realpower
ups.status
ups.timer.shutdown
ups.timer.start
```
