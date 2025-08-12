// internal/midi/player.go (Basic MIDI playback structure)
package midi

import (
	"strconv"
	"time"

	"github.com/Cod-e-Codes/tuitar/internal/models"
)

type Player struct {
	isPlaying    bool
	position     int
	tempo        int
	notes        []PlayableNote
	highlighted  []models.Position
	stopChan     chan bool
}

type PlayableNote struct {
	MidiNote  int
	Start     time.Duration
	Duration  time.Duration
	Velocity  int
	String    int
	Position  int
}

func NewPlayer() *Player {
	return &Player{
		tempo:    120,
		stopChan: make(chan bool),
	}
}

func (p *Player) PlayTab(tab *models.Tab) error {
	if p.isPlaying {
		return nil
	}
	
	p.notes = p.convertTabToNotes(tab)
	p.isPlaying = true
	p.position = 0
	
	go p.playbackLoop()
	
	return nil
}

func (p *Player) Stop() {
	if p.isPlaying {
		p.isPlaying = false
		p.stopChan <- true
		p.highlighted = nil
	}
}

func (p *Player) IsPlaying() bool {
	return p.isPlaying
}

func (p *Player) GetHighlighted() []models.Position {
	return p.highlighted
}

func (p *Player) convertTabToNotes(tab *models.Tab) []PlayableNote {
	var notes []PlayableNote
	
	// Standard guitar tuning MIDI notes (low to high)
	stringMidiNotes := [6]int{40, 45, 50, 55, 59, 64} // E A D G B e
	
	maxLength := 0
	for _, line := range tab.Content {
		if len(line) > maxLength {
			maxLength = len(line)
		}
	}
	
	beatDuration := time.Minute / time.Duration(tab.Tempo*4) // 16th notes
	
	for pos := 0; pos < maxLength; pos++ {
		for stringIdx, line := range tab.Content {
			if pos < len(line) && line[pos] != '-' && line[pos] != '|' {
				if fret, err := strconv.Atoi(string(line[pos])); err == nil {
					midiNote := stringMidiNotes[stringIdx] + fret
					
					note := PlayableNote{
						MidiNote: midiNote,
						Start:    time.Duration(pos) * beatDuration,
						Duration: beatDuration * 3 / 4, // Note length
						Velocity: 127,
						String:   stringIdx,
						Position: pos,
					}
					notes = append(notes, note)
				}
			}
		}
	}
	
	return notes
}

func (p *Player) playbackLoop() {
	defer func() {
		p.isPlaying = false
		p.highlighted = nil
	}()
	
	beatDuration := time.Minute / time.Duration(p.tempo*4)
	ticker := time.NewTicker(beatDuration)
	defer ticker.Stop()
	
	for {
		select {
		case <-p.stopChan:
			return
		case <-ticker.C:
			// Update highlighted positions
			p.highlighted = nil
			for _, note := range p.notes {
				if note.Position == p.position {
					p.highlighted = append(p.highlighted, models.Position{
						String:   note.String,
						Position: note.Position,
					})
				}
			}
			
			p.position++
			
			// Check if we've reached the end
			maxPos := 0
			for _, note := range p.notes {
				if note.Position > maxPos {
					maxPos = note.Position
				}
			}
			
			if p.position > maxPos {
				return
			}
		}
	}
}
