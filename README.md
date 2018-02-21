# mackerel-plugin-nature-remo

[Nature Remo](https://nature.global/) custom metrics plugin for mackerel.io agent.
This plugin can send temperature and humidity detected by your Nature Remo.

## Synopsis

```
mackerel-plugin-nature-remo -access-key=<access-key>
```

## Example of mackerel-agent.conf

```
[plugin.metrics.NatureRemo]
command = "/path/to/mackerel-plugin-nature-remo -access-key=..."
```
