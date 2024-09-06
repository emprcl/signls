# cykl
:loop: a generative midi sequencer in the terminal

 - Represent port with quarter blocks characters
 - Each node can have 2 behaviors :
 	- An emitter behavior
 	- A trigger behavior
 	- A moving behavior ?

- can't represent 8 directions nicely in unicode on one character
 - focus on teleport and zone

TODO:
 - signal brings context
   - store initial bang position to respawn?
   - note relative to the grid !!!!
 - add random note octave spread
 - add random based on perlin noise
 - add lfo on every param
 - add randomness on every note param
 - add note probability
 - add keyboard mode for note input
 - add new notes behaviors:
   - note relative to the previous one (set true intervals? or relative to scale?)
   - note relative to the grid
 - add new emitters
  - Accumulation (trigger after x)
  - Opposite direction
  - Zone for chord
  - Cell for cellular automata rules ? (meh unsure maybe a separate project)
 - add midi cc (note behavior?)
 - add commands (tempo, scale, root, others?)
 - add theme
 - add key mapping
 - add file system

Unsure:
 - signals with different speeds

Names:
  - Nyll
  - KYVL
  - HYDR
  - VYRN
