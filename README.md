[![Go Reference](https://pkg.go.dev/badge/github.com/tommzn/hdb-renderer-syncsign.svg)](https://pkg.go.dev/github.com/tommzn/hdb-renderer-syncsign)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/tommzn/hdb-renderer-syncsign)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/tommzn/hdb-renderer-syncsign)
[![Go Report Card](https://goreportcard.com/badge/github.com/tommzn/hdb-renderer-syncsign)](https://goreportcard.com/report/github.com/tommzn/hdb-renderer-syncsign)
[![Actions Status](https://github.com/tommzn/hdb-renderer-syncsign/actions/workflows/go.pkg.auto-ci.yml/badge.svg)](https://github.com/tommzn/hdb-renderer-syncsign/actions)

# HomeDashboard Renderer for SyncSign® eInk Displays
Renders listen to a data source and generates content for SyncSign® eInk displays.


# Renderers
All renderers implement the [Renderer interface](https://github.com/tommzn/hdb-renderer-core/blob/main/interfaces.go) to provide generic way for content generation.

## Response Renderer
Response renderers generates response payload used in SyncSign template servers, in JSON format. It's template used default structure required for SyncSign displays and provides a posibility for other renderers to add items.
### Config
Defines path to template file.
```yaml
hdb:
  response:
    template: "response.json"
```

## Item Renderers
Item renderes generates items which will be picked up by response renderer to gnereate a complete response for displays. This can be simple text, geometric shapes or icons.

### Timestamp
A timestamp renderer generate a single item with current timestamp. By default it's position is in the lower left corner. Uee NewTimestampRenderer to generate such a renderer.
#### Config
Defines path to template file.
```yaml
hdb:
  response:
    template: "response.json"
```

### Error
In case something went wrong during content genration, error renderer can be used to generate a suitable server response for an error. Use NewErrorRenderer for initialization.
#### Config
Defines path to template file.
```yaml
hdb:
  error
    template: "error.json"
```

### Indoor Climate
This renderer listen to a data source for indoo climate, which includes temperature, humidity and, depending on used sensor, battery status. Indoor climate data
can be processed for diferent devices and can be assigned by config to seperate rooms.
Same template is used for each room and all rooms will be displayed in a row until scrren width exceeds.
Initialized by NewIndoorClimateRenderer.
#### Config
Following example config contains all available config options for indoor climate renderer.
```yaml
hdb:
  indoorclimate:
    template: "indoorclimate.json"
    anchor: 
      x: 10
      y: 10
    size:
      height: 200
      width: 200
    border: 5
    rooms:
      - id: "1"
        name: "Room1"
        displayIndex: "0"
      - id: "2"
        name: "Room2"
        displayIndex: "1"
    devices:
      - id: "Device2"
        roomId: "1"
      - id: "Device1"
        roomId: "2"
```
##### Template
Config option to set template file which should be used to generate a single room element. This file will be reused for all rooms.
##### Anchor
An anchor defines the upper left corner of element for first room.
##### Size
Defines the entire size of a romm element which includes temperature, humidity and battery status icon.
##### Border
Defines a space in pixel between each room element. Border can be set in general for top, right, bottom and left or for each attribute separately.
##### Rooms 
List of room which should be displayed as single element on screen, DisplayIndex defines the order of rooms on the screen from left to right. Name will be displayed on screen and id 
is used to assign devices.
##### Devices
Each room needs at least one assigned device to be displayed on screen.

## General Config
### Tempalte Directory
Use following config to set directory of templates for all renderers. Default value is folder "templates" at runtime location.
```yaml
hdb:
  template_dir: "templates"
```

# Supported Display
Only 7.5 inch display is supported for HomeDashboard project.

# Links
- [SyncSign](https://sync-sign.com)
- [HomeDashboard Documentation](https://github.com/tommzn/hdb-docs/wiki)
