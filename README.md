# godac
[![CI](https://github.com/tmuelli/godac/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/tmuelli/godac/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/tmuelli/godac.svg)](https://pkg.go.dev/github.com/tmuelli/godac)

A terminal damage calculator for Pokémon battles for Generation 4 written in Go.

godac is built for people who are comfortable in a terminal. Learn a handful of shortcuts and you can set up a [Smogon damage-calc](https://github.com/smogon/damage-calc)-style calculation for Gen 4 battles without touching the mouse.

![godac demo: Poochyena vs. Poochyena](docs/demo.gif)

## Features

- Damage calculation following generation 4 formula with damage range and target's HP percentage
- Autofill default values on selecting a Pokémon (types, base stats, max IVS, ability, nature)
- Autofill move's data on select (base power, type, category and hits)
- Custom configurations available, e.g. different types, base stats or move power
- Simple field settings configuration covering general settings (weather, level or battle mode) and per-side effects (e.g. Reflect, Light Screen, etc.)
- Autocomplete on input typing
- Theme adaption for light and dark terminal backgrounds

## Quick Start

You need [Go](https://go.dev/dl/) 1.26.3 or newer and a terminal with Unicode support.

```bash
git clone https://github.com/tmuelli/godac
cd godac
go run .
```

or install the binary:

```bash
go install github.com/tmuelli/godac@latest
godac
```

For the best experience use a terminal with a size of around 110x36.

## Usage

Example calculation for Poochyena against Poochyena using Bite:

1. Use `l` to get into the config view of the left Pokémon. Press `enter` on the `Name` field to start editing, type `Poochyena`, and press `enter` again to confirm.
2. Navigate to the `Moves` tab using `tab` key. Edit the first move like in step 1 and confirm.
3. Use `r` to get into the config view of the right Pokémon. Follow same instructions as in step 1 to set `Poochyena` for the right side.
4. Return to the calc view pressing `c`. To start the calculation, select the move `Bite` followed by `enter`.

### Key bindings

| Key | Action |
| --- | --- |
| `c` | Calc view: choose a move and calculate |
| `l` / `r` | Config view for the left / right Pokémon |
| `f` | Field view: weather, gravity, side effects |
| `↑` `↓` `←` `→` | Move around |
| `enter` | Edit the selected field, toggle it, or calculate (Calc view) |
| `tab` / `shift+tab` | Next / previous tab (Config view) |
| `g` / `s` | General / Sides section (Field view) |
| `?` | Show all key bindings |
| `q` | Quit |

While editing a field, your input filters the suggestions. Press `tab` to complete a suggestion, `enter` to confirm, or `esc` to cancel.

## Motivation

This project is the capstone project for the [boot.dev](https://www.boot.dev) course for the backend developer path. While already working as a software developer, I had only ever done tutorials in Go. To learn the language and for using it as a capstone project I wanted to build something real, but which could also be useful for me.

Growing up with Pokémon and still enjoying to play it from time to time I came across the Smogon's web damage calculator. As a (lazy) computer user being way faster with the keyboard than the mouse, I decided to rebuild the damage calculator as terminal application with convenient shortcuts available, allowing much faster calculations once one learned the few shortcuts.

## How it works

```
main.go          main entry point, loading the data and starting TUI
calc/          calculation logic with damage formula
models/        types and embedded game data (in JSON format)
internal/
  tui/           terminal user interface
  utils/         helpers
```

- Calculation logic is separate from the UI. `calc` takes plain input structs and returns a result. It can be tested without a terminal.
- Game master data is embedded with `go:embed`. Species, moves, abilities, items, natures and types as JSON files in `models`.
- UI follows the Elm architecture of [Bubble Tea](https://github.com/charmbracelet/bubbletea): a model, an update function and view function for rendering. Styling uses [Lip Gloss](https://github.com/charmbracelet/lipgloss).
- Footer and text input come from [Bubbles](https://github.com/charmbracelet/bubbles).

## Testing

```bash
go test ./...
```

Tests cover calculation package only (basic damage, type immunities, abilities, items, weather, field and status effects, fixed-damage moves and multi-hit moves, custom move and type overrides and input validation).

TUI package does not have automated tests yet.

## Limitations and future features

Not supported yet:

- Other generations than gen 4
- Chip damage (e.g. hail, Leech Seed) is not included in the damage result yet
- In-game terrain can not be set yet (used for nature power base power)
- Saving / loading configuration settings
- Saving / loading custom overrides
- Multi-turn and KO calculation
- Convenience shortcuts for easy setup of max or min IVs / EVs

## Use as a library

godac's damage calculation behind the TUI can be imported as a Go library.

- `calc`: Prepare, validate and run calculations
- `models`: Contains data and input types and starting values (`DefaultUserInput`, `DefaultMoveInput`)

Install:

```bash
go get github.com/tmuelli/godac@v0.2.0
```

Example:

```go
package calc_test

import (
	"fmt"
	"log"

	"github.com/tmuelli/godac/calc"
	"github.com/tmuelli/godac/models"
)

func ExampleNewCalc() {
	data, err := models.NewData()
	if err != nil {
		log.Fatal(err)
	}

	poochyena := data.Species["Poochyena"]
	bite := data.Moves["Bite"]

	attacker := data.DefaultUserInput(&poochyena)
	attacker.AttackMove = bite.Name
	attacker.Moves = []string{bite.Name}
	attacker.InputSide = string(models.Attacker)

	defender := data.DefaultUserInput(&poochyena)
	defender.InputSide = string(models.Defender)

	field := models.UserFieldInput{
		Type:      string(models.Singles),
		Level:     int(models.Level100),
		Terrain:   string(models.Clear),
		IsGravity: false,
	}

	c, err := calc.NewCalc(data, &attacker, &defender, &field)
	if err != nil {
		log.Fatal(err)
	}
	if err := c.Calc(); err != nil {
		log.Fatal(err)
	}

	fmt.Println(c.Result.Range.Min, c.Result.Range.Max)
	// prints: 45 53
}

```


The result holds the damage range (`Result.Range.Min`/`Max`, in HP and in %)
and all damage rolls (`Result.DmgMatrix`).

See the
[package docs](https://pkg.go.dev/github.com/tmuelli/godac) for details.
The API may still change in `v0.x`, so pin the version.

## How AI was used

All the code was written by myself and has not been generated by AI. During the development process, LLMs were used for guiding, sketching ideas, making suggestions (especially UI look and feel), research and best practices on golang project structure, building bubbletea apps and styling UI.

One exception to the above is `calc/calc_test.go`, which was AI generated for creating test cases.

## Contributing

If you would like to contribute, please for the repository and open a pull request to the `main` branch.

## Credits and disclaimer

The game data and the damage mechanics are taken from [Smogon's damage-calc](https://github.com/smogon/damage-calc), which is MIT licensed. Generation 4 game data was exported from that repository to JSON files using a small helper script. Sincere thanks to the maintainers of this outstanding project.

godac is an unofficial fan project. It is not affiliated with or endorsed by Smogon, Nintendo, Creatures Inc., GAME FREAK inc. or The Pokémon Company. Pokémon and all related names are trademarks of their respective owners.

## License

[MIT](LICENSE)
