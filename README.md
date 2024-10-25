# Signls

![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/emprcl/signls) ![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/emprcl/signls/build.yml)

Signls (_pronounced signals_) is a non-linear, generative midi sequencer designed for music composition and live performance. It's cross-platform (Linux, macOS et Windows) and it runs in the terminal (TUI).

It takes inspiration from [Orca](https://100r.co/site/orca.html) and [Nodal](https://nodalmusic.com/).

**_Signls is still in development, but it is expected to be relatively stable._**

:notebook: **_[User Manual](https://empr.cl/signls/)_**

_Feel free to [open an issue](https://github.com/emprcl/signls/issues/new)._

![signls screenshot](/docs/screenshot.png)

## Installation

[Download the last release](https://github.com/emprcl/signls/releases) for your platform.

Then:
```sh
# Extract files
mkdir -p signls && tar -zxvf signls_VERSION_PLATFORM.tar.gz -C signls

# Run signls
./signls
```

### Build it yourself

You'll need [go 1.23](https://go.dev/dl/) minimum.
Although you should be able to build it for either **linux**, **macOS** or **Windows**, it has only been tested on **linux**.

```sh
# Linux
sudo apt-get install libasound2-dev
make GOLANG_OS=linux build

# macOS
make GOLANG_OS=darwin build

# Windows
make GOLANG_OS=windows build

# Raspberry Pi OS
sudo apt install libasound2-dev
make GOLANG_OS=linux GOLANG_ARCH=arm64 build
```

## Usage

```sh
# Run signls
./signls

# Display current version
./signls --version
```

Hit `?` to see all keybindings. `esc` to quit.

Some companion apps that receive midi for testing signls:
 - [Enfer](https://neauoire.github.io/Enfer/) ([github](https://github.com/neauoire/Enfer))
 - [QSynth](https://qsynth.sourceforge.io/)

### Keyboard mapping

Keys mapping is fully customizable. After running signls for the first time, a `config.json` is created.
You can edit all the keys inside it.

You can select one of the default keyboard layouts available:
```sh
# QWERTY
./signls --keyboard qwerty

# AZERTY
./signls --keyboard azerty

# QWERTY MAC
./signls --keyboard qwerty-mac

# AZERTY MAC
./signls --keyboard azerty-mac
```

### Default keyboard mapping

For qwerty keyboards, here's the default mapping:

 - `space` **play** or **stop**
 - `tab` **show bank**
 - `1` ... `9` **add nodes**
 - `↑` `↓` `←` `→` **move cursor**
 - `shift`+`↑` `↓` `←` `→` **multiple selection (or modify alt parameter mode in edit mode)**
 - `ctrl`+`↑` `↓` `←` `→` **modify selected node direction (modify parameter or alt parameter value)**
 - `backspace` **remove selected nodes (or grid in bank)**
 - `enter` **edit selected nodes**
 - `m` **toggle selected nodes mute**
 - `M` **mute/unmute all selected nodes**
 - `/` **trigger selected node**
 - `-` `=` **modify tempo**
 - `'` `;` **modify root note**
 - `"` `:` **modify scale**
 - `ctrl`+`c` `x` `v`  **copy, cut, paste selection**
 - `escape` **exit parameter edit or bank selection**
 - `f2` **select midi device**
 - `f10` **fit grid to window**
 - `?` **show help**
 - `ctrl`+`q` **quit**

### Bank management

Each time you start Signls, a json file (default: `default.json`) containing 32 grid slots is loaded.
For selecting a different file, use the `--bank` flag:
```sh
./signls --bank my-grids.json
```

Each time you change grid or quit the program, the current grid is saved to the file.

## Acknowledgments

Signls uses a few awesome packages:
 - [gomidi/midi](https://gitlab.com/gomidi/midi) for all midi communication
 - [charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea) as the main TUI framework
 - [charmbracelet/lipgloss](https://github.com/charmbracelet/lipgloss) for making things beautiful
