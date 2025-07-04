# Custom Column Order

The `g` command supports customizing the column order in long format display mode through both command-line options and configuration files.

## Supported Column Names

- `Permissions` - File permissions (e.g., drwxr-xr-x)
- `Size` - File size (e.g., 4.0 KiB)
- `Owner` - File owner (e.g., runner)
- `Group` - File group (e.g., docker)
- `Time` or `Time Modified` - Last modified time (e.g., Jul 04 02:04)
- `Time Created` - Creation time (when using --create flag)
- `Time Accessed` - Access time (when using --access flag)
- `Name` - File/directory name

## Command Line Usage

Use the `--order` flag with a comma-separated list of column names:

```bash
# Show only size and name
g --order="Size,Name" --long

# Reorder columns: name first, then size and permissions
g --order="Name,Size,Permissions" --long

# Full custom order
g --order="Permissions,Size,Owner,Group,Time Modified,Name" --long
```

## Configuration File Usage

Add an `Order` field to your `g.yaml` config file:

```yaml
Args:
  - long
  - no-update

Order:
  - Size
  - Name  
  - Permissions
```

## Behavior

- Command-line `--order` flag takes precedence over config file settings
- Invalid column names are silently ignored
- Flags like `-O` (no owner) and `-G` (no group) are respected and will exclude those columns even if specified in the order
- If no custom order is specified, the default order is used: `Permissions,Size,Owner,Group,Time,Name`
- Empty order specification falls back to default order

## Examples

```bash
# Default long format
g --long
# Output: drwxr-xr-x  4.0 KiB runner docker Jul 04 02:04 asset

# Size and name only
g --order="Size,Name" --long  
# Output:  4.0 KiB asset

# Name first
g --order="Name,Size,Permissions" --long
# Output: asset  4.0 KiB drwxr-xr-x

# Exclude owner with -O flag
g --order="Permissions,Size,Owner,Group,Name" -O --long
# Output: drwxr-xr-x  4.0 KiB docker asset
```