# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

A Dwarf Fortress-inspired colony sim written purely in Go, using the [Ebiten](https://ebitengine.org/) 2D game engine and the [mizu](https://github.com/sedyh/mizu) ECS framework. Ships as a native binary and as WASM for the browser.

## Commands

```bash
make build              # build native binary -> ./DwarvenFortresses
make wasm                # build web/game.wasm + copy wasm_exec.js into web/
go run .                 # run directly without building
go vet ./...
staticcheck ./...        # golint ./... is also run in CI
go test -race -vet=off ./...
```

CI (`.github/workflows/audit.yml`) runs `go mod verify`, `go build`, `go vet`, `staticcheck`, `golint`, and `go test` on every push/PR to `master`. There are currently no `_test.go` files in the repo, so `go test` is a no-op until tests are added.

No linter config beyond default `go vet`/`staticcheck`/`golint`.

## Architecture

### ECS via mizu

`mizu` is a plain-struct ECS: components are structs registered with `w.AddComponents(...)`, entities are Go structs that embed the components they have (see `entities/*.go`), and systems are structs implementing `Update(w engine.World)` and/or `Draw(w engine.World, screen *ebiten.Image)`.

Every system in this codebase is an **empty struct** (no component-pointer fields), which means mizu calls its `Update`/`Draw` once per frame (not per-entity) and the system does its own entity iteration via `w.View(ComponentA{}, ComponentB{}, ...).Each(func(e engine.Entity) {...})`. Inside `Each`, call `e.Get(&ptrA, &ptrB, ...)` to fetch pointers to that entity's components for that view.

All systems and entities are wired up once, in `scenes/game.go` (`Game.Setup`):
- `w.AddComponents(...)` registers every component type up front.
- `w.AddSystems(...)` registers systems **in the exact order they run**. Both `Update` and `Draw` are called on every system, every frame, in this registration order (see `world.Update`/`world.Draw` in the mizu package) — so a system later in the list sees state already mutated by an earlier one in the same frame. `systems.NewPathfinder()` runs before `systems.NewActor()`, for example, so actor position for the frame is already advanced by the time the actor system's own `Update`/`Draw` runs.
- Starting entities (admin/world singletons, dwarves, mouse, camera, GUI buttons) are created here too.

There's a single scene (`scenes/game.go`); `scenes/mainmenu.go` is the entry scene before it.

### Singleton components

Global/shared state is modeled as components on a single "Admin" entity (`entities/admin.go`), not as package-level state: `GameMapSingleton`, `RenderSingleton`, `InputSingleton`, `NatureSingleton`. Fetch them via helpers like `helpers.GetGameMapSingleton(w)` rather than passing state around manually.

### Map, grid and pathfinding

- The game world is voxel-like: one 2D grid per Z-level. `GameMapSingleton.Grids` is `[]*paths.Grid` (one entry per level), backed by the external `github.com/OpenSauce/paths` library (see `go.mod` `replace` — it's actually pointed at `github.com/tomknightdev/paths`, a fork). That library provides `Grid`, `Cell` (with `Walkable`/`Cost`), and A*-based `Path`.
- `components.Path` wraps `*paths.Path` plus a `Level` (Z) field, since a route can cross floors (e.g. via stairs) and is represented as `[]components.Path`, one segment per level.
- `helpers/pathfinding.go` builds paths (`GetPath`, `GetPathToAdjacent`) using the grid for the relevant Z-level(s).
- `helpers/region.go` maintains a **region index**: a flood-filled partition of each level's walkable cells into connected-component IDs (`GameMapSingleton.RegionIDs`), used to cheaply reject unreachable jobs without running full A*. It's lazily rebuilt (`EnsureRegions`) whenever `MarkRegionDirty` is called (any code path that flips a cell's `Walkable` bit must call this — it also clears `Blocked` on open tasks so previously-unreachable jobs get re-evaluated).
- `assets.TileSize` (loaded from the tileset definition) is the pixel size of one grid cell; world/pixel position = `components.Position{X,Y,Z}` (tile coords, `Z` = level) times `TileSize`.

### Movement

`systems/pathfinder.go` drives per-entity movement for anything with `Move` + `Position` + `Inventory`. It's tile-stepped and energy-gated, not continuous:
- Each `Move` has `CurrentEnergy`/`TotalEnergy`; energy increments once per frame and a step (one grid cell) is only taken once `CurrentEnergy >= TotalEnergy + inventory.Weight` (carrying items slows movement).
- `Move.CurrentPaths []components.Path` holds the remaining route (current segment at index 0); a step calls `.Next()`/`.Advance()` on the embedded `*paths.Path` and snaps `Position` straight to the next cell (no interpolation currently — there's no sub-tile/visual position, `Position` *is* the render position).
- `Move.Adjacent` distinguishes "path to the exact target tile" vs "path to any tile adjacent to it" (used for jobs like pick-up/drop where the actor works from beside the target rather than standing on it).
- If a route becomes unreachable mid-way, the pathfinder system unclaims the job (`Job.ClaimedByID = 0`) and re-adds it as a fresh entity so another worker can pick it up.

### Jobs / tasks

- A `Job` (`components.Job`) is an ordered list of `*Task`s (`components.Task`), each with a `TaskTypeEnum` (see `enums`), a target `Position`, and a required action count.
- `systems/actor.go` (`Actor.Update`) is the worker AI: if idle, it asks `helpers.GetJob` for the nearest claimable job + route; once "arrived" (per `Move.Arrived`) it ticks `Task.ActionsComplete` up to `Task.RequiredActions`, then completes the task and advances to the job's next task (re-routing `Move` to the next task's position), or clears the job when all tasks are done.
- `systems/job.go`, `systems/task.go`, `systems/designations.go`, `systems/building.go`, `systems/item.go` create/mutate jobs for specific player actions (designating chop/mine/stair/stockpile areas, construction, item pickup/storage).

### Rendering

- `systems/render.go` + `helpers/gamemap.go` (`DrawImage`/`DrawImages`) draw entities with a `Sprite` at `Position.X/Y * TileSize`, skipping anything whose `Position.Z` doesn't match the camera's current Z (only one level is visible at a time). Ebiten `DrawImageOptions` + an off-screen buffer (`RenderSingleton.OffScreen`) are used rather than drawing straight to `screen`.
- Camera (`entities/camera.go`, `systems/camera.go`) has its own `Position` (its `Z` is "current visible level") and `Zoom`.
- GUI (`gui/`, `components/gui.go`, `systems/gui.go`) is built on `github.com/yohamta/furex/v2` flexbox layout, driven by `Flex`/`Gui` components rather than immediate-mode drawing.

### World generation

`worldgen/worldgen.go` procedurally generates the initial map (uses `github.com/ojrac/opensimplex-go` for noise).

### Assets

`assets/bundle.go` loads a tileset definition (JSON) plus spritesheet image(s) at startup and slices them into per-tile `*ebiten.Image`s in two variants: `OpaqueImages` and `TransImages` (`enums.TileTypeEnum` -> image), keyed off the tileset JSON's per-tile pixel coordinates.
