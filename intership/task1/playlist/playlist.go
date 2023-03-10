package playlist

import (
    "sync"
    "time"
)

type Song struct {
    Name     string
    Duration time.Duration
}

type node struct {
    song *Song
    next *node
    prev *node
}

type Playlist struct {
    head   *node
    tail   *node
    current *node
    mutex  sync.Mutex
}

func (p *Playlist) AddSong(song *Song) {
    p.mutex.Lock()
    defer p.mutex.Unlock()

    newNode := &node{song: song}
    if p.head == nil {
        p.head = newNode
        p.tail = newNode
    } else {
        newNode.prev = p.tail
        p.tail.next = newNode
        p.tail = newNode
    }
}

func (p *Playlist) Play() {
    p.mutex.Lock()
    defer p.mutex.Unlock()

    if p.current != nil {
        return
    }

    p.current = p.head
    for p.current != nil {
        go func(current *node) {
            time.Sleep(current.song.Duration)
            if current == p.current {
                p.current = p.current.next
                p.Play()
            }
        }(p.current)

        return
    }
}

func (p *Playlist) Pause() {
    p.mutex.Lock()
    defer p.mutex.Unlock()

    if p.current == nil {
        return
    }

    p.current = nil
}

func (p *Playlist) Next() {
    p.mutex.Lock()
    defer p.mutex.Unlock()

    if p.current == nil {
        return
    }

    p.current = p.current.next
    p.Pause()
    p.Play()
}

func (p *Playlist) Prev() {
    p.mutex.Lock()
    defer p.mutex.Unlock()

    if p.current == nil {
        return
    }

    if p.current.prev != nil {
        p.current = p.current.prev
    }
    p.Pause()
    p.Play()
}