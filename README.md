# GoEcsEngine
A simple game engine using Ebiten with ECS.


## Description

### Components
This package contains engine components used for displaying sprites and text and managing UI.

### Loader
This package contains functions for loading entities with components from a TOML file.

### Resources
This package contains engine resources. It includes screen dimensions, fonts, spritesheets and controls.

### States
This package contains functions for managing a state machine.

A state machine has a stack of game states, which can be changed via transitions. Only the systems of the top state are executed in the game loop.

This is useful for pausing game or changing a game level for example.

### Systems
This package contains engine systems used for displaying sprites and text and managing UI. They are run automatically on each frame.

### World
This package defines the world, a global structure containing game data (ECS manager, components and resources).

It is passed as a parameter in all system and state functions.


## Deserialization from a TOML file
The engine uses [a TOML parser](https://github.com/BurntSushi/toml) for reading TOML files. It uses the [TOML v0.4.0](https://github.com/toml-lang/toml/blob/master/versions/en/toml-v0.4.0.md) specification.

Deserialization is relatively straightforward, with TOML fields corresponding directly to components fields, with the exception of Text and SpriteRender components which need to load data dynamically.

See [examples/transform/assets/start.toml](examples/transform/assets/start.toml) or [loader/entity.go](loader/entity.go) for more details.


## Examples
Examples are included in the [examples](examples) directory.

For a more complete example, see the [Arkanoid](https://github.com/x-hgg-x/arkanoid-go) game or the [Space invaders](https://github.com/x-hgg-x/space-invaders-go) game.
