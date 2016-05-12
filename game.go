package main

import (
  "image"
  _ "image/color"
  "fmt"
  "math"
  "math/cmplx"

  "azul3d.org/engine/gfx"
  "azul3d.org/engine/gfx/window"
  "azul3d.org/engine/keyboard"
  _ "azul3d.org/engine/mouse"
)

type p struct {

  x   float64
  y   float64
  dx  float64
  dy  float64

  a   bool

}

type frame struct {

  Image chan image.Image
  bounds image.Rectangle

}

func makep() *p {
  p := &p{
    x:  100,
    y:  400,
    dx: 0,
    dy: 0,
    a:  false,
  }
  return p
}

func gfxLoop(w window.Window, d gfx.Device) {
  p := makep()

  events := make(chan window.Event, 1)
  w.Notify(events, window.KeyboardButtonEvents)

  watcher := keyboard.NewWatcher()

  for {
    window.Poll(events, func(e window.Event) {
      switch ev := e.(type) {
      case keyboard.ButtonEvent:
        watcher.SetState(ev.Key, ev.State)
      }
    })

    fmt.Println(p.x, p.y, p.dx, p.dy, p.a)

    //Important shit: y coordinate's sign is FLIPPED here
    //negative y is UP and positive y is DOWN due to how the engine draws things

    //Accelerate left
    if watcher.Down(keyboard.ArrowLeft) && watcher.Up(keyboard.ArrowRight) {
      p.dx -= 1
    }

    //Accelerate right
    if watcher.Down(keyboard.ArrowRight) && watcher.Up(keyboard.ArrowLeft) {
      p.dx += 1
    }

    //Inertia after stopping
    if watcher.Up(keyboard.ArrowLeft) && watcher.Up(keyboard.ArrowRight) {
      delta := math.Log(math.Abs(p.dx)+1)/4
      if p.dx > 0 {
        p.dx -= delta
      } else
      if p.dx < 0 {
        p.dx += delta
      }
    }

    if watcher.Down(keyboard.ArrowLeft) && watcher.Down(keyboard.ArrowRight) {
      delta := math.Log(math.Abs(p.dx)+1)/4
      if p.dx > 0 {
        p.dx -= delta
      } else
      if p.dx < 0 {
        p.dx += delta
      }
    }

    //Jumping impulse
    if watcher.Down(keyboard.ArrowUp) && !p.a {
      p.dy -= 7
    }

    //If x velocity is low enough, it becomes zero
    if math.Abs(p.dx) < 0.1 {
      p.dx = 0
    }

    //Same with y velocity
    if math.Abs(p.dy) < 0.1 {
      p.dy = 0
    }

    //Maximum x velocity
    if math.Abs(p.dx) > 7 {
      if p.dx > 0 {
        p.dx = 7
      } else
      if p.dx < 0 {
        p.dx = -7
      }
    }

/*    if math.Abs(p.dy) > 7 {
      if p.dx > 0 {
        p.dx = 7
      } else
      if p.dx < 0 {
        p.dx = -7
      }
    }
*/

    //Move position
    p.x += p.dx
    p.dy += 0.3
    p.y += p.dy

    //Ground bounding box
    if p.y >= 400 {
      p.y = 400
      p.dy = 0
      p.a = false
    }

    if p.y < 400 {
      p.a = true
    }

    if p.x < -20 {
      p.x = float64(d.Bounds().Max.X)
    }

    if p.x > float64(d.Bounds().Max.X) {
      p.x = -20
    }

    d.Clear(d.Bounds(), gfx.Color{0, 0, 0, 0})
    x := int(p.x)
    y := int(p.y)

    dabs := int(cmplx.Abs(complex(p.dx, p.dy)))

    d.Clear(image.Rect(0, 400, d.Bounds().Max.X, d.Bounds().Max.Y), gfx.Color{0.7, 0.7, 0.7, 1})

    d.Clear(image.Rect(x-20, y-20, x, y), gfx.Color{1, 1, 1, 1})

    for i := 0; i < dabs; i++ {
      x := int(100+p.dx*float64(i))
      y := int(100+p.dy*float64(i))
      d.Clear(image.Rect(x-1, y-1, x, y), gfx.Color{1, 1, 1, 1})
    }

    d.Render()
  }
}

func main() {
  window.Run(gfxLoop, nil)
}
