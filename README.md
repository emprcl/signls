# cykl
:loop: a generative midi sequencer in the terminal

 - Represent port with quarter blocks characters
 - Each node can have 2 behaviors :
 	- An emitter behavior
 	- A trigger behavior
 	- A moving behavior ?

TODO:
 - signal brings context
   - store initial bang position to respawn?
   - note relative to the grid !!!!
 - add random note octave spread
 - add randomness on every note param
 - add note probability
 - add new notes behaviors:
   - note relative to the previous one (set true intervals? or relative to scale?)
   - note relative to the grid
 - add new emitters
  - Random (toss or dice)
  - Accumulation (trigger after x)
  - Teleport !!Does not trigger not, jump and triggers target right away
  - Zone for chord
  - Cell for cellular automata rules ? (meh unsure maybe a separate project)
 - add midi cc (note behavior?)
 - add key mapping
 - add file system

Unsure:
 - signals with different speeds

Names:
  - Nyll
  - KYVL
  - HYDR
  - VYRN
